<template>
  <div class="landing-page">

    <!-- ─── Navbar ──────────────────────────────────────────────────────────── -->
    <nav class="landing-nav">
      <div class="nav-inner container">
        <div class="nav-logo">
          <div class="logo-icon">S</div>
          <span class="logo-text">Shortlink</span>
        </div>

        <div class="nav-links d-none d-md-flex align-items-center gap-3">
          <a href="#features" class="nav-link text-muted text-decoration-none">Features</a>
          <a href="#pricing" class="nav-link text-muted text-decoration-none">Pricing</a>
        </div>

        <div class="nav-actions d-flex align-items-center gap-2">
          <router-link to="/login" class="btn btn-outline-secondary btn-sm">Login</router-link>
          <router-link to="/signup" class="btn btn-sm btn-primary">Sign Up</router-link>
        </div>
      </div>
    </nav>

    <!-- ─── Hero ─────────────────────────────────────────────────────────────── -->
    <section class="hero-section">
      <div class="hero-inner container text-center">
        <div class="hero-badge mb-3">
          <span class="badge rounded-pill px-3 py-2" style="background: rgba(99,91,255,0.12); color: #635bff; font-size: 0.8rem; font-weight: 600;">
            Fast · Reliable · Analytics-powered
          </span>
        </div>
        <h1 class="hero-title">
          Shorten URLs,<br />
          <span class="hero-highlight">Amplify Reach</span>
        </h1>
        <p class="hero-subtitle text-muted">
          Create short, memorable links in seconds. Track every click with powerful analytics.
          Share smarter with Shortlink.
        </p>

        <!-- Demo Widget -->
        <div class="demo-widget mx-auto">
          <div class="card shadow border-0" style="border-radius: 16px;">
            <div class="card-body p-4">

              <template v-if="!demoResult">
                <div class="d-flex gap-2 flex-column flex-sm-row">
                  <input
                    v-model="demoUrl"
                    type="url"
                    class="form-control form-control-lg flex-fill"
                    placeholder="Paste your long URL here..."
                    @keyup.enter="handleDemoShorten"
                  />
                  <button
                    class="btn btn-primary btn-lg px-4"
                    style="white-space: nowrap;"
                    :disabled="demoLoading || !demoUrl"
                    @click="handleDemoShorten"
                  >
                    <span v-if="demoLoading" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                    Shorten
                  </button>
                </div>
                <div v-if="demoError" class="text-danger small mt-2 text-start">{{ demoError }}</div>
              </template>

              <!-- Result state -->
              <div v-else class="demo-result-box">
                <div class="text-muted small mb-2">Your short link is ready 🎉</div>
                <div class="d-flex align-items-center gap-3 flex-wrap">
                  <a :href="demoResult" target="_blank" rel="noopener noreferrer" class="demo-short-url fw-semibold flex-fill">
                    {{ demoResult }}
                  </a>
                  <button
                    class="btn btn-sm flex-shrink-0"
                    :class="demoCopied ? 'btn-success' : 'btn-outline-secondary'"
                    @click="copyDemoResult"
                  >
                    {{ demoCopied ? '✓ Copied!' : 'Copy' }}
                  </button>
                </div>
                <div class="mt-3 text-start">
                  <button class="btn btn-sm btn-link p-0 text-muted text-decoration-none" @click="resetDemo">
                    ← Shorten another URL
                  </button>
                </div>
              </div>

            </div>
          </div>
          <p class="text-muted small mt-2">No account required to try it out.</p>
        </div>
      </div>
    </section>

    <!-- ─── Features ─────────────────────────────────────────────────────────── -->
    <section id="features" class="features-section">
      <div class="container">
        <div class="text-center mb-5">
          <h2 class="section-title">Everything you need</h2>
          <p class="text-muted">Powerful features to help you share links smarter.</p>
        </div>
        <div class="row g-4">
          <div v-for="feature in features" :key="feature.title" class="col-md-6 col-lg-4">
            <div class="feature-card card h-100 border-0 shadow-sm">
              <div class="card-body p-4">
                <div class="feature-icon-wrap mb-3">{{ feature.icon }}</div>
                <h5 class="fw-semibold mb-2">{{ feature.title }}</h5>
                <p class="text-muted small mb-0">{{ feature.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ─── Pricing ───────────────────────────────────────────────────────────── -->
    <section id="pricing" class="pricing-section">
      <div class="container">
        <div class="text-center mb-5">
          <h2 class="section-title">Simple, transparent pricing</h2>
          <p class="text-muted">Start free, upgrade as you grow. Cancel anytime.</p>
        </div>
        <div class="row g-4 justify-content-center">
          <div v-for="plan in pricingPlans" :key="plan.name" class="col-sm-6 col-xl-3">
            <div
              class="pricing-card card h-100 border-0 shadow-sm"
              :class="{ 'pricing-card--popular': plan.popular }"
            >
              <div v-if="plan.popular" class="popular-badge">Most Popular</div>
              <div class="card-body p-4 d-flex flex-column">
                <div class="mb-4">
                  <h5 class="fw-semibold mb-1">{{ plan.name }}</h5>
                  <div class="d-flex align-items-end gap-1 my-2">
                    <span class="price-amount">${{ plan.price }}</span>
                    <span class="text-muted small pb-1" v-if="plan.price > 0">/mo</span>
                    <span class="text-muted small pb-1" v-else>/ forever</span>
                  </div>
                  <p class="text-muted small mb-0">{{ plan.description }}</p>
                </div>
                <ul class="list-unstyled mb-4 flex-fill">
                  <li v-for="feat in plan.features" :key="feat" class="d-flex align-items-start gap-2 mb-2">
                    <span class="mt-1 flex-shrink-0" style="color: #635bff; font-size: 0.85rem;">✓</span>
                    <span class="small text-muted">{{ feat }}</span>
                  </li>
                </ul>
                <router-link
                  to="/signup"
                  class="btn w-100"
                  :class="plan.popular ? 'btn-primary' : 'btn-outline-secondary'"
                >
                  Get started
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ─── CTA ──────────────────────────────────────────────────────────────── -->
    <section class="cta-section text-center">
      <div class="container">
        <h2 class="section-title text-white mb-3">Start shortening today</h2>
        <p class="mb-4" style="color: rgba(255,255,255,0.75);">
          Join thousands of users who trust Shortlink to power their links.
        </p>
        <router-link to="/signup" class="btn btn-light btn-lg px-5 fw-semibold" style="color: #635bff;">
          Get started for free
        </router-link>
      </div>
    </section>

    <!-- ─── Footer ───────────────────────────────────────────────────────────── -->
    <footer class="landing-footer">
      <div class="container">
        <div class="d-flex flex-column flex-md-row justify-content-between align-items-center gap-3">
          <div class="d-flex align-items-center gap-2">
            <div class="logo-icon-sm">S</div>
            <span class="fw-bold" style="color: #1a1f36;">Shortlink</span>
          </div>
          <p class="text-muted small mb-0">&copy; 2025 Shortlink. All rights reserved.</p>
          <div class="d-flex gap-3">
            <a href="#" class="text-muted small text-decoration-none">Privacy</a>
            <a href="#" class="text-muted small text-decoration-none">Terms</a>
            <a href="#" class="text-muted small text-decoration-none">Contact</a>
          </div>
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
    icon: '⚡',
    title: 'Fast Redirects',
    description: 'Lightning-fast redirects powered by edge caching. Every click reaches its destination in milliseconds.',
  },
  {
    icon: '🎯',
    title: 'Custom Slugs',
    description: 'Create branded short links with custom slugs that reflect your brand and are easy to remember.',
  },
  {
    icon: '📱',
    title: 'QR Codes',
    description: 'Automatically generate QR codes for any short link. Perfect for offline materials and print media.',
  },
  {
    icon: '📊',
    title: 'Click Analytics',
    description: 'Track every click with detailed analytics — referrers, devices, locations, and time-series data.',
  },
  {
    icon: '🔒',
    title: 'Password Protection',
    description: 'Protect sensitive links with a password. Only users with the right password can reach the destination.',
  },
  {
    icon: '⏳',
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
$primary: #635bff;
$primary-dark: #5851db;

.landing-page {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

// ─── Navbar ──────────────────────────────────────────────────────────────────
.landing-nav {
  position: sticky;
  top: 0;
  background: rgba(255, 255, 255, 0.96);
  backdrop-filter: blur(8px);
  border-bottom: 1px solid #e3e8ee;
  z-index: 100;
}

.nav-inner {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.nav-logo {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  flex-shrink: 0;

  .logo-icon {
    width: 34px;
    height: 34px;
    background: $primary;
    color: #fff;
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
    color: #1a1f36;
  }
}

.nav-link {
  font-size: 0.9rem;
  font-weight: 500;
  transition: color 0.15s;

  &:hover {
    color: $primary !important;
  }
}

// ─── Hero ─────────────────────────────────────────────────────────────────────
.hero-section {
  background: linear-gradient(160deg, #f7f9fc 0%, #eef0ff 60%, #f0faff 100%);
  padding: 5.5rem 1rem 4.5rem;
}

.hero-inner {
  max-width: 720px;
  margin: 0 auto;
}

.hero-title {
  font-size: clamp(2.25rem, 5vw, 3.5rem);
  font-weight: 800;
  color: #1a1f36;
  line-height: 1.15;
  margin-bottom: 1.25rem;
  letter-spacing: -0.025em;
}

.hero-highlight {
  background: linear-gradient(135deg, $primary, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-subtitle {
  font-size: 1.1rem;
  line-height: 1.7;
  max-width: 520px;
  margin-left: auto;
  margin-right: auto;
  margin-bottom: 2.5rem;
}

.demo-widget {
  max-width: 580px;
}

.demo-result-box {
  text-align: left;
}

.demo-short-url {
  color: $primary;
  text-decoration: none;
  font-size: 1.05rem;
  word-break: break-all;

  &:hover {
    text-decoration: underline;
  }
}

// ─── Features ─────────────────────────────────────────────────────────────────
.features-section {
  padding: 5rem 0;
  background: #fff;
}

.feature-card {
  border-radius: 12px;
  transition: transform 0.2s, box-shadow 0.2s;

  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1) !important;
  }
}

.feature-icon-wrap {
  font-size: 2.25rem;
  line-height: 1;
}

// ─── Pricing ──────────────────────────────────────────────────────────────────
.pricing-section {
  padding: 5rem 0;
  background: #f7f9fc;
}

.pricing-card {
  border-radius: 12px;
  position: relative;
  transition: transform 0.2s, box-shadow 0.2s;

  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 12px 36px rgba(0, 0, 0, 0.12) !important;
  }

  &--popular {
    border: 2px solid $primary !important;
    transform: scale(1.02);

    &:hover {
      transform: scale(1.02) translateY(-3px);
    }
  }
}

.popular-badge {
  position: absolute;
  top: -12px;
  left: 50%;
  transform: translateX(-50%);
  background: $primary;
  color: #fff;
  font-size: 0.7rem;
  font-weight: 600;
  padding: 4px 14px;
  border-radius: 20px;
  white-space: nowrap;
  letter-spacing: 0.5px;
}

.price-amount {
  font-size: 2.25rem;
  font-weight: 800;
  color: #1a1f36;
  line-height: 1;
}

// ─── CTA ──────────────────────────────────────────────────────────────────────
.cta-section {
  padding: 5.5rem 1rem;
  background: linear-gradient(135deg, $primary, #8b5cf6);
}

// ─── Footer ───────────────────────────────────────────────────────────────────
.landing-footer {
  padding: 2rem 0;
  background: #fff;
  border-top: 1px solid #e3e8ee;
}

.logo-icon-sm {
  width: 28px;
  height: 28px;
  background: $primary;
  color: #fff;
  font-weight: 700;
  font-size: 0.875rem;
  border-radius: 7px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.section-title {
  font-size: clamp(1.5rem, 3vw, 2.25rem);
  font-weight: 800;
  color: #1a1f36;
  letter-spacing: -0.02em;
  margin-bottom: 0.5rem;
}
</style>
