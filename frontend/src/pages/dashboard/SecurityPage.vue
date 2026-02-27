<template>
  <div class="container-fluid py-4">

    <!-- Header -->
    <div class="mb-4">
      <h1 class="h4 fw-bold mb-0">Security</h1>
      <p class="text-muted small mb-0">Manage two-factor authentication and account security.</p>
    </div>

    <!-- Loading -->
    <div v-if="loadingStatus" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading…</span>
      </div>
    </div>

    <template v-else>

      <!-- 2FA Card -->
      <div class="card mb-4">
        <div class="card-body">
          <div class="d-flex align-items-start justify-content-between gap-3">
            <div>
              <h5 class="fw-semibold mb-1">Two-Factor Authentication (2FA)</h5>
              <p class="text-muted small mb-0">
                Add an extra layer of security. When enabled, you'll need a 6-digit code from
                an authenticator app (e.g. Google Authenticator, Authy) at each login.
              </p>
              <span v-if="totpEnabled" class="badge text-bg-success mt-2">Enabled</span>
              <span v-else class="badge text-bg-secondary mt-2">Disabled</span>
            </div>
            <div class="flex-shrink-0">
              <button v-if="!totpEnabled" class="btn btn-primary btn-sm" @click="startSetup">
                Enable 2FA
              </button>
              <button v-else class="btn btn-outline-danger btn-sm" @click="showDisableModal = true">
                Disable 2FA
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Setup Wizard (shown after clicking Enable) -->
      <div v-if="setupStep && !totpEnabled" class="card">
        <div class="card-body">

          <!-- Step 1: Scan QR -->
          <div v-if="setupStep === 'scan'">
            <h6 class="fw-semibold mb-3">Step 1 — Scan this QR code</h6>
            <p class="text-muted small">
              Open your authenticator app and scan the QR code below. If you can't scan,
              manually enter the secret key.
            </p>

            <!-- QR image rendered via qrserver.com public API (client-side only, no server dependency) -->
            <div class="text-center mb-3">
              <img
                :src="qrImageUrl"
                alt="TOTP QR Code"
                width="180"
                height="180"
                class="border rounded"
              />
            </div>

            <div class="mb-3">
              <label class="form-label small fw-medium">Manual secret key</label>
              <div class="input-group input-group-sm">
                <input type="text" class="form-control font-monospace" :value="totpSecret" readonly />
                <button class="btn btn-outline-secondary" type="button" @click="copySecret">
                  {{ copied ? '✓' : 'Copy' }}
                </button>
              </div>
            </div>

            <button class="btn btn-primary" @click="setupStep = 'confirm'">
              I've scanned it →
            </button>
          </div>

          <!-- Step 2: Confirm code -->
          <div v-else-if="setupStep === 'confirm'">
            <h6 class="fw-semibold mb-3">Step 2 — Enter the code to confirm</h6>
            <p class="text-muted small">
              Enter the 6-digit code shown in your authenticator app to verify that setup is correct.
            </p>

            <div v-if="setupError" class="alert alert-danger py-2 small">{{ setupError }}</div>

            <div class="mb-3" style="max-width: 200px;">
              <input
                v-model="confirmCode"
                type="text"
                class="form-control form-control-lg text-center font-monospace"
                placeholder="000000"
                maxlength="6"
                autocomplete="one-time-code"
              />
            </div>

            <div class="d-flex gap-2">
              <button class="btn btn-outline-secondary btn-sm" @click="setupStep = 'scan'">← Back</button>
              <button class="btn btn-primary" :disabled="confirmLoading" @click="confirmSetup">
                <span v-if="confirmLoading" class="spinner-border spinner-border-sm me-2" role="status"></span>
                Confirm &amp; Enable
              </button>
            </div>
          </div>

          <!-- Step 3: Backup codes -->
          <div v-else-if="setupStep === 'backup'">
            <h6 class="fw-semibold mb-2">2FA Enabled!</h6>
            <div class="alert alert-warning py-2 small">
              <strong>Save these backup codes now.</strong> They are shown only once and can be used
              to access your account if you lose your authenticator app.
            </div>
            <div class="row g-2 mb-3">
              <div v-for="code in backupCodes" :key="code" class="col-6">
                <code class="d-block bg-light border rounded px-2 py-1 text-center small">{{ code }}</code>
              </div>
            </div>
            <button class="btn btn-outline-secondary btn-sm me-2" @click="copyBackupCodes">
              {{ copiedBackup ? '✓ Copied' : 'Copy all' }}
            </button>
            <button class="btn btn-primary btn-sm" @click="finishSetup">Done</button>
          </div>

        </div>
      </div>

    </template>

    <!-- Disable Modal -->
    <div v-if="showDisableModal" class="modal d-block" tabindex="-1" style="background: rgba(0,0,0,0.5);">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content border-0 shadow">
          <div class="modal-header border-0 pb-0">
            <h5 class="modal-title fw-semibold">Disable Two-Factor Authentication</h5>
            <button type="button" class="btn-close" @click="closeDisableModal"></button>
          </div>
          <div class="modal-body">
            <p class="text-muted small mb-3">
              Enter your account password to confirm you want to disable 2FA.
            </p>
            <div v-if="disableError" class="alert alert-danger py-2 small">{{ disableError }}</div>
            <div class="mb-3">
              <label class="form-label">Password</label>
              <input
                v-model="disablePassword"
                type="password"
                class="form-control"
                placeholder="••••••••"
                autocomplete="current-password"
              />
            </div>
          </div>
          <div class="modal-footer border-0 pt-0">
            <button type="button" class="btn btn-outline-secondary" @click="closeDisableModal">Cancel</button>
            <button
              type="button"
              class="btn btn-danger"
              :disabled="disableLoading"
              @click="disableTOTP"
            >
              <span v-if="disableLoading" class="spinner-border spinner-border-sm me-2" role="status"></span>
              Disable 2FA
            </button>
          </div>
        </div>
      </div>
    </div>

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
