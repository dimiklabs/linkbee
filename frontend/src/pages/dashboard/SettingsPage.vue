<template>
  <div class="settings-page">

    <!-- Page header -->
    <div class="page-header">
      <h1 class="page-title">Profile &amp; Settings</h1>
      <p class="page-subtitle">Manage your account information and security preferences.</p>
    </div>

    <!-- ── Profile ──────────────────────────────────────────────────── -->
    <div class="an-card">
      <div class="an-card-header">
        <div class="an-card-header__left">
          <div class="an-card-icon an-card-icon--primary">
            <span class="material-symbols-outlined">person</span>
          </div>
          <div class="an-card-title">Profile Information</div>
        </div>
      </div>
      <div class="an-card-body">

        <!-- Avatar row -->
        <div class="avatar-row">
          <div class="avatar-wrapper">
            <div
              class="avatar-circle"
              :style="{
                backgroundImage: avatarPreviewUrl ? `url(${avatarPreviewUrl})` : undefined,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
                backgroundColor: avatarPreviewUrl ? 'transparent' : 'var(--md-sys-color-primary)',
              }"
            >
              <span v-if="!avatarPreviewUrl" class="avatar-initials">
                {{ authStore.userInitials }}
              </span>
            </div>
            <div class="avatar-status-ring"></div>
          </div>
          <div class="avatar-controls">
            <div class="avatar-controls__label">Profile photo</div>
            <div class="avatar-controls__row">
              <label class="choose-photo-label">
                <span class="material-symbols-outlined choose-photo-label__icon">upload</span>
                Choose photo
                <input
                  ref="avatarFileInputRef"
                  type="file"
                  accept="image/*"
                  class="visually-hidden"
                  @change="onAvatarFileChange"
                />
              </label>
              <span class="avatar-filename">{{ avatarFileName || 'No file chosen' }}</span>
            </div>
            <div class="avatar-url-toggle">
              <button class="btn-text" @click="showUrlInput = !showUrlInput">
                {{ showUrlInput ? 'Hide URL input' : 'Or enter image URL manually' }}
              </button>
              <div v-if="showUrlInput" class="avatar-url-field">
                <label class="form-field__label">Image URL</label>
                <input
                  type="url"
                  class="form-input"
                  :value="profileForm.profile_picture"
                  @input="profileForm.profile_picture = ($event.target as HTMLInputElement).value"
                  placeholder="https://example.com/photo.jpg"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Name row -->
        <div class="form-row-2">
          <div class="form-field">
            <label class="form-field__label">First name</label>
            <input
              type="text"
              class="form-input"
              :value="profileForm.first_name"
              @input="profileForm.first_name = ($event.target as HTMLInputElement).value"
              maxlength="50"
              placeholder="First name"
            />
          </div>
          <div class="form-field">
            <label class="form-field__label">Last name</label>
            <input
              type="text"
              class="form-input"
              :value="profileForm.last_name"
              @input="profileForm.last_name = ($event.target as HTMLInputElement).value"
              maxlength="50"
              placeholder="Last name"
            />
          </div>
        </div>

        <!-- Email (readonly) -->
        <div class="form-field">
          <label class="form-field__label">Email</label>
          <input
            type="email"
            class="form-input form-input--readonly"
            :value="authStore.profile?.email ?? ''"
            readonly
          />
          <span class="form-field__hint">Email cannot be changed here.</span>
        </div>

        <!-- Phone -->
        <div class="form-field">
          <label class="form-field__label">Phone</label>
          <input
            type="tel"
            class="form-input"
            :value="profileForm.phone"
            @input="profileForm.phone = ($event.target as HTMLInputElement).value"
            maxlength="20"
            placeholder="+1 555 000 0000"
          />
        </div>

        <!-- Feedback -->
        <div v-if="profileSuccess" class="feedback-success">
          <span class="material-symbols-outlined feedback__icon">check_circle</span>
          {{ profileSuccess }}
        </div>
        <div v-if="profileError" class="feedback-error">
          <span class="material-symbols-outlined feedback__icon">error</span>
          {{ profileError }}
        </div>

        <button class="btn-filled" :disabled="savingProfile" @click="saveProfile">
          <div v-if="savingProfile" class="css-spinner css-spinner--sm css-spinner--white"></div>
          Save profile
        </button>
      </div>
    </div>

    <!-- ── Change Password ─────────────────────────────────────────── -->
    <div class="an-card">
      <div class="an-card-header">
        <div class="an-card-header__left">
          <div class="an-card-icon an-card-icon--neutral">
            <span class="material-symbols-outlined">lock</span>
          </div>
          <div class="an-card-title">Change Password</div>
        </div>
      </div>
      <div class="an-card-body">
        <div class="form-field">
          <label class="form-field__label">Current password</label>
          <input
            type="password"
            class="form-input"
            :value="passwordForm.old_password"
            @input="passwordForm.old_password = ($event.target as HTMLInputElement).value"
            autocomplete="current-password"
            placeholder="••••••••"
          />
        </div>
        <div class="form-field">
          <label class="form-field__label">New password <span class="form-field__hint-inline">(at least 8 characters)</span></label>
          <input
            type="password"
            class="form-input"
            :value="passwordForm.new_password"
            @input="passwordForm.new_password = ($event.target as HTMLInputElement).value"
            autocomplete="new-password"
            placeholder="••••••••"
          />
        </div>
        <div class="form-field">
          <label class="form-field__label">Confirm new password</label>
          <input
            type="password"
            class="form-input"
            :class="{ 'form-input--error': passwordMismatch }"
            :value="passwordForm.confirm_password"
            @input="passwordForm.confirm_password = ($event.target as HTMLInputElement).value"
            autocomplete="new-password"
            placeholder="••••••••"
          />
          <span v-if="passwordMismatch" class="form-field__error">Passwords do not match.</span>
        </div>

        <div v-if="passwordSuccess" class="feedback-success">
          <span class="material-symbols-outlined feedback__icon">check_circle</span>
          {{ passwordSuccess }}
        </div>
        <div v-if="passwordError" class="feedback-error">
          <span class="material-symbols-outlined feedback__icon">error</span>
          {{ passwordError }}
        </div>

        <button class="btn-filled" :disabled="savingPassword || !canChangePassword" @click="changePassword">
          <div v-if="savingPassword" class="css-spinner css-spinner--sm css-spinner--white"></div>
          Update password
        </button>
      </div>
    </div>

    <!-- ── Notification Preferences ──────────────────────────────── -->
    <div class="an-card">
      <div class="an-card-header">
        <div class="an-card-header__left">
          <div class="an-card-icon an-card-icon--primary">
            <span class="material-symbols-outlined">notifications</span>
          </div>
          <div class="an-card-title">Notification Preferences</div>
        </div>
      </div>

      <div class="switch-rows">
        <div class="switch-row">
          <div class="switch-row__content">
            <div class="switch-row__title">Link expiry warnings</div>
            <div class="switch-row__desc">Get emailed 7 days before a link expires.</div>
          </div>
          <label class="m3-switch">
            <input type="checkbox" class="m3-switch__input"
              :checked="notifPrefs.expiry_warnings"
              @change="notifPrefs.expiry_warnings = ($event.target as HTMLInputElement).checked"
            />
            <span class="m3-switch__track"></span>
          </label>
        </div>

        <div class="switch-row">
          <div class="switch-row__content">
            <div class="switch-row__title">Health alert emails</div>
            <div class="switch-row__desc">Notify when a link becomes unhealthy.</div>
          </div>
          <label class="m3-switch">
            <input type="checkbox" class="m3-switch__input"
              :checked="notifPrefs.health_alerts"
              @change="notifPrefs.health_alerts = ($event.target as HTMLInputElement).checked"
            />
            <span class="m3-switch__track"></span>
          </label>
        </div>

        <div class="switch-row">
          <div class="switch-row__content">
            <div class="switch-row__title">Weekly analytics digest</div>
            <div class="switch-row__desc">Summary of clicks sent every Monday.</div>
          </div>
          <label class="m3-switch">
            <input type="checkbox" class="m3-switch__input"
              :checked="notifPrefs.weekly_digest"
              @change="notifPrefs.weekly_digest = ($event.target as HTMLInputElement).checked"
            />
            <span class="m3-switch__track"></span>
          </label>
        </div>

        <div class="switch-row switch-row--last">
          <div class="switch-row__content">
            <div class="switch-row__title">Click milestone alerts</div>
            <div class="switch-row__desc">Notified at 100, 1K, 10K clicks per link.</div>
          </div>
          <label class="m3-switch">
            <input type="checkbox" class="m3-switch__input"
              :checked="notifPrefs.milestone_alerts"
              @change="notifPrefs.milestone_alerts = ($event.target as HTMLInputElement).checked"
            />
            <span class="m3-switch__track"></span>
          </label>
        </div>
      </div>

      <div class="an-card-body">
        <p class="notif-note">
          Email preferences affect automated system emails. Transactional emails (password reset, verification) are always sent.
        </p>
        <div v-if="notifSaved" class="feedback-success">
          <span class="material-symbols-outlined feedback__icon">check_circle</span>
          Preferences saved.
        </div>
        <button class="btn-filled" @click="saveNotifPrefs">Save preferences</button>
      </div>
    </div>

    <!-- ── Active Sessions ────────────────────────────────────────── -->
    <div class="an-card">
      <div class="an-card-header">
        <div class="an-card-header__left">
          <div class="an-card-icon an-card-icon--neutral">
            <span class="material-symbols-outlined">devices</span>
          </div>
          <div class="an-card-title">Active Sessions</div>
        </div>
        <button class="btn-outlined btn-danger" :disabled="revokingAll" @click="logoutAllOther">
          <div v-if="revokingAll" class="css-spinner css-spinner--sm"></div>
          Logout all other devices
        </button>
      </div>

      <div v-if="sessionsLoading" class="sessions-loading">
        <div class="css-spinner"></div>
      </div>
      <div v-else-if="sessions.length === 0" class="sessions-empty">
        No active sessions found.
      </div>
      <div v-else>
        <div v-for="session in sessions" :key="session.id" class="session-row">
          <span class="material-symbols-outlined session-device-icon">
            {{ session.device_type?.toLowerCase() === 'mobile' ? 'smartphone' : session.device_type?.toLowerCase() === 'tablet' ? 'tablet' : 'computer' }}
          </span>
          <div class="session-body">
            <div class="session-label-row">
              <span class="session-label">{{ sessionLabel(session) }}</span>
              <span v-if="session.is_current" class="m3-badge m3-badge--success">Current</span>
              <span v-if="session.login_method" class="m3-badge m3-badge--neutral">{{ session.login_method }}</span>
            </div>
            <div class="session-meta">
              <span v-if="session.ip_address">{{ session.ip_address }}</span>
              <span v-if="session.ip_address && session.last_activity_at"> · </span>
              <span>Last active {{ formatRelative(session.last_activity_at) }}</span>
            </div>
            <div class="session-meta">Signed in {{ formatDate(session.created_at) }}</div>
          </div>
          <button
            class="btn-text"
            v-if="!session.is_current"
            :disabled="revokingId === session.id"
            @click="revokeSession(session)"
          >
            <div v-if="revokingId === session.id" class="css-spinner css-spinner--sm"></div>
            <span v-else>Revoke</span>
          </button>
        </div>
      </div>
    </div>

    <!-- ── Data Export ────────────────────────────────────────────── -->
    <div class="an-card">
      <div class="an-card-header">
        <div class="an-card-header__left">
          <div class="an-card-icon an-card-icon--neutral">
            <span class="material-symbols-outlined">download</span>
          </div>
          <div class="an-card-title">Data Export</div>
        </div>
      </div>
      <div class="an-card-body">
        <div class="export-row">
          <div class="export-row__text">
            <div class="export-row__title">Download your data</div>
            <p class="export-row__desc">
              Export a JSON file containing your profile information and all your short links.
              Your data will be ready to download immediately.
            </p>
          </div>
          <button class="btn-tonal export-btn" :disabled="exportingData" @click="downloadExport">
            <div v-if="exportingData" class="css-spinner css-spinner--sm"></div>
            <span v-else class="material-symbols-outlined">download</span>
            Download my data
          </button>
        </div>
        <div v-if="exportError" class="feedback-error feedback--mt">
          <span class="material-symbols-outlined feedback__icon">error</span>
          {{ exportError }}
        </div>
      </div>
    </div>

    <!-- ── Danger Zone ────────────────────────────────────────────── -->
    <div class="an-card danger-zone-card">
      <div class="an-card-header danger-zone-header">
        <div class="an-card-header__left">
          <div class="an-card-icon an-card-icon--error">
            <span class="material-symbols-outlined">warning</span>
          </div>
          <div class="an-card-title danger-zone-title">Danger Zone</div>
        </div>
      </div>
      <div class="an-card-body">
        <div class="export-row">
          <div class="export-row__text">
            <div class="export-row__title">Delete account</div>
            <p class="export-row__desc">
              Permanently delete your account and all associated links, analytics, and data.
              This action <strong>cannot be undone</strong>.
            </p>
          </div>
          <button class="btn-outlined btn-danger" @click="showDeleteConfirm = true">
            Delete account
          </button>
        </div>

        <!-- Delete confirmation -->
        <div v-if="showDeleteConfirm" class="delete-confirm-box">
          <p class="delete-confirm-box__prompt">
            To confirm, type <code class="inline-code">DELETE</code> in the box below:
          </p>
          <div class="delete-confirm-box__actions">
            <div class="form-field form-field--inline">
              <input
                type="text"
                class="form-input"
                :value="deleteConfirmText"
                @input="deleteConfirmText = ($event.target as HTMLInputElement).value"
                placeholder="DELETE"
              />
            </div>
            <button
              class="btn-filled btn-danger"
              :disabled="deleteConfirmText !== 'DELETE' || deletingAccount"
              @click="deleteAccount"
            >
              <div v-if="deletingAccount" class="css-spinner css-spinner--sm css-spinner--white"></div>
              Confirm delete
            </button>
            <button class="btn-outlined" @click="showDeleteConfirm = false; deleteConfirmText = ''">Cancel</button>
          </div>
          <div v-if="deleteError" class="feedback-error feedback--mt">
            <span class="material-symbols-outlined feedback__icon">error</span>
            {{ deleteError }}
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import authApi from '@/api/auth';
import type { SessionResponse } from '@/types/auth';

