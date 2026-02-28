<template>
  <div class="teams-page">

    <!-- Toast notifications -->
    <Transition name="snack">
      <div v-if="toasts.length > 0" class="m3-snackbar">
        <span class="material-symbols-outlined snack-icon">
          {{ toasts[0]?.type === 'error' ? 'error' : 'check_circle' }}
        </span>
        <span class="snack-text">{{ toasts[0]?.message }}</span>
        <button class="btn-text" @click="removeToast(toasts[0]?.id ?? 0)">Dismiss</button>
      </div>
    </Transition>

    <!-- Page Header -->
    <div class="page-header">
      <div class="page-header__left">
        <h1 class="page-title">Teams</h1>
        <p class="page-subtitle">Manage your teams and collaborate with others.</p>
      </div>
      <div class="page-header__actions">
        <button class="btn-filled" @click="openCreateModal">
          <span class="material-symbols-outlined">add</span>
          Create Team
        </button>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="loading-center">
      <div class="css-spinner"></div>
    </div>

    <!-- Empty state -->
    <div v-else-if="teams.length === 0" class="an-card empty-state">
      <div class="empty-icon">
        <span class="material-symbols-outlined">group</span>
      </div>
      <div class="empty-title">No teams yet</div>
      <p class="empty-desc">Create a team to collaborate with your colleagues.</p>
      <button class="btn-filled" @click="openCreateModal">
        <span class="material-symbols-outlined">add</span>
        Create Your First Team
      </button>
    </div>

    <!-- Teams grid -->
    <div v-else class="teams-grid">
      <div
        v-for="team in teams"
        :key="team.id"
        class="team-card"
        :class="{ 'team-card--selected': selectedTeam?.id === team.id }"
        @click="selectTeam(team)"
      >
        <div class="team-card__body">
          <div class="team-card__top">
            <div class="team-avatar">{{ team.name.charAt(0).toUpperCase() }}</div>
            <span class="m3-badge" :class="getRoleBadgeM3Class(getMyRole(team))">
              {{ getMyRole(team) }}
            </span>
          </div>
          <div class="team-name">{{ team.name }}</div>
          <div class="team-slug">@{{ team.slug }}</div>
          <p v-if="team.description" class="team-description">{{ team.description }}</p>
          <div class="team-members-count">
            <span class="material-symbols-outlined member-icon">group</span>
            <span>{{ team.member_count }} member{{ team.member_count !== 1 ? 's' : '' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Team Detail Panel ──────────────────────────────────────────────── -->
    <div v-if="selectedTeam" class="an-card team-detail">
      <div class="an-card-header">
        <div class="an-card-header__left">
          <div class="team-avatar-sm">{{ selectedTeam.name.charAt(0).toUpperCase() }}</div>
          <div>
            <div class="team-detail-name">{{ selectedTeam.name }}</div>
            <div class="team-detail-slug">@{{ selectedTeam.slug }}</div>
          </div>
        </div>
        <div class="team-detail-actions">
          <button class="btn-filled" v-if="canManageMembers(selectedTeam)" @click="openInviteModal">
            <span class="material-symbols-outlined">person_add</span>
            Invite Member
          </button>
          <button class="btn-outlined" v-if="isTeamOwner(selectedTeam)" @click="openEditModal">
            <span class="material-symbols-outlined">edit</span>
            Edit
          </button>
          <button class="btn-outlined btn-danger" v-if="!isTeamOwner(selectedTeam)" @click="confirmLeave(selectedTeam)">
            Leave Team
          </button>
          <button class="btn-outlined btn-danger" v-if="isTeamOwner(selectedTeam)" @click="confirmDelete(selectedTeam)">
            Delete Team
          </button>
        </div>
      </div>

      <!-- Members loading -->
      <div v-if="membersLoading" class="loading-center loading-center--sm">
        <div class="css-spinner"></div>
      </div>

      <!-- Members table -->
      <div v-else class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>Member</th>
              <th>Role</th>
              <th>Status</th>
              <th>Joined</th>
              <th v-if="canManageMembers(selectedTeam)" class="th-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="members.length === 0">
              <td :colspan="canManageMembers(selectedTeam) ? 5 : 4" class="td-empty">No members yet.</td>
            </tr>
            <tr v-for="member in members" :key="member.id">
              <td>
                <div class="member-row">
                  <div class="member-avatar">{{ getMemberInitial(member) }}</div>
                  <div>
                    <div class="member-name">{{ getMemberName(member) }}</div>
                    <div class="member-email">{{ member.email }}</div>
                  </div>
                </div>
              </td>
              <td>
                <span class="m3-badge" :class="getRoleBadgeM3Class(member.role)">{{ member.role }}</span>
              </td>
              <td>
                <span class="m3-badge" :class="getStatusBadgeM3Class(member.status)">{{ member.status }}</span>
              </td>
              <td class="cell-muted cell-sm">
                {{ member.joined_at ? formatDate(member.joined_at) : '—' }}
              </td>
              <td v-if="canManageMembers(selectedTeam)" class="td-right">
                <div class="member-actions">
                  <AppSelect
                    v-if="member.role !== 'owner' && member.user_id !== authStore.profile?.id"
                    :model-value="member.role"
                    @update:model-value="changeMemberRole(member, $event)"
                    label="Role"
                  >
                    <option value="admin">Admin</option>
                    <option value="member">Member</option>
                  </AppSelect>
                  <button class="btn-icon btn-icon--danger"
                    v-if="member.role !== 'owner' && member.user_id !== authStore.profile?.id"
                    @click="confirmRemoveMember(member)"
                    title="Remove member"
                  >
                    <span class="material-symbols-outlined">person_remove</span>
                  </button>
                  <span v-if="member.role === 'owner'" class="owner-label">Owner</span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ── Create Team Dialog ─────────────────────────────────────────────── -->
    <BaseModal v-model="showCreateModal" size="md" @closed="closeCreateModal">
      <template #headline>
        <span class="material-symbols-outlined modal-icon modal-icon--primary">group_add</span>
        Create Team
      </template>
      <div class="modal-form">
        <label class="form-field">
          <span class="form-field__label">Team Name *</span>
          <input type="text" class="form-input" :value="createForm.name"
            @input="createForm.name = ($event.target as HTMLInputElement).value; autoSlug()"
            required />
        </label>
        <div class="slug-row">
          <span class="slug-prefix">@</span>
          <label class="form-field form-field--grow">
            <span class="form-field__label">Slug *</span>
            <input type="text" class="form-input" :value="createForm.slug"
              @input="createForm.slug = ($event.target as HTMLInputElement).value"
              pattern="[a-z0-9-]+" required />
          </label>
        </div>
        <p class="form-hint">Lowercase letters, numbers, and hyphens only.</p>
        <label class="form-field">
          <span class="form-field__label">Description (optional)</span>
          <textarea class="form-textarea" :value="createForm.description"
            @input="createForm.description = ($event.target as HTMLTextAreaElement).value"
            rows="2"></textarea>
        </label>
        <div v-if="createError" class="feedback-error">
          <span class="material-symbols-outlined feedback-icon">error</span>
          {{ createError }}
        </div>
      </div>
      <template #actions>
        <button class="btn-text" @click="closeCreateModal">Cancel</button>
        <button class="btn-filled" :disabled="createLoading" @click="submitCreateTeam">
          <div v-if="createLoading" class="css-spinner css-spinner--sm css-spinner--white"></div>
          Create Team
        </button>
      </template>
    </BaseModal>

    <!-- ── Edit Team Dialog ───────────────────────────────────────────────── -->
    <BaseModal v-model="showEditModal" size="md" @closed="showEditModal = false">
      <template #headline>
        <span class="material-symbols-outlined modal-icon modal-icon--primary">edit</span>
        Edit Team
      </template>
      <div class="modal-form">
        <label class="form-field">
          <span class="form-field__label">Team Name</span>
          <input type="text" class="form-input" :value="editForm.name"
            @input="editForm.name = ($event.target as HTMLInputElement).value" />
        </label>
        <label class="form-field">
          <span class="form-field__label">Description</span>
          <textarea class="form-textarea" :value="editForm.description"
            @input="editForm.description = ($event.target as HTMLTextAreaElement).value"
            rows="2"></textarea>
        </label>
        <div v-if="editError" class="feedback-error">
          <span class="material-symbols-outlined feedback-icon">error</span>
          {{ editError }}
        </div>
      </div>
      <template #actions>
        <button class="btn-text" @click="showEditModal = false">Cancel</button>
        <button class="btn-filled" :disabled="editLoading" @click="submitEditTeam">
          <div v-if="editLoading" class="css-spinner css-spinner--sm css-spinner--white"></div>
          Save Changes
        </button>
      </template>
    </BaseModal>

    <!-- ── Invite Member Dialog ───────────────────────────────────────────── -->
    <BaseModal v-model="showInviteModal" size="md" @closed="closeInviteModal">
      <template #headline>
        <span class="material-symbols-outlined modal-icon modal-icon--primary">person_add</span>
        Invite Member
      </template>
      <div class="modal-form">
        <label class="form-field">
          <span class="form-field__label">Email Address *</span>
          <input type="email" class="form-input" :value="inviteForm.email"
            @input="inviteForm.email = ($event.target as HTMLInputElement).value"
            required />
        </label>
        <AppSelect :model-value="inviteForm.role"
          @update:model-value="inviteForm.role = $event as 'admin' | 'member'"
          label="Role">
          <option value="member">Member</option>
          <option value="admin">Admin</option>
        </AppSelect>
        <div v-if="inviteError" class="feedback-error">
          <span class="material-symbols-outlined feedback-icon">error</span>
          {{ inviteError }}
        </div>
      </div>
      <template #actions>
        <button class="btn-text" @click="closeInviteModal">Cancel</button>
        <button class="btn-filled" :disabled="inviteLoading" @click="submitInvite">
          <div v-if="inviteLoading" class="css-spinner css-spinner--sm css-spinner--white"></div>
          Send Invitation
        </button>
      </template>
    </BaseModal>

    <!-- ── Confirm Delete Dialog ──────────────────────────────────────────── -->
    <BaseModal v-model="showDeleteConfirm" size="sm" :persistent="false" @closed="showDeleteConfirm = false">
      <template #headline>
        <span class="material-symbols-outlined modal-icon modal-icon--danger">delete</span>
        Delete Team
      </template>
      <p class="modal-text">
        Are you sure you want to delete <strong>{{ teamToDelete?.name }}</strong>? This action cannot be undone.
      </p>
      <template #actions>
        <button class="btn-text" @click="showDeleteConfirm = false">Cancel</button>
        <button class="btn-filled btn-danger" :disabled="deleteLoading" @click="deleteTeam">
          <div v-if="deleteLoading" class="css-spinner css-spinner--sm css-spinner--white"></div>
          Delete Team
        </button>
      </template>
    </BaseModal>

    <!-- ── Confirm Leave Dialog ───────────────────────────────────────────── -->
    <BaseModal v-model="showLeaveConfirm" size="sm" :persistent="false" @closed="showLeaveConfirm = false">
      <template #headline>
        <span class="material-symbols-outlined modal-icon">logout</span>
        Leave Team
      </template>
      <p class="modal-text">Are you sure you want to leave <strong>{{ teamToLeave?.name }}</strong>?</p>
      <template #actions>
        <button class="btn-text" @click="showLeaveConfirm = false">Cancel</button>
        <button class="btn-filled btn-danger" :disabled="leaveLoading" @click="leaveTeam">
          <div v-if="leaveLoading" class="css-spinner css-spinner--sm css-spinner--white"></div>
          Leave Team
        </button>
      </template>
    </BaseModal>

    <!-- ── Confirm Remove Member Dialog ──────────────────────────────────── -->
    <BaseModal v-model="showRemoveConfirm" size="sm" :persistent="false" @closed="showRemoveConfirm = false">
      <template #headline>
        <span class="material-symbols-outlined modal-icon modal-icon--danger">person_remove</span>
        Remove Member
      </template>
      <p class="modal-text">
        Remove <strong>{{ memberToRemove?.email }}</strong> from the team? This cannot be undone.
      </p>
      <template #actions>
        <button class="btn-text" @click="showRemoveConfirm = false">Cancel</button>
        <button class="btn-filled btn-danger" :disabled="removeLoading" @click="removeMember">
          <div v-if="removeLoading" class="css-spinner css-spinner--sm css-spinner--white"></div>
          Remove
        </button>
      </template>
    </BaseModal>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import BaseModal from '@/components/BaseModal.vue';
import AppSelect from '@/components/AppSelect.vue';
import { useAuthStore } from '@/stores/auth';
import teamsApi from '@/api/teams';
import type { Team, TeamMember } from '@/types/teams';

const authStore = useAuthStore();

// ── State ─────────────────────────────────────────────────────────────────────
const teams = ref<Team[]>([]);
const loading = ref(false);
const selectedTeam = ref<Team | null>(null);
const members = ref<TeamMember[]>([]);
const membersLoading = ref(false);

// Toasts
const toasts = ref<{ id: number; message: string; type: 'success' | 'error' }[]>([]);
let toastCounter = 0;

// Create modal
const showCreateModal = ref(false);
const createForm = ref({ name: '', slug: '', description: '' });
const createLoading = ref(false);
const createError = ref('');

// Edit modal
const showEditModal = ref(false);
const editForm = ref({ name: '', description: '' });
const editLoading = ref(false);
const editError = ref('');

// Invite modal
const showInviteModal = ref(false);
const inviteForm = ref({ email: '', role: 'member' as 'admin' | 'member' });
const inviteLoading = ref(false);
const inviteError = ref('');

// Delete confirm
const showDeleteConfirm = ref(false);
const teamToDelete = ref<Team | null>(null);
const deleteLoading = ref(false);

// Leave confirm
const showLeaveConfirm = ref(false);
const teamToLeave = ref<Team | null>(null);
const leaveLoading = ref(false);

// Remove member confirm
const showRemoveConfirm = ref(false);
const memberToRemove = ref<TeamMember | null>(null);
const removeLoading = ref(false);

// ── Toast helpers ─────────────────────────────────────────────────────────────
function showToast(message: string, type: 'success' | 'error' = 'success') {
  const id = ++toastCounter;
  toasts.value.push({ id, message, type });
  setTimeout(() => removeToast(id), 4000);
}

function removeToast(id: number) {
  toasts.value = toasts.value.filter((t) => t.id !== id);
}

// ── Data loading ──────────────────────────────────────────────────────────────
async function loadTeams() {
  loading.value = true;
  try {
    const res = await teamsApi.listMyTeams();
    teams.value = res.data.data ?? [];
  } catch (e: any) {
    showToast(e?.response?.data?.error?.description ?? 'Failed to load teams', 'error');
  } finally {
    loading.value = false;
  }
}

async function loadMembers(teamId: string) {
  membersLoading.value = true;
  members.value = [];
  try {
    const res = await teamsApi.listMembers(teamId);
    members.value = res.data.data ?? [];
  } catch (e: any) {
    showToast(e?.response?.data?.error?.description ?? 'Failed to load members', 'error');
  } finally {
    membersLoading.value = false;
  }
}

// ── Team selection ────────────────────────────────────────────────────────────
function selectTeam(team: Team) {
  if (selectedTeam.value?.id === team.id) {
    selectedTeam.value = null;
    members.value = [];
    return;
  }
  selectedTeam.value = team;
  loadMembers(team.id);
}

// ── Role helpers ──────────────────────────────────────────────────────────────
function getMyRole(team: Team): string {
  if (team.owner_id === authStore.profile?.id) return 'owner';
  const me = members.value.find((m) => m.user_id === authStore.profile?.id);
  return me?.role ?? 'member';
}

function isTeamOwner(team: Team): boolean {
  return team.owner_id === authStore.profile?.id;
}

function canManageMembers(team: Team): boolean {
  if (isTeamOwner(team)) return true;
  const me = members.value.find((m) => m.user_id === authStore.profile?.id);
  return me?.role === 'admin';
}

function getRoleBadgeClass(role: string): string {
  switch (role) {
    case 'owner': return 'text-bg-primary';
    case 'admin': return 'text-bg-info';
    default: return 'text-bg-secondary';
  }
}

function getRoleBadgeM3Class(role: string): string {
  switch (role) {
    case 'owner': return 'm3-badge--primary';
    case 'admin': return 'm3-badge--info';
    default: return 'm3-badge--neutral';
  }
}

function getStatusBadgeClass(status: string): string {
  switch (status) {
    case 'active': return 'text-bg-success';
    case 'pending': return 'text-bg-warning text-dark';
    case 'declined': return 'text-bg-danger';
    default: return 'text-bg-secondary';
  }
}

function getStatusBadgeM3Class(status: string): string {
  switch (status) {
    case 'active': return 'm3-badge--success';
    case 'pending': return 'm3-badge--warning';
    case 'declined': return 'm3-badge--error';
    default: return 'm3-badge--neutral';
  }
}

function getMemberInitial(member: TeamMember): string {
  if (member.first_name) return member.first_name.charAt(0).toUpperCase();
  if (member.email) return member.email.charAt(0).toUpperCase();
  return '?';
}

function getMemberName(member: TeamMember): string {
  if (member.first_name || member.last_name) {
    return `${member.first_name ?? ''} ${member.last_name ?? ''}`.trim();
  }
  return member.invite_email ?? member.email ?? '—';
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' });
}

// ── Create Team ───────────────────────────────────────────────────────────────
function openCreateModal() {
  createForm.value = { name: '', slug: '', description: '' };
  createError.value = '';
  showCreateModal.value = true;
}

function closeCreateModal() {
  showCreateModal.value = false;
  createError.value = '';
}

function autoSlug() {
  createForm.value.slug = createForm.value.name
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '');
}

async function submitCreateTeam() {
  createError.value = '';
  createLoading.value = true;
  try {
    const res = await teamsApi.createTeam(createForm.value);
    const newTeam = res.data.data;
    if (newTeam) {
      teams.value.unshift(newTeam);
      showToast('Team created successfully!');
      closeCreateModal();
    }
  } catch (e: any) {
    createError.value = e?.response?.data?.error?.description ?? 'Failed to create team';
  } finally {
    createLoading.value = false;
  }
}

// ── Edit Team ─────────────────────────────────────────────────────────────────
function openEditModal() {
  if (!selectedTeam.value) return;
  editForm.value = {
    name: selectedTeam.value.name,
    description: selectedTeam.value.description ?? '',
  };
  editError.value = '';
  showEditModal.value = true;
}

async function submitEditTeam() {
  if (!selectedTeam.value) return;
  editError.value = '';
  editLoading.value = true;
  try {
    const res = await teamsApi.updateTeam(selectedTeam.value.id, editForm.value);
    const updated = res.data.data;
    if (updated) {
      const idx = teams.value.findIndex((t) => t.id === updated.id);
      if (idx !== -1) teams.value[idx] = updated;
      selectedTeam.value = updated;
      showToast('Team updated successfully!');
      showEditModal.value = false;
    }
  } catch (e: any) {
    editError.value = e?.response?.data?.error?.description ?? 'Failed to update team';
  } finally {
    editLoading.value = false;
  }
}

// ── Invite Member ─────────────────────────────────────────────────────────────
function openInviteModal() {
  inviteForm.value = { email: '', role: 'member' };
  inviteError.value = '';
  showInviteModal.value = true;
}

function closeInviteModal() {
  showInviteModal.value = false;
  inviteError.value = '';
}

async function submitInvite() {
  if (!selectedTeam.value) return;
  inviteError.value = '';
  inviteLoading.value = true;
  try {
    const res = await teamsApi.inviteMember(selectedTeam.value.id, inviteForm.value);
    const member = res.data.data;
    if (member) {
      members.value.push(member);
      showToast('Invitation sent successfully!');
      closeInviteModal();
    }
  } catch (e: any) {
    inviteError.value = e?.response?.data?.error?.description ?? 'Failed to send invitation';
  } finally {
    inviteLoading.value = false;
  }
}

// ── Change Member Role ────────────────────────────────────────────────────────
async function changeMemberRole(member: TeamMember, newRole: string) {
  if (!selectedTeam.value) return;
  try {
    await teamsApi.updateMemberRole(selectedTeam.value.id, member.user_id, newRole);
    const idx = members.value.findIndex((m) => m.id === member.id);
    if (idx !== -1) members.value[idx] = { ...members.value[idx]!, role: newRole as 'admin' | 'member' };
    showToast('Role updated successfully!');
  } catch (e: any) {
    showToast(e?.response?.data?.error?.description ?? 'Failed to update role', 'error');
  }
}

// ── Remove Member ─────────────────────────────────────────────────────────────
function confirmRemoveMember(member: TeamMember) {
  memberToRemove.value = member;
  showRemoveConfirm.value = true;
}

async function removeMember() {
  if (!selectedTeam.value || !memberToRemove.value) return;
  removeLoading.value = true;
  try {
    await teamsApi.removeMember(selectedTeam.value.id, memberToRemove.value.user_id);
    members.value = members.value.filter((m) => m.id !== memberToRemove.value!.id);
    showToast('Member removed successfully!');
    showRemoveConfirm.value = false;
    memberToRemove.value = null;
    // Update member count
    const idx = teams.value.findIndex((t) => t.id === selectedTeam.value?.id);
    if (idx !== -1) teams.value[idx]!.member_count = Math.max(0, teams.value[idx]!.member_count - 1);
  } catch (e: any) {
    showToast(e?.response?.data?.error?.description ?? 'Failed to remove member', 'error');
  } finally {
    removeLoading.value = false;
  }
}

// ── Delete Team ───────────────────────────────────────────────────────────────
function confirmDelete(team: Team) {
  teamToDelete.value = team;
  showDeleteConfirm.value = true;
}

async function deleteTeam() {
  if (!teamToDelete.value) return;
  deleteLoading.value = true;
  try {
    await teamsApi.deleteTeam(teamToDelete.value.id);
    teams.value = teams.value.filter((t) => t.id !== teamToDelete.value!.id);
    if (selectedTeam.value?.id === teamToDelete.value.id) {
      selectedTeam.value = null;
      members.value = [];
    }
    showToast('Team deleted successfully!');
    showDeleteConfirm.value = false;
    teamToDelete.value = null;
  } catch (e: any) {
    showToast(e?.response?.data?.error?.description ?? 'Failed to delete team', 'error');
  } finally {
    deleteLoading.value = false;
  }
}

// ── Leave Team ────────────────────────────────────────────────────────────────
function confirmLeave(team: Team) {
  teamToLeave.value = team;
  showLeaveConfirm.value = true;
}

async function leaveTeam() {
  if (!teamToLeave.value) return;
  leaveLoading.value = true;
  try {
    await teamsApi.leaveTeam(teamToLeave.value.id);
    teams.value = teams.value.filter((t) => t.id !== teamToLeave.value!.id);
    if (selectedTeam.value?.id === teamToLeave.value.id) {
      selectedTeam.value = null;
      members.value = [];
    }
    showToast('Left team successfully!');
    showLeaveConfirm.value = false;
    teamToLeave.value = null;
  } catch (e: any) {
    showToast(e?.response?.data?.error?.description ?? 'Failed to leave team', 'error');
  } finally {
    leaveLoading.value = false;
  }
}

// ── Lifecycle ─────────────────────────────────────────────────────────────────
onMounted(() => {
  loadTeams();
});
</script>

<style scoped lang="scss">
.teams-page {
  max-width: 1100px;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* ── Page Header ─────────────────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;

  &__left { display: flex; flex-direction: column; gap: 4px; }
  &__actions { display: flex; align-items: center; gap: 8px; flex-shrink: 0; }
}

.page-title {
  font-size: 1.375rem;
  font-weight: 700;
  margin: 0;
  color: var(--md-sys-color-on-surface);
}

.page-subtitle {
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

  &--sm { width: 16px; height: 16px; border-width: 2px; }
  &--white { border-color: rgba(255,255,255,0.35); border-top-color: #fff; }
}

@keyframes spin { to { transform: rotate(360deg); } }

/* ── Loading ─────────────────────────────────────────────────────────────── */
.loading-center {
  display: flex;
  justify-content: center;
  padding: 64px;

  &--sm { padding: 32px; }
}

/* ── Cards ───────────────────────────────────────────────────────────────── */
.an-card {
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  overflow: hidden;
}

.an-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  background: var(--md-sys-color-surface-container-low);
  flex-wrap: wrap;

  &__left {
    display: flex;
    align-items: center;
    gap: 12px;
  }
}

/* ── Empty state ─────────────────────────────────────────────────────────── */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 64px 24px;
  text-align: center;
}

