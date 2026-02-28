<template>
  <div class="landing-page">

    <!-- ─── Navbar ──────────────────────────────────────────────────────────── -->
    <nav class="landing-nav" :class="{ 'landing-nav--scrolled': isScrolled }">
      <div class="nav-inner">
        <!-- Logo -->
        <div class="nav-logo">
          <div class="logo-icon">
            <span class="material-symbols-outlined">link</span>
          </div>
          <span class="logo-text">Linkbee</span>
        </div>

        <!-- Desktop nav links -->
        <div class="nav-links">
          <a href="#features" class="nav-link" @click.prevent="scrollTo('features')">Features</a>
          <a href="#pricing" class="nav-link" @click.prevent="scrollTo('pricing')">Pricing</a>
        </div>

        <!-- Nav actions (desktop) -->
        <div class="nav-actions">
          <router-link to="/login" class="nav-action-link">
            <button class="btn-outlined nav-btn">Login</button>
          </router-link>
          <router-link to="/signup" class="nav-action-link">
            <button class="btn-filled nav-btn">Sign Up</button>
          </router-link>
          <!-- Hamburger (mobile) -->
          <button
            class="hamburger-btn"
            :class="{ 'hamburger-btn--open': mobileMenuOpen }"
            @click="mobileMenuOpen = !mobileMenuOpen"
            aria-label="Toggle navigation"
          >
            <span class="material-symbols-outlined">
              {{ mobileMenuOpen ? 'close' : 'menu' }}
            </span>
          </button>
        </div>
      </div>

      <!-- Mobile slide-down menu -->
      <div class="mobile-menu" :class="{ 'mobile-menu--open': mobileMenuOpen }">
        <div class="mobile-menu-inner">
          <a href="#features" class="mobile-nav-link" @click="handleMobileNav('features')">Features</a>
          <a href="#pricing" class="mobile-nav-link" @click="handleMobileNav('pricing')">Pricing</a>
          <div class="mobile-nav-divider"></div>
          <router-link to="/login" class="mobile-action-link" @click="mobileMenuOpen = false">
            <button class="btn-outlined mobile-action-btn">Login</button>
          </router-link>
          <router-link to="/signup" class="mobile-action-link" @click="mobileMenuOpen = false">
            <button class="btn-filled mobile-action-btn">Sign Up</button>
          </router-link>
        </div>
      </div>
    </nav>

    <!-- ─── Hero ─────────────────────────────────────────────────────────────── -->
    <section class="hero-section">
      <div class="hero-inner">
        <!-- Badge pill -->
        <div class="hero-badge">
          <span class="hero-badge-pill">✨ Trusted by 10,000+ teams</span>
        </div>

        <!-- Headline -->
        <h1 class="hero-title">
          Create Short Links, Track<br class="hero-br" />
          <span class="hero-highlight">Every Click</span>
        </h1>

        <!-- Subtitle -->
        <p class="hero-subtitle">
          Transform long URLs into powerful short links. Gain deep insights with real-time analytics,
          custom slugs, and QR codes — all in one platform.
        </p>

        <!-- CTA buttons -->
        <div class="hero-cta-row">
          <router-link to="/signup" class="hero-cta-link">
            <button class="btn-filled hero-cta-btn hero-cta-btn--primary">
              Get started for free
            </button>
          </router-link>
          <a href="#features" class="hero-cta-link" @click.prevent="scrollTo('features')">
            <button class="btn-outlined hero-cta-btn">
              <span class="material-symbols-outlined hero-play-icon">play_circle</span>
              See how it works
            </button>
          </a>
        </div>

        <!-- Demo widget card -->
        <div class="demo-widget">
          <div class="demo-card">
            <div class="demo-card-header">
              <span class="material-symbols-outlined demo-card-icon">link</span>
              <span class="demo-card-label">Paste any URL to shorten it</span>
            </div>

            <div class="demo-input-row">
              <md-outlined-text-field
                :value="demoUrl"
                @input="demoUrl = ($event.target as HTMLInputElement).value"
                label="Paste your long URL here..."
                type="url"
                class="demo-text-field"
                @keyup.enter="handleDemoShorten"
              />
              <button class="btn-filled demo-shorten-btn" :disabled="!demoUrl" @click="handleDemoShorten">
                Shorten
              </button>
            </div>
          </div>
          <p class="demo-hint">Free to sign up · No credit card needed</p>
        </div>

        <LoginModal ref="loginModalRef" />
      </div>
    </section>

    <!-- ─── Stats / Trust section ─────────────────────────────────────────────── -->
    <section class="stats-section">
      <div class="stats-inner">
        <div v-for="stat in stats" :key="stat.label" class="stat-item">
          <span class="stat-number">{{ stat.number }}</span>
          <span class="stat-label">{{ stat.label }}</span>
        </div>
      </div>
    </section>

    <!-- ─── Features ─────────────────────────────────────────────────────────── -->
    <section id="features" class="features-section">
      <div class="section-inner">
        <div class="section-header">
          <h2 class="section-title">Everything you need</h2>
          <p class="section-subtitle">
            Powerful features to help you share links smarter and track performance effortlessly.
          </p>
        </div>
        <div class="features-grid">
          <div
            v-for="feature in features"
            :key="feature.title"
            class="feature-card"
          >
            <div class="feature-icon-wrap">
              <span class="material-symbols-outlined feature-icon">{{ feature.materialIcon }}</span>
            </div>
            <h3 class="feature-title">{{ feature.title }}</h3>
            <p class="feature-desc">{{ feature.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ─── Pricing ───────────────────────────────────────────────────────────── -->
    <section id="pricing" class="pricing-section">
      <div class="section-inner">
        <div class="section-header">
          <h2 class="section-title">Simple, transparent pricing</h2>
          <p class="section-subtitle">
            Start free, upgrade as you grow. No hidden fees, cancel anytime.
          </p>
        </div>
        <div class="pricing-grid">
          <div
            v-for="plan in pricingPlans"
            :key="plan.name"
            class="pricing-card"
            :class="{ 'pricing-card--popular': plan.popular }"
          >
            <div v-if="plan.popular" class="popular-badge-wrap">
              <span class="popular-badge">Most Popular</span>
            </div>
            <div class="pricing-card-body">
              <div class="pricing-header">
                <h3 class="pricing-plan-name">{{ plan.name }}</h3>
                <div class="price-row">
                  <span class="price-amount">${{ plan.price }}</span>
                  <span class="price-period">
                    {{ plan.price > 0 ? '/mo' : '/ forever' }}
                  </span>
                </div>
                <p class="pricing-description">{{ plan.description }}</p>
              </div>
              <ul class="pricing-features">
                <li
                  v-for="feat in plan.features"
                  :key="feat"
                  class="pricing-feature-item"
                >
                  <span class="material-symbols-outlined pricing-check-icon">check_circle</span>
                  <span class="pricing-feat-text">{{ feat }}</span>
                </li>
              </ul>
              <router-link to="/signup" class="pricing-cta-link">
                <button class="btn-filled pricing-cta-btn" v-if="plan.popular" >Get started</button>
                <button class="btn-outlined pricing-cta-btn" v-else >Get started</button>
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ─── CTA ──────────────────────────────────────────────────────────────── -->
    <section class="cta-section">
      <div class="cta-inner">
        <h2 class="cta-title">Ready to grow your reach?</h2>
        <p class="cta-subtitle">
          Join thousands of teams who trust Linkbee to power their links and analytics.
          Start in seconds — no credit card required.
        </p>
        <router-link to="/signup" class="cta-action-link">
          <button class="btn-filled cta-btn">
            Get started for free
          </button>
        </router-link>
      </div>
    </section>

    <!-- ─── Footer ───────────────────────────────────────────────────────────── -->
    <footer class="landing-footer">
      <div class="footer-inner">
        <div class="footer-logo">
          <div class="footer-logo-icon">
            <span class="material-symbols-outlined footer-logo-symbol">link</span>
          </div>
          <span class="footer-logo-text">Linkbee</span>
        </div>
        <p class="footer-copyright">
          &copy; 2026 Linkbee. All rights reserved.
        </p>
        <div class="footer-links">
          <a href="#" class="footer-link">Privacy</a>
          <a href="#" class="footer-link">Terms</a>
          <a href="#" class="footer-link">Contact</a>
        </div>
      </div>
    </footer>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import LoginModal from '@/components/LoginModal.vue';

// ── Scroll-aware navbar ────────────────────────────────────────────────────────
const isScrolled = ref(false);

function handleScroll() {
  isScrolled.value = window.scrollY > 10;
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll, { passive: true });
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});

