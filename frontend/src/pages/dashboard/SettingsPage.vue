<template>
  <div class="settings-page">
    <!-- Page header -->
    <div class="page-header">
      <h1 class="md-headline-small page-title">Profile &amp; Settings</h1>
      <p class="md-body-medium page-subtitle">
        Manage your account information and security preferences.
      </p>
    </div>

    <!-- ── Profile ──────────────────────────────────────────────────────── -->
    <div class="m3-card m3-card--outlined section-card">
      <div class="card-section-header">
        <span class="material-symbols-outlined card-section-header__icon">person</span>
        <span class="md-title-medium">Profile Information</span>
      </div>
      <div class="card-section-body">
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
            <div class="md-label-large" style="margin-bottom:8px;">Profile photo</div>
            <div style="display:flex;align-items:center;gap:8px;flex-wrap:wrap;">
              <label class="choose-photo-label">
                <span class="material-symbols-outlined" style="font-size:18px;vertical-align:middle;margin-right:4px;">upload</span>
                Choose photo
                <input
                  ref="avatarFileInputRef"
                  type="file"
                  accept="image/*"
                  style="display:none"
                  @change="onAvatarFileChange"
                />
              </label>
              <span class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);max-width:200px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">
                {{ avatarFileName || 'No file chosen' }}
              </span>
            </div>
            <div style="margin-top:8px;">
              <button class="btn-text" @click="showUrlInput = !showUrlInput">
                {{ showUrlInput ? 'Hide URL input' : 'Or enter image URL manually' }}
              </button>
              <div v-if="showUrlInput" style="margin-top:8px;">
                <md-outlined-text-field
                  :value="profileForm.profile_picture"
                  @input="profileForm.profile_picture = ($event.target as HTMLInputElement).value"
                  label="Image URL"
                  type="url"
                  style="width:100%;"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Name row -->
        <div class="form-row-2">
          <md-outlined-text-field
            :value="profileForm.first_name"
            @input="profileForm.first_name = ($event.target as HTMLInputElement).value"
            label="First name"
            maxlength="50"
            style="width:100%;"
          />
          <md-outlined-text-field
            :value="profileForm.last_name"
            @input="profileForm.last_name = ($event.target as HTMLInputElement).value"
            label="Last name"
            maxlength="50"
            style="width:100%;"
          />
        </div>

        <!-- Email (readonly) -->
        <md-outlined-text-field
          :value="authStore.profile?.email ?? ''"
          label="Email"
          type="email"
          readonly
          style="width:100%;margin-bottom:16px;"
        >
          <span slot="supporting-text">Email cannot be changed here.</span>
        </md-outlined-text-field>

        <!-- Phone -->
        <md-outlined-text-field
          :value="profileForm.phone"
          @input="profileForm.phone = ($event.target as HTMLInputElement).value"
          label="Phone"
          type="tel"
          maxlength="20"
          style="width:100%;margin-bottom:20px;"
        />

        <!-- Feedback -->
        <div v-if="profileSuccess" class="feedback-success">
          <span class="material-symbols-outlined" style="font-size:18px;">check_circle</span>
          {{ profileSuccess }}
        </div>
        <div v-if="profileError" class="feedback-error">
          <span class="material-symbols-outlined" style="font-size:18px;">error</span>
          {{ profileError }}
        </div>

        <button class="btn-filled" :disabled="savingProfile" @click="saveProfile">
          <md-circular-progress v-if="savingProfile" indeterminate style="margin-right:8px;" />
          Save profile
        </button>
      </div>
    </div>

    <!-- ── Change Password ──────────────────────────────────────────────── -->
    <div class="m3-card m3-card--outlined section-card">
      <div class="card-section-header">
        <span class="material-symbols-outlined card-section-header__icon">lock</span>
        <span class="md-title-medium">Change Password</span>
      </div>
      <div class="card-section-body">
        <md-outlined-text-field
          :value="passwordForm.old_password"
          @input="passwordForm.old_password = ($event.target as HTMLInputElement).value"
          label="Current password"
          type="password"
          autocomplete="current-password"
          style="width:100%;margin-bottom:16px;"
        />
        <md-outlined-text-field
          :value="passwordForm.new_password"
          @input="passwordForm.new_password = ($event.target as HTMLInputElement).value"
          label="New password (at least 8 characters)"
          type="password"
          autocomplete="new-password"
          style="width:100%;margin-bottom:16px;"
        />
        <md-outlined-text-field
          :value="passwordForm.confirm_password"
          @input="passwordForm.confirm_password = ($event.target as HTMLInputElement).value"
          label="Confirm new password"
          type="password"
          autocomplete="new-password"
          style="width:100%;margin-bottom:4px;"
          :error="passwordMismatch"
        >
          <span v-if="passwordMismatch" slot="error-text">Passwords do not match.</span>
        </md-outlined-text-field>
        <div style="margin-bottom:20px;" />

        <div v-if="passwordSuccess" class="feedback-success">
          <span class="material-symbols-outlined" style="font-size:18px;">check_circle</span>
          {{ passwordSuccess }}
        </div>
        <div v-if="passwordError" class="feedback-error">
          <span class="material-symbols-outlined" style="font-size:18px;">error</span>
          {{ passwordError }}
        </div>

        <button class="btn-filled" :disabled="savingPassword || !canChangePassword" @click="changePassword">
          <md-circular-progress v-if="savingPassword" indeterminate style="margin-right:8px;" />
          Update password
        </button>
      </div>
    </div>

    <!-- ── Notification Preferences ────────────────────────────────────── -->
    <div class="m3-card m3-card--outlined section-card">
      <div class="card-section-header">
        <span class="material-symbols-outlined card-section-header__icon">notifications</span>
        <span class="md-title-medium">Notification Preferences</span>
      </div>
      <div class="card-section-body">
        <!-- Link expiry warnings -->
        <div class="switch-row">
          <div>
            <div class="md-title-small">Link expiry warnings</div>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Get emailed 7 days before a link expires.</div>
          </div>
          <md-switch
            :selected="notifPrefs.expiry_warnings"
            @change="notifPrefs.expiry_warnings = ($event.target as HTMLInputElement).checked"
          />
        </div>

        <!-- Health alert emails -->
        <div class="switch-row">
          <div>
            <div class="md-title-small">Health alert emails</div>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Notify when a link becomes unhealthy.</div>
          </div>
          <md-switch
            :selected="notifPrefs.health_alerts"
            @change="notifPrefs.health_alerts = ($event.target as HTMLInputElement).checked"
          />
        </div>

        <!-- Weekly analytics digest -->
        <div class="switch-row">
          <div>
            <div class="md-title-small">Weekly analytics digest</div>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Summary of clicks sent every Monday.</div>
          </div>
          <md-switch
            :selected="notifPrefs.weekly_digest"
            @change="notifPrefs.weekly_digest = ($event.target as HTMLInputElement).checked"
          />
        </div>

        <!-- Click milestone alerts -->
        <div class="switch-row" style="border-bottom:none;">
          <div>
            <div class="md-title-small">Click milestone alerts</div>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Notified at 100, 1K, 10K clicks per link.</div>
          </div>
          <md-switch
            :selected="notifPrefs.milestone_alerts"
            @change="notifPrefs.milestone_alerts = ($event.target as HTMLInputElement).checked"
          />
        </div>

        <p class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin:16px 0 20px;">
          Email preferences affect automated system emails. Transactional emails (password reset, verification) are always sent.
        </p>

        <div v-if="notifSaved" class="feedback-success" style="margin-bottom:16px;">
          <span class="material-symbols-outlined" style="font-size:18px;">check_circle</span>
          Preferences saved.
        </div>

        <button class="btn-filled" @click="saveNotifPrefs">Save preferences</button>
      </div>
    </div>

    <!-- ── Active Sessions ──────────────────────────────────────────────── -->
    <div class="m3-card m3-card--outlined section-card">
      <div class="card-section-header card-section-header--row">
        <div style="display:flex;align-items:center;gap:10px;">
          <span class="material-symbols-outlined card-section-header__icon">devices</span>
          <span class="md-title-medium">Active Sessions</span>
        </div>
        <button class="btn-outlined btn-danger"
          :disabled="revokingAll"
          @click="logoutAllOther"
        >
          <md-circular-progress v-if="revokingAll" indeterminate style="margin-right:6px;" />
          Logout all other devices
        </button>
      </div>
      <div>
        <div v-if="sessionsLoading" style="text-align:center;padding:32px 24px;">
          <md-circular-progress indeterminate />
        </div>
        <div v-else-if="sessions.length === 0" style="text-align:center;padding:32px 24px;" class="md-body-medium">
          <span style="color:var(--md-sys-color-on-surface-variant);">No active sessions found.</span>
        </div>
        <div v-else>
          <div
            v-for="session in sessions"
            :key="session.id"
            class="session-row"
          >
            <span class="material-symbols-outlined session-device-icon">
              {{ session.device_type?.toLowerCase() === 'mobile' ? 'smartphone' : session.device_type?.toLowerCase() === 'tablet' ? 'tablet' : 'computer' }}
            </span>
            <div style="flex:1;min-width:0;">
              <div style="display:flex;align-items:center;gap:8px;flex-wrap:wrap;">
                <span class="md-label-large">{{ sessionLabel(session) }}</span>
                <span v-if="session.is_current" class="m3-badge m3-badge--success">Current</span>
                <span v-if="session.login_method" class="m3-badge m3-badge--neutral">{{ session.login_method }}</span>
              </div>
              <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin-top:2px;">
                <span v-if="session.ip_address">{{ session.ip_address }}</span>
                <span v-if="session.ip_address && session.last_activity_at"> · </span>
                <span>Last active {{ formatRelative(session.last_activity_at) }}</span>
              </div>
              <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">
                Signed in {{ formatDate(session.created_at) }}
              </div>
            </div>
            <button class="btn-text"
              v-if="!session.is_current"
              :disabled="revokingId === session.id"
              @click="revokeSession(session)"
             
            >
              <md-circular-progress v-if="revokingId === session.id" indeterminate />
              <span v-else>Revoke</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Data Export ──────────────────────────────────────────────────── -->
    <div class="m3-card m3-card--outlined section-card">
      <div class="card-section-header">
        <span class="material-symbols-outlined card-section-header__icon">download</span>
        <span class="md-title-medium">Data Export</span>
      </div>
      <div class="card-section-body">
        <div style="display:flex;align-items:flex-start;justify-content:space-between;gap:16px;flex-wrap:wrap;">
          <div>
            <div class="md-label-large" style="margin-bottom:4px;">Download your data</div>
            <p class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin:0;max-width:480px;">
              Export a JSON file containing your profile information and all your short links.
              Your data will be ready to download immediately.
            </p>
          </div>
          <button class="btn-tonal" :disabled="exportingData" @click="downloadExport">
            <md-circular-progress v-if="exportingData" indeterminate style="margin-right:8px;" />
            <span v-else>
              <span class="material-symbols-outlined" style="font-size:18px;vertical-align:middle;margin-right:4px;">download</span>
            </span>
            Download my data
          </button>
        </div>
        <div v-if="exportError" class="feedback-error" style="margin-top:16px;">
          <span class="material-symbols-outlined" style="font-size:18px;">error</span>
          {{ exportError }}
        </div>
      </div>
    </div>

    <!-- ── Danger Zone ──────────────────────────────────────────────────── -->
    <div class="m3-card m3-card--outlined section-card danger-zone-card">
      <div class="card-section-header danger-zone-header">
        <span class="material-symbols-outlined card-section-header__icon" style="color:var(--md-sys-color-error);">warning</span>
        <span class="md-title-medium" style="color:var(--md-sys-color-error);">Danger Zone</span>
      </div>
      <div class="card-section-body">
        <div style="display:flex;align-items:flex-start;justify-content:space-between;gap:16px;flex-wrap:wrap;">
          <div>
            <div class="md-label-large" style="margin-bottom:4px;">Delete account</div>
            <p class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin:0;max-width:480px;">
              Permanently delete your account and all associated links, analytics, and data.
              This action <strong>cannot be undone</strong>.
            </p>
          </div>
          <button class="btn-outlined btn-danger"
            @click="showDeleteConfirm = true"
          >
            Delete account
          </button>
        </div>

        <!-- Delete confirmation -->
        <div v-if="showDeleteConfirm" class="delete-confirm-box">
          <p class="md-body-medium" style="margin-bottom:12px;">
            To confirm, type <code>DELETE</code> in the box below:
          </p>
          <div style="display:flex;gap:12px;align-items:center;flex-wrap:wrap;">
            <md-outlined-text-field
              :value="deleteConfirmText"
              @input="deleteConfirmText = ($event.target as HTMLInputElement).value"
              label="Type DELETE"
              style="max-width:200px;"
            />
            <button class="btn-filled btn-danger"
              :disabled="deleteConfirmText !== 'DELETE' || deletingAccount"
              @click="deleteAccount"
            >
              <md-circular-progress v-if="deletingAccount" indeterminate style="margin-right:8px;" />
              Confirm delete
            </button>
            <button class="btn-outlined" @click="showDeleteConfirm = false; deleteConfirmText = ''">Cancel</button>
          </div>
          <div v-if="deleteError" class="feedback-error" style="margin-top:12px;">
            <span class="material-symbols-outlined" style="font-size:18px;">error</span>
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
.settings-page {
  max-width: 780px;
  padding: 24px 0;
  display: flex;
  flex-direction: column;
  gap: 24px;

  @media (max-width: 575px) {
    padding: 16px 0;
    gap: 16px;
  }
}