.empty-icon {
  width: 72px;
  height: 72px;
  border-radius: 20px;
  background: var(--md-sys-color-surface-container-low);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;

  .material-symbols-outlined {
    font-size: 2rem;
    color: var(--md-sys-color-on-surface-variant);
    opacity: 0.6;
  }
}

.empty-title {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--md-sys-color-on-surface);
}

.empty-desc {
  margin: 0 0 20px;
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Teams grid ──────────────────────────────────────────────────────────── */
.teams-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.team-card {
  background: var(--md-sys-color-surface);
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 14px;
  overflow: hidden;
  cursor: pointer;
  transition: box-shadow 0.15s, border-color 0.15s;

  &:hover {
    box-shadow: 0 4px 16px rgba(99, 91, 255, 0.15);
    border-color: var(--md-sys-color-primary);
  }

  &--selected {
    border-color: var(--md-sys-color-primary);
    border-width: 2px;
    box-shadow: 0 0 0 4px rgba(99, 91, 255, 0.12);
  }

  &__body { padding: 20px; }

  &__top {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 12px;
  }
}

.team-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
  margin-bottom: 4px;
}

.team-slug {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
  margin-bottom: 6px;
}

.team-description {
  font-size: 0.85rem;
  color: var(--md-sys-color-on-surface-variant);
  margin: 0 0 8px;
  line-height: 1.4;
}

