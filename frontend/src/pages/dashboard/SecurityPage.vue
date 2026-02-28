<template>
  <div class="security-page">

    <!-- Page header -->
    <div class="page-header">
      <h1 class="page-title">Security</h1>
      <p class="page-subtitle">Manage two-factor authentication and account security.</p>
    </div>

    <!-- Loading -->
    <div v-if="loadingStatus" class="loading-center">
      <div class="css-spinner"></div>
      <p class="loading-label">Loading security settings…</p>
    </div>

    <template v-else>

      <!-- ── 2FA Card ───────────────────────────────────────────────────── -->
      <div class="an-card">
        <div class="an-card-header">
          <div class="an-card-icon an-card-icon--primary">
            <span class="material-symbols-outlined">shield</span>
          </div>
          <span class="an-card-title">Two-Factor Authentication (2FA)</span>
        </div>
        <div class="an-card-body">
          <div class="twofa-status-row">
            <div class="twofa-status-info">
              <div class="twofa-status-badge-row">
                <div :class="['twofa-status-dot', totpEnabled ? 'twofa-status-dot--on' : 'twofa-status-dot--off']"></div>
                <span v-if="totpEnabled" class="m3-badge m3-badge--success">Enabled</span>
                <span v-else class="m3-badge m3-badge--neutral">Disabled</span>
              </div>
              <p class="twofa-description">
                Add an extra layer of security to your account. When enabled, you'll need a 6-digit
                code from an authenticator app (e.g. Google Authenticator, Authy) at every login.
              </p>
            </div>
            <div class="twofa-status-action">
              <button class="btn-filled" v-if="!totpEnabled" @click="startSetup">
                <span class="material-symbols-outlined">lock</span>
                Enable 2FA
              </button>
              <button class="btn-outlined btn-danger" v-else @click="showDisableModal = true">
                <span class="material-symbols-outlined">lock_open</span>
                Disable 2FA
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- ── Setup Wizard ──────────────────────────────────────────────── -->
      <div v-if="setupStep && !totpEnabled" class="an-card">
        <div class="an-card-header">
          <div class="an-card-icon an-card-icon--neutral">
            <span class="material-symbols-outlined">manage_accounts</span>
          </div>
          <span class="an-card-title">Setup Two-Factor Authentication</span>
        </div>

        <!-- Step indicator -->
        <div class="wizard-steps">
          <div class="wizard-step" :class="{ 'wizard-step--active': setupStep === 'scan', 'wizard-step--done': setupStep === 'confirm' || setupStep === 'backup' }">
            <div class="wizard-step__dot">
              <span v-if="setupStep === 'confirm' || setupStep === 'backup'" class="material-symbols-outlined wizard-check-icon">check</span>
              <span v-else>1</span>
            </div>
            <span class="wizard-step__label">Scan QR</span>
          </div>
          <div class="wizard-step__connector" :class="{ 'wizard-step__connector--done': setupStep === 'confirm' || setupStep === 'backup' }"></div>
          <div class="wizard-step" :class="{ 'wizard-step--active': setupStep === 'confirm', 'wizard-step--done': setupStep === 'backup' }">
            <div class="wizard-step__dot">
              <span v-if="setupStep === 'backup'" class="material-symbols-outlined wizard-check-icon">check</span>
              <span v-else>2</span>
            </div>
            <span class="wizard-step__label">Verify Code</span>
          </div>
          <div class="wizard-step__connector" :class="{ 'wizard-step__connector--done': setupStep === 'backup' }"></div>
          <div class="wizard-step" :class="{ 'wizard-step--active': setupStep === 'backup' }">
            <div class="wizard-step__dot">3</div>
            <span class="wizard-step__label">Save Codes</span>
          </div>
        </div>

        <div class="an-card-body">

          <!-- Step 1: Scan QR -->
          <div v-if="setupStep === 'scan'" class="wizard-panel">
            <h2 class="wizard-panel__title">Scan this QR code</h2>
            <p class="wizard-panel__subtitle">
              Open your authenticator app and scan the QR code below. If you can't scan,
              use the manual secret key instead.
            </p>

            <div class="qr-container">
              <div class="qr-frame">
                <img :src="qrImageUrl" alt="TOTP QR Code" width="180" height="180" class="qr-image" />
              </div>
              <p class="qr-helper">Scan with Google Authenticator, Authy, or similar</p>
            </div>

            <div class="secret-key-block">
              <div class="secret-key-label">
                <span class="material-symbols-outlined secret-key-icon">vpn_key</span>
                Manual secret key
              </div>
              <div class="secret-key-row">
                <label class="form-field form-field--mono">
                  <span class="form-field__label">Secret key</span>
                  <input type="text" class="form-input form-input--mono" :value="totpSecret" readonly />
                </label>
                <button class="btn-outlined" @click="copySecret">
                  <span class="material-symbols-outlined">{{ copied ? 'check' : 'content_copy' }}</span>
                  {{ copied ? 'Copied' : 'Copy' }}
                </button>
              </div>
            </div>

            <div class="wizard-actions">
              <button class="btn-filled" @click="setupStep = 'confirm'">
                I've scanned it
                <span class="material-symbols-outlined">arrow_forward</span>
              </button>
            </div>
          </div>

          <!-- Step 2: Confirm code -->
          <div v-else-if="setupStep === 'confirm'" class="wizard-panel">
            <h2 class="wizard-panel__title">Enter the verification code</h2>
            <p class="wizard-panel__subtitle">
              Enter the 6-digit code shown in your authenticator app to verify the setup is working correctly.
            </p>

            <div v-if="setupError" class="feedback-error">
              <span class="material-symbols-outlined feedback-icon">error</span>
              {{ setupError }}
            </div>

            <div class="code-input-wrap">
              <label class="form-field form-field--narrow">
                <span class="form-field__label">6-digit code</span>
                <input
                  type="text"
                  class="form-input form-input--mono"
                  :value="confirmCode"
                  @input="confirmCode = ($event.target as HTMLInputElement).value"
                  maxlength="6"
                  autocomplete="one-time-code"
                  placeholder="000000"
                />
              </label>
              <p class="code-hint">The code refreshes every 30 seconds.</p>
            </div>

            <div class="wizard-actions">
              <button class="btn-outlined" @click="setupStep = 'scan'">
                <span class="material-symbols-outlined">arrow_back</span>
                Back
              </button>
              <button class="btn-filled" :disabled="confirmLoading" @click="confirmSetup">
                <div v-if="confirmLoading" class="css-spinner css-spinner--sm css-spinner--white"></div>
                <span class="material-symbols-outlined" v-else>verified</span>
                Confirm &amp; Enable
              </button>
            </div>
          </div>

          <!-- Step 3: Backup codes -->
          <div v-else-if="setupStep === 'backup'" class="wizard-panel">
            <div class="success-banner">
              <span class="material-symbols-outlined success-banner__icon">verified</span>
              <div>
                <div class="success-banner__title">2FA Successfully Enabled!</div>
                <div class="success-banner__subtitle">Your account is now protected with two-factor authentication.</div>
              </div>
            </div>

            <div class="backup-warning">
              <span class="material-symbols-outlined backup-warning__icon">warning</span>
              <div class="backup-warning__text">
                <strong>Save these backup codes now.</strong> They are shown only once. Use them to
                access your account if you ever lose your authenticator app.
              </div>
            </div>

            <div class="backup-codes-grid">
              <code v-for="code in backupCodes" :key="code" class="backup-code">{{ code }}</code>
            </div>

            <div class="wizard-actions">
              <button class="btn-outlined" @click="copyBackupCodes">
                <span class="material-symbols-outlined">{{ copiedBackup ? 'check' : 'content_copy' }}</span>
                {{ copiedBackup ? 'Copied all!' : 'Copy all codes' }}
              </button>
              <button class="btn-filled" @click="finishSetup">
                <span class="material-symbols-outlined">done</span>
                Done
              </button>
            </div>
          </div>

        </div>
      </div>

    </template>

    <!-- ── Disable 2FA Dialog ────────────────────────────────────────────── -->
    <BaseModal v-model="showDisableModal" size="sm" @closed="closeDisableModal">
      <template #headline>Disable Two-Factor Authentication</template>
      <div class="disable-modal-body">
        <p class="disable-modal-desc">
          Enter your account password to confirm you want to disable 2FA.
        </p>
        <div v-if="disableError" class="feedback-error disable-modal-error">
          <span class="material-symbols-outlined feedback-icon">error</span>
          {{ disableError }}
        </div>
        <label class="form-field">
          <span class="form-field__label">Password</span>
          <input
            type="password"
            class="form-input"
            :value="disablePassword"
            @input="disablePassword = ($event.target as HTMLInputElement).value"
            autocomplete="current-password"
          />
        </label>
      </div>
      <template #actions>
        <button class="btn-text" @click="closeDisableModal">Cancel</button>
        <button class="btn-filled btn-danger" :disabled="disableLoading" @click="disableTOTP">
          <div v-if="disableLoading" class="css-spinner css-spinner--sm css-spinner--white"></div>
          Disable 2FA
        </button>
      </template>
    </BaseModal>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import BaseModal from '@/components/BaseModal.vue';