/* ── Page header ─────────────────────────────────────────────────────────── */
.page-header {
  margin-bottom: 4px;
}

.page-title {
  margin: 0 0 4px;
  font-weight: 700;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
  margin: 0;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Section card ────────────────────────────────────────────────────────── */
.section-card {
  border-radius: 16px;
  overflow: hidden;
}

/* ── Card section header ─────────────────────────────────────────────────── */
.card-section-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 24px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface-container-low);

  &--row {
    justify-content: space-between;
  }
}

.card-section-header__icon {
  font-size: 20px;
  color: var(--md-sys-color-primary);
  flex-shrink: 0;
}

.card-section-body {
  padding: 24px;
}

/* ── Avatar ──────────────────────────────────────────────────────────────── */
.avatar-row {
  display: flex;
  align-items: flex-start;
  gap: 24px;
  margin-bottom: 24px;

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
}

.choose-photo-label {
  display: inline-flex;
  align-items: center;
  padding: 8px 16px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--md-sys-color-on-surface);
  transition: background 0.15s, border-color 0.15s;
  gap: 4px;

  &:hover {
    background: var(--md-sys-color-surface-container-low);
    border-color: var(--md-sys-color-outline);
  }
}

/* ── Form layout ─────────────────────────────────────────────────────────── */
.form-row-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 16px;

  @media (max-width: 540px) {
    grid-template-columns: 1fr;
  }
}

