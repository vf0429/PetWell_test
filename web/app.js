const API_BASE = "http://localhost:8090";
const SESSION_KEY = "petwell_merchant_session";
const GUIDE_SEEN_KEY = "petwell_merchant_guide_seen";
const LANG_KEY = "petwell_merchant_lang";

const I18N = {
  en: {
    auth_title: "Welcome back",
    auth_subtitle: "Use email, phone, or Google to sign in.",
    auth_google_btn: "Continue with Google",
    auth_or: "or continue with",
    auth_tab_email: "Email",
    auth_tab_phone: "Phone",
    auth_email_label: "Email",
    auth_email_placeholder: "owner@petwell.com",
    auth_password_label: "Password",
    auth_password_placeholder: "********",
    auth_email_btn: "Sign in with Email",
    auth_country_code_label: "Country code",
    auth_phone_label: "Phone",
    auth_phone_placeholder: "13800138000",
    auth_phone_btn: "Sign in with Phone",
    auth_demo_hint: "Demo: shop_admin@petwell.com / Shop123456, clinic_admin@petwell.com / Clinic123456",
    app_title: "Merchant Console",
    menu_dashboard: "Dashboard",
    menu_orders: "Orders / Appointments",
    menu_catalog: "Products / Visits",
    menu_customers: "Customers / Pet Profiles",
    menu_analytics: "Analytics",
    menu_settings: "Settings",
    mode_shop: "Shop",
    mode_clinic: "Clinic",
    logout_btn: "Logout",
    topbar_desc_shop: "Manage sales, inventory, customers, and Shopify sync in one place",
    topbar_desc_clinic: "Manage appointments, visits, prescriptions, and follow-ups in one place",
    dashboard_system_status: "System Status",
    dashboard_quick_actions: "Quick Actions",
    quick_orders: "Go to Orders/Appointments",
    quick_catalog: "Go to Products/Visits",
    quick_reload: "Refresh Data",
    orders_create_order: "Create Order",
    orders_customer_id: "Customer ID",
    orders_customer_id_placeholder: "cus_xxx",
    orders_amount_hkd: "Amount (HKD)",
    orders_status: "Status",
    orders_create_order_btn: "Create Order",
    orders_create_appt: "Create Appointment",
    orders_pet_id: "Pet ID",
    orders_pet_placeholder: "pet_001",
    orders_doctor_id: "Doctor ID",
    orders_doctor_placeholder: "doc_001",
    orders_scheduled_at: "Schedule Time (ISO)",
    orders_schedule_placeholder: "2026-03-01T10:00:00Z",
    orders_chief_complaint: "Chief Complaint",
    orders_complaint_placeholder: "vomiting",
    orders_create_appt_btn: "Create Appointment",
    catalog_create_product: "Create Product",
    catalog_product_name: "Product Name",
    catalog_sku: "SKU",
    catalog_stock: "Stock",
    catalog_price_hkd: "Price (HKD)",
    catalog_create_product_btn: "Create Product",
    catalog_create_visit: "Create Visit",
    catalog_pet_id: "Pet ID",
    catalog_doctor_id: "Doctor ID",
    catalog_appt_id: "Appointment ID",
    catalog_appt_placeholder: "Optional",
    catalog_status: "Status",
    catalog_create_visit_btn: "Create Visit",
    customers_create_customer: "Create Customer",
    customers_name: "Name",
    customers_phone: "Phone",
    customers_create_customer_btn: "Create Customer",
    customers_create_rx: "Create Prescription",
    customers_visit_id: "Visit ID",
    customers_medicine: "Medicine",
    customers_dosage: "Dosage",
    customers_frequency: "Frequency",
    customers_days: "Duration Days",
    customers_notes: "Notes",
    customers_create_rx_btn: "Create Prescription",
    analytics_title: "KPI JSON Preview (for integration)",
    settings_oauth_title: "Shopify OAuth Placeholder Test",
    settings_oauth_btn: "Send OAuth Placeholder Request",
    settings_webhook_title: "Shopify Webhook Placeholder Test",
    settings_webhook_btn: "Send Webhook Placeholder Event",
    settings_last_resp: "Last Response",
    settings_guide_title: "Onboarding Guide",
    settings_guide_desc: "Replay the function guide for each page.",
    settings_guide_btn: "Replay Guide",
    guide_skip: "Skip guide",
    guide_next: "Next",
    guide_get_started: "Get Started",
    not_signed_in: "Not signed in",
    dashboard_status_mode: "Mode",
    dashboard_status_backend: "Backend",
    dashboard_status_user: "Current User",
    dashboard_status_time: "Updated At",
    list_orders: "Orders",
    list_appointments: "Appointments",
    list_products: "Products",
    list_visits: "Visits",
    list_customers: "Customers",
    list_prescriptions: "Prescriptions",
    table_col_id: "ID",
    table_col_customer: "Customer",
    table_col_status: "Status",
    table_col_amount: "Amount",
    table_col_pet: "Pet",
    table_col_doctor: "Doctor",
    table_col_scheduled_at: "Scheduled At",
    table_col_name: "Name",
    table_col_sku: "SKU",
    table_col_stock: "Stock",
    table_col_price: "Price",
    table_col_checkin_at: "Check In",
    table_col_email: "Email",
    table_col_phone: "Phone",
    table_col_visit: "Visit",
    table_col_medicine: "Medicine",
    table_col_dosage: "Dosage",
    table_col_frequency: "Frequency",
    login_success: "Login successful. Opening console...",
    login_fail_email: "Email login failed",
    login_fail_phone: "Phone login failed",
    login_fail_google: "Google login failed",
    logout_success: "Signed out.",
    backend_fail_1: "Cannot connect to backend:",
    backend_fail_2: "Please run: cd backend && go run .",
    session_expired: "Session expired, please sign in again.",
  },
  zh: {
    auth_title: "欢迎回来",
    auth_subtitle: "使用邮箱、手机号或 Google 登录。",
    auth_google_btn: "使用 Google 登录",
    auth_or: "或使用以下方式",
    auth_tab_email: "邮箱",
    auth_tab_phone: "手机号",
    auth_email_label: "邮箱",
    auth_email_placeholder: "owner@petwell.com",
    auth_password_label: "密码",
    auth_password_placeholder: "********",
    auth_email_btn: "邮箱登录",
    auth_country_code_label: "区号",
    auth_phone_label: "手机号",
    auth_phone_placeholder: "13800138000",
    auth_phone_btn: "手机号登录",
    auth_demo_hint: "测试账号：shop_admin@petwell.com / Shop123456，clinic_admin@petwell.com / Clinic123456",
    app_title: "商家工作台",
    menu_dashboard: "Dashboard",
    menu_orders: "订单 / 预约",
    menu_catalog: "商品 / 就诊",
    menu_customers: "客户 / 宠物档案",
    menu_analytics: "分析",
    menu_settings: "设置",
    mode_shop: "商品商家",
    mode_clinic: "医院商家",
    logout_btn: "退出",
    topbar_desc_shop: "统一管理销售、库存、客户和 Shopify 同步",
    topbar_desc_clinic: "统一管理预约、就诊、处方与回访",
    dashboard_system_status: "系统状态",
    dashboard_quick_actions: "快速动作",
    quick_orders: "去订单/预约页",
    quick_catalog: "去商品/就诊页",
    quick_reload: "刷新数据",
    orders_create_order: "新增订单",
    orders_customer_id: "客户ID",
    orders_customer_id_placeholder: "cus_xxx",
    orders_amount_hkd: "金额(HKD)",
    orders_status: "状态",
    orders_create_order_btn: "创建订单",
    orders_create_appt: "新增预约",
    orders_pet_id: "宠物ID",
    orders_pet_placeholder: "pet_001",
    orders_doctor_id: "医生ID",
    orders_doctor_placeholder: "doc_001",
    orders_scheduled_at: "预约时间(ISO)",
    orders_schedule_placeholder: "2026-03-01T10:00:00Z",
    orders_chief_complaint: "主诉",
    orders_complaint_placeholder: "vomiting",
    orders_create_appt_btn: "创建预约",
    catalog_create_product: "新增商品",
    catalog_product_name: "商品名称",
    catalog_sku: "SKU",
    catalog_stock: "库存",
    catalog_price_hkd: "价格(HKD)",
    catalog_create_product_btn: "创建商品",
    catalog_create_visit: "新增就诊记录",
    catalog_pet_id: "宠物ID",
    catalog_doctor_id: "医生ID",
    catalog_appt_id: "关联预约ID",
    catalog_appt_placeholder: "可选",
    catalog_status: "状态",
    catalog_create_visit_btn: "创建就诊",
    customers_create_customer: "新增客户",
    customers_name: "姓名",
    customers_phone: "电话",
    customers_create_customer_btn: "创建客户",
    customers_create_rx: "新增处方",
    customers_visit_id: "就诊ID",
    customers_medicine: "药品",
    customers_dosage: "剂量",
    customers_frequency: "频次",
    customers_days: "天数",
    customers_notes: "备注",
    customers_create_rx_btn: "创建处方",
    analytics_title: "关键指标 JSON 预览（便于联调）",
    settings_oauth_title: "Shopify OAuth 占位测试",
    settings_oauth_btn: "发起 OAuth 占位请求",
    settings_webhook_title: "Shopify Webhook 占位测试",
    settings_webhook_btn: "发送 Webhook 占位事件",
    settings_last_resp: "最后响应",
    settings_guide_title: "新手引导",
    settings_guide_desc: "重新播放每一页功能介绍。",
    settings_guide_btn: "重新播放 Guide",
    guide_skip: "跳过 guide",
    guide_next: "下一步",
    guide_get_started: "开始使用",
    not_signed_in: "未登录",
    dashboard_status_mode: "模式",
    dashboard_status_backend: "后端",
    dashboard_status_user: "当前用户",
    dashboard_status_time: "数据更新时间",
    list_orders: "订单列表",
    list_appointments: "预约列表",
    list_products: "商品列表",
    list_visits: "就诊列表",
    list_customers: "客户列表",
    list_prescriptions: "处方列表",
    table_col_id: "ID",
    table_col_customer: "客户",
    table_col_status: "状态",
    table_col_amount: "金额",
    table_col_pet: "宠物",
    table_col_doctor: "医生",
    table_col_scheduled_at: "预约时间",
    table_col_name: "名称",
    table_col_sku: "SKU",
    table_col_stock: "库存",
    table_col_price: "价格",
    table_col_checkin_at: "到诊时间",
    table_col_email: "邮箱",
    table_col_phone: "电话",
    table_col_visit: "就诊",
    table_col_medicine: "药品",
    table_col_dosage: "剂量",
    table_col_frequency: "频次",
    login_success: "登录成功，正在进入控制台...",
    login_fail_email: "邮箱登录失败",
    login_fail_phone: "手机号登录失败",
    login_fail_google: "Google 登录失败",
    logout_success: "已退出登录",
    backend_fail_1: "无法连接后端：",
    backend_fail_2: "请先运行：cd backend && go run .",
    session_expired: "登录会话已过期，请重新登录。",
  },
};