import authApi from '@/api/auth';

const loadingStatus = ref(true);
const totpEnabled = ref(false);

// Setup wizard
const setupStep = ref<'scan' | 'confirm' | 'backup' | null>(null);
const totpSecret = ref('');
const qrCodeURL = ref('');
const confirmCode = ref('');
const confirmLoading = ref(false);
const setupError = ref('');
const backupCodes = ref<string[]>([]);
const copied = ref(false);
const copiedBackup = ref(false);

// Disable modal
const showDisableModal = ref(false);
const disablePassword = ref('');
const disableLoading = ref(false);
const disableError = ref('');

// Build a QR image URL via a public API (avoids needing a backend QR generator)
const qrImageUrl = computed(() => {
  if (!qrCodeURL.value) return '';
  const encoded = encodeURIComponent(qrCodeURL.value);
  return `https://api.qrserver.com/v1/create-qr-code/?data=${encoded}&size=180x180&margin=4`;
});

onMounted(async () => {
  try {
    const res = await authApi.getTOTPStatus();
    totpEnabled.value = res.data?.enabled ?? false;
  } catch {
    // ignore
  } finally {
    loadingStatus.value = false;
  }
});

async function startSetup() {
  setupError.value = '';
  try {
    const res = await authApi.setupTOTP();
    if (res.data) {
      totpSecret.value = res.data.secret;
      qrCodeURL.value = res.data.qr_code_url;
      setupStep.value = 'scan';
    }
  } catch (err: any) {
    const data = err?.response?.data;
    setupError.value = data?.message || data?.description || 'Failed to start 2FA setup.';
  }
}

