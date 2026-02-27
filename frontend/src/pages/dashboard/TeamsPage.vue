<template>
  <div class="teams-page">

    <!-- Toast notifications -->
    <Transition name="snack">
      <div
        v-if="toasts.length > 0"
        class="m3-snackbar"
      >
        <span class="material-symbols-outlined" style="font-size:20px;">
          {{ toasts[0].type === 'error' ? 'error' : 'check_circle' }}
        </span>
        <span style="flex:1;">{{ toasts[0].message }}</span>
        <md-text-button @click="removeToast(toasts[0].id)" style="--md-text-button-label-text-color:#CFBCFF;">Dismiss</md-text-button>
      </div>
    </Transition>

    <!-- Page header -->
    <div class="page-header-row">
      <div>
        <h1 class="md-headline-small">Teams</h1>
        <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin:4px 0 0;">
          Manage your teams and collaborate with others.
        </p>
      </div>
      <md-filled-button @click="openCreateModal">
        <span class="material-symbols-outlined" slot="icon">add</span>
        Create Team
      </md-filled-button>
    </div>

    <!-- Loading state -->
    <div v-if="loading" style="text-align:center;padding:64px 0;">
      <md-circular-progress indeterminate />
    </div>

    <!-- Empty state -->
    <div v-else-if="teams.length === 0" class="m3-card m3-card--elevated m3-empty-state">
      <div class="m3-empty-state__icon">
        <span class="material-symbols-outlined">group</span>
      </div>
      <div class="md-title-medium" style="margin-bottom:8px;">No teams yet</div>
      <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin:0 0 20px;">
        Create a team to collaborate with your colleagues.
      </p>
      <md-filled-button @click="openCreateModal">
        <span class="material-symbols-outlined" slot="icon">add</span>
        Create Your First Team
      </md-filled-button>
    </div>

    <!-- Teams grid -->
    <div v-else class="teams-grid">
      <div
        v-for="team in teams"
        :key="team.id"
        class="m3-card m3-card--elevated team-card"
        :class="{ selected: selectedTeam?.id === team.id }"
        @click="selectTeam(team)"
      >
        <div class="team-card-body">
          <div style="display:flex;align-items:flex-start;justify-content:space-between;margin-bottom:12px;">
            <div class="team-avatar">{{ team.name.charAt(0).toUpperCase() }}</div>
            <span class="m3-badge" :class="getRoleBadgeM3Class(getMyRole(team))">
              {{ getMyRole(team) }}
            </span>
          </div>
          <div class="md-title-medium" style="margin-bottom:4px;">{{ team.name }}</div>
          <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin-bottom:6px;">@{{ team.slug }}</div>
          <p v-if="team.description" class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin-bottom:8px;line-height:1.4;">
            {{ team.description }}
          </p>
          <div style="display:flex;align-items:center;gap:6px;" class="md-body-small">
            <span class="material-symbols-outlined" style="font-size:16px;color:var(--md-sys-color-on-surface-variant);">group</span>
            <span style="color:var(--md-sys-color-on-surface-variant);">
              {{ team.member_count }} member{{ team.member_count !== 1 ? 's' : '' }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Team Detail Panel ──────────────────────────────────────────────── -->
    <div v-if="selectedTeam" class="m3-card m3-card--elevated section-card" style="margin-top:20px;">
      <div class="card-section-header">
        <div style="display:flex;align-items:center;gap:12px;">
          <div class="team-avatar-sm">{{ selectedTeam.name.charAt(0).toUpperCase() }}</div>
          <div>
            <div class="md-title-medium">{{ selectedTeam.name }}</div>
            <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">@{{ selectedTeam.slug }}</div>
          </div>
        </div>
        <div style="display:flex;gap:8px;flex-wrap:wrap;">
          <md-filled-button
            v-if="canManageMembers(selectedTeam)"
            @click="openInviteModal"
          >
            <span class="material-symbols-outlined" slot="icon">person_add</span>
            Invite Member
          </md-filled-button>
          <md-outlined-button
            v-if="isTeamOwner(selectedTeam)"
            @click="openEditModal"
          >
            <span class="material-symbols-outlined" slot="icon">edit</span>
            Edit
          </md-outlined-button>
          <md-outlined-button
            v-if="!isTeamOwner(selectedTeam)"
            @click="confirmLeave(selectedTeam)"
            style="--md-outlined-button-outline-color:var(--md-sys-color-error);--md-outlined-button-label-text-color:var(--md-sys-color-error);"
          >
            Leave Team
          </md-outlined-button>
          <md-outlined-button
            v-if="isTeamOwner(selectedTeam)"
            @click="confirmDelete(selectedTeam)"
            style="--md-outlined-button-outline-color:var(--md-sys-color-error);--md-outlined-button-label-text-color:var(--md-sys-color-error);"
          >
            Delete Team
          </md-outlined-button>
        </div>
      </div>

      <div>
        <!-- Members loading -->
        <div v-if="membersLoading" style="text-align:center;padding:32px;">
          <md-circular-progress indeterminate />
        </div>

        <!-- Members table -->
        <div v-else class="m3-table-wrapper">
          <table class="m3-table members-table">
            <thead>
              <tr>
                <th>Member</th>
                <th>Role</th>
                <th>Status</th>
                <th>Joined</th>
                <th v-if="canManageMembers(selectedTeam)" style="text-align:right;">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="members.length === 0">
                <td :colspan="canManageMembers(selectedTeam) ? 5 : 4" style="text-align:center;color:var(--md-sys-color-on-surface-variant);padding:32px;">
                  No members yet.
                </td>
              </tr>
              <tr v-for="member in members" :key="member.id">
                <td>
                  <div style="display:flex;align-items:center;gap:10px;">
                    <div class="member-avatar">{{ getMemberInitial(member) }}</div>
                    <div>
                      <div class="md-label-large">{{ getMemberName(member) }}</div>
                      <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">{{ member.email }}</div>
                    </div>
                  </div>
                </td>
                <td>
                  <span class="m3-badge" :class="getRoleBadgeM3Class(member.role)">{{ member.role }}</span>
                </td>
                <td>
                  <span class="m3-badge" :class="getStatusBadgeM3Class(member.status)">{{ member.status }}</span>
                </td>
                <td class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">
                  {{ member.joined_at ? formatDate(member.joined_at) : '—' }}
                </td>
                <td v-if="canManageMembers(selectedTeam)" style="text-align:right;">
                  <div style="display:flex;align-items:center;justify-content:flex-end;gap:8px;">
                    <md-outlined-select
                      v-if="member.role !== 'owner' && member.user_id !== authStore.profile?.id"
                      :value="member.role"
                      @change="changeMemberRole(member, ($event.target as HTMLSelectElement).value)"
                      label="Role"
                      style="min-width:110px;"
                    >
                      <md-select-option value="admin"><div slot="headline">Admin</div></md-select-option>
                      <md-select-option value="member"><div slot="headline">Member</div></md-select-option>
                    </md-outlined-select>
                    <md-icon-button
                      v-if="member.role !== 'owner' && member.user_id !== authStore.profile?.id"
                      @click="confirmRemoveMember(member)"
                      title="Remove member"
                    >
                      <span class="material-symbols-outlined" style="color:var(--md-sys-color-error);">person_remove</span>
                    </md-icon-button>
                    <span v-if="member.role === 'owner'" class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);">Owner</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ── Create Team Dialog ─────────────────────────────────────────────── -->
    <md-dialog :open="showCreateModal" @closed="closeCreateModal">
      <div slot="headline">Create Team</div>
      <div slot="content" style="min-width:480px;max-width:100%;">
        <form @submit.prevent="submitCreateTeam">
          <md-outlined-text-field
            :value="createForm.name"
            @input="createForm.name = ($event.target as HTMLInputElement).value; autoSlug()"
            label="Team Name"
            required
            style="width:100%;margin-bottom:16px;"
          />
          <div style="display:flex;align-items:center;gap:8px;margin-bottom:4px;">
            <span class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);">@</span>
            <md-outlined-text-field
              :value="createForm.slug"
              @input="createForm.slug = ($event.target as HTMLInputElement).value"
              label="Slug"
              required
              pattern="[a-z0-9-]+"
              title="Only lowercase letters, numbers, and hyphens"
              style="flex:1;"
            />
          </div>
          <div class="md-body-small" style="color:var(--md-sys-color-on-surface-variant);margin-bottom:16px;">
            Lowercase letters, numbers, and hyphens only.
          </div>
          <md-outlined-text-field
            :value="createForm.description"
            @input="createForm.description = ($event.target as HTMLInputElement).value"
            label="Description (optional)"
            type="textarea"
            rows="2"
            style="width:100%;margin-bottom:16px;"
          />
          <div v-if="createError" class="feedback-error">
            <span class="material-symbols-outlined" style="font-size:18px;">error</span>
            {{ createError }}
          </div>
        </form>
      </div>
      <div slot="actions">
        <md-text-button @click="closeCreateModal">Cancel</md-text-button>
        <md-filled-button :disabled="createLoading" @click="submitCreateTeam">
          <md-circular-progress v-if="createLoading" indeterminate style="--md-circular-progress-size:20px;margin-right:8px;" />
          Create Team
        </md-filled-button>
      </div>
    </md-dialog>

    <!-- ── Edit Team Dialog ───────────────────────────────────────────────── -->
    <md-dialog :open="showEditModal && !!selectedTeam" @closed="showEditModal = false">
      <div slot="headline">Edit Team</div>
      <div slot="content" style="min-width:480px;max-width:100%;">
        <md-outlined-text-field
          :value="editForm.name"
          @input="editForm.name = ($event.target as HTMLInputElement).value"
          label="Team Name"
          style="width:100%;margin-bottom:16px;"
        />
        <md-outlined-text-field
          :value="editForm.description"
          @input="editForm.description = ($event.target as HTMLInputElement).value"
          label="Description"
          type="textarea"
          rows="2"
          style="width:100%;"
        />
        <div v-if="editError" class="feedback-error" style="margin-top:16px;">
          <span class="material-symbols-outlined" style="font-size:18px;">error</span>
          {{ editError }}
        </div>
      </div>
      <div slot="actions">
        <md-text-button @click="showEditModal = false">Cancel</md-text-button>
        <md-filled-button :disabled="editLoading" @click="submitEditTeam">
          <md-circular-progress v-if="editLoading" indeterminate style="--md-circular-progress-size:20px;margin-right:8px;" />
          Save Changes
        </md-filled-button>
      </div>
    </md-dialog>

    <!-- ── Invite Member Dialog ───────────────────────────────────────────── -->
    <md-dialog :open="showInviteModal" @closed="closeInviteModal">
      <div slot="headline">Invite Member</div>
      <div slot="content" style="min-width:480px;max-width:100%;">
        <md-outlined-text-field
          :value="inviteForm.email"
          @input="inviteForm.email = ($event.target as HTMLInputElement).value"
          label="Email Address"
          type="email"
          required
          style="width:100%;margin-bottom:16px;"
        />
        <md-outlined-select
          :value="inviteForm.role"
          @change="inviteForm.role = ($event.target as HTMLSelectElement).value as 'admin' | 'member'"
          label="Role"
          style="width:100%;"
        >
          <md-select-option value="member"><div slot="headline">Member</div></md-select-option>
          <md-select-option value="admin"><div slot="headline">Admin</div></md-select-option>
        </md-outlined-select>
        <div v-if="inviteError" class="feedback-error" style="margin-top:16px;">
          <span class="material-symbols-outlined" style="font-size:18px;">error</span>
          {{ inviteError }}
        </div>
      </div>
      <div slot="actions">
        <md-text-button @click="closeInviteModal">Cancel</md-text-button>
        <md-filled-button :disabled="inviteLoading" @click="submitInvite">
          <md-circular-progress v-if="inviteLoading" indeterminate style="--md-circular-progress-size:20px;margin-right:8px;" />
          Send Invitation
        </md-filled-button>
      </div>
    </md-dialog>

    <!-- ── Confirm Delete Dialog ──────────────────────────────────────────── -->
    <md-dialog :open="showDeleteConfirm" @closed="showDeleteConfirm = false">
      <div slot="headline">Delete Team</div>
      <div slot="content" style="min-width:360px;max-width:100%;text-align:center;">
        <span class="material-symbols-outlined" style="font-size:48px;color:var(--md-sys-color-error);">warning</span>
        <p class="md-body-medium" style="margin:16px 0 0;color:var(--md-sys-color-on-surface-variant);">
          Are you sure you want to delete <strong>{{ teamToDelete?.name }}</strong>?
          This action cannot be undone.
        </p>
      </div>
      <div slot="actions">
        <md-text-button @click="showDeleteConfirm = false">Cancel</md-text-button>
        <md-filled-button
          :disabled="deleteLoading"
          @click="deleteTeam"
          style="--md-filled-button-container-color:var(--md-sys-color-error);"
        >
          <md-circular-progress v-if="deleteLoading" indeterminate style="--md-circular-progress-size:20px;margin-right:8px;" />
          Delete Team
        </md-filled-button>
      </div>
    </md-dialog>

    <!-- ── Confirm Leave Dialog ───────────────────────────────────────────── -->
    <md-dialog :open="showLeaveConfirm" @closed="showLeaveConfirm = false">
      <div slot="headline">Leave Team</div>
      <div slot="content" style="min-width:360px;max-width:100%;text-align:center;">
        <span class="material-symbols-outlined" style="font-size:48px;color:var(--md-sys-color-on-surface-variant);">logout</span>
        <p class="md-body-medium" style="margin:16px 0 0;color:var(--md-sys-color-on-surface-variant);">
          Are you sure you want to leave <strong>{{ teamToLeave?.name }}</strong>?
        </p>
      </div>
      <div slot="actions">
        <md-text-button @click="showLeaveConfirm = false">Cancel</md-text-button>
        <md-filled-button
          :disabled="leaveLoading"
          @click="leaveTeam"
        >
          <md-circular-progress v-if="leaveLoading" indeterminate style="--md-circular-progress-size:20px;margin-right:8px;" />
          Leave Team
        </md-filled-button>
      </div>
    </md-dialog>

    <!-- ── Confirm Remove Member Dialog ──────────────────────────────────── -->
    <md-dialog :open="showRemoveConfirm" @closed="showRemoveConfirm = false">
      <div slot="headline">Remove Member</div>
      <div slot="content" style="min-width:360px;max-width:100%;">
        <p class="md-body-medium" style="color:var(--md-sys-color-on-surface-variant);margin:0;">
          Remove <strong>{{ memberToRemove?.email }}</strong> from the team?
        </p>
      </div>
      <div slot="actions">
        <md-text-button @click="showRemoveConfirm = false">Cancel</md-text-button>
        <md-filled-button
          :disabled="removeLoading"
          @click="removeMember"
          style="--md-filled-button-container-color:var(--md-sys-color-error);"
        >
          <md-circular-progress v-if="removeLoading" indeterminate style="--md-circular-progress-size:20px;margin-right:8px;" />
          Remove
        </md-filled-button>
      </div>
    </md-dialog>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
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
    if (idx !== -1) members.value[idx] = { ...members.value[idx], role: newRole as 'admin' | 'member' };
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
    if (idx !== -1) teams.value[idx].member_count = Math.max(0, teams.value[idx].member_count - 1);
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
  padding: 24px 0;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.page-header-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

