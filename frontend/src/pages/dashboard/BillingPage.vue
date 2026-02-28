<template>
  <div class="billing-page">

    <!-- Page Header -->
    <div class="page-header">
      <div class="page-header__left">
        <h1 class="page-title">Billing &amp; Plan</h1>
        <p class="page-subtitle">Manage your subscription and track usage.</p>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-center">
      <div class="css-spinner"></div>
    </div>

    <div v-else class="billing-content">

      <!-- ── Current Plan + Usage row ──────────────────────────────────── -->
      <div class="billing-grid">

        <!-- Current Plan card -->
        <div class="an-card">
          <div class="an-card-header">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">workspace_premium</span>
            </div>
            <span class="an-card-title">Current Plan</span>
            <span class="m3-badge" :class="statusClass">{{ statusLabel }}</span>
          </div>
          <div class="an-card-body">
            <div class="plan-name">{{ planLabel }}</div>

            <div v-if="sub?.current_period_end" class="meta-row">
              <span class="meta-row__label">Renews</span>
              <span class="meta-row__value">{{ formatDate(sub.current_period_end) }}</span>
            </div>
            <div v-if="sub?.cancelled_at" class="meta-row">
              <span class="meta-row__label">Cancelled</span>
              <span class="meta-row__value meta-row__value--danger">{{ formatDate(sub.cancelled_at) }}</span>
            </div>

            <!-- Upgrade CTAs -->
            <div class="upgrade-actions">
              <button class="btn-filled"
                v-if="currentPlanID === 'free'"
                :disabled="checkoutLoading === 'pro'"
                @click="goCheckout('pro')"
              >
                <div v-if="checkoutLoading === 'pro'" class="css-spinner css-spinner--sm css-spinner--white"></div>
                <span v-else class="material-symbols-outlined">arrow_upward</span>
                Upgrade to Pro
              </button>
              <button class="btn-filled"
                v-if="currentPlanID !== 'business'"
                :disabled="checkoutLoading === 'business'"
                @click="goCheckout('business')"
              >
                <div v-if="checkoutLoading === 'business'" class="css-spinner css-spinner--sm css-spinner--white"></div>
                <span v-else class="material-symbols-outlined">rocket_launch</span>
                Upgrade to Business
              </button>
              <span v-if="currentPlanID === 'business'" class="plan-top-badge">
                <span class="material-symbols-outlined">check_circle</span>
                You're on the highest plan
              </span>
            </div>

            <p v-if="checkoutError" class="checkout-error">{{ checkoutError }}</p>
          </div>
        </div>

        <!-- Usage card -->
        <div class="an-card">
          <div class="an-card-header">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">analytics</span>
            </div>
            <span class="an-card-title">Usage This Period</span>
          </div>
          <div class="an-card-body">

            <!-- Links -->
            <div class="usage-item">
              <div class="usage-item__header">
                <div class="usage-item__name">
                  <span class="material-symbols-outlined usage-icon">link</span>
                  <span class="usage-label">Links</span>
                </div>
                <div class="usage-item__stats">
                  <span class="usage-count">{{ usage?.links ?? 0 }} / {{ limitLabel(plan?.max_links) }} used</span>
                  <span class="m3-usage-badge" :class="pct(usage?.links ?? 0, plan?.max_links ?? 0) >= 100 ? 'm3-usage-badge--danger' : pct(usage?.links ?? 0, plan?.max_links ?? 0) >= 80 ? 'm3-usage-badge--warning' : 'm3-usage-badge--ok'">{{ pct(usage?.links ?? 0, plan?.max_links ?? 0) }}%</span>
                </div>
              </div>
              <div class="prog-bar">
                <div class="prog-fill" :class="barClass(usage?.links ?? 0, plan?.max_links ?? 0)" :style="{ width: pct(usage?.links ?? 0, plan?.max_links ?? 0) + '%' }"></div>
              </div>
            </div>

            <!-- API Keys -->
            <div class="usage-item">
              <div class="usage-item__header">
                <div class="usage-item__name">
                  <span class="material-symbols-outlined usage-icon">key</span>
                  <span class="usage-label">API Keys</span>
                </div>
                <div class="usage-item__stats">
                  <span class="usage-count">{{ usage?.api_keys ?? 0 }} / {{ limitLabel(plan?.max_api_keys) }} used</span>
                  <span class="m3-usage-badge" :class="pct(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0) >= 100 ? 'm3-usage-badge--danger' : pct(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0) >= 80 ? 'm3-usage-badge--warning' : 'm3-usage-badge--ok'">{{ pct(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0) }}%</span>
                </div>
              </div>
              <div class="prog-bar">
                <div class="prog-fill" :class="barClass(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0)" :style="{ width: pct(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0) + '%' }"></div>
              </div>
            </div>

            <!-- Webhooks -->
            <div class="usage-item usage-item--last">
              <div class="usage-item__header">
                <div class="usage-item__name">
                  <span class="material-symbols-outlined usage-icon">notifications</span>
                  <span class="usage-label">Webhooks</span>
                </div>
                <div class="usage-item__stats">
                  <span class="usage-count">{{ usage?.webhooks ?? 0 }} / {{ limitLabel(plan?.max_webhooks) }} used</span>
                  <span class="m3-usage-badge" :class="pct(usage?.webhooks ?? 0, plan?.max_webhooks ?? 0) >= 100 ? 'm3-usage-badge--danger' : pct(usage?.webhooks ?? 0, plan?.max_webhooks ?? 0) >= 80 ? 'm3-usage-badge--warning' : 'm3-usage-badge--ok'">{{ pct(usage?.webhooks ?? 0, plan?.max_webhooks ?? 0) }}%</span>
                </div>
              </div>
              <div class="prog-bar">
                <div class="prog-fill" :class="barClass(usage?.webhooks ?? 0, plan?.max_webhooks ?? 0)" :style="{ width: pct(usage?.webhooks ?? 0, plan?.max_webhooks ?? 0) + '%' }"></div>
              </div>
              <p v-if="!plan?.has_webhooks" class="usage-item__note">
                Webhooks are not available on the {{ planLabel }} plan.
                <button class="btn-link-inline" @click="goCheckout('pro')">Upgrade to Pro</button>
              </p>
            </div>

          </div>
        </div>
      </div>

      <!-- ── Plan Comparison table ─────────────────────────────────────── -->
      <div class="an-card">
        <div class="an-card-header">
          <div class="an-card-icon an-card-icon--primary">
            <span class="material-symbols-outlined">compare</span>
          </div>
          <span class="an-card-title">Plan Comparison</span>
        </div>
        <div class="table-wrapper">
          <table class="compare-table">
            <thead>
              <tr>
                <th>Feature</th>
                <th class="text-center" :class="{ 'current-col': currentPlanID === 'free' }">
                  Free
                  <span v-if="currentPlanID === 'free'" class="current-plan-tag">Current</span>
                </th>
                <th class="text-center" :class="{ 'current-col': currentPlanID === 'pro' }">
                  Pro
                  <span v-if="currentPlanID === 'pro'" class="current-plan-tag">Current</span>
                </th>
                <th class="text-center" :class="{ 'current-col': currentPlanID === 'business' }">
                  Business
                  <span v-if="currentPlanID === 'business'" class="current-plan-tag">Current</span>
                </th>
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
                <td class="text-center" :class="{ 'current-col': currentPlanID === 'free' }">
                  <span class="material-symbols-outlined x-icon">close</span>
                </td>
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
                <td class="text-center" :class="{ 'current-col': currentPlanID === 'free' }">
                  <span class="material-symbols-outlined x-icon">close</span>
                </td>
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
  if (limit === -1) return 'prog-fill--success';
  const p = pct(used, limit);
  if (p >= 100) return 'prog-fill--danger';
  if (p >= 80) return 'prog-fill--warning';
  return '';
}