async function confirmSetup() {
  if (!confirmCode.value.trim()) return;
  confirmLoading.value = true;
  setupError.value = '';
  try {
    const res = await authApi.confirmTOTP(confirmCode.value.trim());
    if (res.data?.backup_codes) {
      backupCodes.value = res.data.backup_codes;
      totpEnabled.value = true;
      setupStep.value = 'backup';
    }
  } catch (err: any) {
    const data = err?.response?.data;
    setupError.value = data?.message || data?.description || 'Invalid code. Please try again.';
  } finally {
    confirmLoading.value = false;
  }
}

function finishSetup() {
  setupStep.value = null;
  confirmCode.value = '';
  totpSecret.value = '';
  qrCodeURL.value = '';
  backupCodes.value = [];
}

async function copySecret() {
  await navigator.clipboard.writeText(totpSecret.value);
  copied.value = true;
  setTimeout(() => { copied.value = false; }, 2000);
}

async function copyBackupCodes() {
  await navigator.clipboard.writeText(backupCodes.value.join('\n'));
  copiedBackup.value = true;
  setTimeout(() => { copiedBackup.value = false; }, 2000);
}

function closeDisableModal() {
  showDisableModal.value = false;
  disablePassword.value = '';
  disableError.value = '';
}

async function disableTOTP() {
  if (!disablePassword.value) {
    disableError.value = 'Password is required.';
    return;
  }
  disableLoading.value = true;
  disableError.value = '';
  try {
    await authApi.disableTOTP(disablePassword.value);
    totpEnabled.value = false;
    closeDisableModal();
  } catch (err: any) {
    const data = err?.response?.data;
    disableError.value = data?.message || data?.description || 'Failed to disable 2FA.';
  } finally {
    disableLoading.value = false;
  }
}
</script>

