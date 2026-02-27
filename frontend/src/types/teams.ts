export interface Team {
  id: string;
  name: string;
  slug: string;
  owner_id: string;
  description?: string;
  avatar_url?: string;
  member_count: number;
  created_at: string;
}

export interface TeamMember {
  id: string;
  user_id: string;
  email: string;
  first_name?: string;
  last_name?: string;
  role: 'owner' | 'admin' | 'member';
  status: 'pending' | 'active' | 'declined';
  joined_at?: string;
  invite_email?: string;
}

export interface CreateTeamRequest {
  name: string;
  slug: string;
  description?: string;
}

export interface InviteMemberRequest {
  email: string;
  role: 'admin' | 'member';
}
