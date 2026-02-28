<template>
  <div class="legal-page">
    <nav class="legal-nav">
      <div class="legal-nav-inner">
        <router-link to="/" class="legal-logo">
          <img src="/logo.png" alt="Linkbee" class="legal-logo-img" />
        </router-link>
        <div class="legal-nav-links">
          <router-link to="/terms" class="legal-nav-link">Terms of Service</router-link>
          <router-link to="/login" class="legal-nav-link">Login</router-link>
        </div>
      </div>
    </nav>

    <main class="legal-main">
      <div class="legal-container">
        <header class="legal-header">
          <h1 class="legal-title">Privacy Policy</h1>
          <p class="legal-meta">Effective date: February 28, 2026 &nbsp;·&nbsp; Last updated: February 28, 2026</p>
        </header>

        <div class="legal-body">

          <section>
            <h2>1. Introduction</h2>
            <p>
              Linkbee ("we," "us," or "our") is committed to protecting your personal information. This Privacy Policy explains what data we collect, how we use it, how we share it, and your rights in relation to it when you use the Linkbee platform and related services (the "Service").
            </p>
            <p>
              By using the Service you agree to the collection and use of information in accordance with this policy.
            </p>
          </section>

          <section>
            <h2>2. Information We Collect</h2>

            <h3>2.1 Information you provide directly</h3>
            <ul>
              <li><strong>Account data</strong> — name, email address, password (stored as a hashed value), and profile details provided during registration or in account settings.</li>
              <li><strong>Content data</strong> — destination URLs, custom slugs, link titles, bio page content, and any other content you create within the Service.</li>
              <li><strong>Payment data</strong> — billing address and payment method details collected by our payment processor on our behalf. We do not store full card numbers on our servers.</li>
              <li><strong>Communications</strong> — messages you send us via email or support channels.</li>
            </ul>

            <h3>2.2 Information collected automatically</h3>
            <ul>
              <li><strong>Usage data</strong> — pages visited, features used, session duration, and actions performed within the Service.</li>
              <li><strong>Click analytics</strong> — when a visitor clicks one of your shortened links, we record a hashed (anonymised) version of their IP address, approximate geographic location (country/city), device type, browser, operating system, and referrer URL. Raw IP addresses are not stored.</li>
              <li><strong>Log data</strong> — server logs including request timestamps, HTTP method, URL path, response status, and error information, retained for security and debugging purposes.</li>
              <li><strong>Cookies and local storage</strong> — we use browser local storage to persist your authentication session and UI preferences (e.g., dark mode). We do not use third-party advertising cookies.</li>
            </ul>
          </section>

          <section>
            <h2>3. How We Use Your Information</h2>
            <p>We use the information we collect to:</p>
            <ul>
              <li>Provide, operate, and improve the Service</li>
              <li>Authenticate your account and maintain security</li>
              <li>Process payments and manage your subscription</li>
              <li>Provide click analytics and reporting features to you</li>
              <li>Send transactional emails (account verification, password resets, billing receipts)</li>
              <li>Respond to support requests and communications</li>
              <li>Monitor and enforce compliance with our Terms of Service</li>
              <li>Comply with legal obligations</li>
            </ul>
            <p>We do not sell your personal data to third parties.</p>
          </section>

          <section>
            <h2>4. Legal Bases for Processing (GDPR)</h2>
            <p>If you are located in the European Economic Area (EEA) or United Kingdom, we process your personal data on the following legal bases:</p>
            <ul>
              <li><strong>Contract performance</strong> — processing necessary to provide the Service you have subscribed to.</li>
              <li><strong>Legitimate interests</strong> — operating and improving the Service, preventing fraud, and ensuring security.</li>
              <li><strong>Legal obligation</strong> — complying with applicable laws and regulations.</li>
              <li><strong>Consent</strong> — where we have asked for your consent (e.g., optional marketing communications), which you may withdraw at any time.</li>
            </ul>
          </section>

          <section>
            <h2>5. Sharing of Information</h2>
            <p>We may share your information with:</p>
            <ul>
              <li><strong>Payment processors</strong> — to handle subscription billing. They are bound by their own privacy policies and data processing agreements.</li>
              <li><strong>Infrastructure and hosting providers</strong> — cloud services used to host and operate the Service, under data processing agreements.</li>
              <li><strong>Email service providers</strong> — to deliver transactional emails on our behalf.</li>
              <li><strong>Law enforcement or regulators</strong> — where required by applicable law, court order, or legal process.</li>
              <li><strong>Business transfers</strong> — in connection with a merger, acquisition, or sale of assets, your data may be transferred as part of that transaction.</li>
            </ul>
            <p>We do not share your personal data with advertisers or data brokers.</p>
          </section>

          <section>
            <h2>6. Click Analytics Data</h2>
            <p>
              When visitors click your shortened links, we collect analytics data to power your dashboard reports. This data is associated with your account (as the link owner) and is never shared with third parties for advertising purposes.
            </p>
            <p>
              IP addresses of link visitors are hashed using a one-way cryptographic function before storage. We cannot reverse this hash to identify individual visitors.
            </p>
          </section>

          <section>
            <h2>7. Data Retention</h2>
            <p>
              We retain your account data for as long as your account is active or as needed to provide the Service. Upon account deletion, we will delete or anonymise your personal data within 30 days, except where retention is required by law.
            </p>
            <p>
              Anonymised and aggregated analytics data (from which individual visitors cannot be identified) may be retained indefinitely for product improvement purposes.
            </p>
            <p>Server logs are retained for up to 90 days.</p>
          </section>

          <section>
            <h2>8. Security</h2>
            <p>
              We implement industry-standard security measures including TLS encryption in transit, bcrypt password hashing, JWT-based authentication with short-lived access tokens, and token blacklisting on logout.
            </p>
            <p>
              While we take reasonable steps to protect your information, no method of transmission over the internet or method of electronic storage is 100% secure. We cannot guarantee absolute security.
            </p>
          </section>

          <section>
            <h2>9. Your Rights</h2>
            <p>Depending on your location, you may have the following rights regarding your personal data:</p>
            <ul>
              <li><strong>Access</strong> — request a copy of the personal data we hold about you.</li>
              <li><strong>Correction</strong> — request that we correct inaccurate or incomplete data.</li>
              <li><strong>Deletion</strong> — request that we delete your personal data ("right to be forgotten").</li>
              <li><strong>Portability</strong> — request a machine-readable export of your data.</li>
              <li><strong>Restriction</strong> — request that we restrict processing of your data in certain circumstances.</li>
              <li><strong>Objection</strong> — object to processing based on legitimate interests.</li>
              <li><strong>Withdraw consent</strong> — where processing is based on consent, withdraw it at any time without affecting prior lawful processing.</li>
            </ul>
            <p>
              To exercise any of these rights, contact us at <a href="mailto:privacy@linkbee.click">privacy@linkbee.click</a>. We will respond within 30 days.
            </p>
          </section>

          <section>
            <h2>10. Cookies and Tracking</h2>
            <p>
              Linkbee uses browser <strong>local storage</strong> (not cookies) to store your authentication tokens and UI preferences. We do not use tracking pixels, third-party analytics scripts, or advertising cookies.
            </p>
            <p>
              Our public redirect infrastructure may set no cookies on link visitors' browsers.
            </p>
          </section>

          <section>
            <h2>11. Children's Privacy</h2>
            <p>
              The Service is not directed to individuals under the age of 18. We do not knowingly collect personal data from children. If you become aware that a child has provided us with personal data, please contact us and we will take steps to delete such information.
            </p>
          </section>

          <section>
            <h2>12. International Data Transfers</h2>
            <p>
              Your information may be transferred to and processed in countries other than your own. We take appropriate safeguards to ensure that your data is treated securely and in accordance with this policy and applicable data protection law.
            </p>
          </section>

          <section>
            <h2>13. Changes to This Policy</h2>
            <p>
              We may update this Privacy Policy from time to time. When we make material changes, we will update the "Last updated" date and notify you via email or an in-app notice. Your continued use of the Service after changes constitutes acceptance of the revised policy.
            </p>
          </section>

          <section>
            <h2>14. Contact Us</h2>
            <p>If you have questions or concerns about this Privacy Policy or our data practices, please contact us:</p>
            <ul>
              <li>Email: <a href="mailto:privacy@linkbee.click">privacy@linkbee.click</a></li>
              <li>Support: <a href="mailto:support@linkbee.click">support@linkbee.click</a></li>
              <li>Website: <a href="https://linkbee.click">https://linkbee.click</a></li>
            </ul>
          </section>

        </div>
      </div>
    </main>

    <footer class="legal-footer">
      <p>© {{ new Date().getFullYear() }} Linkbee. All rights reserved. &nbsp;·&nbsp;
        <router-link to="/terms">Terms of Service</router-link> &nbsp;·&nbsp;
        <router-link to="/privacy">Privacy Policy</router-link>
      </p>
    </footer>
  </div>