const router = useRouter();
const authStore = useAuthStore();

// ── Profile ───────────────────────────────────────────────────────────────────

const profileForm = ref({
  first_name: '',
  last_name: '',
  phone: '',
  profile_picture: '',
});
const savingProfile = ref(false);
const profileSuccess = ref('');
const profileError = ref('');

// Avatar file picker state
const avatarFileInputRef = ref<HTMLInputElement | null>(null);
const avatarFileName = ref('');
const showUrlInput = ref(false);

const avatarPreviewUrl = computed(() => profileForm.value.profile_picture.trim() || '');

function onAvatarFileChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  avatarFileName.value = file.name;

  const reader = new FileReader();
  reader.onload = (e) => {
    const dataUrl = e.target?.result as string;
    const img = new Image();
    img.onload = () => {
      const MAX = 200;
      let { width, height } = img;

      // Preserve aspect ratio, cap at 200x200
      if (width > height) {
        if (width > MAX) { height = Math.round(height * MAX / width); width = MAX; }
      } else {
        if (height > MAX) { width = Math.round(width * MAX / height); height = MAX; }
      }

      const canvas = document.createElement('canvas');
      canvas.width = width;
      canvas.height = height;
      const ctx = canvas.getContext('2d');
      if (!ctx) return;
      ctx.drawImage(img, 0, 0, width, height);

      // Convert to JPEG base64, 0.85 quality
      const compressedDataUri = canvas.toDataURL('image/jpeg', 0.85);
      profileForm.value.profile_picture = compressedDataUri;
    };
    img.src = dataUrl;
  };
  reader.readAsDataURL(file);
}

