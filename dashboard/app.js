const API_BASE_URL = "http://localhost:8000/api";

function formatCurrency(value) {
    const numberValue = Number(value || 0);

    return new Intl.NumberFormat("id-ID", {
        style: "currency",
        currency: "IDR",
        maximumFractionDigits: 0,
    }).format(numberValue);
}

function formatNumber(value) {
    return new Intl.NumberFormat("id-ID").format(Number(value || 0));
}

async function fetchJson(endpoint) {
    const response = await fetch(`${API_BASE_URL}${endpoint}`);

    if (!response.ok) {
        throw new Error(`Failed to fetch ${endpoint}`);
    }

    return response.json();
}

async function loadSummary() {
    const data = await fetchJson("/summary");

    document.getElementById("totalRevenue").textContent = formatCurrency(data.total_revenue);
    document.getElementById("totalTransactions").textContent = formatNumber(data.total_transactions);
    document.getElementById("totalCustomers").textContent = formatNumber(data.total_customers);
    document.getElementById("averageBasketSize").textContent = formatCurrency(data.average_basket_size);
}

function createBarChart(canvasId, labels, values, label) {
    const context = document.getElementById(canvasId);

    new Chart(context, {
        type: "bar",
        data: {
            labels,
            datasets: [
                {
                    label,
                    data: values,
                },
            ],
        },
        options: {
            responsive: true,
            plugins: {
                legend: {
                    display: false,
                },
            },
            scales: {
                y: {
                    beginAtZero: true,
                },
            },
        },
    });
}

function createLineChart(canvasId, labels, values, label) {
    const context = document.getElementById(canvasId);

    new Chart(context, {
        type: "line",
        data: {
            labels,
            datasets: [
                {
                    label,
                    data: values,
                    tension: 0.3,
                },
            ],
        },
        options: {
            responsive: true,
            plugins: {
                legend: {
                    display: false,
                },
            },
            scales: {
                y: {
                    beginAtZero: true,
                },
            },
        },
    });
}

async function loadRevenueByStore() {
    const data = await fetchJson("/revenue-by-store");

    const topStores = data.slice(0, 10);

    createBarChart(
        "revenueByStoreChart",
        topStores.map((row) => `${row.store_name} (${row.tenant_id})`),
        topStores.map((row) => Number(row.total_revenue)),
        "Revenue"
    );
}

async function loadPromotionEffectiveness() {
    const data = await fetchJson("/promotion-effectiveness");

    const topPromotions = data.slice(0, 10);

    createBarChart(
        "promotionEffectivenessChart",
        topPromotions.map((row) => `${row.promo_name} (${row.tenant_id})`),
        topPromotions.map((row) => Number(row.promoted_revenue)),
        "Promoted Revenue"
    );
}

async function loadTopProductsByCity() {
    const data = await fetchJson("/top-products-by-city");

    const topProducts = data.slice(0, 10);

    createBarChart(
        "topProductsChart",
        topProducts.map((row) => `${row.product_name} - ${row.city}`),
        topProducts.map((row) => Number(row.total_revenue)),
        "Revenue"
    );
}

async function loadCustomerSegments() {
    const data = await fetchJson("/customer-segments");

    createBarChart(
        "customerSegmentsChart",
        data.map((row) => `${row.customer_segment} (${row.tenant_id})`),
        data.map((row) => Number(row.segment_revenue)),
        "Segment Revenue"
    );
}

async function loadTransactionsByDay() {
    const data = await fetchJson("/transactions-by-day");

    createLineChart(
        "transactionsByDayChart",
        data.map((row) => row.day_name),
        data.map((row) => Number(row.total_transactions)),
        "Transactions"
    );
}

async function main() {
    try {
        await loadSummary();
        await loadRevenueByStore();
        await loadPromotionEffectiveness();
        await loadTopProductsByCity();
        await loadCustomerSegments();
        await loadTransactionsByDay();
    } catch (error) {
        console.error(error);
        alert("Failed to load dashboard data. Please check the FastAPI service.");
    }
}

main();