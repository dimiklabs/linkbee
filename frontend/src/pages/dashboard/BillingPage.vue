<template>
  <div class="billing-page">
    <div class="page-header mb-4">
      <h1 class="page-title">Billing &amp; Plan</h1>
      <p class="page-subtitle">Manage your subscription and track usage.</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading…</span>
      </div>
    </div>

    <div v-else class="billing-grid">

      <!-- ── Current Plan ───────────────────────────────────────────────── -->
      <div class="card plan-card">
        <div class="card-body">
          <div class="d-flex align-items-start justify-content-between mb-3">
            <div>
              <p class="section-label mb-1">Current plan</p>
              <h2 class="plan-name mb-1">{{ planLabel }}</h2>
            </div>
            <span class="status-badge" :class="statusClass">{{ statusLabel }}</span>
          </div>

          <div v-if="sub?.current_period_end" class="meta-row">
            <span class="meta-label">Renews</span>
            <span class="meta-value">{{ formatDate(sub.current_period_end) }}</span>
          </div>
          <div v-if="sub?.cancelled_at" class="meta-row">
            <span class="meta-label">Cancelled</span>
            <span class="meta-value text-danger">{{ formatDate(sub.cancelled_at) }}</span>
          </div>

          <!-- Upgrade CTAs -->
          <div class="mt-4 d-flex flex-wrap gap-2">
            <button
              v-if="currentPlanID === 'free'"
              class="btn btn-primary"
              :disabled="checkoutLoading === 'pro'"
              @click="goCheckout('pro')"
            >
              <span v-if="checkoutLoading === 'pro'" class="spinner-border spinner-border-sm me-2"></span>
              Upgrade to Pro
            </button>
            <button
              v-if="currentPlanID !== 'business'"
              class="btn"
              :class="currentPlanID === 'free' ? 'btn-outline-primary' : 'btn-primary'"
              :disabled="checkoutLoading === 'business'"
              @click="goCheckout('business')"
            >
              <span v-if="checkoutLoading === 'business'" class="spinner-border spinner-border-sm me-2"></span>
              Upgrade to Business
            </button>
            <span v-if="currentPlanID === 'business'" class="text-success fw-medium">
              You're on the highest plan
            </span>
          </div>

          <p v-if="checkoutError" class="text-danger small mt-2 mb-0">{{ checkoutError }}</p>
        </div>
      </div>

      <!-- ── Usage Meters ───────────────────────────────────────────────── -->
      <div class="card usage-card">
        <div class="card-body">
          <p class="section-label mb-3">Usage this period</p>

          <div class="usage-list">
            <!-- Links -->
            <div class="usage-item">
              <div class="usage-header">
                <span class="usage-resource">Links</span>
                <span class="usage-count">
                  {{ usage?.links ?? 0 }}
                  <span class="usage-limit">/ {{ limitLabel(plan?.max_links) }}</span>
                </span>
              </div>
              <div class="progress usage-bar" :title="`${usage?.links ?? 0} of ${plan?.max_links === -1 ? '∞' : plan?.max_links} links`">
                <div
                  class="progress-bar"
                  :class="barClass(usage?.links ?? 0, plan?.max_links ?? 0)"
                  :style="{ width: pct(usage?.links ?? 0, plan?.max_links ?? 0) + '%' }"
                ></div>
              </div>
            </div>

            <!-- API Keys -->
            <div class="usage-item">
              <div class="usage-header">
                <span class="usage-resource">API Keys</span>
                <span class="usage-count">
                  {{ usage?.api_keys ?? 0 }}
                  <span class="usage-limit">/ {{ limitLabel(plan?.max_api_keys) }}</span>
                </span>
              </div>
              <div class="progress usage-bar">
                <div
                  class="progress-bar"
                  :class="barClass(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0)"
                  :style="{ width: pct(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0) + '%' }"
                ></div>
              </div>
            </div>

            <!-- Webhooks -->
            <div class="usage-item">
              <div class="usage-header">
                <span class="usage-resource">Webhooks</span>
                <span class="usage-count">
                  {{ usage?.webhooks ?? 0 }}
                  <span class="usage-limit">/ {{ limitLabel(plan?.max_webhooks) }}</span>
                </span>
              </div>
              <div class="progress usage-bar">
                <div
                  class="progress-bar"
                  :class="barClass(usage?.webhooks ?? 0, plan?.max_webhooks ?? 0)"
                  :style="{ width: pct(usage?.webhooks ?? 0, plan?.max_webhooks ?? 0) + '%' }"
                ></div>
              </div>
              <p v-if="!plan?.has_webhooks" class="feature-note mt-1">
                Webhooks are not available on the {{ planLabel }} plan.
                <button class="btn-link-inline" @click="goCheckout('pro')">Upgrade to Pro</button>
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- ── Plan Comparison ───────────────────────────────────────────── -->
      <div class="card compare-card">
        <div class="card-body">
          <p class="section-label mb-3">Plan comparison</p>
          <div class="table-responsive">
            <table class="table compare-table mb-0">
              <thead>
                <tr>
                  <th>Feature</th>
                  <th class="text-center" :class="{ 'current-col': currentPlanID === 'free' }">Free</th>
                  <th class="text-center" :class="{ 'current-col': currentPlanID === 'pro' }">Pro</th>
                  <th class="text-center" :class="{ 'current-col': currentPlanID === 'business' }">Business</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Links</td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'free' }">5</td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'pro' }">100</td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'business' }">Unlimited</td>
                </tr>
                <tr>
                  <td>API Keys</td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'free' }">1</td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'pro' }">5</td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'business' }">Unlimited</td>
                </tr>
                <tr>
                  <td>Webhooks</td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'free' }">—</td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'pro' }">5</td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'business' }">Unlimited</td>
                </tr>
                <tr>
                  <td>Analytics</td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'free' }">
                    <span class="check">✓</span>
                  </td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'pro' }">
                    <span class="check">✓</span>
                  </td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'business' }">
                    <span class="check">✓</span>
                  </td>
                </tr>
                <tr>
                  <td>QR Codes</td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'free' }">
                    <span class="check">✓</span>
                  </td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'pro' }">
                    <span class="check">✓</span>
                  </td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'business' }">
                    <span class="check">✓</span>
                  </td>
                </tr>
                <tr>
                  <td>Priority support</td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'free' }">—</td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'pro' }">
                    <span class="check">✓</span>
                  </td>
                  <td class="text-center" :class="{ 'current-col': currentPlanID === 'business' }">
                    <span class="check">✓</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import billingApi from '@/api/billing';