function loadProfile() {
  const p = authStore.profile;
  if (!p) return;
  profileForm.value = {
    first_name: p.first_name ?? '',
    last_name: p.last_name ?? '',
    phone: p.phone ?? '',
    profile_picture: p.profile_picture ?? '',
  };
}

watch(() => authStore.profile, loadProfile, { immediate: true });

async function saveProfile() {
  savingProfile.value = true;
  profileSuccess.value = '';
  profileError.value = '';
  try {
    const picValue = profileForm.value.profile_picture;
    await authApi.updateProfile({
      first_name: profileForm.value.first_name.trim() || undefined,
      last_name: profileForm.value.last_name.trim() || undefined,
      phone: profileForm.value.phone.trim() || undefined,
      // Send base64 data URIs as-is; only omit if truly empty
      profile_picture: picValue || undefined,
    });
    await authStore.fetchProfile();
    profileSuccess.value = 'Profile updated successfully.';
  } catch (err: any) {
    profileError.value = err?.response?.data?.description ?? 'Failed to update profile.';
  } finally {
    savingProfile.value = false;
  }
}

// ── Change password ───────────────────────────────────────────────────────────

const passwordForm = ref({
  old_password: '',
  new_password: '',
  confirm_password: '',
});
const savingPassword = ref(false);
const passwordSuccess = ref('');
const passwordError = ref('');