<style scoped lang="scss">
.security-page {
  max-width: 780px;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* ── Page header ─────────────────────────────────────────────────────────── */
.page-header {
  margin-bottom: 4px;
}

.page-title {
  margin: 0 0 4px;
  font-size: 1.375rem;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
  margin: 0;
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Loading ─────────────────────────────────────────────────────────────── */
.loading-center {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 0;
  gap: 16px;
}

.loading-label {
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
    border-color: rgba(255, 255, 255, 0.35);
    border-top-color: #fff;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Cards ───────────────────────────────────────────────────────────────── */
.an-card {
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 16px;
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

  .material-symbols-outlined {
    font-size: 18px;
  }

  &--primary {
    background: color-mix(in srgb, var(--md-sys-color-primary) 12%, transparent);
    color: var(--md-sys-color-primary);
  }

  &--neutral {
    background: var(--md-sys-color-surface-container-high);
    color: var(--md-sys-color-on-surface-variant);
  }
}

.an-card-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.an-card-body {
  padding: 24px;
}

/* ── 2FA status row ──────────────────────────────────────────────────────── */
.twofa-status-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
  flex-wrap: wrap;
}

.twofa-status-info {
  flex: 1;
  min-width: 0;
}

.twofa-status-badge-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.twofa-status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
  animation: pulse 2s infinite;

  &--on {
    background: #1aa563;
    box-shadow: 0 0 0 3px color-mix(in srgb, #1aa563 20%, transparent);
  }

  &--off {
    background: var(--md-sys-color-on-surface-variant);
    opacity: 0.5;
    animation: none;
  }
}

@keyframes pulse {
  0%, 100% { box-shadow: 0 0 0 3px color-mix(in srgb, #1aa563 20%, transparent); }
  50% { box-shadow: 0 0 0 6px color-mix(in srgb, #1aa563 8%, transparent); }
}

.twofa-description {
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
  font-size: 0.9rem;
  max-width: 500px;
}

.twofa-status-action {
  flex-shrink: 0;
}

/* ── Badges ──────────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  font-size: 0.72rem;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: 6px;

  &--success {
    background: color-mix(in srgb, #1aa563 14%, transparent);
    color: #0a7040;
  }

  &--neutral {
    background: var(--md-sys-color-surface-container-high);
    color: var(--md-sys-color-on-surface-variant);
  }
}

/* ── Wizard steps ────────────────────────────────────────────────────────── */
.wizard-steps {
  display: flex;
  align-items: center;
  padding: 20px 24px 16px;
  overflow-x: auto;
}

.wizard-step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.wizard-step__dot {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  font-weight: 700;
  background: var(--md-sys-color-surface-container);
  color: var(--md-sys-color-on-surface-variant);
  border: 2px solid var(--md-sys-color-outline-variant);
  transition: all 0.25s;

  .wizard-step--active & {
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 4px color-mix(in srgb, var(--md-sys-color-primary) 15%, transparent);
  }

  .wizard-step--done & {
    background: color-mix(in srgb, #1aa563 15%, transparent);
    color: #0a7040;
    border-color: #1aa563;
  }
}

.wizard-check-icon {
  font-size: 16px;
}

.wizard-step__label {
  font-size: 0.72rem;
  color: var(--md-sys-color-on-surface-variant);
  white-space: nowrap;

  .wizard-step--active & {
    color: var(--md-sys-color-primary);
    font-weight: 600;
  }

  .wizard-step--done & {
    color: #0a7040;
  }
}

.wizard-step__connector {
  height: 2px;
  flex: 1;
  min-width: 32px;
  background: var(--md-sys-color-outline-variant);
  margin: 0 4px;
  margin-bottom: 22px;
  transition: background 0.25s;

  &--done {
    background: #1aa563;
  }
}

/* ── Wizard panel ────────────────────────────────────────────────────────── */
.wizard-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.wizard-panel__title {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.wizard-panel__subtitle {
  margin: 0;
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── QR code ─────────────────────────────────────────────────────────────── */
.qr-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.qr-frame {
  padding: 12px;
  border: 2px solid var(--md-sys-color-outline-variant);
  border-radius: 16px;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.qr-image {
  display: block;
  border-radius: 4px;
}

.qr-helper {
  margin: 0;
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Secret key block ────────────────────────────────────────────────────── */
.secret-key-block {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.secret-key-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface-variant);
}

.secret-key-icon {
  font-size: 16px;
}

.secret-key-row {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  flex-wrap: wrap;
}

/* ── Form fields ─────────────────────────────────────────────────────────── */
.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;

  &--narrow {
    max-width: 220px;
  }

  &--mono .form-input {
    font-family: 'Courier New', monospace;
    letter-spacing: 0.04em;
  }
}

.form-field__label {
  font-size: 0.8rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface-variant);
}

.form-input {
  height: 40px;
  padding: 0 12px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface);
  font-size: 0.9rem;
  outline: none;
  transition: border-color 0.15s;
  width: 100%;
  box-sizing: border-box;

  &:focus {
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary) 12%, transparent);
  }

  &[readonly] {
    background: var(--md-sys-color-surface-container-low);
    color: var(--md-sys-color-on-surface-variant);
    cursor: default;
  }
}

/* ── Code input ──────────────────────────────────────────────────────────── */
.code-input-wrap {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.code-hint {
  margin: 0;
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Wizard actions ──────────────────────────────────────────────────────── */
.wizard-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

/* ── Success banner ──────────────────────────────────────────────────────── */
.success-banner {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px 20px;
  border-radius: 12px;
  background: color-mix(in srgb, #1aa563 10%, transparent);
  border: 1px solid color-mix(in srgb, #1aa563 30%, transparent);
}

.success-banner__icon {
  font-size: 28px;
  color: #1aa563;
  flex-shrink: 0;
}

.success-banner__title {
  font-size: 0.9375rem;
  font-weight: 600;
  margin-bottom: 2px;
  color: var(--md-sys-color-on-surface);
}

.success-banner__subtitle {
  font-size: 0.85rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Backup warning ──────────────────────────────────────────────────────── */
.backup-warning {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 14px 18px;
  border-radius: 10px;
  background: color-mix(in srgb, #f59e0b 12%, transparent);
  border: 1px solid color-mix(in srgb, #f59e0b 35%, transparent);
}

.backup-warning__icon {
  font-size: 20px;
  color: #d97706;
  flex-shrink: 0;
  margin-top: 1px;
}

.backup-warning__text {
  font-size: 0.9rem;
  color: #78350f;
}

/* ── Backup codes grid ───────────────────────────────────────────────────── */
.backup-codes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 8px;
}

.backup-code {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--md-sys-color-surface-container-low);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 10px;
  padding: 10px 14px;
  text-align: center;
  font-family: 'Courier New', monospace;
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  letter-spacing: 0.06em;
  transition: background 0.15s;

  &:hover {
    background: var(--md-sys-color-surface-container);
  }
}

/* ── Disable modal ───────────────────────────────────────────────────────── */
.disable-modal-body {
  min-width: 320px;
  max-width: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.disable-modal-desc {
  margin: 0;
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
}

.disable-modal-error {
  margin: 0;
}

/* ── Feedback ────────────────────────────────────────────────────────────── */
.feedback-error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 10px;
  background: var(--md-sys-color-error-container);
  border: 1px solid color-mix(in srgb, var(--md-sys-color-error) 30%, transparent);
  color: var(--md-sys-color-on-error-container, var(--md-sys-color-error));
  font-size: 0.875rem;
}

.feedback-icon {
  font-size: 18px;
  flex-shrink: 0;
}

/* ── Buttons ─────────────────────────────────────────────────────────────── */
.btn-danger {
  --btn-color: var(--md-sys-color-error, #dc2626);
  background: var(--btn-color) !important;
  border-color: var(--btn-color) !important;

  &.btn-outlined {
    background: transparent !important;
    color: var(--btn-color) !important;
    border-color: var(--btn-color) !important;
  }
}
</style>
