<template>
  <div class="teams-page">

    <!-- Toast notifications -->
    <div class="toast-container position-fixed top-0 end-0 p-3" style="z-index: 1100;">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="toast show align-items-center border-0"
        :class="toast.type === 'success' ? 'text-bg-success' : 'text-bg-danger'"
        role="alert"
      >
        <div class="d-flex">
          <div class="toast-body">{{ toast.message }}</div>
          <button type="button" class="btn-close btn-close-white me-2 m-auto" @click="removeToast(toast.id)"></button>
        </div>
      </div>
    </div>

    <!-- Page header -->
    <div class="d-flex align-items-center justify-content-between mb-4">
      <div>
        <h4 class="mb-1 fw-semibold">Teams</h4>
        <p class="text-muted mb-0 small">Manage your teams and collaborate with others.</p>
      </div>
      <button class="btn btn-primary d-flex align-items-center gap-2" @click="openCreateModal">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
          <path d="M8 2a.5.5 0 0 1 .5.5v5h5a.5.5 0 0 1 0 1h-5v5a.5.5 0 0 1-1 0v-5h-5a.5.5 0 0 1 0-1h5v-5A.5.5 0 0 1 8 2z"/>
        </svg>
        Create Team
      </button>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else-if="teams.length === 0" class="text-center py-5">
      <div class="mb-3" style="font-size: 3rem;">👥</div>
      <h5 class="fw-semibold mb-1">No teams yet</h5>
      <p class="text-muted mb-3">Create a team to collaborate with your colleagues.</p>
      <button class="btn btn-primary" @click="openCreateModal">Create Your First Team</button>
    </div>

    <!-- Teams grid -->
    <div v-else class="row g-3 mb-4">
      <div v-for="team in teams" :key="team.id" class="col-12 col-md-6 col-xl-4">
        <div
          class="card h-100 team-card"
          :class="{ selected: selectedTeam?.id === team.id }"
          @click="selectTeam(team)"
          style="cursor: pointer;"
        >
          <div class="card-body">
            <div class="d-flex align-items-start justify-content-between mb-2">
              <div class="team-avatar me-3">
                {{ team.name.charAt(0).toUpperCase() }}
              </div>
              <span class="badge" :class="getRoleBadgeClass(getMyRole(team))">
                {{ getMyRole(team) }}
              </span>
            </div>
            <h6 class="fw-semibold mb-1">{{ team.name }}</h6>
            <p class="text-muted small mb-2">@{{ team.slug }}</p>
            <p v-if="team.description" class="text-muted small mb-2" style="line-height: 1.4;">
              {{ team.description }}
            </p>
            <div class="d-flex align-items-center gap-1 text-muted small">
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" fill="currentColor" viewBox="0 0 16 16">
                <path d="M15 14s1 0 1-1-1-4-5-4-5 3-5 4 1 1 1 1zm-7.978-1A.261.261 0 0 1 7 12.996c.001-.264.167-1.03.76-1.72C8.312 10.629 9.282 10 11 10c1.717 0 2.687.63 3.24 1.276.593.69.758 1.457.76 1.72l-.008.002A.274.274 0 0 1 15 13H7zM11 7a2 2 0 1 0 0-4 2 2 0 0 0 0 4zm-3.5 0a3.5 3.5 0 1 1 7 0 3.5 3.5 0 0 1-7 0zM5 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0m0-7.5a2 2 0 1 0 0 4 2 2 0 0 0 0-4M1 6.5a4 4 0 1 1 8 0 4 4 0 0 1-8 0"/>
              </svg>
              {{ team.member_count }} member{{ team.member_count !== 1 ? 's' : '' }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Team detail panel -->
    <div v-if="selectedTeam" class="card mb-4">
      <div class="card-header d-flex align-items-center justify-content-between">
        <div class="d-flex align-items-center gap-2">
          <div class="team-avatar-sm">{{ selectedTeam.name.charAt(0).toUpperCase() }}</div>
          <div>
            <h6 class="mb-0 fw-semibold">{{ selectedTeam.name }}</h6>
            <span class="text-muted small">@{{ selectedTeam.slug }}</span>
          </div>
        </div>
        <div class="d-flex gap-2">
          <button
            v-if="canManageMembers(selectedTeam)"
            class="btn btn-sm btn-outline-primary d-flex align-items-center gap-1"
            @click="openInviteModal"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
              <path d="M6 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6zm2-3a2 2 0 1 1-4 0 2 2 0 0 1 4 0zm4 8c0 1-1 1-1 1H1s-1 0-1-1 1-4 6-4 6 3 6 4zm-1-.004c-.001-.246-.154-.986-.832-1.664C9.516 10.68 8.289 10 6 10c-2.29 0-3.516.68-4.168 1.332-.678.678-.83 1.418-.832 1.664h10z"/>
              <path fill-rule="evenodd" d="M13.5 5a.5.5 0 0 1 .5.5V7h1.5a.5.5 0 0 1 0 1H14v1.5a.5.5 0 0 1-1 0V8h-1.5a.5.5 0 0 1 0-1H13V5.5a.5.5 0 0 1 .5-.5z"/>
            </svg>
            Invite Member
          </button>
          <button
            v-if="isTeamOwner(selectedTeam)"
            class="btn btn-sm btn-outline-secondary"
            @click="openEditModal"
          >
            Edit
          </button>
          <button
            v-if="!isTeamOwner(selectedTeam)"
            class="btn btn-sm btn-outline-warning"
            @click="confirmLeave(selectedTeam)"
          >
            Leave Team
          </button>
          <button
            v-if="isTeamOwner(selectedTeam)"
            class="btn btn-sm btn-outline-danger"
            @click="confirmDelete(selectedTeam)"
          >
            Delete Team
          </button>
        </div>
      </div>

      <div class="card-body p-0">
        <!-- Members loading -->
        <div v-if="membersLoading" class="text-center py-4">
          <div class="spinner-border spinner-border-sm text-primary"></div>
        </div>

        <!-- Members table -->
        <div v-else class="table-responsive">
          <table class="table table-hover mb-0">
            <thead class="table-light">
              <tr>
                <th class="ps-4">Member</th>
                <th>Role</th>
                <th>Status</th>
                <th>Joined</th>
                <th v-if="canManageMembers(selectedTeam)" class="text-end pe-4">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="members.length === 0">
                <td :colspan="canManageMembers(selectedTeam) ? 5 : 4" class="text-center text-muted py-4">
                  No members yet.
                </td>
              </tr>
              <tr v-for="member in members" :key="member.id">
                <td class="ps-4">
                  <div class="d-flex align-items-center gap-2">
                    <div class="member-avatar">
                      {{ getMemberInitial(member) }}
                    </div>
                    <div>
                      <div class="fw-medium small">
                        {{ getMemberName(member) }}
                      </div>
                      <div class="text-muted" style="font-size: 0.75rem;">{{ member.email }}</div>
                    </div>
                  </div>
                </td>
                <td class="align-middle">
                  <span class="badge" :class="getRoleBadgeClass(member.role)">
                    {{ member.role }}
                  </span>
                </td>
                <td class="align-middle">
                  <span class="badge" :class="getStatusBadgeClass(member.status)">
                    {{ member.status }}
                  </span>
                </td>
                <td class="align-middle text-muted small">
                  {{ member.joined_at ? formatDate(member.joined_at) : '—' }}
                </td>
                <td v-if="canManageMembers(selectedTeam)" class="align-middle text-end pe-4">
                  <div class="d-flex align-items-center justify-content-end gap-2">
                    <!-- Change role (not for owner) -->
                    <select
                      v-if="member.role !== 'owner' && member.user_id !== authStore.profile?.id"
                      class="form-select form-select-sm"
                      style="width: auto;"
                      :value="member.role"
                      @change="changeMemberRole(member, ($event.target as HTMLSelectElement).value)"
                    >
                      <option value="admin">Admin</option>
                      <option value="member">Member</option>
                    </select>
                    <!-- Remove (not for owner or self) -->
                    <button
                      v-if="member.role !== 'owner' && member.user_id !== authStore.profile?.id"
                      class="btn btn-sm btn-outline-danger"
                      @click="confirmRemoveMember(member)"
                    >
                      Remove
                    </button>
                    <span v-if="member.role === 'owner'" class="text-muted small">Owner</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ── Create Team Modal ──────────────────────────────────────────────── -->
    <div v-if="showCreateModal" class="modal-backdrop-custom" @click.self="closeCreateModal">
      <div class="modal-dialog-custom card shadow-lg">
        <div class="card-header d-flex align-items-center justify-content-between">
          <h6 class="mb-0 fw-semibold">Create Team</h6>
          <button class="btn-close" @click="closeCreateModal"></button>
        </div>
        <div class="card-body">
          <form @submit.prevent="submitCreateTeam">
            <div class="mb-3">
              <label class="form-label fw-medium">Team Name <span class="text-danger">*</span></label>
              <input
                v-model="createForm.name"
                type="text"
                class="form-control"
                placeholder="e.g. Engineering"
                required
                @input="autoSlug"
              />
            </div>
            <div class="mb-3">
              <label class="form-label fw-medium">Slug <span class="text-danger">*</span></label>
              <div class="input-group">
                <span class="input-group-text text-muted">@</span>
                <input
                  v-model="createForm.slug"
                  type="text"
                  class="form-control"
                  placeholder="engineering"
                  required
                  pattern="[a-z0-9-]+"
                  title="Only lowercase letters, numbers, and hyphens"
                />
              </div>
              <div class="form-text">Lowercase letters, numbers, and hyphens only.</div>
            </div>
            <div class="mb-3">
              <label class="form-label fw-medium">Description</label>
              <textarea
                v-model="createForm.description"
                class="form-control"
                rows="2"
                placeholder="Optional description..."
              ></textarea>
            </div>
            <div v-if="createError" class="alert alert-danger py-2 small">{{ createError }}</div>
            <div class="d-flex justify-content-end gap-2 mt-3">
              <button type="button" class="btn btn-outline-secondary" @click="closeCreateModal">Cancel</button>
              <button type="submit" class="btn btn-primary" :disabled="createLoading">
                <span v-if="createLoading" class="spinner-border spinner-border-sm me-2"></span>
                Create Team
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- ── Edit Team Modal ────────────────────────────────────────────────── -->
    <div v-if="showEditModal && selectedTeam" class="modal-backdrop-custom" @click.self="showEditModal = false">
      <div class="modal-dialog-custom card shadow-lg">
        <div class="card-header d-flex align-items-center justify-content-between">
          <h6 class="mb-0 fw-semibold">Edit Team</h6>
          <button class="btn-close" @click="showEditModal = false"></button>
        </div>
        <div class="card-body">
          <form @submit.prevent="submitEditTeam">
            <div class="mb-3">
              <label class="form-label fw-medium">Team Name</label>
              <input v-model="editForm.name" type="text" class="form-control" />
            </div>
            <div class="mb-3">
              <label class="form-label fw-medium">Description</label>
              <textarea v-model="editForm.description" class="form-control" rows="2"></textarea>
            </div>
            <div v-if="editError" class="alert alert-danger py-2 small">{{ editError }}</div>
            <div class="d-flex justify-content-end gap-2 mt-3">
              <button type="button" class="btn btn-outline-secondary" @click="showEditModal = false">Cancel</button>
              <button type="submit" class="btn btn-primary" :disabled="editLoading">
                <span v-if="editLoading" class="spinner-border spinner-border-sm me-2"></span>
                Save Changes
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- ── Invite Member Modal ────────────────────────────────────────────── -->
    <div v-if="showInviteModal" class="modal-backdrop-custom" @click.self="closeInviteModal">
      <div class="modal-dialog-custom card shadow-lg">
        <div class="card-header d-flex align-items-center justify-content-between">
          <h6 class="mb-0 fw-semibold">Invite Member</h6>
          <button class="btn-close" @click="closeInviteModal"></button>
        </div>
        <div class="card-body">
          <form @submit.prevent="submitInvite">
            <div class="mb-3">
              <label class="form-label fw-medium">Email Address <span class="text-danger">*</span></label>
              <input
                v-model="inviteForm.email"
                type="email"
                class="form-control"
                placeholder="colleague@example.com"
                required
              />
            </div>
            <div class="mb-3">
              <label class="form-label fw-medium">Role <span class="text-danger">*</span></label>
              <select v-model="inviteForm.role" class="form-select" required>
                <option value="member">Member</option>
                <option value="admin">Admin</option>
              </select>
            </div>
            <div v-if="inviteError" class="alert alert-danger py-2 small">{{ inviteError }}</div>
            <div class="d-flex justify-content-end gap-2 mt-3">
              <button type="button" class="btn btn-outline-secondary" @click="closeInviteModal">Cancel</button>
              <button type="submit" class="btn btn-primary" :disabled="inviteLoading">
                <span v-if="inviteLoading" class="spinner-border spinner-border-sm me-2"></span>
                Send Invitation
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- ── Confirm Delete Modal ───────────────────────────────────────────── -->
    <div v-if="showDeleteConfirm" class="modal-backdrop-custom" @click.self="showDeleteConfirm = false">
      <div class="modal-dialog-custom modal-dialog-sm card shadow-lg">
        <div class="card-body text-center py-4">
          <div class="mb-3" style="font-size: 2.5rem;">⚠️</div>
          <h6 class="fw-semibold mb-2">Delete Team</h6>
          <p class="text-muted small mb-4">
            Are you sure you want to delete <strong>{{ teamToDelete?.name }}</strong>?
            This action cannot be undone.
          </p>
          <div class="d-flex justify-content-center gap-2">
            <button class="btn btn-outline-secondary" @click="showDeleteConfirm = false">Cancel</button>
            <button class="btn btn-danger" :disabled="deleteLoading" @click="deleteTeam">
              <span v-if="deleteLoading" class="spinner-border spinner-border-sm me-2"></span>
              Delete Team
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Confirm Leave Modal ────────────────────────────────────────────── -->
    <div v-if="showLeaveConfirm" class="modal-backdrop-custom" @click.self="showLeaveConfirm = false">
      <div class="modal-dialog-custom modal-dialog-sm card shadow-lg">
        <div class="card-body text-center py-4">
          <div class="mb-3" style="font-size: 2.5rem;">🚪</div>
          <h6 class="fw-semibold mb-2">Leave Team</h6>
          <p class="text-muted small mb-4">
            Are you sure you want to leave <strong>{{ teamToLeave?.name }}</strong>?
          </p>
          <div class="d-flex justify-content-center gap-2">
            <button class="btn btn-outline-secondary" @click="showLeaveConfirm = false">Cancel</button>
            <button class="btn btn-warning" :disabled="leaveLoading" @click="leaveTeam">
              <span v-if="leaveLoading" class="spinner-border spinner-border-sm me-2"></span>
              Leave Team
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Confirm Remove Member Modal ────────────────────────────────────── -->
    <div v-if="showRemoveConfirm" class="modal-backdrop-custom" @click.self="showRemoveConfirm = false">
      <div class="modal-dialog-custom modal-dialog-sm card shadow-lg">
        <div class="card-body text-center py-4">
          <h6 class="fw-semibold mb-2">Remove Member</h6>
          <p class="text-muted small mb-4">
            Remove <strong>{{ memberToRemove?.email }}</strong> from the team?
          </p>
          <div class="d-flex justify-content-center gap-2">
            <button class="btn btn-outline-secondary" @click="showRemoveConfirm = false">Cancel</button>
            <button class="btn btn-danger" :disabled="removeLoading" @click="removeMember">
              <span v-if="removeLoading" class="spinner-border spinner-border-sm me-2"></span>
              Remove
            </button>
          </div>
        </div>
      </div>
    </div>

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

function getStatusBadgeClass(status: string): string {
  switch (status) {
    case 'active': return 'text-bg-success';
    case 'pending': return 'text-bg-warning text-dark';
    case 'declined': return 'text-bg-danger';
    default: return 'text-bg-secondary';
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
$primary: #635bff;

.teams-page {
  max-width: 1100px;
}

.team-card {
  border: 1.5px solid #e3e8ee;
  transition: border-color 0.15s, box-shadow 0.15s;

  &:hover {
    border-color: $primary;
    box-shadow: 0 2px 8px rgba(99, 91, 255, 0.1);
  }

  &.selected {
    border-color: $primary;
    box-shadow: 0 0 0 3px rgba(99, 91, 255, 0.15);
  }
}

.team-avatar {
  width: 42px;
  height: 42px;
  background: $primary;
  color: #fff;
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
  background: $primary;
  color: #fff;
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
  background: #e3e8ee;
  color: #697386;
  font-weight: 600;
  font-size: 0.8rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.modal-backdrop-custom {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  z-index: 1050;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.modal-dialog-custom {
  width: 100%;
  max-width: 480px;
  border-radius: 12px;
  overflow: hidden;
}

.modal-dialog-sm {
  max-width: 380px;
}
</style>
