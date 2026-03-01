<template>
  <div class="billing-page">

    <!-- Page Header -->
    <div class="page-header">
      <h1 class="page-title">Billing &amp; Plan</h1>
      <p class="page-subtitle">Manage your subscription and track usage.</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-center">
      <div class="css-spinner"></div>
    </div>

    <div v-else class="billing-content">

      <!-- ── Plan Cards ──────────────────────────────────────────────────── -->
      <div class="plan-cards">

        <!-- Free -->
        <div class="plan-card" :class="{ 'plan-card--current': currentPlanID === 'free' }">
          <div class="plan-card-top">
            <div class="plan-tier-row">
              <span class="plan-tier-name">Free</span>
              <span v-if="currentPlanID === 'free'" class="plan-badge plan-badge--current">Current Plan</span>
            </div>
            <div class="plan-price-row">
              <span class="plan-price">$0</span>
              <span class="plan-price-period">/ forever</span>
            </div>
            <p class="plan-tagline">Get started at no cost</p>
          </div>

          <ul class="plan-features">
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              5 short links
            </li>
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              QR code generation
            </li>
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              Basic click counts
            </li>
            <li class="plan-feature plan-feature--no">
              <span class="material-symbols-outlined feat-icon feat-icon--no">remove_circle</span>
              Analytics &amp; reports
            </li>
            <li class="plan-feature plan-feature--no">
              <span class="material-symbols-outlined feat-icon feat-icon--no">remove_circle</span>
              Custom slugs
            </li>
            <li class="plan-feature plan-feature--no">
              <span class="material-symbols-outlined feat-icon feat-icon--no">remove_circle</span>
              API keys
            </li>
            <li class="plan-feature plan-feature--no">
              <span class="material-symbols-outlined feat-icon feat-icon--no">remove_circle</span>
              Bio link click tracking
            </li>
          </ul>

          <div class="plan-card-footer">
            <div v-if="currentPlanID === 'free'" class="plan-active-state">
              <span class="material-symbols-outlined">check_circle</span>
              Active plan
            </div>
          </div>
        </div>

        <!-- Pro -->
        <div class="plan-card plan-card--pro" :class="{ 'plan-card--current': currentPlanID === 'pro' }">
          <div class="plan-card-top">
            <div class="plan-tier-row">
              <span class="plan-tier-name">Pro</span>
              <span v-if="currentPlanID === 'pro'" class="plan-badge plan-badge--current">Current Plan</span>
              <span v-else class="plan-badge plan-badge--recommended">Recommended</span>
            </div>
            <div class="plan-price-row">
              <span class="plan-price">$7</span>
              <span class="plan-price-period">/ month</span>
            </div>
            <p class="plan-tagline">Everything you need to grow</p>
          </div>

          <ul class="plan-features">
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              100 short links
            </li>
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              QR code generation
            </li>
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              Full analytics &amp; charts
            </li>
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              Scheduled email reports
            </li>
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              Custom slugs
            </li>
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              5 API keys
            </li>
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              Bio link click tracking
            </li>
          </ul>

          <div class="plan-card-footer">
            <div v-if="currentPlanID === 'pro'" class="plan-active-state plan-active-state--pro">
              <span class="material-symbols-outlined">check_circle</span>
              Active plan
            </div>
            <button
              v-else
              class="btn-upgrade"
              :disabled="checkoutLoading === 'pro'"
              @click="goCheckout('pro')"
            >
              <div v-if="checkoutLoading === 'pro'" class="css-spinner css-spinner--sm css-spinner--white"></div>
              <span v-else class="material-symbols-outlined">workspace_premium</span>
              Upgrade to Pro
            </button>
            <p v-if="checkoutError && checkoutLoading !== 'growth'" class="checkout-error">{{ checkoutError }}</p>
            <p v-if="currentPlanID !== 'pro'" class="billing-terms-note">
              By upgrading you agree to our
              <router-link to="/terms" target="_blank">Terms of Service</router-link>
              and
              <router-link to="/privacy" target="_blank">Privacy Policy</router-link>.
            </p>
          </div>
        </div>

        <!-- Growth -->
        <div class="plan-card plan-card--growth" :class="{ 'plan-card--current': currentPlanID === 'growth' }">
          <div class="plan-card-top">
            <div class="plan-tier-row">
              <span class="plan-tier-name">Growth</span>
              <span v-if="currentPlanID === 'growth'" class="plan-badge plan-badge--current">Current Plan</span>
            </div>
            <div class="plan-price-row">
              <span class="plan-price">$15</span>
              <span class="plan-price-period">/ month</span>
            </div>
            <p class="plan-tagline">For teams and growing businesses</p>
          </div>

          <ul class="plan-features">
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              Unlimited short links
            </li>
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              Everything in Pro
            </li>
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              10 API keys
            </li>
            <li class="plan-feature plan-feature--yes">
              <span class="material-symbols-outlined feat-icon feat-icon--yes">check_circle</span>
              Priority support
            </li>
          </ul>

          <div class="plan-card-footer">
            <div v-if="currentPlanID === 'growth'" class="plan-active-state plan-active-state--pro">
              <span class="material-symbols-outlined">check_circle</span>
              Active plan
            </div>
            <button
              v-else
              class="btn-upgrade btn-upgrade--growth"
              :disabled="checkoutLoading === 'growth'"
              @click="goCheckout('growth')"
            >
              <div v-if="checkoutLoading === 'growth'" class="css-spinner css-spinner--sm css-spinner--white"></div>
              <span v-else class="material-symbols-outlined">rocket_launch</span>
              Upgrade to Growth
            </button>
            <p v-if="checkoutError && checkoutLoading !== 'pro'" class="checkout-error">{{ checkoutError }}</p>
            <p v-if="currentPlanID !== 'growth'" class="billing-terms-note">
              By upgrading you agree to our
              <router-link to="/terms" target="_blank">Terms of Service</router-link>
              and
              <router-link to="/privacy" target="_blank">Privacy Policy</router-link>.
            </p>
          </div>
        </div>

      </div><!-- /plan-cards -->

      <!-- ── Subscription status + Usage ────────────────────────────────── -->
      <div class="billing-grid">

        <!-- Subscription details -->
        <div class="an-card" v-if="sub">
          <div class="an-card-header">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">receipt_long</span>
            </div>
            <span class="an-card-title">Subscription</span>
            <span class="m3-badge" :class="statusClass">{{ statusLabel }}</span>
          </div>
          <div class="an-card-body">
            <div class="sub-detail-row">
              <span class="sub-detail-label">Plan</span>
              <span class="sub-detail-value">{{ planLabel }}</span>
            </div>
            <div v-if="sub.current_period_end" class="sub-detail-row">
              <span class="sub-detail-label">Renews</span>
              <span class="sub-detail-value">{{ formatDate(sub.current_period_end) }}</span>
            </div>
            <div v-if="sub.cancelled_at" class="sub-detail-row">
              <span class="sub-detail-label">Cancelled</span>
              <span class="sub-detail-value sub-detail-value--danger">{{ formatDate(sub.cancelled_at) }}</span>
            </div>
          </div>
        </div>

        <!-- Usage -->
        <div class="an-card">
          <div class="an-card-header">
            <div class="an-card-icon an-card-icon--primary">
              <span class="material-symbols-outlined">analytics</span>
            </div>
            <span class="an-card-title">Usage This Period</span>
          </div>
          <div class="an-card-body">

            <div class="usage-item">
              <div class="usage-item__header">
                <div class="usage-item__name">
                  <span class="material-symbols-outlined usage-icon">link</span>
                  <span class="usage-label">Links</span>
                </div>
                <div class="usage-item__stats">
                  <span class="usage-count">{{ usage?.links ?? 0 }} / {{ limitLabel(plan?.max_links) }}</span>
                  <span class="m3-usage-badge" :class="badgeClass(usage?.links ?? 0, plan?.max_links ?? 0)">
                    {{ pct(usage?.links ?? 0, plan?.max_links ?? 0) }}%
                  </span>
                </div>
              </div>
              <div class="prog-bar">
                <div class="prog-fill" :class="barClass(usage?.links ?? 0, plan?.max_links ?? 0)" :style="{ width: pct(usage?.links ?? 0, plan?.max_links ?? 0) + '%' }"></div>
              </div>
            </div>

            <div class="usage-item usage-item--last">
              <div class="usage-item__header">
                <div class="usage-item__name">
                  <span class="material-symbols-outlined usage-icon">key</span>
                  <span class="usage-label">API Keys</span>
                </div>
                <div class="usage-item__stats">
                  <span v-if="(plan?.max_api_keys ?? 0) === 0" class="usage-locked">
                    <span class="material-symbols-outlined" style="font-size:13px;">lock</span>
                    Paid feature
                  </span>
                  <template v-else>
                    <span class="usage-count">{{ usage?.api_keys ?? 0 }} / {{ limitLabel(plan?.max_api_keys) }}</span>
                    <span class="m3-usage-badge" :class="badgeClass(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0)">
                      {{ pct(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0) }}%
                    </span>
                  </template>
                </div>
              </div>
              <div v-if="(plan?.max_api_keys ?? 0) > 0" class="prog-bar">
                <div class="prog-fill" :class="barClass(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0)" :style="{ width: pct(usage?.api_keys ?? 0, plan?.max_api_keys ?? 0) + '%' }"></div>
              </div>
            </div>

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
const checkoutLoading = ref<'pro' | 'growth' | false>(false);
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
  if (limit <= 0) return 0;
  return Math.min(100, Math.round((used / limit) * 100));
}

