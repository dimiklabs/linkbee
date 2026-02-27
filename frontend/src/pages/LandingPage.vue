<template>
  <div class="landing-page">

    <!-- ─── Navbar ──────────────────────────────────────────────────────────── -->
    <nav class="landing-nav">
      <div class="nav-inner">
        <div class="nav-logo">
          <div class="logo-icon">S</div>
          <span class="logo-text">Shortlink</span>
        </div>

        <div class="nav-links">
          <a href="#features" class="nav-link">Features</a>
          <a href="#pricing" class="nav-link">Pricing</a>
        </div>

        <div class="nav-actions">
          <router-link to="/login">
            <md-outlined-button>Login</md-outlined-button>
          </router-link>
          <router-link to="/signup">
            <md-filled-button>Sign Up</md-filled-button>
          </router-link>
        </div>
      </div>
    </nav>

    <!-- ─── Hero ─────────────────────────────────────────────────────────────── -->
    <section class="hero-section">
      <div class="hero-inner">
        <div class="hero-badge">
          <span class="m3-badge m3-badge--primary">Fast · Reliable · Analytics-powered</span>
        </div>
        <h1 class="hero-title md-display-medium">
          Shorten URLs,<br />
          <span class="hero-highlight">Amplify Reach</span>
        </h1>
        <p class="hero-subtitle md-body-large">
          Create short, memorable links in seconds. Track every click with powerful analytics.
          Share smarter with Shortlink.
        </p>

        <!-- Demo Widget -->
        <div class="demo-widget">
          <div class="m3-card m3-card--elevated demo-card">

            <template v-if="!demoResult">
              <div class="demo-input-row">
                <md-outlined-text-field
                  :value="demoUrl"
                  @input="demoUrl = ($event.target as HTMLInputElement).value"
                  label="Paste your long URL here..."
                  type="url"
                  style="flex: 1;"
                  @keyup.enter="handleDemoShorten"
                />
                <md-filled-button
                  :disabled="demoLoading || !demoUrl"
                  @click="handleDemoShorten"
                  style="height: 56px; white-space: nowrap;"
                >
                  <md-circular-progress v-if="demoLoading" indeterminate style="--md-circular-progress-size:20px; margin-right:8px;" />
                  Shorten
                </md-filled-button>
              </div>
              <div v-if="demoError" class="demo-error md-body-small">{{ demoError }}</div>
            </template>

            <!-- Result state -->
            <div v-else class="demo-result-box">
              <p class="md-body-small" style="color: var(--md-sys-color-on-surface-variant); margin-bottom: 8px;">
                Your short link is ready!
              </p>
              <div class="demo-result-row">
                <a :href="demoResult" target="_blank" rel="noopener noreferrer" class="demo-short-url md-title-medium">
                  {{ demoResult }}
                </a>
                <md-filled-tonal-button @click="copyDemoResult">
                  <span class="material-symbols-outlined" style="font-size:18px; margin-right:4px;">
                    {{ demoCopied ? 'check_circle' : 'content_copy' }}
                  </span>
                  {{ demoCopied ? 'Copied!' : 'Copy' }}
                </md-filled-tonal-button>
              </div>
              <div style="margin-top: 16px;">
                <md-text-button @click="resetDemo">
                  <span class="material-symbols-outlined" style="font-size:18px; margin-right:4px;">arrow_back</span>
                  Shorten another URL
                </md-text-button>
              </div>
            </div>

          </div>
          <p class="demo-hint md-body-small">No account required to try it out.</p>
        </div>
      </div>
    </section>

    <!-- ─── Features ─────────────────────────────────────────────────────────── -->
    <section id="features" class="features-section">
      <div class="section-inner">
        <div class="section-header">
          <h2 class="md-headline-medium section-title">Everything you need</h2>
          <p class="md-body-large" style="color: var(--md-sys-color-on-surface-variant);">
            Powerful features to help you share links smarter.
          </p>
        </div>
        <div class="features-grid">
          <div v-for="feature in features" :key="feature.title" class="m3-card m3-card--outlined feature-card">
            <div class="feature-icon-wrap">
              <span class="material-symbols-outlined feature-icon">{{ feature.materialIcon }}</span>
            </div>
            <h3 class="md-title-medium feature-title">{{ feature.title }}</h3>
            <p class="md-body-medium feature-desc">{{ feature.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ─── Pricing ───────────────────────────────────────────────────────────── -->
    <section id="pricing" class="pricing-section">
      <div class="section-inner">
        <div class="section-header">
          <h2 class="md-headline-medium section-title">Simple, transparent pricing</h2>
          <p class="md-body-large" style="color: var(--md-sys-color-on-surface-variant);">
            Start free, upgrade as you grow. Cancel anytime.
          </p>
        </div>
        <div class="pricing-grid">
          <div
            v-for="plan in pricingPlans"
            :key="plan.name"
            class="pricing-card"
            :class="plan.popular ? 'pricing-card--popular' : 'm3-card--outlined'"
          >
            <div v-if="plan.popular" class="popular-badge">
              <span class="m3-badge m3-badge--primary">Most Popular</span>
            </div>
            <div class="pricing-card-body">
              <div class="pricing-header">
                <h3 class="md-title-large">{{ plan.name }}</h3>
                <div class="price-row">
                  <span class="price-amount">${{ plan.price }}</span>
                  <span class="md-body-medium price-period" style="color: var(--md-sys-color-on-surface-variant);">
                    {{ plan.price > 0 ? '/mo' : '/ forever' }}
                  </span>
                </div>
                <p class="md-body-medium" style="color: var(--md-sys-color-on-surface-variant);">{{ plan.description }}</p>
              </div>
              <ul class="pricing-features">
                <li v-for="feat in plan.features" :key="feat" class="pricing-feature-item">
                  <span class="material-symbols-outlined" style="font-size:18px; color: var(--md-sys-color-primary); flex-shrink: 0;">check_circle</span>
                  <span class="md-body-medium">{{ feat }}</span>
                </li>
              </ul>
              <router-link to="/signup" style="text-decoration: none; display: block; margin-top: auto;">
                <md-filled-button v-if="plan.popular" style="width: 100%;">Get started</md-filled-button>
                <md-outlined-button v-else style="width: 100%;">Get started</md-outlined-button>
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ─── CTA ──────────────────────────────────────────────────────────────── -->
    <section class="cta-section">
      <div class="cta-inner">
        <h2 class="md-headline-medium cta-title">Start shortening today</h2>
        <p class="md-body-large cta-subtitle">
          Join thousands of users who trust Shortlink to power their links.
        </p>
        <router-link to="/signup" style="text-decoration: none;">
          <md-filled-button style="--md-filled-button-container-color: #fff; --md-filled-button-label-text-color: var(--md-sys-color-primary); height: 48px; font-size: 1rem; padding: 0 32px;">
            Get started for free
          </md-filled-button>
        </router-link>
      </div>
    </section>

    <!-- ─── Footer ───────────────────────────────────────────────────────────── -->
    <footer class="landing-footer">
      <div class="footer-inner">
        <div class="footer-logo">
          <div class="logo-icon-sm">S</div>
          <span class="md-label-large" style="color: var(--md-sys-color-on-surface);">Shortlink</span>
        </div>
        <p class="md-body-small" style="color: var(--md-sys-color-on-surface-variant); margin: 0;">
          &copy; 2025 Shortlink. All rights reserved.
        </p>
        <div class="footer-links">
          <a href="#" class="footer-link md-body-small">Privacy</a>
          <a href="#" class="footer-link md-body-small">Terms</a>
          <a href="#" class="footer-link md-body-small">Contact</a>
        </div>
      </div>
    </footer>

  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import linksApi from '@/api/links';

const router = useRouter();

// Demo widget
const demoUrl = ref('');
const demoLoading = ref(false);
const demoError = ref('');
const demoResult = ref('');
const demoCopied = ref(false);

async function handleDemoShorten() {
  if (!demoUrl.value) return;

  demoLoading.value = true;
  demoError.value = '';

  try {
    const response = await linksApi.demoShorten({ destination_url: demoUrl.value });
    demoResult.value = response.data?.short_url || response.short_url || '';
  } catch (err: any) {
    const data = err?.response?.data;
    demoError.value = data?.message || data?.description || 'Failed to shorten URL. Please try again.';
  } finally {
    demoLoading.value = false;
  }
}

async function copyDemoResult() {
  try {
    await navigator.clipboard.writeText(demoResult.value);
    demoCopied.value = true;
    setTimeout(() => { demoCopied.value = false; }, 2500);
  } catch {
    // clipboard not available
  }
}

function resetDemo() {
  demoResult.value = '';
  demoUrl.value = '';
  demoError.value = '';
  demoCopied.value = false;
}

// Features list (6 features, 3 cols × 2 rows)
const features = [
  {
    materialIcon: 'bolt',
    title: 'Fast Redirects',
    description: 'Lightning-fast redirects powered by edge caching. Every click reaches its destination in milliseconds.',
  },
  {
    materialIcon: 'ads_click',
    title: 'Custom Slugs',
    description: 'Create branded short links with custom slugs that reflect your brand and are easy to remember.',
  },
  {
    materialIcon: 'qr_code_2',
    title: 'QR Codes',
    description: 'Automatically generate QR codes for any short link. Perfect for offline materials and print media.',
  },
  {
    materialIcon: 'bar_chart',
    title: 'Click Analytics',
    description: 'Track every click with detailed analytics — referrers, devices, locations, and time-series data.',
  },
  {
    materialIcon: 'lock',
    title: 'Password Protection',
    description: 'Protect sensitive links with a password. Only users with the right password can reach the destination.',
  },
  {
    materialIcon: 'schedule',
    title: 'Link Expiry',
    description: 'Set expiration dates for your links. Automatically deactivate them after a set time or click limit.',
  },
];

// Pricing plans
const pricingPlans = [
  {
    name: 'Free',
    price: 0,
    description: 'Perfect for personal use',
    popular: false,
    features: [
      '50 links/month',
      'Basic analytics',
      'QR codes',
      '30-day history',
    ],
  },
  {
    name: 'Starter',
    price: 9,
    description: 'For growing creators',
    popular: true,
    features: [
      '500 links/month',
      'Full analytics',
      'Custom slugs',
      'Password protection',
      '1 year history',
    ],
  },
  {
    name: 'Pro',
    price: 29,
    description: 'For teams and businesses',
    popular: false,
    features: [
      'Unlimited links',
      'Advanced analytics',
      'Team collaboration',
      'API access',
      'Link expiry',
      'Priority support',
    ],
  },
  {
    name: 'Business',
    price: 79,
    description: 'Enterprise-grade solution',
    popular: false,
    features: [
      'Everything in Pro',
      'Custom domain',
      'SSO / SAML',
      'SLA guarantee',
      'Dedicated support',
      'Audit logs',
    ],
  },
];
</script>

<style scoped lang="scss">
.landing-page {
  background: var(--md-sys-color-background);
  min-height: 100vh;
}

// ─── Navbar ──────────────────────────────────────────────────────────────────
.landing-nav {
  position: sticky;
  top: 0;
  background: var(--md-sys-color-surface);
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  z-index: 100;
}

.nav-inner {
  max-width: 1200px;
  margin: 0 auto;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0 24px;
}

.nav-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
  text-decoration: none;
}

.logo-icon {
  width: 34px;
  height: 34px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  font-weight: 700;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
}

.logo-text {
  font-weight: 700;
  font-size: 1.0625rem;
  color: var(--md-sys-color-on-surface);
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 8px;

  @media (max-width: 640px) {
    display: none;
  }
}

.nav-link {
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface-variant);
  text-decoration: none;
  padding: 6px 12px;
  border-radius: 20px;
  transition: background 0.15s, color 0.15s;

  &:hover {
    background: var(--md-sys-color-surface-container);
    color: var(--md-sys-color-on-surface);
  }
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 8px;

  a {
    text-decoration: none;
  }
}