// ── Smooth scroll helper ───────────────────────────────────────────────────────
function scrollTo(id: string) {
  const el = document.getElementById(id);
  if (el) {
    el.scrollIntoView({ behavior: 'smooth' });
  }
}

// ── Mobile menu ────────────────────────────────────────────────────────────────
const mobileMenuOpen = ref(false);

function handleMobileNav(id: string) {
  mobileMenuOpen.value = false;
  scrollTo(id);
}

// ── Demo widget ────────────────────────────────────────────────────────────────
const demoUrl = ref('');
const loginModalRef = ref<InstanceType<typeof LoginModal> | null>(null);

function handleDemoShorten() {
  if (!demoUrl.value) return;
  loginModalRef.value?.show();
}

// ── Stats ──────────────────────────────────────────────────────────────────────
const stats = [
  { number: '10M+', label: 'Links Created' },
  { number: '50M+', label: 'Clicks Tracked' },
  { number: '99.9%', label: 'Uptime' },
];

// ── Features list ──────────────────────────────────────────────────────────────
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

// ── Pricing plans ──────────────────────────────────────────────────────────────
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
// ─── Global ────────────────────────────────────────────────────────────────────
.landing-page {
  background: var(--md-sys-color-background);
  min-height: 100vh;
  font-family: inherit;
}