/* ── Toggle / Switch rows ────────────────────────────────────────────────── */
.switch-row,
.toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 0;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  transition: background 0.15s;

  &:hover {
    background: var(--md-sys-color-surface-container-lowest, transparent);
  }
}

/* ── Session rows ────────────────────────────────────────────────────────── */
.session-row {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 16px 24px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  transition: background 0.15s;

  &:last-child {
    border-bottom: none;
  }

  &:hover {
    background: var(--md-sys-color-surface-container-low);
  }
}

.session-device-icon {
  color: var(--md-sys-color-on-surface-variant);
  margin-top: 2px;
  flex-shrink: 0;
  font-size: 28px;
}

/* ── M3 Badges ───────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  font-size: 0.72rem;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: 999px;
}

.m3-badge--success {
  background: color-mix(in srgb, #1aa563 14%, transparent);
  color: #0a7040;
}

.m3-badge--neutral {
  background: var(--md-sys-color-surface-container-high);
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Feedback ────────────────────────────────────────────────────────────── */
.feedback-success {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 10px;
  background: color-mix(in srgb, #1aa563 12%, transparent);
  border: 1px solid color-mix(in srgb, #1aa563 30%, transparent);
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
  background: var(--md-sys-color-error-container);
  border: 1px solid color-mix(in srgb, var(--md-sys-color-error) 30%, transparent);
  color: var(--md-sys-color-on-error-container, var(--md-sys-color-error));
  font-size: 0.875rem;
  margin-bottom: 16px;
}

/* ── Danger zone ─────────────────────────────────────────────────────────── */
.danger-zone-card {
  border-left: 4px solid var(--md-sys-color-error) !important;
  border-color: var(--md-sys-color-error) !important;
}

.danger-zone-header {
  border-bottom-color: color-mix(in srgb, var(--md-sys-color-error) 30%, transparent) !important;
  background: color-mix(in srgb, var(--md-sys-color-error) 6%, transparent) !important;
}

.delete-confirm-box {
  margin-top: 24px;
  padding: 16px 20px;
  border-radius: 10px;
  border: 1px solid color-mix(in srgb, var(--md-sys-color-error) 40%, transparent);
  background: color-mix(in srgb, var(--md-sys-color-error) 5%, transparent);
}

/* ── M3 Cards ────────────────────────────────────────────────────────────── */
.m3-card {
  border-radius: 16px;
  overflow: hidden;

  &--outlined {
    background: var(--md-sys-color-surface);
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}
</style>
