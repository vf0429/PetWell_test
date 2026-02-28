const API_BASE = "http://localhost:8090";

const state = {
  mode: "shop",
  section: "dashboard",
};

function el(id) {
  return document.getElementById(id);
}

async function api(path, options = {}) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`${res.status} ${text}`);
  }
  return res.json();
}

function setModeUI() {
  document.querySelectorAll(".mode-btn").forEach((btn) => {
    btn.classList.toggle("active", btn.dataset.mode === state.mode);
  });
  document.querySelectorAll("[data-mode-only]").forEach((node) => {
    node.style.display = node.dataset.modeOnly === state.mode ? "block" : "none";
  });
  el("topbar-desc").textContent =
    state.mode === "shop"
      ? "统一管理销售、库存、客户和 Shopify 同步"
      : "统一管理预约、就诊、处方与回访";
}

function setSectionUI(section) {
  state.section = section;
  document.querySelectorAll(".menu-item").forEach((btn) => {
    btn.classList.toggle("active", btn.dataset.section === section);
  });
  document.querySelectorAll(".panel").forEach((panel) => {
    panel.classList.toggle("active", panel.id === section);
  });
}

function fillTable(headId, bodyId, columns, rows) {
  el(headId).innerHTML = `<tr>${columns.map((c) => `<th>${c.title}</th>`).join("")}</tr>`;
  el(bodyId).innerHTML = rows
    .map(
      (row) =>
        `<tr>${columns
          .map((c) => `<td>${row[c.key] ?? ""}</td>`)
          .join("")}</tr>`
    )
    .join("");
}

async function renderDashboard() {
  const payload = await api(`/api/merchant/dashboard?mode=${state.mode}`);
  const kpis = payload.kpis || [];

  el("kpi-grid").innerHTML = kpis
    .map(
      (k) => `<article class="kpi"><div class="label">${k.label}</div><div class="value">${k.value}</div></article>`
    )
    .join("");

  const statusItems = [];
  statusItems.push(`模式：${state.mode}`);
  statusItems.push(`后端：${API_BASE}`);
  statusItems.push(`数据更新时间：${new Date().toLocaleString()}`);
  el("dashboard-status").innerHTML = statusItems.map((s) => `<li>${s}</li>`).join("");
  el("analytics-json").textContent = JSON.stringify(payload, null, 2);
}

async function renderOrdersOrAppointments() {
  if (state.mode === "shop") {
    const rows = await api("/api/merchant/orders");
    el("orders-list-title").textContent = "订单列表";
    fillTable(
      "orders-table-head",
      "orders-table-body",
      [
        { key: "id", title: "ID" },
        { key: "customer_id", title: "Customer" },
        { key: "status", title: "Status" },
        { key: "amount_hkd", title: "Amount" },
      ],
      rows
    );
  } else {
    const rows = await api("/api/merchant/appointments");
    el("orders-list-title").textContent = "预约列表";
    fillTable(
      "orders-table-head",
      "orders-table-body",
      [
        { key: "id", title: "ID" },
        { key: "pet_id", title: "Pet" },
        { key: "doctor_id", title: "Doctor" },
        { key: "status", title: "Status" },
        { key: "scheduled_at", title: "ScheduledAt" },
      ],
      rows
    );
  }
}

async function renderCatalogOrVisits() {
  if (state.mode === "shop") {
    const rows = await api("/api/merchant/products");
    el("catalog-list-title").textContent = "商品列表";
    fillTable(
      "catalog-table-head",
      "catalog-table-body",
      [
        { key: "id", title: "ID" },
        { key: "name", title: "Name" },
        { key: "sku", title: "SKU" },
        { key: "stock", title: "Stock" },
        { key: "price_hkd", title: "Price" },
      ],
      rows
    );
  } else {
    const rows = await api("/api/merchant/visits");
    el("catalog-list-title").textContent = "就诊列表";
    fillTable(
      "catalog-table-head",
      "catalog-table-body",
      [
        { key: "id", title: "ID" },
        { key: "pet_id", title: "Pet" },
        { key: "doctor_id", title: "Doctor" },
        { key: "status", title: "Status" },
        { key: "checkin_at", title: "CheckIn" },
      ],
      rows
    );
  }
}