const GUIDE_STEPS = [
  {
    section: "dashboard",
    selector: "#dashboard .kpi-grid",
    title: { en: "Dashboard Overview", zh: "Dashboard 总览" },
    desc: {
      en: "See key KPIs, system health, and quick actions first.",
      zh: "这里看核心指标、系统状态和快速动作，先快速了解业务健康度。",
    },
  },
  {
    section: "orders",
    selector: "#orders",
    title: { en: "Orders / Appointments", zh: "订单 / 预约" },
    desc: {
      en: "Shop mode handles orders; clinic mode handles appointments.",
      zh: "商品商家管理订单；医院商家管理预约。这里是日常操作高频入口。",
    },
  },
  {
    section: "catalog",
    selector: "#catalog",
    title: { en: "Products / Visits", zh: "商品 / 就诊" },
    desc: {
      en: "Maintain products/stock in shop mode, visits/status in clinic mode.",
      zh: "商品模式维护商品和库存；医院模式维护就诊记录和状态流转。",
    },
  },
  {
    section: "customers",
    selector: "#customers",
    title: { en: "Customers / Prescriptions", zh: "客户 / 处方" },
    desc: {
      en: "Manage customers in shop mode and prescriptions in clinic mode.",
      zh: "商品模式维护客户信息；医院模式维护处方与用药信息。",
    },
  },
  {
    section: "analytics",
    selector: "#analytics",
    title: { en: "Analytics", zh: "分析" },
    desc: {
      en: "Preview KPI JSON here for frontend-backend integration and debugging.",
      zh: "这里展示关键数据 JSON，便于前后端联调和接口排查。",
    },
  },
  {
    section: "settings",
    selector: "#settings",
    title: { en: "Settings", zh: "设置" },
    desc: {
      en: "Run Shopify OAuth/Webhook placeholder tests and replay onboarding guide.",
      zh: "这里做 Shopify OAuth / Webhook 测试，也可以重新播放新手引导。",
    },
  },
];