async function goCheckout(planID: 'pro' | 'business') {
  checkoutError.value = '';
  checkoutLoading.value = planID;
  try {
    const res = await billingApi.getCheckoutURL(planID);
    if (res.data.data) window.location.href = res.data.data.checkout_url;
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
    if (subRes.data.data) {
      sub.value = subRes.data.data.subscription;
      plan.value = subRes.data.data.plan;
    }
    if (usageRes.data.data) usage.value = usageRes.data.data;
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped lang="scss">
.billing-page {
  max-width: 900px;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* ── Page Header ─────────────────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
}

.page-header__left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.page-title {
  font-size: 1.375rem;
  font-weight: 700;
  margin: 0;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
  margin: 0;
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── CSS Spinner ─────────────────────────────────────────────────────────── */
.css-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--md-sys-color-outline-variant);
  border-top-color: var(--md-sys-color-primary);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;

  &--sm {
    width: 16px;
    height: 16px;
    border-width: 2px;
  }

  &--white {
    border-color: rgba(255,255,255,0.35);
    border-top-color: #fff;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Loading ─────────────────────────────────────────────────────────────── */
.loading-center {
  display: flex;
  justify-content: center;
  padding: 64px;
}

/* ── Layout ──────────────────────────────────────────────────────────────── */
.billing-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.billing-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;

  @media (max-width: 767px) {
    grid-template-columns: 1fr;
  }
}

/* ── Cards ───────────────────────────────────────────────────────────────── */
.an-card {
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  overflow: hidden;
}

.an-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface-container-low);
}