.team-members-count {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.85rem;
  color: var(--md-sys-color-on-surface-variant);
}

.member-icon {
  font-size: 16px;
}

/* ── Avatars ─────────────────────────────────────────────────────────────── */
.team-avatar {
  width: 42px;
  height: 42px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  font-weight: 700;
  font-size: 1.125rem;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.team-avatar-sm {
  width: 34px;
  height: 34px;
  background: var(--md-sys-color-primary);
  color: var(--md-sys-color-on-primary);
  font-weight: 700;
  font-size: 0.9rem;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

/* ── Team detail ─────────────────────────────────────────────────────────── */
.team-detail-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--md-sys-color-on-surface);
}

.team-detail-slug {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
}

.team-detail-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  align-items: center;
}

/* ── Table ───────────────────────────────────────────────────────────────── */
.table-wrapper {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;

  th {
    padding: 12px 16px;
    text-align: left;
    font-weight: 600;
    font-size: 0.8rem;
    color: var(--md-sys-color-on-surface-variant);
    background: var(--md-sys-color-surface-container-low);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    white-space: nowrap;
  }

  td {
    padding: 12px 16px;
    vertical-align: middle;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface);
  }

  tbody tr:last-child td { border-bottom: none; }
  tbody tr:hover td { background: var(--md-sys-color-surface-container-low); }
}