const state = {
  mode: "shop",
  section: "dashboard",
  lang: localStorage.getItem(LANG_KEY) || "en",
  session: null,
  guideIndex: 0,
  guideActive: false,
  guideHighlightEl: null,
};

function el(id) {
  return document.getElementById(id);
}

function t(key) {
  const dict = I18N[state.lang] || I18N.en;
  return dict[key] || I18N.en[key] || key;
}

function applyI18n() {
  document.documentElement.lang = state.lang === "zh" ? "zh-Hans" : "en";
  document.querySelectorAll("[data-i18n]").forEach((node) => {
    node.textContent = t(node.dataset.i18n);
  });
  document.querySelectorAll("[data-i18n-placeholder]").forEach((node) => {
    node.placeholder = t(node.dataset.i18nPlaceholder);
  });
  document.querySelectorAll(".lang-btn").forEach((btn) => {
    btn.classList.toggle("active", btn.dataset.lang === state.lang);
  });
}

function showMessage(msg, isError = false) {
  const node = el("auth-message");
  node.textContent = msg;
  node.style.color = isError ? "#dc2626" : "#16a34a";
}

function setAppVisible(isVisible) {
  el("auth-screen").classList.toggle("hidden", isVisible);
  el("app-shell").classList.toggle("hidden", !isVisible);
}