// ─── Hero ─────────────────────────────────────────────────────────────────────
.hero-section {
  background: var(--md-sys-color-surface-container-low);
  padding: 80px 24px 64px;
}

.hero-inner {
  max-width: 720px;
  margin: 0 auto;
  text-align: center;
}

.hero-badge {
  margin-bottom: 20px;
}

.hero-title {
  color: var(--md-sys-color-on-surface);
  line-height: 1.15;
  margin-bottom: 20px;
  letter-spacing: -0.025em;
}

.hero-highlight {
  color: var(--md-sys-color-primary);
}

.hero-subtitle {
  color: var(--md-sys-color-on-surface-variant);
  max-width: 520px;
  margin: 0 auto 40px;
  line-height: 1.7;
}

.demo-widget {
  max-width: 580px;
  margin: 0 auto;
}

.demo-card {
  padding: 24px;
}

.demo-input-row {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  flex-wrap: wrap;
}

.demo-error {
  color: var(--md-sys-color-error);
  margin-top: 8px;
  text-align: left;
}

.demo-result-box {
  text-align: left;
}

.demo-result-row {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.demo-short-url {
  color: var(--md-sys-color-primary);
  text-decoration: none;
  word-break: break-all;
  flex: 1;

  &:hover {
    text-decoration: underline;
  }
}

.demo-hint {
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 12px;
}

// ─── Features ─────────────────────────────────────────────────────────────────
.features-section {
  padding: 80px 24px;
  background: var(--md-sys-color-surface);
}

.section-inner {
  max-width: 1200px;
  margin: 0 auto;
}

.section-header {
  text-align: center;
  margin-bottom: 48px;

  p {
    margin-top: 8px;
  }
}

.section-title {
  color: var(--md-sys-color-on-surface);
  margin-bottom: 8px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;

  @media (max-width: 900px) {
    grid-template-columns: repeat(2, 1fr);
  }

  @media (max-width: 580px) {
    grid-template-columns: 1fr;
  }
}

.feature-card {
  padding: 24px;
  border-radius: 12px;
  background: var(--md-sys-color-surface-container-low);
  transition: transform 0.2s;

  &:hover {
    transform: translateY(-3px);
  }
}

.feature-icon-wrap {
  width: 48px;
  height: 48px;
  background: var(--md-sys-color-primary-container);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}

.feature-icon {
  font-size: 24px;
  color: var(--md-sys-color-on-primary-container);
}

.feature-title {
  color: var(--md-sys-color-on-surface);
  margin-bottom: 8px;
}

.feature-desc {
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1.6;
}

// ─── Pricing ──────────────────────────────────────────────────────────────────
.pricing-section {
  padding: 80px 24px;
  background: var(--md-sys-color-surface-container-low);
}

.pricing-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;

  @media (max-width: 1100px) {
    grid-template-columns: repeat(2, 1fr);
  }

  @media (max-width: 580px) {
    grid-template-columns: 1fr;
  }
}