// ─── Navbar ────────────────────────────────────────────────────────────────────
.landing-nav {
  position: sticky;
  top: 0;
  z-index: 200;
  background: var(--md-sys-color-surface);
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  transition: backdrop-filter 0.3s ease, box-shadow 0.3s ease, background 0.3s ease;

  &--scrolled {
    backdrop-filter: blur(16px) saturate(180%);
    -webkit-backdrop-filter: blur(16px) saturate(180%);
    background: rgba(var(--md-sys-color-surface-rgb, 255 255 255) / 0.85);
    box-shadow: 0 1px 8px rgba(0, 0, 0, 0.08);
  }
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
  width: 36px;
  height: 36px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;

  .material-symbols-outlined {
    font-size: 20px;
  }
}

.logo-text {
  font-weight: 700;
  font-size: 1.0625rem;
  color: var(--md-sys-color-on-surface);
  letter-spacing: -0.01em;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 4px;

  @media (max-width: 640px) {
    display: none;
  }
}

.nav-link {
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface-variant);
  text-decoration: none;
  padding: 6px 14px;
  border-radius: 20px;
  transition: background 0.15s ease, color 0.15s ease;

  &:hover {
    background: var(--md-sys-color-surface-container);
    color: var(--md-sys-color-on-surface);
  }
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nav-action-link {
  text-decoration: none;

  @media (max-width: 640px) {
    display: none;
  }
}

.nav-btn {
  white-space: nowrap;
}

// Hamburger button
.hamburger-btn {
  display: none;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 50%;
  color: var(--md-sys-color-on-surface);
  transition: background 0.15s ease;

  &:hover {
    background: var(--md-sys-color-surface-container);
  }

  .material-symbols-outlined {
    font-size: 24px;
  }

  @media (max-width: 640px) {
    display: flex;
  }
}

// Mobile menu
.mobile-menu {
  max-height: 0;
  overflow: hidden;
  transition: max-height 0.3s ease;
  background: var(--md-sys-color-surface);
  border-top: 1px solid transparent;

  &--open {
    max-height: 320px;
    border-top-color: var(--md-sys-color-outline-variant);
  }
}

.mobile-menu-inner {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px 24px 20px;
}

.mobile-nav-link {
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface-variant);
  text-decoration: none;
  padding: 10px 14px;
  border-radius: 10px;
  transition: background 0.15s ease, color 0.15s ease;

  &:hover {
    background: var(--md-sys-color-surface-container);
    color: var(--md-sys-color-on-surface);
  }
}

.mobile-nav-divider {
  height: 1px;
  background: var(--md-sys-color-outline-variant);
  margin: 8px 0;
}

.mobile-action-link {
  text-decoration: none;
  display: block;
}

.mobile-action-btn {
  width: 100%;
}

// ─── Hero ──────────────────────────────────────────────────────────────────────
.hero-section {
  background:
    radial-gradient(ellipse at 30% 50%, rgba(99, 91, 255, 0.12), transparent 70%),
    var(--md-sys-color-surface-container-low);
  padding: 96px 24px 72px;

  @media (max-width: 640px) {
    padding: 64px 20px 48px;
  }
}

.hero-inner {
  max-width: 760px;
  margin: 0 auto;
  text-align: center;
}

// Badge pill
.hero-badge {
  margin-bottom: 24px;
}

.hero-badge-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: var(--md-sys-color-primary-container);
  color: var(--md-sys-color-on-primary-container);
  padding: 6px 16px;
  border-radius: 100px;
  font-size: 0.875rem;
  font-weight: 500;
  letter-spacing: 0.01em;
}