function formatUser(session) {
  if (!session || !session.user) return t("not_signed_in");
  const user = session.user;
  return `${user.display_name} (${user.method})`;
}

function setAuthTab(tab) {
  document.querySelectorAll(".auth-tab").forEach((btn) => {
    btn.classList.toggle("active", btn.dataset.authTab === tab);
  });
  el("email-login-form").classList.toggle("active", tab === "email");
  el("phone-login-form").classList.toggle("active", tab === "phone");
  showMessage("");
}

async function api(path, options = {}) {
  const headers = { "Content-Type": "application/json" };
  if (state.session?.session_id) {
    headers["X-Session-ID"] = state.session.session_id;
  }
  const res = await fetch(`${API_BASE}${path}`, {
    headers,
    ...options,
  });
  if (!res.ok) {
    const text = await res.text();
    if (res.status === 401 && state.session?.session_id) {
      clearSession();
      setAppVisible(false);
      showMessage(t("session_expired"), true);
    }
    throw new Error(`${res.status} ${text}`);
  }
  return res.json();
}

function formToObject(form) {
  return Object.fromEntries(new FormData(form).entries());
}

async function submitJSON(path, payload) {
  return api(path, { method: "POST", body: JSON.stringify(payload) });
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
      ? t("topbar_desc_shop")
      : t("topbar_desc_clinic");
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
    .map((row) => `<tr>${columns.map((c) => `<td>${row[c.key] ?? ""}</td>`).join("")}</tr>`)
    .join("");
}

async function renderDashboard() {
  const payload = await api(`/api/merchant/dashboard?mode=${state.mode}`);
  const kpis = payload.kpis || [];
  el("kpi-grid").innerHTML = kpis
    .map((k) => `<article class="kpi"><div class="label">${k.label}</div><div class="value">${k.value}</div></article>`)
    .join("");

  const statusItems = [
    `${t("dashboard_status_mode")}: ${state.mode}`,
    `${t("dashboard_status_backend")}: ${API_BASE}`,
    `${t("dashboard_status_user")}: ${formatUser(state.session)}`,
    `${t("dashboard_status_time")}: ${new Date().toLocaleString()}`,
  ];
  el("dashboard-status").innerHTML = statusItems.map((s) => `<li>${s}</li>`).join("");
  el("analytics-json").textContent = JSON.stringify(payload, null, 2);
}

