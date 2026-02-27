<template>
  <div class="billing-page">
    <!-- Page header -->
    <div class="page-header">
      <h1 class="md-headline-small">Billing &amp; Plan</h1>
      <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin:4px 0 0;">
        Manage your subscription and track usage.
      </p>
    </div>

    <!-- Loading -->
    <div v-if="loading" style="text-align:center;padding:64px 0;">
      <md-circular-progress indeterminate />
    </div>

    <div v-else class="billing-content">

      <!-- ── Current Plan + Usage row ──────────────────────────────────── -->
      <div class="billing-grid">

        <!-- Current Plan card -->
        <div class="m3-card m3-card--elevated section-card">
          <div class="card-section-body">
            <div style="display:flex;align-items:flex-start;justify-content:space-between;margin-bottom:12px;">
              <div>
                <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);text-transform:uppercase;letter-spacing:0.05em;font-weight:600;margin-bottom:4px;">
                  Current plan
                </div>
                <div class="md-headline-small" style="color:var(--md-sys-color-on-surface);">{{ planLabel }}</div>
              </div>
              <span class="m3-badge" :class="statusClass">{{ statusLabel }}</span>
            </div>

            <div v-if="sub?.current_period_end" class="meta-row">
              <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Renews</span>
              <span class="md-body-small" style="font-weight:500;">{{ formatDate(sub.current_period_end) }}</span>
            </div>
            <div v-if="sub?.cancelled_at" class="meta-row">
              <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Cancelled</span>
              <span class="md-body-small" style="font-weight:500;color:var(--md-sys-color-error);">{{ formatDate(sub.cancelled_at) }}</span>
            </div>

            <!-- Upgrade CTAs -->
            <div style="margin-top:20px;display:flex;flex-wrap:wrap;gap:8px;">
              <md-filled-button
                v-if="currentPlanID === 'free'"
                :disabled="checkoutLoading === 'pro'"
                @click="goCheckout('pro')"
              >
                <md-circular-progress v-if="checkoutLoading === 'pro'" indeterminate style="--md-circular-progress-size:20px;margin-right:8px;" />
                Upgrade to Pro
              </md-filled-button>
              <md-filled-button
                v-if="currentPlanID !== 'business'"
                :disabled="checkoutLoading === 'business'"
                @click="goCheckout('business')"
              >
                <md-circular-progress v-if="checkoutLoading === 'business'" indeterminate style="--md-circular-progress-size:20px;margin-right:8px;" />
                Upgrade to Business
              </md-filled-button>
              <span v-if="currentPlanID === 'business'" class="md-body-medium" style="color:var(--md-sys-color-primary);font-weight:500;align-self:center;">
                <span class="material-symbols-outlined" style="font-size:18px;vertical-align:middle;margin-right:4px;">check_circle</span>
                You're on the highest plan
              </span>
            </div>

            <p v-if="checkoutError" class="md-body-small" style="color:var(--md-sys-color-error);margin-top:8px;">{{ checkoutError }}</p>
          </div>
        </div>

        <!-- Usage card -->
        <div class="m3-card m3-card--elevated section-card">
          <div class="card-section-header">
            <div style="display:flex;align-items:center;gap:8px;">
              <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">analytics</span>
              <span class="md-label-large" style="text-transform:uppercase;letter-spacing:0.05em;color:var(--md-sys-color-on-surface-variant);">Usage this period</span>
            </div>
          </div>
          <div class="card-section-body">

            <!-- Links -->
            <div class="usage-item">
              <div class="usage-item__header">
                <div style="display:flex;align-items:center;gap:6px;">
                  <span class="material-symbols-outlined" style="font-size:16px;color:var(--md-sys-color-on-surface-variant);">link</span>
                  <span class="md-label-large">Links</span>
                </div>
                <div style="display:flex;align-items:center;gap:8px;">
                  <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">
                    {{ usage?.links ?? 0 }} / {{ limitLabel(plan?.max_links) }}
                  </span>
                  <span
                    class="m3-badge"
                    :class="pct(usage?.links ?? 0, plan?.max_links ?? 0) >= 100 ? 'm3-usage-badge--danger' : pct(usage?.links ?? 0, plan?.max_links ?? 0) >= 80 ? 'm3-usage-badge--warning' : 'm3-usage-badge--ok'"
                  >{{ pct(usage?.links ?? 0, plan?.max_links ?? 0) }}%</span>
                </div>
              </div>
              <md-linear-progress :value="pct(usage?.links ?? 0, plan?.max_links ?? 0) / 100" />
            </div>

            <!-- API Keys -->
            <div class="usage-item">
              <div class="usage-item__header">
                <div style="display:flex;align-items:center;gap:6px;">
                  <span class="material-symbols-outlined" style="font-size:16px;color:var(--md-sys-color-on-surface-variant);">key</span>
                  <span class="md-label-large">API Keys</span>
                </div>
                <div style="display:flex;align-items:center;gap:8px;">
                  <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">
                    {{ usage?.api_keys ?? 0 }} / {{ limitLabel(plan?.max_api_keys) }}
                  </span>
                  <span
                    class="m3-badge"
                    :class="pct(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0) >= 100 ? 'm3-usage-badge--danger' : pct(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0) >= 80 ? 'm3-usage-badge--warning' : 'm3-usage-badge--ok'"
                  >{{ pct(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0) }}%</span>
                </div>
              </div>
              <md-linear-progress :value="pct(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0) / 100" />
            </div>

            <!-- Webhooks -->
            <div class="usage-item" style="margin-bottom:0;">
              <div class="usage-item__header">
                <div style="display:flex;align-items:center;gap:6px;">
                  <span class="material-symbols-outlined" style="font-size:16px;color:var(--md-sys-color-on-surface-variant);">notifications</span>
                  <span class="md-label-large">Webhooks</span>
                </div>
                <div style="display:flex;align-items:center;gap:8px;">
                  <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">
                    {{ usage?.webhooks ?? 0 }} / {{ limitLabel(plan?.max_webhooks) }}
                  </span>
                  <span
                    class="m3-badge"
                    :class="pct(usage?.webhooks ?? 0, plan?.max_webhooks ?? 0) >= 100 ? 'm3-usage-badge--danger' : pct(usage?.webhooks ?? 0, plan?.max_webhooks ?? 0) >= 80 ? 'm3-usage-badge--warning' : 'm3-usage-badge--ok'"
                  >{{ pct(usage?.webhooks ?? 0, plan?.max_webhooks ?? 0) }}%</span>
                </div>
              </div>
              <md-linear-progress :value="pct(usage?.webhooks ?? 0, plan?.max_webhooks ?? 0) / 100" />
              <p v-if="!plan?.has_webhooks" class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin:6px 0 0;">
                Webhooks are not available on the {{ planLabel }} plan.
                <button class="btn-link-inline" @click="goCheckout('pro')">Upgrade to Pro</button>
              </p>
            </div>

          </div>
        </div>
      </div>

      <!-- ── Plan Comparison table ─────────────────────────────────────── -->
      <div class="m3-card m3-card--elevated section-card" style="margin-top:20px;">
        <div class="card-section-header">
          <div style="display:flex;align-items:center;gap:8px;">
            <span class="material-symbols-outlined" style="font-size:18px;color:var(--md-sys-color-primary);">compare</span>
            <span class="md-label-large" style="text-transform:uppercase;letter-spacing:0.05em;color:var(--md-sys-color-on-surface-variant);">Plan comparison</span>
          </div>
        </div>
        <div class="m3-table-wrapper">
          <table class="m3-table compare-table">
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
                  <span class="material-symbols-outlined check-icon">check_circle</span>
                </td>
                <td class="text-center" :class="{ 'current-col': currentPlanID === 'pro' }">
                  <span class="material-symbols-outlined check-icon">check_circle</span>
                </td>
                <td class="text-center" :class="{ 'current-col': currentPlanID === 'business' }">
                  <span class="material-symbols-outlined check-icon">check_circle</span>
                </td>
              </tr>
              <tr>
                <td>QR Codes</td>
                <td class="text-center" :class="{ 'current-col': currentPlanID === 'free' }">
                  <span class="material-symbols-outlined check-icon">check_circle</span>
                </td>
                <td class="text-center" :class="{ 'current-col': currentPlanID === 'pro' }">
                  <span class="material-symbols-outlined check-icon">check_circle</span>
                </td>
                <td class="text-center" :class="{ 'current-col': currentPlanID === 'business' }">
                  <span class="material-symbols-outlined check-icon">check_circle</span>
                </td>
              </tr>
              <tr>
                <td>Priority support</td>
                <td class="text-center" :class="{ 'current-col': currentPlanID === 'free' }">—</td>
                <td class="text-center" :class="{ 'current-col': currentPlanID === 'pro' }">
                  <span class="material-symbols-outlined check-icon">check_circle</span>
                </td>
                <td class="text-center" :class="{ 'current-col': currentPlanID === 'business' }">
                  <span class="material-symbols-outlined check-icon">check_circle</span>
                </td>
              </tr>
            </tbody>
          </table>
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
.billing-page {
  max-width: 900px;
  padding: 24px 0;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.page-header {
  margin-bottom: 24px;
}

.billing-content {
  display: flex;
  flex-direction: column;
}

.billing-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;

  @media (max-width: 767px) {
    grid-template-columns: 1fr;
  }
}

