<template>
  <div class="container-fluid py-4" style="max-width: 780px;">
    <div class="mb-4">
      <h4 class="mb-1 fw-bold">Profile & Settings</h4>
      <p class="text-muted small mb-0">Manage your account information and security preferences.</p>
    </div>

    <!-- ── Profile ──────────────────────────────────────────────────────── -->
    <div class="card border-0 shadow-sm mb-4">
      <div class="card-header bg-white border-bottom py-3 px-4">
        <span class="fw-semibold">Profile Information</span>
      </div>
      <div class="card-body px-4 py-4">
        <div class="d-flex align-items-center gap-4 mb-4">
          <!-- Avatar preview -->
          <div class="flex-shrink-0">
            <div
              class="rounded-circle d-flex align-items-center justify-content-center fw-bold text-white"
              :style="{ width: '80px', height: '80px', fontSize: '1.75rem', backgroundImage: avatarPreviewUrl ? `url(${avatarPreviewUrl})` : undefined, backgroundSize: 'cover', backgroundPosition: 'center', backgroundColor: avatarPreviewUrl ? 'transparent' : '#635bff' }"
            >
              <span v-if="!avatarPreviewUrl">{{ authStore.userInitials }}</span>
            </div>
          </div>
          <div class="flex-fill">
            <!-- Primary: file picker -->
            <label class="form-label small fw-medium mb-1">Profile photo</label>
            <div class="d-flex align-items-center gap-2 mb-2">
              <label class="btn btn-sm btn-outline-secondary choose-photo-btn" style="cursor:pointer;white-space:nowrap;">
                Choose photo
                <input
                  ref="avatarFileInputRef"
                  type="file"
                  accept="image/*"
                  class="visually-hidden"
                  @change="onAvatarFileChange"
                />
              </label>
              <span class="text-muted small text-truncate" style="max-width:180px;">
                {{ avatarFileName || 'No file chosen' }}
              </span>
            </div>
            <!-- Secondary: URL input (collapsed by default) -->
            <div>
              <button
                class="btn btn-link btn-sm p-0 text-muted small"
                type="button"
                @click="showUrlInput = !showUrlInput"
                style="text-decoration:none;"
              >
                {{ showUrlInput ? 'Hide URL input' : 'Or enter image URL manually' }}
              </button>
              <div v-if="showUrlInput" class="mt-1">
                <input
                  v-model="profileForm.profile_picture"
                  type="url"
                  class="form-control form-control-sm"
                  placeholder="https://example.com/avatar.png"
                />
              </div>
            </div>
          </div>
        </div>

        <div class="row g-3 mb-3">
          <div class="col-sm-6">
            <label class="form-label small fw-medium">First name</label>
            <input v-model="profileForm.first_name" type="text" class="form-control" placeholder="First name" maxlength="50" />
          </div>
          <div class="col-sm-6">
            <label class="form-label small fw-medium">Last name</label>
            <input v-model="profileForm.last_name" type="text" class="form-control" placeholder="Last name" maxlength="50" />
          </div>
        </div>

        <div class="mb-3">
          <label class="form-label small fw-medium">Email</label>
          <input
            :value="authStore.profile?.email"
            type="email"
            class="form-control bg-light text-muted"
            readonly
          />
          <div class="form-text">Email cannot be changed here.</div>
        </div>

        <div class="mb-3">
          <label class="form-label small fw-medium">Phone</label>
          <input v-model="profileForm.phone" type="tel" class="form-control" placeholder="+1 555 000 0000" maxlength="20" />
        </div>

        <div v-if="profileSuccess" class="alert alert-success py-2 small">{{ profileSuccess }}</div>
        <div v-if="profileError" class="alert alert-danger py-2 small">{{ profileError }}</div>

        <button class="btn btn-primary" :disabled="savingProfile" @click="saveProfile">
          <span v-if="savingProfile" class="spinner-border spinner-border-sm me-2"></span>
          Save profile
        </button>
      </div>
    </div>

    <!-- ── Change Password ──────────────────────────────────────────────── -->
    <div class="card border-0 shadow-sm mb-4">
      <div class="card-header bg-white border-bottom py-3 px-4">
        <span class="fw-semibold">Change Password</span>
      </div>
      <div class="card-body px-4 py-4">
        <div class="mb-3">
          <label class="form-label small fw-medium">Current password</label>
          <input
            v-model="passwordForm.old_password"
            type="password"
            class="form-control"
            autocomplete="current-password"
            placeholder="Enter current password"
          />
        </div>
        <div class="mb-3">
          <label class="form-label small fw-medium">New password</label>
          <input
            v-model="passwordForm.new_password"
            type="password"
            class="form-control"
            autocomplete="new-password"
            placeholder="At least 8 characters"
          />
        </div>
        <div class="mb-3">
          <label class="form-label small fw-medium">Confirm new password</label>
          <input
            v-model="passwordForm.confirm_password"
            type="password"
            class="form-control"
            autocomplete="new-password"
            placeholder="Repeat new password"
            :class="{ 'is-invalid': passwordMismatch }"
          />
          <div v-if="passwordMismatch" class="invalid-feedback">Passwords do not match.</div>
        </div>

        <div v-if="passwordSuccess" class="alert alert-success py-2 small">{{ passwordSuccess }}</div>
        <div v-if="passwordError" class="alert alert-danger py-2 small">{{ passwordError }}</div>

        <button class="btn btn-primary" :disabled="savingPassword || !canChangePassword" @click="changePassword">
          <span v-if="savingPassword" class="spinner-border spinner-border-sm me-2"></span>
          Update password
        </button>
      </div>
    </div>

    <!-- ── Notification Preferences ────────────────────────────────────── -->
    <div class="card border-0 shadow-sm mb-4">
      <div class="card-header bg-white border-bottom py-3 px-4">
        <span class="fw-semibold">Notification Preferences</span>
      </div>
      <div class="card-body px-4 py-4">

        <!-- Link expiry warnings -->
        <div class="d-flex align-items-start gap-3 mb-4">
          <div class="form-check form-switch flex-shrink-0 mt-1 mb-0">
            <input
              id="notif-expiry"
              v-model="notifPrefs.expiry_warnings"
              type="checkbox"
              role="switch"
              class="form-check-input"
            />
          </div>
          <div>
            <label for="notif-expiry" class="fw-medium mb-1 d-block" style="cursor: pointer;">Link expiry warnings</label>
            <div class="text-muted small">Get emailed 7 days before a link expires.</div>
          </div>
        </div>

        <!-- Health alert emails -->
        <div class="d-flex align-items-start gap-3 mb-4">
          <div class="form-check form-switch flex-shrink-0 mt-1 mb-0">
            <input
              id="notif-health"
              v-model="notifPrefs.health_alerts"
              type="checkbox"
              role="switch"
              class="form-check-input"
            />
          </div>
          <div>
            <label for="notif-health" class="fw-medium mb-1 d-block" style="cursor: pointer;">Health alert emails</label>
            <div class="text-muted small">Notify when a link becomes unhealthy.</div>
          </div>
        </div>

        <!-- Weekly analytics digest -->
        <div class="d-flex align-items-start gap-3 mb-4">
          <div class="form-check form-switch flex-shrink-0 mt-1 mb-0">
            <input
              id="notif-digest"
              v-model="notifPrefs.weekly_digest"
              type="checkbox"
              role="switch"
              class="form-check-input"
            />
          </div>
          <div>
            <label for="notif-digest" class="fw-medium mb-1 d-block" style="cursor: pointer;">Weekly analytics digest</label>
            <div class="text-muted small">Summary of clicks sent every Monday.</div>
          </div>
        </div>

        <!-- Click milestone alerts -->
        <div class="d-flex align-items-start gap-3 mb-4">
          <div class="form-check form-switch flex-shrink-0 mt-1 mb-0">
            <input
              id="notif-milestone"
              v-model="notifPrefs.milestone_alerts"
              type="checkbox"
              role="switch"
              class="form-check-input"
            />
          </div>
          <div>
            <label for="notif-milestone" class="fw-medium mb-1 d-block" style="cursor: pointer;">Click milestone alerts</label>
            <div class="text-muted small">Notified at 100, 1K, 10K clicks per link.</div>
          </div>
        </div>

        <p class="text-muted small mb-3">
          Email preferences affect automated system emails. Transactional emails (password reset, verification) are always sent.
        </p>

        <div v-if="notifSaved" class="alert alert-success py-2 small mb-3">Preferences saved.</div>

        <button class="btn btn-primary" @click="saveNotifPrefs">Save preferences</button>
      </div>
    </div>

    <!-- ── Active Sessions ──────────────────────────────────────────────── -->
    <div class="card border-0 shadow-sm mb-4">
      <div class="card-header bg-white border-bottom py-3 px-4 d-flex align-items-center justify-content-between">
        <span class="fw-semibold">Active Sessions</span>
        <button
          class="btn btn-sm btn-outline-danger"
          :disabled="revokingAll"
          @click="logoutAllOther"
        >
          <span v-if="revokingAll" class="spinner-border spinner-border-sm me-1"></span>
          Logout all other devices
        </button>
      </div>
      <div class="card-body px-0 py-0">
        <div v-if="sessionsLoading" class="text-center py-4">
          <div class="spinner-border spinner-border-sm text-primary"></div>
        </div>
        <div v-else-if="sessions.length === 0" class="text-center text-muted py-4 small">
          No active sessions found.
        </div>
        <div v-else class="list-group list-group-flush">
          <div
            v-for="session in sessions"
            :key="session.id"
            class="list-group-item py-3 px-4"
          >
            <div class="d-flex align-items-start gap-3">
              <!-- Device icon -->
              <div class="flex-shrink-0 text-muted mt-1" style="font-size: 1.2rem;">
                {{ deviceIcon(session.device_type) }}
              </div>
              <div class="flex-fill min-w-0">
                <div class="d-flex align-items-center gap-2 flex-wrap">
                  <span class="fw-medium small">{{ sessionLabel(session) }}</span>
                  <span v-if="session.is_current" class="badge text-bg-success" style="font-size: 0.7rem;">Current</span>
                  <span v-if="session.login_method" class="badge bg-light text-secondary border" style="font-size: 0.7rem;">{{ session.login_method }}</span>
                </div>
                <div class="text-muted" style="font-size: 0.78rem;">
                  <span v-if="session.ip_address">{{ session.ip_address }}</span>
                  <span v-if="session.ip_address && session.last_activity_at"> &middot; </span>
                  <span>Last active {{ formatRelative(session.last_activity_at) }}</span>
                </div>
                <div class="text-muted" style="font-size: 0.75rem;">
                  Signed in {{ formatDate(session.created_at) }}
                </div>
              </div>
              <button
                v-if="!session.is_current"
                class="btn btn-sm btn-outline-danger flex-shrink-0"
                :disabled="revokingId === session.id"
                @click="revokeSession(session)"
              >
                <span v-if="revokingId === session.id" class="spinner-border spinner-border-sm"></span>
                <span v-else>Revoke</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Data Export ──────────────────────────────────────────────────── -->
    <div class="card border-0 shadow-sm mb-4">
      <div class="card-header bg-white border-bottom py-3 px-4">
        <span class="fw-semibold">Data Export</span>
      </div>
      <div class="card-body px-4 py-4">
        <div class="d-flex align-items-start justify-content-between gap-3 flex-wrap">
          <div>
            <div class="fw-medium mb-1">Download your data</div>
            <p class="text-muted small mb-0">
              Export a JSON file containing your profile information and all your short links.
              Your data will be ready to download immediately.
            </p>
          </div>
          <button class="btn btn-outline-primary flex-shrink-0" :disabled="exportingData" @click="downloadExport">
            <span v-if="exportingData" class="spinner-border spinner-border-sm me-2"></span>
            <span v-else>Download my data</span>
          </button>
        </div>
        <div v-if="exportError" class="alert alert-danger py-2 small mt-3 mb-0">{{ exportError }}</div>
      </div>
    </div>

    <!-- ── Danger Zone ──────────────────────────────────────────────────── -->
    <div class="card border-0 shadow-sm border border-danger border-opacity-25">
      <div class="card-header bg-danger bg-opacity-10 border-bottom border-danger border-opacity-25 py-3 px-4">
        <span class="fw-semibold text-danger">Danger Zone</span>
      </div>
      <div class="card-body px-4 py-4">
        <div class="d-flex align-items-start justify-content-between gap-3 flex-wrap">
          <div>
            <div class="fw-medium mb-1">Delete account</div>
            <p class="text-muted small mb-0">
              Permanently delete your account and all associated links, analytics, and data.
              This action <strong>cannot be undone</strong>.
            </p>
          </div>
          <button class="btn btn-outline-danger flex-shrink-0" @click="showDeleteConfirm = true">
            Delete account
          </button>
        </div>

        <!-- Delete confirmation -->
        <div v-if="showDeleteConfirm" class="mt-4 p-3 border border-danger border-opacity-25 rounded bg-danger bg-opacity-5">
          <p class="small fw-medium mb-2">
            To confirm, type <code>DELETE</code> in the box below:
          </p>
          <div class="d-flex gap-2 align-items-center">
            <input
              v-model="deleteConfirmText"
              type="text"
              class="form-control form-control-sm"
              placeholder="DELETE"
              style="max-width: 200px;"
            />
            <button
              class="btn btn-sm btn-danger"
              :disabled="deleteConfirmText !== 'DELETE' || deletingAccount"
              @click="deleteAccount"
            >
              <span v-if="deletingAccount" class="spinner-border spinner-border-sm me-1"></span>
              Confirm delete
            </button>
            <button class="btn btn-sm btn-outline-secondary" @click="showDeleteConfirm = false; deleteConfirmText = ''">
              Cancel
            </button>
          </div>
          <div v-if="deleteError" class="alert alert-danger py-2 small mt-2 mb-0">{{ deleteError }}</div>
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

<style scoped>
.btn-primary {
  background-color: #635bff;
  border-color: #635bff;
}
.btn-primary:hover:not(:disabled) {
  background-color: #5249e0;
  border-color: #5249e0;
}
.choose-photo-btn {
  border-color: #d1d5db;
  color: #374151;
}
.choose-photo-btn:hover {
  background-color: #f3f4f6;
  border-color: #9ca3af;
}
</style>