// Headline
.hero-title {
  font-size: clamp(2rem, 5vw, 3.25rem);
  font-weight: 800;
  color: var(--md-sys-color-on-surface);
  line-height: 1.12;
  letter-spacing: -0.03em;
  margin: 0 0 20px;
}

.hero-br {
  @media (max-width: 480px) {
    display: none;
  }
}

.hero-highlight {
  color: var(--md-sys-color-primary);
}

// Subtitle
.hero-subtitle {
  font-size: 1.0625rem;
  color: var(--md-sys-color-on-surface-variant);
  max-width: 540px;
  margin: 0 auto 36px;
  line-height: 1.7;
}

// CTA row
.hero-cta-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 48px;
}

.hero-cta-link {
  text-decoration: none;
}

.hero-cta-btn {
  height: 48px;
  font-size: 0.9375rem;
  padding: 0 24px;
  white-space: nowrap;

  &--primary {
      }
}

.hero-play-icon {
  font-size: 18px;
  margin-right: 6px;
  vertical-align: middle;
}

// Demo widget
.demo-widget {
  max-width: 600px;
  margin: 0 auto;
}

.demo-card {
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06);
}

.demo-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.demo-card-icon {
  font-size: 18px;
  color: var(--md-sys-color-primary);
}

.demo-card-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface-variant);
}

.demo-input-row {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  flex-wrap: wrap;
}

.demo-text-field {
  flex: 1;
  min-width: 0;
}

.demo-shorten-btn {
  height: 56px;
  white-space: nowrap;
  flex-shrink: 0;
}

.demo-spinner {
    margin-right: 8px;
}

.demo-error {
  color: var(--md-sys-color-error);
  font-size: 0.8125rem;
  margin-top: 8px;
  text-align: left;
}

.demo-result-box {
  text-align: left;
}

.demo-result-success-row {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 10px;
}

.demo-result-check {
  font-size: 18px;
  color: var(--md-sys-color-primary);
}

.demo-result-ready {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
}

.demo-result-row {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.demo-short-url {
  color: var(--md-sys-color-primary);
  font-size: 1.0625rem;
  font-weight: 600;
  text-decoration: none;
  word-break: break-all;
  flex: 1;

  &:hover {
    text-decoration: underline;
  }
}

.demo-copy-btn {
  flex-shrink: 0;
}

.demo-copy-icon {
  font-size: 18px;
  margin-right: 4px;
}

.demo-reset-row {
  margin-top: 14px;
}

.demo-back-icon {
  font-size: 18px;
  margin-right: 4px;
}

.demo-hint {
  color: var(--md-sys-color-on-surface-variant);
  font-size: 0.8125rem;
  margin-top: 12px;
  text-align: center;
}

// ─── Stats / Trust ─────────────────────────────────────────────────────────────
.stats-section {
  background: var(--md-sys-color-surface-container);
  border-top: 1px solid var(--md-sys-color-outline-variant);
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  padding: 40px 24px;
}

.stats-inner {
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  text-align: center;

  @media (max-width: 480px) {
    grid-template-columns: 1fr;
    gap: 28px;
  }
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-number {
  font-size: 2rem;
  font-weight: 800;
  color: var(--md-sys-color-primary);
  letter-spacing: -0.03em;
  line-height: 1;
}

.stat-label {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
  font-weight: 500;
}

// ─── Features ──────────────────────────────────────────────────────────────────
.features-section {
  padding: 96px 24px;
  background: var(--md-sys-color-surface);

  @media (max-width: 640px) {
    padding: 64px 20px;
  }
}

.section-inner {
  max-width: 1200px;
  margin: 0 auto;
}

.section-header {
  text-align: center;
  margin-bottom: 56px;
}

.section-title {
  font-size: clamp(1.5rem, 3vw, 2rem);
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
  letter-spacing: -0.02em;
  margin: 0 0 12px;
}

.section-subtitle {
  font-size: 1.0625rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
  line-height: 1.6;
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
  padding: 28px 24px;
  border-radius: 16px;
  border: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface);
  transition: transform 0.2s ease, box-shadow 0.2s ease;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  }
}

.feature-icon-wrap {
  width: 52px;
  height: 52px;
  background: var(--md-sys-color-primary-container);
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 18px;
}

.feature-icon {
  font-size: 26px;
  color: var(--md-sys-color-on-primary-container);
}

.feature-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin: 0 0 8px;
}

.feature-desc {
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1.65;
  margin: 0;
}