.section-card {
  border-radius: 12px;
  overflow: hidden;
}

.m3-card {
  border-radius: 12px;
  background: var(--md-sys-color-surface);
  overflow: hidden;

  &--elevated {
    box-shadow: 0 1px 3px rgba(0,0,0,0.10), 0 2px 6px rgba(0,0,0,0.07);
  }
}

.card-section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 14px 24px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.card-section-body {
  padding: 24px;
}

.meta-row {
  display: flex;
  gap: 8px;
  margin-top: 6px;
}

/* Usage items */
.usage-item {
  margin-bottom: 20px;

  &:last-child {
    margin-bottom: 0;
  }
}

.usage-item__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
  flex-wrap: wrap;
  gap: 4px;
}

/* Status badges */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;

  &.badge-success {
    background: rgba(22, 163, 74, 0.12);
    color: #16a34a;
  }

  &.badge-danger {
    background: var(--md-sys-color-error-container);
    color: var(--md-sys-color-error);
  }

  &.badge-warning {
    background: rgba(245, 158, 11, 0.15);
    color: #92400e;
  }

  &.badge-secondary {
    background: var(--md-sys-color-surface-container-low);
    color: var(--md-sys-color-on-surface-variant);
  }
}

/* Usage percentage badges */
.m3-usage-badge--danger {
  background: rgba(220, 38, 38, 0.12);
  color: #dc2626;
}

.m3-usage-badge--warning {
  background: rgba(245, 158, 11, 0.12);
  color: #b45309;
}

.m3-usage-badge--ok {
  background: rgba(22, 163, 74, 0.12);
  color: #16a34a;
}

.btn-link-inline {
  background: none;
  border: none;
  padding: 0;
  color: var(--md-sys-color-primary);
  font-size: inherit;
  font-weight: 500;
  cursor: pointer;
  text-decoration: underline;
  text-underline-offset: 2px;
}

/* Compare table */
.m3-table-wrapper {
  overflow-x: auto;
}

.m3-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}

.compare-table {
  width: 100%;

  th,
  td {
    padding: 10px 16px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  tbody tr:last-child td {
    border-bottom: none;
  }

  th {
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--md-sys-color-on-surface-variant);
    font-weight: 600;
    background: var(--md-sys-color-surface-container-low);
  }

  .text-center {
    text-align: center;
  }

  .current-col {
    background: rgba(99, 91, 255, 0.08);
    font-weight: 600;
  }
}

.check-icon {
  font-size: 18px;
  color: #16a34a;
  vertical-align: middle;
}
</style>