async function renderCustomersOrPrescriptions() {
  if (state.mode === "shop") {
    const rows = await api("/api/merchant/customers");
    el("customers-list-title").textContent = "客户列表";
    fillTable(
      "customers-table-head",
      "customers-table-body",
      [
        { key: "id", title: "ID" },
        { key: "name", title: "Name" },
        { key: "email", title: "Email" },
        { key: "phone", title: "Phone" },
      ],
      rows
    );
  } else {
    const rows = await api("/api/merchant/prescriptions");
    el("customers-list-title").textContent = "处方列表";
    fillTable(
      "customers-table-head",
      "customers-table-body",
      [
        { key: "id", title: "ID" },
        { key: "visit_id", title: "Visit" },
        { key: "medicine_name", title: "Medicine" },
        { key: "dosage", title: "Dosage" },
        { key: "frequency", title: "Frequency" },
      ],
      rows
    );
  }
}

async function renderCurrentSection() {
  await renderDashboard();
  await renderOrdersOrAppointments();
  await renderCatalogOrVisits();
  await renderCustomersOrPrescriptions();
}

function formToObject(form) {
  return Object.fromEntries(new FormData(form).entries());
}

async function submitJSON(path, payload) {
  const resp = await api(path, { method: "POST", body: JSON.stringify(payload) });
  return resp;
}

function bindForms() {
  el("order-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = formToObject(e.target);
    payload.tenant_id = "tenant_demo";
    payload.amount_hkd = Number(payload.amount_hkd);
    await submitJSON("/api/merchant/orders", payload);
    e.target.reset();
    await renderCurrentSection();
  });

  el("appointment-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = formToObject(e.target);
    payload.tenant_id = "tenant_demo";
    await submitJSON("/api/merchant/appointments", payload);
    e.target.reset();
    await renderCurrentSection();
  });

  el("product-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = formToObject(e.target);
    payload.tenant_id = "tenant_demo";
    payload.stock = Number(payload.stock);
    payload.price_hkd = Number(payload.price_hkd);
    await submitJSON("/api/merchant/products", payload);
    e.target.reset();
    await renderCurrentSection();
  });

  el("visit-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = formToObject(e.target);
    payload.tenant_id = "tenant_demo";
    await submitJSON("/api/merchant/visits", payload);
    e.target.reset();
    await renderCurrentSection();
  });

  el("customer-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = formToObject(e.target);
    payload.tenant_id = "tenant_demo";
    await submitJSON("/api/merchant/customers", payload);
    e.target.reset();
    await renderCurrentSection();
  });

  el("prescription-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = formToObject(e.target);
    payload.duration_days = Number(payload.duration_days);
    await submitJSON("/api/merchant/prescriptions", payload);
    e.target.reset();
    await renderCurrentSection();
  });

  el("shopify-oauth-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = formToObject(e.target);
    const resp = await submitJSON("/api/merchant/shopify/oauth/start", payload);
    el("settings-response").textContent = JSON.stringify(resp, null, 2);
  });

  el("shopify-webhook-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = formToObject(e.target);
    const resp = await submitJSON("/api/merchant/shopify/webhooks", {
      topic: payload.topic,
      payload: { source: "manual_test" },
    });
    el("settings-response").textContent = JSON.stringify(resp, null, 2);
  });
}

function bindNavigation() {
  document.querySelectorAll(".menu-item").forEach((btn) => {
    btn.addEventListener("click", () => {
      setSectionUI(btn.dataset.section);
    });
  });

  document.querySelectorAll(".mode-btn").forEach((btn) => {
    btn.addEventListener("click", async () => {
      state.mode = btn.dataset.mode;
      setModeUI();
      await renderCurrentSection();
    });
  });

  document.querySelectorAll(".quick-actions .btn").forEach((btn) => {
    btn.addEventListener("click", async () => {
      const action = btn.dataset.action;
      if (action === "goto-orders") setSectionUI("orders");
      if (action === "goto-catalog") setSectionUI("catalog");
      if (action === "reload") await renderCurrentSection();
    });
  });
}

async function init() {
  bindNavigation();
  bindForms();
  setModeUI();
  setSectionUI("dashboard");

  try {
    await renderCurrentSection();
  } catch (err) {
    console.error(err);
    el("dashboard-status").innerHTML = `<li>无法连接后端：${err.message}</li><li>请先运行：cd backend && go run .</li>`;
  }
}

init();