.pricing-card {
  border-radius: 12px;
  border: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface);
  position: relative;
  display: flex;
  flex-direction: column;
  transition: transform 0.2s;

  &:hover {
    transform: translateY(-3px);
  }

  &--popular {
    background: var(--md-sys-color-primary-container);
    border-color: var(--md-sys-color-primary);
    border-width: 2px;
    transform: scale(1.02);

    &:hover {
      transform: scale(1.02) translateY(-3px);
    }
  }
}

.pricing-card-body {
  padding: 28px 24px;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.pricing-header {
  margin-bottom: 20px;
}

.price-row {
  display: flex;
  align-items: flex-end;
  gap: 4px;
  margin: 8px 0;
}

.price-amount {
  font-size: 2.25rem;
  font-weight: 800;
  color: var(--md-sys-color-on-surface);
  line-height: 1;
}

.price-period {
  padding-bottom: 4px;
}

.pricing-features {
  list-style: none;
  padding: 0;
  margin: 0 0 24px 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
}

.pricing-feature-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  color: var(--md-sys-color-on-surface);
}

.popular-badge {
  position: absolute;
  top: -14px;
  left: 50%;
  transform: translateX(-50%);
  white-space: nowrap;
}

// ─── CTA ──────────────────────────────────────────────────────────────────────
.cta-section {
  padding: 80px 24px;
  background: linear-gradient(135deg, var(--md-sys-color-primary), #8b5cf6);
  text-align: center;
}

.cta-inner {
  max-width: 600px;
  margin: 0 auto;
}

.cta-title {
  color: #fff;
  margin-bottom: 16px;
}

.cta-subtitle {
  color: rgba(255, 255, 255, 0.85);
  margin-bottom: 32px;
}

// ─── Footer ───────────────────────────────────────────────────────────────────
.landing-footer {
  padding: 24px;
  background: var(--md-sys-color-surface);
  border-top: 1px solid var(--md-sys-color-outline-variant);
}

.footer-inner {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.footer-logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo-icon-sm {
  width: 28px;
  height: 28px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  font-weight: 700;
  font-size: 0.875rem;
  border-radius: 7px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.footer-links {
  display: flex;
  gap: 24px;
}

.footer-link {
  color: var(--md-sys-color-on-surface-variant);
  text-decoration: none;

  &:hover {
    color: var(--md-sys-color-on-surface);
  }
}
</style>