.th-right { text-align: right; }
.td-right { text-align: right; }
.td-empty { text-align: center; color: var(--md-sys-color-on-surface-variant); padding: 32px 16px !important; }
.cell-muted { color: var(--md-sys-color-on-surface-variant); }
.cell-sm { font-size: 0.8rem; }

/* ── Member row ──────────────────────────────────────────────────────────── */
.member-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.member-avatar {
  width: 32px;
  height: 32px;
  background: var(--md-sys-color-surface-container-low);
  color: var(--md-sys-color-on-surface-variant);
  font-weight: 600;
  font-size: 0.8rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.member-name {
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--md-sys-color-on-surface);
}

.member-email {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
}

.member-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
}

.owner-label {
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
}

/* ── Badges ──────────────────────────────────────────────────────────────── */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 6px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;
  text-transform: capitalize;

  &--primary { background: rgba(99, 91, 255, 0.12); color: var(--md-sys-color-primary); }
  &--info { background: rgba(2, 136, 209, 0.12); color: #0277bd; }
  &--success { background: rgba(22, 163, 74, 0.12); color: #16a34a; }
  &--warning { background: rgba(245, 158, 11, 0.12); color: #92400e; }
  &--error { background: rgba(220, 38, 38, 0.12); color: var(--md-sys-color-error); }
  &--neutral {
    background: var(--md-sys-color-surface-container-low);
    color: var(--md-sys-color-on-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}

/* ── Danger button ───────────────────────────────────────────────────────── */
.btn-danger {
  color: var(--md-sys-color-error) !important;
  border-color: var(--md-sys-color-error) !important;
}

.btn-icon--danger {
  color: var(--md-sys-color-error);
  &:hover { background: rgba(220,38,38,0.08); }
}

/* ── Modal ───────────────────────────────────────────────────────────────── */
.modal-form {
  min-width: 480px;
  max-width: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.modal-text {
  color: var(--md-sys-color-on-surface-variant);
  margin: 0;
}

.modal-icon {
  font-size: 20px;
  vertical-align: middle;
  margin-right: 4px;

  &--primary { color: var(--md-sys-color-primary); }
  &--danger { color: var(--md-sys-color-error); }
}

/* ── Form fields ─────────────────────────────────────────────────────────── */
.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;

  &--grow { flex: 1; }
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
  width: 100%;
  box-sizing: border-box;

  &:focus {
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary) 12%, transparent);
  }
}

.form-textarea {
  padding: 8px 12px;
  border: 1px solid var(--md-sys-color-outline-variant);
  border-radius: 8px;
  background: var(--md-sys-color-surface);
  color: var(--md-sys-color-on-surface);
  font-size: 0.9rem;
  outline: none;
  width: 100%;
  box-sizing: border-box;
  resize: vertical;
  font-family: inherit;

  &:focus {
    border-color: var(--md-sys-color-primary);
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary) 12%, transparent);
  }
}

.form-hint {
  margin: -8px 0 0;
  font-size: 0.8rem;
  color: var(--md-sys-color-on-surface-variant);
}

.slug-row {
  display: flex;
  align-items: flex-end;
  gap: 8px;
}

.slug-prefix {
  font-size: 0.9rem;
  color: var(--md-sys-color-on-surface-variant);
  padding-bottom: 8px;
  flex-shrink: 0;
}

/* ── Feedback ────────────────────────────────────────────────────────────── */
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

.feedback-icon { font-size: 18px; flex-shrink: 0; }

/* ── Snackbar ────────────────────────────────────────────────────────────── */
.m3-snackbar {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  background: #313033;
  color: #fff;
  border-radius: 4px;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 280px;
  max-width: 560px;
  z-index: 1000;
  box-shadow: 0 4px 12px rgba(0,0,0,0.24);
}

.snack-icon { font-size: 20px; flex-shrink: 0; }
.snack-text { flex: 1; font-size: 0.875rem; }

.snack-enter-active, .snack-leave-active { transition: all 0.25s; }
.snack-enter-from, .snack-leave-to { transform: translateY(80px); opacity: 0; }
</style>