</template>

<script setup lang="ts">
// Static page — no script logic needed
</script>

<style scoped lang="scss">
.legal-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--md-sys-color-background);
  color: var(--md-sys-color-on-background);
}

/* ── Nav ───────────────────────────────────────────────────────────────────── */
.legal-nav {
  position: sticky;
  top: 0;
  z-index: 100;
  background: var(--md-sys-color-surface);
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.legal-nav-inner {
  max-width: 860px;
  margin: 0 auto;
  padding: 0 24px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.legal-logo {
  display: flex;
  align-items: center;
  text-decoration: none;
}

.legal-logo-img {
  height: 32px;
}

.legal-nav-links {
  display: flex;
  align-items: center;
  gap: 24px;
}

.legal-nav-link {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
  text-decoration: none;
  font-weight: 500;

  &:hover { color: var(--md-sys-color-primary); }
}

/* ── Main ──────────────────────────────────────────────────────────────────── */
.legal-main {
  flex: 1;
  padding: 48px 24px 64px;
}

.legal-container {
  max-width: 780px;
  margin: 0 auto;
}

/* ── Header ────────────────────────────────────────────────────────────────── */
.legal-header {
  margin-bottom: 40px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.legal-title {
  font-size: 2rem;
  font-weight: 700;
  margin: 0 0 8px;
  color: var(--md-sys-color-on-surface);
}

.legal-meta {
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
}

/* ── Body ──────────────────────────────────────────────────────────────────── */
.legal-body {
  display: flex;
  flex-direction: column;
  gap: 36px;
  font-size: 0.9375rem;
  line-height: 1.75;
  color: var(--md-sys-color-on-surface);

  section { display: flex; flex-direction: column; gap: 12px; }

  h2 {
    font-size: 1.125rem;
    font-weight: 700;
    margin: 0;
    color: var(--md-sys-color-on-surface);
    padding-bottom: 6px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  h3 {
    font-size: 0.9375rem;
    font-weight: 600;
    margin: 4px 0 0;
    color: var(--md-sys-color-on-surface);
  }

  p { margin: 0; }

  ul {
    margin: 0;
    padding-left: 24px;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  a {
    color: var(--md-sys-color-primary);
    text-decoration: underline;
    text-underline-offset: 2px;
  }

  strong { font-weight: 600; }
}

/* ── Footer ────────────────────────────────────────────────────────────────── */
.legal-footer {
  padding: 20px 24px;
  border-top: 1px solid var(--md-sys-color-outline-variant);
  text-align: center;
  font-size: 0.8125rem;
  color: var(--md-sys-color-on-surface-variant);

  a {
    color: var(--md-sys-color-on-surface-variant);
    text-decoration: none;

    &:hover { color: var(--md-sys-color-primary); }
  }
}
</style>