import type { PlanInfo, Subscription, UsageCounts } from '@/types/billing';
import { PLAN_LABELS, SUB_STATUS_LABELS } from '@/types/billing';

const loading = ref(true);
const checkoutLoading = ref<'pro' | 'business' | null>(null);
const checkoutError = ref('');

const sub = ref<Subscription | null>(null);
const plan = ref<PlanInfo | null>(null);
const usage = ref<UsageCounts | null>(null);

const currentPlanID = computed(() => sub.value?.plan_id ?? 'free');
const planLabel = computed(() => PLAN_LABELS[currentPlanID.value] ?? currentPlanID.value);
const statusLabel = computed(() => SUB_STATUS_LABELS[sub.value?.status ?? 'active'] ?? sub.value?.status ?? 'Active');
const statusClass = computed(() => {
  switch (sub.value?.status) {
    case 'active':
    case 'on_trial':
      return 'badge-success';
    case 'cancelled':
    case 'expired':
      return 'badge-danger';
    case 'past_due':
      return 'badge-warning';
    default:
      return 'badge-secondary';
  }
});

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' });
}

function limitLabel(limit: number | undefined): string {
  if (limit === undefined || limit === null) return '?';
  return limit === -1 ? '∞' : String(limit);
}

function pct(used: number, limit: number): number {
  if (limit === -1 || limit === 0) return 0;
  return Math.min(100, Math.round((used / limit) * 100));
}

function barClass(used: number, limit: number): string {
  if (limit === -1) return 'bg-success';
  const p = pct(used, limit);
  if (p >= 100) return 'bg-danger';
  if (p >= 80) return 'bg-warning';
  return 'bg-primary';
}

async function goCheckout(planID: 'pro' | 'business') {
  checkoutError.value = '';
  checkoutLoading.value = planID;
  try {
    const res = await billingApi.getCheckoutURL(planID);
    window.location.href = res.data.data.checkout_url;
  } catch {
    checkoutError.value = 'Could not generate checkout link. Please try again.';
    checkoutLoading.value = null;
  }
}