async function renderOrdersOrAppointments() {
  if (state.mode === "shop") {
    const rows = await api("/api/merchant/orders");
    el("orders-list-title").textContent = t("list_orders");
    fillTable(
      "orders-table-head",
      "orders-table-body",
      [
        { key: "id", title: t("table_col_id") },
        { key: "customer_id", title: t("table_col_customer") },
        { key: "status", title: t("table_col_status") },
        { key: "amount_hkd", title: t("table_col_amount") },
      ],
      rows
    );
    return;
  }

  const rows = await api("/api/merchant/appointments");
  el("orders-list-title").textContent = t("list_appointments");
  fillTable(
    "orders-table-head",
    "orders-table-body",
    [
      { key: "id", title: t("table_col_id") },
      { key: "pet_id", title: t("table_col_pet") },
      { key: "doctor_id", title: t("table_col_doctor") },
      { key: "status", title: t("table_col_status") },
      { key: "scheduled_at", title: t("table_col_scheduled_at") },
    ],
    rows
  );
}

async function renderCatalogOrVisits() {
  if (state.mode === "shop") {
    const rows = await api("/api/merchant/products");
    el("catalog-list-title").textContent = t("list_products");
    fillTable(
      "catalog-table-head",
      "catalog-table-body",
      [
        { key: "id", title: t("table_col_id") },
        { key: "name", title: t("table_col_name") },
        { key: "sku", title: t("table_col_sku") },
        { key: "stock", title: t("table_col_stock") },
        { key: "price_hkd", title: t("table_col_price") },
      ],
      rows
    );
    return;
  }

  const rows = await api("/api/merchant/visits");
  el("catalog-list-title").textContent = t("list_visits");
  fillTable(
    "catalog-table-head",
    "catalog-table-body",
    [
      { key: "id", title: t("table_col_id") },
      { key: "pet_id", title: t("table_col_pet") },
      { key: "doctor_id", title: t("table_col_doctor") },
      { key: "status", title: t("table_col_status") },
      { key: "checkin_at", title: t("table_col_checkin_at") },
    ],
    rows
  );
}

async function renderCustomersOrPrescriptions() {
  if (state.mode === "shop") {
    const rows = await api("/api/merchant/customers");
    el("customers-list-title").textContent = t("list_customers");
    fillTable(
      "customers-table-head",
      "customers-table-body",
      [
        { key: "id", title: t("table_col_id") },
        { key: "name", title: t("table_col_name") },
        { key: "email", title: t("table_col_email") },
        { key: "phone", title: t("table_col_phone") },
      ],
      rows
    );
    return;
  }

  const rows = await api("/api/merchant/prescriptions");
  el("customers-list-title").textContent = t("list_prescriptions");
  fillTable(
    "customers-table-head",
    "customers-table-body",
    [
      { key: "id", title: t("table_col_id") },
      { key: "visit_id", title: t("table_col_visit") },
      { key: "medicine_name", title: t("table_col_medicine") },
      { key: "dosage", title: t("table_col_dosage") },
      { key: "frequency", title: t("table_col_frequency") },
    ],
    rows
  );
}

async function renderCurrentSection() {
  await renderDashboard();
  await renderOrdersOrAppointments();
  await renderCatalogOrVisits();
  await renderCustomersOrPrescriptions();
}

function clearGuideHighlight() {
  if (!state.guideHighlightEl) return;
  state.guideHighlightEl.classList.remove("guide-focus");
  state.guideHighlightEl = null;
}

function finishGuide() {
  state.guideActive = false;
  clearGuideHighlight();
  el("guide-backdrop").classList.add("hidden");
  el("guide-panel").classList.add("hidden");
  localStorage.setItem(GUIDE_SEEN_KEY, "1");
}

function renderGuideStep() {
  const step = GUIDE_STEPS[state.guideIndex];
  if (!step) {
    finishGuide();
    return;
  }

  setSectionUI(step.section);
  clearGuideHighlight();

  const target = document.querySelector(step.selector);
  if (target) {
    state.guideHighlightEl = target;
    state.guideHighlightEl.classList.add("guide-focus");
    state.guideHighlightEl.scrollIntoView({ behavior: "smooth", block: "center" });
  }

  el("guide-progress").textContent = `${state.guideIndex + 1} / ${GUIDE_STEPS.length}`;
  el("guide-title").textContent = step.title[state.lang] || step.title.en;
  el("guide-desc").textContent = step.desc[state.lang] || step.desc.en;
  el("next-guide-btn").textContent = state.guideIndex === GUIDE_STEPS.length - 1 ? t("guide_get_started") : t("guide_next");
}