function barClass(used: number, limit: number): string {
  if (limit === -1) return 'prog-fill--success';
  const p = pct(used, limit);
  if (p >= 100) return 'prog-fill--danger';
  if (p >= 80) return 'prog-fill--warning';
  return '';
}

function badgeClass(used: number, limit: number): string {
  const p = pct(used, limit);
  if (p >= 100) return 'm3-usage-badge--danger';
  if (p >= 80) return 'm3-usage-badge--warning';
  return 'm3-usage-badge--ok';
}

async function goCheckout(planID: 'pro' | 'growth') {
  checkoutError.value = '';
  checkoutLoading.value = planID;
  try {
    const res = await billingApi.getCheckoutURL(planID);
    if (res.data.data) window.location.href = res.data.data.checkout_url;
  } catch {
    checkoutError.value = 'Could not generate checkout link. Please try again.';
  } finally {
    checkoutLoading.value = false;
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
$ease: cubic-bezier(0.2, 0, 0, 1);

.billing-page {
  max-width: 1060px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* ── Header ──────────────────────────────────────────────────────────────── */
.page-header { display: flex; flex-direction: column; gap: 4px; }

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

/* ── Loading ─────────────────────────────────────────────────────────────── */
.css-spinner {
  width: 32px; height: 32px;
  border: 3px solid var(--md-sys-color-outline-variant);
  border-top-color: var(--md-sys-color-primary);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;

  &--sm { width: 16px; height: 16px; border-width: 2px; }
  &--white { border-color: rgba(255,255,255,0.35); border-top-color: #fff; }
}
@keyframes spin { to { transform: rotate(360deg); } }

.loading-center { display: flex; justify-content: center; padding: 64px; }

.billing-content { display: flex; flex-direction: column; gap: 24px; }

/* ── Plan Cards ──────────────────────────────────────────────────────────── */
.plan-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;

  @media (max-width: 860px) { grid-template-columns: 1fr 1fr; }
  @media (max-width: 560px) { grid-template-columns: 1fr; }
}

.plan-card {
  display: flex;
  flex-direction: column;
  background: var(--md-sys-color-surface);
  border: 2px solid var(--md-sys-color-outline-variant);
  border-radius: 18px;
  overflow: hidden;
  transition: border-color 0.2s $ease, box-shadow 0.2s $ease;

  &--current {
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary) 12%, transparent);
  }

  &--pro {
    background: linear-gradient(160deg,
      color-mix(in srgb, var(--md-sys-color-primary-container) 40%, var(--md-sys-color-surface)) 0%,
      var(--md-sys-color-surface) 60%
    );
  }

  &--growth {
    background: linear-gradient(160deg,
      color-mix(in srgb, var(--md-sys-color-secondary-container) 40%, var(--md-sys-color-surface)) 0%,
      var(--md-sys-color-surface) 60%
    );
  }
}

.plan-card-top {
  padding: 24px 24px 16px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.plan-tier-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.plan-tier-name {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
}

.plan-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 100px;
  font-size: 0.68rem;
  font-weight: 700;
  letter-spacing: 0.3px;
  text-transform: uppercase;

  &--current {
    background: color-mix(in srgb, var(--md-sys-color-primary) 15%, transparent);
    color: var(--md-sys-color-primary);
  }

  &--recommended {
    background: var(--md-sys-color-tertiary-container);
    color: var(--md-sys-color-on-tertiary-container);
  }
}

.plan-price-row {
  display: flex;
  align-items: baseline;
  gap: 4px;
  margin-bottom: 6px;
}

.plan-price {
  font-size: 2.5rem;
  font-weight: 800;
  color: var(--md-sys-color-on-surface);
  line-height: 1;
}

.plan-price-period {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
}

.plan-tagline {
  margin: 0;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Feature list ────────────────────────────────────────────────────────── */
.plan-features {
  list-style: none;
  margin: 0;
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  flex: 1;
}

.plan-feature {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 0.875rem;

  &--no { color: var(--md-sys-color-on-surface-variant); opacity: 0.6; }
  &--yes { color: var(--md-sys-color-on-surface); }
}

.feat-icon {
  font-size: 18px;
  flex-shrink: 0;

  &--yes { color: #16a34a; }
  &--no  { color: var(--md-sys-color-on-surface-variant); opacity: 0.4; }
}

/* ── Card footer ─────────────────────────────────────────────────────────── */
.plan-card-footer {
  padding: 16px 24px 24px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.plan-active-state {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface-variant);

  .material-symbols-outlined { font-size: 18px; }

  &--pro {
    color: var(--md-sys-color-primary);
  }
}

.btn-upgrade {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  height: 48px;
  border: none;
  border-radius: 12px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  font-size: 0.9375rem;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.15s $ease;

  .material-symbols-outlined { font-size: 20px; }

  &:hover:not(:disabled) { opacity: 0.88; }
  &:disabled { opacity: 0.5; cursor: not-allowed; }

  &--growth {
    background: var(--md-sys-color-secondary);
    color: var(--md-sys-color-on-secondary);
  }
}

.checkout-error {
  margin: 0;
  font-size: 0.8rem;
  color: var(--md-sys-color-error);
}

.billing-terms-note {
  margin: 0;
  font-size: 0.75rem;
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1.5;

  a {
    color: var(--md-sys-color-primary);
    text-decoration: underline;
    text-underline-offset: 2px;
  }
}

/* ── Subscription + Usage grid ───────────────────────────────────────────── */
.billing-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;

  @media (max-width: 640px) { grid-template-columns: 1fr; }
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
  width: 36px; height: 36px;
  border-radius: 10px;
  display: flex; align-items: center; justify-content: center;
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

.an-card-body { padding: 20px 24px 24px; }

/* ── Subscription detail rows ────────────────────────────────────────────── */
.sub-detail-row {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  padding: 8px 0;
  font-size: 0.875rem;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);

  &:last-child { border-bottom: none; }
}

.sub-detail-label { color: var(--md-sys-color-on-surface-variant); }
.sub-detail-value { font-weight: 500; }
.sub-detail-value--danger { color: var(--md-sys-color-error); }

/* ── Usage ───────────────────────────────────────────────────────────────── */
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

.usage-locked {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--md-sys-color-tertiary);
  background: var(--md-sys-color-tertiary-container);
  border-radius: 100px;
  padding: 2px 8px;
}

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
  &--danger  { background: #dc2626; }
}

/* ── Badges ──────────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;
  flex-shrink: 0;

  &.badge-success { background: rgba(22,163,74,0.12); color: #16a34a; }
  &.badge-danger  { background: var(--md-sys-color-error-container); color: var(--md-sys-color-error); }
  &.badge-warning { background: rgba(245,158,11,0.15); color: #92400e; }
  &.badge-secondary { background: var(--md-sys-color-surface-container-low); color: var(--md-sys-color-on-surface-variant); }
}

.m3-usage-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;

  &--danger  { background: rgba(220,38,38,0.12); color: #dc2626; }
  &--warning { background: rgba(245,158,11,0.12); color: #b45309; }
  &--ok      { background: rgba(22,163,74,0.12);  color: #16a34a; }
}
</style>
