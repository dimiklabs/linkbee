<template>
  <div class="security-page">

    <!-- Page header -->
    <div class="page-header">
      <h1 class="md-headline-small">Security</h1>
      <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin:4px 0 0;">
        Manage two-factor authentication and account security.
      </p>
    </div>

    <!-- Loading -->
    <div v-if="loadingStatus" style="text-align:center;padding:64px 0;">
      <md-circular-progress indeterminate />
    </div>

    <template v-else>

      <!-- ── 2FA Card ───────────────────────────────────────────────────── -->
      <div class="m3-card m3-card--outlined section-card">
        <div class="card-section-header">
          <span class="md-title-medium">Two-Factor Authentication (2FA)</span>
        </div>
        <div class="card-section-body">
          <div style="display:flex;align-items:flex-start;justify-content:space-between;gap:16px;flex-wrap:wrap;">
            <div style="flex:1;min-width:0;">
              <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin:0 0 12px;">
                Add an extra layer of security. When enabled, you'll need a 6-digit code from
                an authenticator app (e.g. Google Authenticator, Authy) at each login.
              </p>
              <span v-if="totpEnabled" class="m3-badge m3-badge--success">Enabled</span>
              <span v-else class="m3-badge m3-badge--neutral">Disabled</span>
            </div>
            <div style="flex-shrink:0;">
              <md-filled-button v-if="!totpEnabled" @click="startSetup">Enable 2FA</md-filled-button>
              <md-outlined-button
                v-else
                @click="showDisableModal = true"
                style="--md-outlined-button-outline-color:var(--md-sys-color-error);--md-outlined-button-label-text-color:var(--md-sys-color-error);"
              >
                Disable 2FA
              </md-outlined-button>
            </div>
          </div>
        </div>
      </div>

      <!-- ── Setup Wizard ──────────────────────────────────────────────── -->
      <div v-if="setupStep && !totpEnabled" class="m3-card m3-card--outlined section-card">
        <div class="card-section-body">

          <!-- Step 1: Scan QR -->
          <div v-if="setupStep === 'scan'">
            <div class="md-title-medium" style="margin-bottom:12px;">Step 1 — Scan this QR code</div>
            <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin-bottom:20px;">
              Open your authenticator app and scan the QR code below. If you can't scan,
              manually enter the secret key.
            </p>

            <div style="text-align:center;margin-bottom:20px;">
              <img
                :src="qrImageUrl"
                alt="TOTP QR Code"
                width="180"
                height="180"
                style="border:1px solid var(--md-sys-color-outline-variant);border-radius:8px;"
              />
            </div>

            <div style="margin-bottom:20px;">
              <div class="md-label-large" style="margin-bottom:8px;">Manual secret key</div>
              <div style="display:flex;align-items:center;gap:8px;">
                <md-outlined-text-field
                  :value="totpSecret"
                  label="Secret key"
                  readonly
                  style="flex:1;font-family:monospace;"
                />
                <md-outlined-button @click="copySecret">
                  <span class="material-symbols-outlined" style="font-size:18px;vertical-align:middle;margin-right:4px;">
                    {{ copied ? 'check' : 'content_copy' }}
                  </span>
                  {{ copied ? 'Copied' : 'Copy' }}
                </md-outlined-button>
              </div>
            </div>

            <md-filled-button @click="setupStep = 'confirm'">
              I've scanned it
              <span class="material-symbols-outlined" slot="icon">arrow_forward</span>
            </md-filled-button>
          </div>

          <!-- Step 2: Confirm code -->
          <div v-else-if="setupStep === 'confirm'">
            <div class="md-title-medium" style="margin-bottom:12px;">Step 2 — Enter the code to confirm</div>
            <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin-bottom:20px;">
              Enter the 6-digit code shown in your authenticator app to verify that setup is correct.
            </p>

            <div v-if="setupError" class="feedback-error" style="margin-bottom:16px;">
              <span class="material-symbols-outlined" style="font-size:18px;">error</span>
              {{ setupError }}
            </div>

            <md-outlined-text-field
              :value="confirmCode"
              @input="confirmCode = ($event.target as HTMLInputElement).value"
              label="6-digit code"
              maxlength="6"
              autocomplete="one-time-code"
              style="max-width:200px;margin-bottom:20px;font-family:monospace;font-size:1.25rem;"
            />

            <div style="display:flex;gap:12px;flex-wrap:wrap;">
              <md-outlined-button @click="setupStep = 'scan'">
                <span class="material-symbols-outlined" slot="icon">arrow_back</span>
                Back
              </md-outlined-button>
              <md-filled-button :disabled="confirmLoading" @click="confirmSetup">
                <md-circular-progress v-if="confirmLoading" indeterminate style="--md-circular-progress-size:20px;margin-right:8px;" />
                Confirm &amp; Enable
              </md-filled-button>
            </div>
          </div>

          <!-- Step 3: Backup codes -->
          <div v-else-if="setupStep === 'backup'">
            <div class="md-title-medium" style="margin-bottom:12px;">2FA Enabled!</div>
            <div class="feedback-warning" style="margin-bottom:20px;">
              <span class="material-symbols-outlined" style="font-size:18px;">warning</span>
              <span><strong>Save these backup codes now.</strong> They are shown only once and can be used
              to access your account if you lose your authenticator app.</span>
            </div>
            <div class="backup-codes-grid">
              <code
                v-for="code in backupCodes"
                :key="code"
                class="backup-code"
              >{{ code }}</code>
            </div>
            <div style="display:flex;gap:12px;flex-wrap:wrap;margin-top:20px;">
              <md-outlined-button @click="copyBackupCodes">
                <span class="material-symbols-outlined" style="font-size:18px;vertical-align:middle;margin-right:4px;">
                  {{ copiedBackup ? 'check' : 'content_copy' }}
                </span>
                {{ copiedBackup ? 'Copied' : 'Copy all' }}
              </md-outlined-button>
              <md-filled-button @click="finishSetup">Done</md-filled-button>
            </div>
          </div>

        </div>
      </div>

    </template>

    <!-- ── Disable 2FA Dialog ────────────────────────────────────────────── -->
    <md-dialog :open="showDisableModal" @closed="closeDisableModal">
      <div slot="headline">Disable Two-Factor Authentication</div>
      <div slot="content" style="min-width:360px;max-width:100%;">
        <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin-bottom:20px;">
          Enter your account password to confirm you want to disable 2FA.
        </p>
        <div v-if="disableError" class="feedback-error" style="margin-bottom:16px;">
          <span class="material-symbols-outlined" style="font-size:18px;">error</span>
          {{ disableError }}
        </div>
        <md-outlined-text-field
          :value="disablePassword"
          @input="disablePassword = ($event.target as HTMLInputElement).value"
          label="Password"
          type="password"
          autocomplete="current-password"
          style="width:100%;"
        />
      </div>
      <div slot="actions">
        <md-text-button @click="closeDisableModal">Cancel</md-text-button>
        <md-filled-button
          :disabled="disableLoading"
          @click="disableTOTP"
          style="--md-filled-button-container-color:var(--md-sys-color-error);"
        >
          <md-circular-progress v-if="disableLoading" indeterminate style="--md-circular-progress-size:20px;margin-right:8px;" />
          Disable 2FA
        </md-filled-button>
      </div>
    </md-dialog>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
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

<style scoped>
.security-page {
  max-width: 780px;
  padding: 24px 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  margin-bottom: 8px;
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

.feedback-error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 8px;
  background: var(--md-sys-color-error-container);
  color: var(--md-sys-color-error);
  font-size: 0.875rem;
}

.feedback-warning {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 12px 16px;
  border-radius: 8px;
  background: color-mix(in srgb, #f59e0b 12%, transparent);
  color: #92400e;
  font-size: 0.875rem;
}

.backup-codes-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.backup-code {
  display: block;
  background: var(--md-sys-color-surface-container-low);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 6px;
  padding: 8px 12px;
  text-align: center;
  font-family: monospace;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);
}
</style>