/* Empty state */
.m3-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 24px;
  text-align: center;

  &__icon {
    width: 72px;
    height: 72px;
    border-radius: 50%;
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
}

.teams-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
  margin-bottom: 20px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.team-card {
  border-radius: 12px;
  cursor: pointer;
  transition: box-shadow 0.15s, outline 0.15s;

  &:hover {
    box-shadow: 0 4px 16px rgba(99, 91, 255, 0.15);
    outline: 1px solid var(--md-sys-color-primary);
  }

  &.selected {
    outline: 2px solid var(--md-sys-color-primary);
    box-shadow: 0 0 0 4px rgba(99, 91, 255, 0.12);
  }
}

.team-card-body {
  padding: 20px;
}

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

.section-card {
  border-radius: 12px;
  overflow: hidden;
}

.card-section-header {
  padding: 14px 20px;
  border-bottom: 1px solid var(--md-sys-color-outline-variant);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.m3-table-wrapper {
  overflow-x: auto;
}

.members-table {
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

  tbody tr:last-child td {
    border-bottom: none;
  }

  tbody tr:hover td {
    background: var(--md-sys-color-surface-container-low);
  }
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

/* M3 badge variants */
.m3-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 0.72rem;
  font-weight: 600;
  white-space: nowrap;
  text-transform: capitalize;

  &--primary {
    background: rgba(99, 91, 255, 0.12);
    color: var(--md-sys-color-primary);
  }

  &--info {
    background: rgba(2, 136, 209, 0.12);
    color: #0277bd;
  }

  &--success {
    background: rgba(22, 163, 74, 0.12);
    color: #16a34a;
  }

  &--warning {
    background: rgba(245, 158, 11, 0.12);
    color: #92400e;
  }

  &--error {
    background: rgba(220, 38, 38, 0.12);
    color: var(--md-sys-color-error);
  }

  &--neutral {
    background: var(--md-sys-color-surface-container-low);
    color: var(--md-sys-color-on-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
  }
}

/* Snackbar transition */
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

.snack-enter-active,
.snack-leave-active {
  transition: all 0.25s;
}

.snack-enter-from,
.snack-leave-to {
  transform: translateY(80px);
  opacity: 0;
}
</style>
