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
        <div class="m3-card m3-card--outlined section-card">
          <div class="card-section-header">
            <span class="md-label-large" style="text-transform:uppercase;letter-spacing:0.05em;color:var(--md-sys-color-on-surface-variant);">Usage this period</span>
          </div>
          <div class="card-section-body">

            <!-- Links -->
            <div style="margin-bottom:20px;">
              <div style="display:flex;justify-content:space-between;margin-bottom:6px;">
                <span class="md-label-large">Links</span>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">
                  {{ usage?.links ?? 0 }} / {{ limitLabel(plan?.max_links) }}
                </span>
              </div>
              <md-linear-progress :value="pct(usage?.links ?? 0, plan?.max_links ?? 0) / 100" />
            </div>

            <!-- API Keys -->
            <div style="margin-bottom:20px;">
              <div style="display:flex;justify-content:space-between;margin-bottom:6px;">
                <span class="md-label-large">API Keys</span>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">
                  {{ usage?.api_keys ?? 0 }} / {{ limitLabel(plan?.max_api_keys) }}
                </span>
              </div>
              <md-linear-progress :value="pct(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0) / 100" />
            </div>

            <!-- Webhooks -->
            <div>
              <div style="display:flex;justify-content:space-between;margin-bottom:6px;">
                <span class="md-label-large">Webhooks</span>
                <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">
                  {{ usage?.webhooks ?? 0 }} / {{ limitLabel(plan?.max_webhooks) }}
                </span>
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
      <div class="m3-card m3-card--outlined section-card" style="margin-top:20px;">
        <div class="card-section-header">
          <span class="md-label-large" style="text-transform:uppercase;letter-spacing:0.05em;color:var(--md-sys-color-on-surface-variant);">Plan comparison</span>
        </div>
        <div class="card-section-body" style="padding:0;overflow-x:auto;">
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

<style scoped>
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
}

@media (max-width: 767px) {
  .billing-grid {
    grid-template-columns: 1fr;
  }
}

.section-card {
  border-radius: 12px;
  overflow: hidden;
}

.card-section-header {
  padding: 16px 24px;
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

/* Status badges */
.m3-badge.badge-success {
  background: color-mix(in srgb, #1e7e34 12%, transparent);
  color: #1e7e34;
}

.m3-badge.badge-danger {
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-error);
}

.m3-badge.badge-warning {
  background: color-mix(in srgb, #f59e0b 15%, transparent);
  color: #92400e;
}

.m3-badge.badge-secondary {
  background: var(--md-sys-color-surface-container-low);
  color: var(--md-sys-color-on-surface-variant);
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
.compare-table {
  width: 100%;
}

.compare-table th,
.compare-table td {
  padding: 10px 16px;
}

.compare-table th {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--md-sys-color-on-surface-variant);
  font-weight: 600;
}

.compare-table .text-center {
  text-align: center;
}

.compare-table .current-col {
  background: color-mix(in srgb, var(--md-sys-color-primary) 8%, transparent);
  font-weight: 600;
}

.check-icon {
  font-size: 18px;
  color: #1e7e34;
  vertical-align: middle;
}
</style>