const passwordMismatch = computed(
  () => passwordForm.value.confirm_password.length > 0 && passwordForm.value.new_password !== passwordForm.value.confirm_password
);

const canChangePassword = computed(
  () =>
    passwordForm.value.old_password.length > 0 &&
    passwordForm.value.new_password.length >= 8 &&
    !passwordMismatch.value
);

async function changePassword() {
  if (!canChangePassword.value) return;
  savingPassword.value = true;
  passwordSuccess.value = '';
  passwordError.value = '';
  try {
    await authApi.changePassword({
      old_password: passwordForm.value.old_password,
      new_password: passwordForm.value.new_password,
    });
    passwordSuccess.value = 'Password changed successfully.';
    passwordForm.value = { old_password: '', new_password: '', confirm_password: '' };
  } catch (err: any) {
    passwordError.value = err?.response?.data?.description ?? 'Failed to change password.';
  } finally {
    savingPassword.value = false;
  }
}

// ── Sessions ──────────────────────────────────────────────────────────────────

const sessions = ref<SessionResponse[]>([]);
const sessionsLoading = ref(false);
const revokingId = ref<string | null>(null);
const revokingAll = ref(false);

async function loadSessions() {
  sessionsLoading.value = true;
  try {
    const res = await authApi.getSessions();
    sessions.value = res.data?.sessions ?? [];
  } finally {
    sessionsLoading.value = false;
  }
}