// ─── Pricing ───────────────────────────────────────────────────────────────────
.pricing-section {
  padding: 96px 24px;
  background: var(--md-sys-color-surface-container-low);

  @media (max-width: 640px) {
    padding: 64px 20px;
  }
}

.pricing-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  align-items: start;

  @media (max-width: 1100px) {
    grid-template-columns: repeat(2, 1fr);
  }

  @media (max-width: 580px) {
    grid-template-columns: 1fr;
  }
}

.pricing-card {
  border-radius: 16px;
  border: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface);
  position: relative;
  display: flex;
  flex-direction: column;
  transition: transform 0.2s ease, box-shadow 0.2s ease;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  }

  &--popular {
    border-color: var(--md-sys-color-primary);
    border-width: 2px;
    background: var(--md-sys-color-surface);
    box-shadow: 0 4px 24px rgba(99, 91, 255, 0.16);

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 12px 40px rgba(99, 91, 255, 0.2);
    }
  }
}

.popular-badge-wrap {
  position: absolute;
  top: -14px;
  left: 50%;
  transform: translateX(-50%);
  white-space: nowrap;
  z-index: 1;
}

.popular-badge {
  display: inline-block;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  font-size: 0.75rem;
  font-weight: 600;
  padding: 4px 14px;
  border-radius: 100px;
  letter-spacing: 0.03em;
}

.pricing-card-body {
  padding: 32px 24px 24px;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.pricing-header {
  margin-bottom: 20px;
}

.pricing-plan-name {
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
  margin: 0 0 8px;
}

.price-row {
  display: flex;
  align-items: flex-end;
  gap: 4px;
  margin: 0 0 8px;
}

.price-amount {
  font-size: 2.5rem;
  font-weight: 800;
  color: var(--md-sys-color-on-surface);
  line-height: 1;
  letter-spacing: -0.03em;
}

.price-period {
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
  padding-bottom: 4px;
}

.pricing-description {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
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
}

.pricing-check-icon {
  font-size: 18px;
  color: var(--md-sys-color-primary);
  flex-shrink: 0;
  margin-top: 1px;
}

.pricing-feat-text {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);
  line-height: 1.5;
}

.pricing-cta-link {
  text-decoration: none;
  display: block;
  margin-top: auto;
}

.pricing-cta-btn {
  width: 100%;
}

// ─── CTA ───────────────────────────────────────────────────────────────────────
.cta-section {
  padding: 96px 24px;
  background: linear-gradient(135deg, var(--md-sys-color-primary) 0%, #8b5cf6 100%);
  text-align: center;

  @media (max-width: 640px) {
    padding: 64px 20px;
  }
}

.cta-inner {
  max-width: 620px;
  margin: 0 auto;
}

.cta-title {
  font-size: clamp(1.5rem, 3vw, 2.25rem);
  font-weight: 800;
  color: #fff;
  letter-spacing: -0.025em;
  margin: 0 0 16px;
}

.cta-subtitle {
  font-size: 1.0625rem;
  color: rgba(255, 255, 255, 0.85);
  margin: 0 0 36px;
  line-height: 1.65;
}

.cta-action-link {
  text-decoration: none;
}

.cta-btn {
        height: 52px;
  font-size: 1rem;
  padding: 0 36px;
  font-weight: 600;
}

// ─── Footer ────────────────────────────────────────────────────────────────────
.landing-footer {
  padding: 32px 24px;
  background: var(--md-sys-color-surface);
  border-top: 1px solid var(--md-sys-color-outline-variant);
}

.footer-inner {
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  gap: 16px;

  @media (max-width: 640px) {
    grid-template-columns: 1fr;
    text-align: center;
    gap: 14px;
    justify-items: center;
  }
}

.footer-logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.footer-logo-icon {
  width: 28px;
  height: 28px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  border-radius: 7px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.footer-logo-symbol {
  font-size: 16px;
}

.footer-logo-text {
  font-weight: 700;
  font-size: 0.9375rem;
  color: var(--md-sys-color-on-surface);
}

.footer-copyright {
  font-size: 0.8125rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
}

.footer-links {
  display: flex;
  gap: 24px;
  justify-content: flex-end;

  @media (max-width: 640px) {
    justify-content: center;
  }
}

.footer-link {
  font-size: 0.8125rem;
  color: var(--md-sys-color-on-surface-variant);
  text-decoration: none;
  transition: color 0.15s ease;

  &:hover {
    color: var(--md-sys-color-on-surface);
  }
}
</style>