.an-card-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  .material-symbols-outlined { font-size: 18px; }

  &--primary {
    background: color-mix(in srgb, var(--md-sys-color-primary) 12%, transparent);
    color: var(--md-sys-color-primary);
  }
}

.an-card-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  flex: 1;
}

.an-card-body {
  padding: 20px 24px 24px;
}

/* ── Plan name ───────────────────────────────────────────────────────────── */
.plan-name {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
  margin-bottom: 8px;
}

.plan-top-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--md-sys-color-primary);
  font-weight: 500;
  font-size: 0.875rem;

  .material-symbols-outlined { font-size: 18px; }
}

.upgrade-actions {
  margin-top: 20px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.checkout-error {
  margin: 8px 0 0;
  font-size: 0.8rem;
  color: var(--md-sys-color-error);
}

/* ── Meta row ────────────────────────────────────────────────────────────── */
.meta-row {
  display: flex;
  gap: 8px;
  margin-top: 6px;
  font-size: 0.875rem;

  &__label { color: var(--md-sys-color-on-surface-variant); }
  &__value { font-weight: 500; }
  &__value--danger { color: var(--md-sys-color-error); }
}

/* ── Usage items ─────────────────────────────────────────────────────────── */
.usage-item {
  margin-bottom: 20px;

  &--last { margin-bottom: 0; }
}

.usage-item__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
  flex-wrap: wrap;
  gap: 4px;
}

.usage-item__name {
  display: flex;
  align-items: center;
  gap: 6px;
}

.usage-icon {
  font-size: 16px;
  color: var(--md-sys-color-on-surface-variant);
}

.usage-label {
  font-size: 0.875rem;
  font-weight: 500;
}

.usage-item__stats {
  display: flex;
  align-items: center;
  gap: 8px;
}

.usage-count {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
}

.usage-item__note {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 6px 0 0;
}

/* ── Progress bar ────────────────────────────────────────────────────────── */
.prog-bar {
  height: 6px;
  background: var(--md-sys-color-surface-container-high);
  border-radius: 999px;
  overflow: hidden;
}

.prog-fill {
  height: 100%;
  background: var(--md-sys-color-primary);
  border-radius: 999px;
  transition: width 0.4s ease;

  &--success { background: #16a34a; }
  &--warning { background: #f59e0b; }
  &--danger { background: #dc2626; }
}

/* ── Status badges ───────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;
  flex-shrink: 0;

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

/* ── Usage percentage badges ─────────────────────────────────────────────── */
.m3-usage-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;

  &--danger { background: rgba(220, 38, 38, 0.12); color: #dc2626; }
  &--warning { background: rgba(245, 158, 11, 0.12); color: #b45309; }
  &--ok { background: rgba(22, 163, 74, 0.12); color: #16a34a; }
}

/* ── Inline link button ──────────────────────────────────────────────────── */
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

/* ── Compare table ───────────────────────────────────────────────────────── */
.table-wrapper {
  overflow-x: auto;
}

.compare-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;

  th,
  td {
    padding: 10px 16px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  tbody tr:last-child td {
    border-bottom: none;
  }

  th {
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
    background: var(--md-sys-color-surface-container-low);
    text-align: left;
  }

  .text-center { text-align: center; }

  .current-col {
    background: rgba(99, 91, 255, 0.08);
    font-weight: 600;
  }
}

.current-plan-tag {
  display: block;
  font-size: 0.65rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--md-sys-color-primary);
  margin-top: 2px;
}

.check-icon {
  font-size: 18px;
  color: #16a34a;
  vertical-align: middle;
}

.x-icon {
  font-size: 18px;
  color: var(--md-sys-color-on-surface-variant);
  opacity: 0.4;
  vertical-align: middle;
}
</style>