async function revokeSession(session: SessionResponse) {
  revokingId.value = session.id;
  try {
    await authApi.deleteSession(session.id);
    sessions.value = sessions.value.filter(s => s.id !== session.id);
  } finally {
    revokingId.value = null;
  }
}

async function logoutAllOther() {
  revokingAll.value = true;
  try {
    await authApi.logoutAll();
    await loadSessions();
  } finally {
    revokingAll.value = false;
  }
}

function deviceIcon(deviceType?: string) {
  switch (deviceType?.toLowerCase()) {
    case 'mobile':  return '📱';
    case 'tablet':  return '📟';
    case 'desktop': return '🖥️';
    default:        return '💻';
  }
}

function sessionLabel(session: SessionResponse) {
  const parts: string[] = [];
  if (session.browser) parts.push(session.browser);
  if (session.os) parts.push(session.os);
  if (parts.length === 0 && session.user_agent) {
    return session.user_agent.length > 60 ? session.user_agent.slice(0, 60) + '…' : session.user_agent;
  }
  return parts.join(' on ') || 'Unknown device';
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString(undefined, {
    year: 'numeric', month: 'short', day: 'numeric',
  });
}

function formatRelative(dateStr: string) {
  const diff = Date.now() - new Date(dateStr).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return 'just now';
  if (mins < 60) return `${mins}m ago`;
  const hrs = Math.floor(mins / 60);
  if (hrs < 24) return `${hrs}h ago`;
  const days = Math.floor(hrs / 24);
  return `${days}d ago`;
}