onMounted(async () => {
  try {
    const [subRes, usageRes] = await Promise.all([
      billingApi.getSubscription(),
      billingApi.getUsage(),
    ]);
    sub.value = subRes.data.data.subscription;
    plan.value = subRes.data.data.plan;
    usage.value = usageRes.data.data;
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped lang="scss">
$primary: #635bff;

.billing-page {
  max-width: 900px;
}

.page-title {
  font-size: 1.375rem;
  font-weight: 700;
  color: #1a1f36;
  margin: 0;
}

.page-subtitle {
  color: #697386;
  font-size: 0.875rem;
  margin: 0.25rem 0 0;
}

.billing-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.25rem;

  @media (max-width: 767px) {
    grid-template-columns: 1fr;
  }
}

.plan-card,
.usage-card {
  border-radius: 12px;
  border: 1px solid #e3e8ee;
}

.compare-card {
  grid-column: 1 / -1;
  border-radius: 12px;
  border: 1px solid #e3e8ee;
}

.section-label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #697386;
  margin: 0;
}

.plan-name {
  font-size: 1.75rem;
  font-weight: 700;
  color: #1a1f36;
  margin: 0;
}

// Status badge
.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
  white-space: nowrap;

  &.badge-success {
    background: #e3f9e5;
    color: #1e7e34;
  }

  &.badge-danger {
    background: #fde8e8;
    color: #b91c1c;
  }

  &.badge-warning {
    background: #fef3c7;
    color: #92400e;
  }

  &.badge-secondary {
    background: #f1f3f5;
    color: #697386;
  }
}

.meta-row {
  display: flex;
  gap: 0.5rem;
  font-size: 0.875rem;
  margin-top: 0.375rem;
}

.meta-label {
  color: #697386;
}

.meta-value {
  color: #1a1f36;
  font-weight: 500;
}

// Usage
.usage-list {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.usage-item {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.usage-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.usage-resource {
  font-size: 0.875rem;
  font-weight: 500;
  color: #1a1f36;
}

.usage-count {
  font-size: 0.875rem;
  font-weight: 700;
  color: #1a1f36;
}

.usage-limit {
  font-weight: 400;
  color: #697386;
}

.usage-bar {
  height: 6px;
  border-radius: 3px;
  background: #e9ecef;
  overflow: hidden;

  .progress-bar {
    transition: width 0.4s ease;
    border-radius: 3px;
  }
}

.feature-note {
  font-size: 0.78rem;
  color: #697386;
  margin: 0;
}

.btn-link-inline {
  background: none;
  border: none;
  padding: 0;
  color: $primary;
  font-size: inherit;
  font-weight: 500;
  cursor: pointer;
  text-decoration: underline;
  text-underline-offset: 2px;
}

// Compare table
.compare-table {
  font-size: 0.875rem;

  th {
    font-weight: 600;
    font-size: 0.8rem;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: #697386;
    border-bottom: 2px solid #e3e8ee;
    padding: 0.625rem 0.75rem;
  }

  td {
    padding: 0.625rem 0.75rem;
    color: #1a1f36;
    vertical-align: middle;
    border-bottom: 1px solid #f1f3f5;
  }

  tr:last-child td {
    border-bottom: none;
  }

  .current-col {
    background: rgba($primary, 0.05);
    font-weight: 600;
  }
}

.check {
  color: #1e7e34;
  font-size: 1rem;
}

// Dark mode
.dark-mode & {
  .page-title { color: #e6edf3; }
  .page-subtitle { color: #8b949e; }

  .plan-card,
  .usage-card,
  .compare-card {
    background: #161b22;
    border-color: #30363d;
  }

  .plan-name { color: #e6edf3; }
  .meta-value { color: #e6edf3; }

  .usage-resource,
  .usage-count { color: #e6edf3; }

  .usage-bar { background: #30363d; }

  .compare-table {
    th { color: #8b949e; border-bottom-color: #30363d; }
    td { color: #e6edf3; border-bottom-color: #21262d; }
    .current-col { background: rgba($primary, 0.12); }
  }
}
</style>