function startGuide() {
  state.guideActive = true;
  state.guideIndex = 0;
  el("guide-backdrop").classList.remove("hidden");
  el("guide-panel").classList.remove("hidden");
  renderGuideStep();
}

function saveSession(session) {
  state.session = session;
  localStorage.setItem(SESSION_KEY, JSON.stringify(session));
  el("current-user").textContent = formatUser(session);
}

function clearSession() {
  state.session = null;
  localStorage.removeItem(SESSION_KEY);
  el("current-user").textContent = t("not_signed_in");
}

function bindAuth() {
  document.querySelectorAll(".auth-tab").forEach((btn) => {
    btn.addEventListener("click", () => setAuthTab(btn.dataset.authTab));
  });

  el("email-login-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = formToObject(e.target);
    try {
      const session = await submitJSON("/api/merchant/auth/login", {
        method: "email",
        email: payload.email,
        password: payload.password,
      });
      saveSession(session);
      showMessage(t("login_success"));
      await enterApp({ playGuide: true });
    } catch (err) {
      showMessage(`${t("login_fail_email")}: ${err.message}`, true);
    }
  });

  el("phone-login-form").addEventListener("submit", async (e) => {
    e.preventDefault();
    const payload = formToObject(e.target);
    try {
      const session = await submitJSON("/api/merchant/auth/login", {
        method: "phone",
        country_code: payload.country_code,
        phone: payload.phone,
        password: payload.password,
      });
      saveSession(session);
      showMessage(t("login_success"));
      await enterApp({ playGuide: true });
    } catch (err) {
      showMessage(`${t("login_fail_phone")}: ${err.message}`, true);
    }
  });

  el("google-login-btn").addEventListener("click", async () => {
    try {
      const session = await submitJSON("/api/merchant/auth/login", {
        method: "google",
        email: "google_shop@petwell.com",
      });
      saveSession(session);
      showMessage(t("login_success"));
      await enterApp({ playGuide: true });
    } catch (err) {
      showMessage(`${t("login_fail_google")}: ${err.message}`, true);
    }
  });
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
  document.querySelectorAll(".lang-btn").forEach((btn) => {
    btn.addEventListener("click", async () => {
      state.lang = btn.dataset.lang;
      localStorage.setItem(LANG_KEY, state.lang);
      applyI18n();
      setModeUI();
      if (state.session) {
        el("current-user").textContent = formatUser(state.session);
        await renderCurrentSection();
        if (state.guideActive) renderGuideStep();
      }
    });
  });

  document.querySelectorAll(".menu-item").forEach((btn) => {
    btn.addEventListener("click", () => setSectionUI(btn.dataset.section));
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

  el("logout-btn").addEventListener("click", () => {
    finishGuide();
    clearSession();
    setAppVisible(false);
    setAuthTab("email");
    showMessage(t("logout_success"));
  });

  el("replay-guide-btn").addEventListener("click", () => {
    localStorage.removeItem(GUIDE_SEEN_KEY);
    startGuide();
  });

  el("next-guide-btn").addEventListener("click", () => {
    if (!state.guideActive) return;
    state.guideIndex += 1;
    renderGuideStep();
  });

  el("skip-guide-btn").addEventListener("click", () => finishGuide());
}

async function enterApp({ playGuide }) {
  setAppVisible(true);
  setModeUI();
  setSectionUI("dashboard");

  try {
    await renderCurrentSection();
  } catch (err) {
    console.error(err);
    el("dashboard-status").innerHTML = `<li>${t("backend_fail_1")} ${err.message}</li><li>${t("backend_fail_2")}</li>`;
  }

  const seen = localStorage.getItem(GUIDE_SEEN_KEY) === "1";
  if (playGuide && !seen) startGuide();
}

function restoreSession() {
  const raw = localStorage.getItem(SESSION_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw);
  } catch {
    return null;
  }
}

async function init() {
  if (!I18N[state.lang]) state.lang = "en";
  applyI18n();
  bindAuth();
  bindForms();
  bindNavigation();
  setAuthTab("email");
  setModeUI();
  setSectionUI("dashboard");
  el("current-user").textContent = t("not_signed_in");

  const saved = restoreSession();
  if (!saved) {
    setAppVisible(false);
    return;
  }

  saveSession(saved);
  await enterApp({ playGuide: false });
}

init();