// ── Data export ───────────────────────────────────────────────────────────────

const exportingData = ref(false);
const exportError = ref('');

async function downloadExport() {
  exportingData.value = true;
  exportError.value = '';
  try {
    await authApi.downloadExport();
  } catch (err: any) {
    exportError.value = err?.response?.data?.description ?? 'Failed to export data. Please try again.';
  } finally {
    exportingData.value = false;
  }
}

// ── Delete account ────────────────────────────────────────────────────────────

const showDeleteConfirm = ref(false);
const deleteConfirmText = ref('');
const deletingAccount = ref(false);
const deleteError = ref('');

async function deleteAccount() {
  if (deleteConfirmText.value !== 'DELETE') return;
  deletingAccount.value = true;
  deleteError.value = '';
  try {
    await authApi.deleteAccount();
    authStore.logout();
    router.push('/login');
  } catch (err: any) {
    deleteError.value = err?.response?.data?.description ?? 'Failed to delete account.';
    deletingAccount.value = false;
  }
}

// ── Notification preferences ──────────────────────────────────────────────────

interface NotificationPrefs {
  expiry_warnings: boolean;
  health_alerts: boolean;
  weekly_digest: boolean;
  milestone_alerts: boolean;
}

const NOTIF_KEY = 'notification_prefs';

const notifPrefs = ref<NotificationPrefs>({
  expiry_warnings: true,
  health_alerts: true,
  weekly_digest: true,
  milestone_alerts: true,
});

const notifSaved = ref(false);

function loadNotifPrefs() {
  try {
    const raw = localStorage.getItem(NOTIF_KEY);
    if (raw) notifPrefs.value = { ...notifPrefs.value, ...JSON.parse(raw) };
  } catch {}
}

function saveNotifPrefs() {
  localStorage.setItem(NOTIF_KEY, JSON.stringify(notifPrefs.value));
  notifSaved.value = true;
  setTimeout(() => (notifSaved.value = false), 2000);
}

// ── Init ──────────────────────────────────────────────────────────────────────

onMounted(async () => {
  if (!authStore.profile) {
    await authStore.fetchProfile();
  }
  loadSessions();
  loadNotifPrefs();
});
</script>

<style scoped lang="scss">
/* ── Root ─────────────────────────────────────────────────────────── */
.settings-page {
  max-width: 780px;
  display: flex;
  flex-direction: column;
  gap: 24px;
  padding: 0;

  @media (max-width: 575px) { gap: 16px; }
}

/* ── Page header ─────────────────────────────────────────────────── */
.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0 0 4px;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
  margin: 0;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Cards ───────────────────────────────────────────────────────── */
.an-card {
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant, #e3e8ee);
  border-radius: 16px;
  overflow: hidden;
}

.an-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  gap: 0.75rem;
  flex-wrap: wrap;
  border-bottom: 1px solid var(--md-sys-color-outline-variant, #e3e8ee);
  background: var(--md-sys-color-surface-container-low, #f8f9fa);

  &__left {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
}

.an-card-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
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
    background: rgba(99, 91, 255, 0.1);
    color: var(--md-sys-color-primary, #635bff);
  }

  &--neutral {
    background: var(--md-sys-color-surface-container, #f3f4f6);
    color: var(--md-sys-color-on-surface-variant);
  }

  &--error {
    background: rgba(220, 38, 38, 0.1);
    color: #dc2626;
  }
}

.an-card-body {
  padding: 20px 24px;
}

/* ── CSS spinner ─────────────────────────────────────────────────── */
.css-spinner {
  width: 24px;
  height: 24px;
  border: 3px solid var(--md-sys-color-outline-variant, #e3e8ee);
  border-top-color: var(--md-sys-color-primary, #635bff);
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

@keyframes spin { to { transform: rotate(360deg); } }

/* ── Avatar ──────────────────────────────────────────────────────── */
.avatar-row {
  display: flex;
  align-items: flex-start;
  gap: 24px;
  margin-bottom: 20px;

  @media (max-width: 500px) {
    flex-direction: column;
    align-items: center;
  }
}

.avatar-wrapper {
  position: relative;
  flex-shrink: 0;
}

.avatar-circle {
  width: 84px;
  height: 84px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 3px solid var(--md-sys-color-surface);
  box-shadow: 0 0 0 2px var(--md-sys-color-outline-variant);
}

.avatar-initials {
  color: var(--md-sys-color-on-primary);
  font-weight: 700;
  font-size: 1.75rem;
  line-height: 1;
  letter-spacing: -0.02em;
}

.avatar-status-ring {
  position: absolute;
  bottom: 4px;
  right: 4px;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: #1aa563;
  border: 2px solid var(--md-sys-color-surface);
}

.avatar-controls {
  flex: 1;
  min-width: 0;

  &__label {
    font-size: 0.8125rem;
    font-weight: 500;
    color: var(--md-sys-color-on-surface);
    margin-bottom: 8px;
  }

  &__row {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
  }
}

.avatar-filename {
  font-size: 0.8125rem;
  color: var(--md-sys-color-on-surface-variant);
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.avatar-url-toggle {
  margin-top: 8px;
}

.avatar-url-field {
  margin-top: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.choose-photo-label {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 7px 14px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
  transition: background 0.15s, border-color 0.15s;

  &:hover {
    background: var(--md-sys-color-surface-container-low);
    border-color: var(--md-sys-color-outline);
  }

  &__icon { font-size: 18px; }
}

/* ── Form fields ─────────────────────────────────────────────────── */
.form-row-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 12px;

  @media (max-width: 540px) { grid-template-columns: 1fr; }
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 12px;

  &:last-of-type { margin-bottom: 0; }

  &--inline { margin-bottom: 0; }

  &__label {
    font-size: 0.75rem;
    font-weight: 500;
    color: var(--md-sys-color-on-surface-variant);
  }

  &__hint {
    font-size: 0.75rem;
    color: var(--md-sys-color-on-surface-variant);
  }

  &__hint-inline {
    font-weight: 400;
    font-size: 0.75rem;
    color: var(--md-sys-color-on-surface-variant);
  }

  &__error {
    font-size: 0.75rem;
    color: var(--md-sys-color-error, #dc2626);
  }
}

.form-input {
  height: 42px;
  padding: 0 12px;
  border: 1px solid var(--md-sys-color-outline-variant, #e3e8ee);
  border-radius: 8px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface);
  font-size: 0.875rem;
  font-family: inherit;
  outline: none;
  transition: border-color 0.15s;

  &::placeholder { color: var(--md-sys-color-on-surface-variant); opacity: 0.6; }
  &:focus { border-color: var(--md-sys-color-primary, #635bff); }

  &--readonly {
    background: var(--md-sys-color-surface-container-low, #f8f9fa);
    color: var(--md-sys-color-on-surface-variant);
    cursor: not-allowed;
  }

  &--error {
    border-color: var(--md-sys-color-error, #dc2626);
  }
}

/* ── Feedback ────────────────────────────────────────────────────── */
.feedback-success {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 10px;
  background: rgba(26, 165, 99, 0.1);
  border: 1px solid rgba(26, 165, 99, 0.25);
  color: #0a6639;
  font-size: 0.875rem;
  margin-bottom: 16px;
}

.feedback-error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 10px;
  background: rgba(220, 38, 38, 0.08);
  border: 1px solid rgba(220, 38, 38, 0.25);
  color: #dc2626;
  font-size: 0.875rem;
  margin-bottom: 16px;
}

.feedback--mt { margin-top: 16px; margin-bottom: 0; }
.feedback__icon { font-size: 18px; flex-shrink: 0; }

/* ── Notification switches ───────────────────────────────────────── */
.switch-rows {
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
}

.switch-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 24px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  transition: background 0.15s;

  &:hover { background: var(--md-sys-color-surface-container-low, #f8f9fa); }

  &--last { border-bottom: none; }

  &__content { flex: 1; min-width: 0; }

  &__title {
    font-size: 0.9rem;
    font-weight: 500;
    color: var(--md-sys-color-on-surface);
    margin-bottom: 2px;
  }

  &__desc {
    font-size: 0.8125rem;
    color: var(--md-sys-color-on-surface-variant);
  }
}

/* ── M3 native toggle switch ─────────────────────────────────────── */
.m3-switch {
  position: relative;
  display: inline-flex;
  align-items: center;
  cursor: pointer;
  flex-shrink: 0;
}

.m3-switch__input {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
}

.m3-switch__track {
  width: 52px;
  height: 32px;
  border-radius: 16px;
  background: var(--md-sys-color-surface-container-highest, #dde1e7);
  border: 2px solid var(--md-sys-color-outline, #888);
  position: relative;
  transition: background 0.2s, border-color 0.2s;

  &::after {
    content: '';
    position: absolute;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: var(--md-sys-color-outline, #888);
    top: 6px;
    left: 6px;
    transition: transform 0.2s, background 0.2s, width 0.15s, height 0.15s, top 0.15s, left 0.15s;
  }
}

.m3-switch__input:checked + .m3-switch__track {
  background: var(--md-sys-color-primary, #635bff);
  border-color: var(--md-sys-color-primary, #635bff);

  &::after {
    background: #fff;
    transform: translateX(20px);
    width: 24px;
    height: 24px;
    top: 2px;
    left: 2px;
  }
}

.m3-switch:hover .m3-switch__track {
  border-color: var(--md-sys-color-on-surface-variant, #666);
}

/* ── Notif note ──────────────────────────────────────────────────── */
.notif-note {
  font-size: 0.8125rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0 0 16px;
}

/* ── Sessions ────────────────────────────────────────────────────── */
.sessions-loading {
  display: flex;
  justify-content: center;
  padding: 32px 24px;
}

.sessions-empty {
  text-align: center;
  padding: 32px 24px;
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
}

.session-row {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 16px 24px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  transition: background 0.15s;

  &:last-child { border-bottom: none; }
  &:hover { background: var(--md-sys-color-surface-container-low); }
}

.session-device-icon {
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 2px;
  flex-shrink: 0;
  font-size: 28px;
}

.session-body {
  flex: 1;
  min-width: 0;
}

.session-label-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 2px;
}

.session-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
}

.session-meta {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  line-height: 1.5;
}

/* ── Export row ──────────────────────────────────────────────────── */
.export-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;

  &__title {
    font-size: 0.9rem;
    font-weight: 500;
    color: var(--md-sys-color-on-surface);
    margin-bottom: 4px;
  }

  &__desc {
    font-size: 0.8125rem;
    color: var(--md-sys-color-on-surface-variant);
    margin: 0;
    max-width: 480px;
  }
}

.export-btn {
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

/* ── Danger zone ─────────────────────────────────────────────────── */
.danger-zone-card {
  border-color: rgba(220, 38, 38, 0.25);
}

.danger-zone-header {
  background: rgba(220, 38, 38, 0.04);
  border-bottom-color: rgba(220, 38, 38, 0.15);
}

.danger-zone-title { color: #dc2626; }

.delete-confirm-box {
  margin-top: 20px;
  padding: 16px 20px;
  border-radius: 12px;
  border: 1px solid rgba(220, 38, 38, 0.3);
  background: rgba(220, 38, 38, 0.04);

  &__prompt {
    font-size: 0.9rem;
    color: var(--md-sys-color-on-surface);
    margin: 0 0 12px;
  }

  &__actions {
    display: flex;
    gap: 12px;
    align-items: center;
    flex-wrap: wrap;
  }
}

.inline-code {
  font-family: 'SFMono-Regular', Consolas, monospace;
  background: var(--md-sys-color-surface-container-low);
  padding: 1px 5px;
  border-radius: 4px;
  font-size: 0.875em;
}

/* ── Badges ──────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  font-size: 0.72rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 6px;
  white-space: nowrap;

  &--success {
    background: rgba(26, 165, 99, 0.12);
    color: #0a7040;
  }

  &--neutral {
    background: var(--md-sys-color-surface-container-high);
    color: var(--md-sys-color-on-surface-variant);
  }
}

/* ── Utility ─────────────────────────────────────────────────────── */
.visually-hidden {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0 0 0 0);
  white-space: nowrap;
}
</style>
