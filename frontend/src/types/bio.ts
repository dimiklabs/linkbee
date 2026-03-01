export interface BioLinkItem {
  id: string;
  title: string;
  url: string;
  is_active: boolean;
  position: number;
  click_count: number;
}

export interface BioPage {
  id: string;
  username: string;
  title: string;
  description: string;
  avatar_url: string;
  theme: 'light' | 'dark';
  is_published: boolean;
  links: BioLinkItem[];
  created_at: string;
}

export interface UpdateBioPageRequest {
  username?: string;
  title?: string;
  description?: string;
  avatar_url?: string;
  theme?: 'light' | 'dark';
  is_published?: boolean;
}

export interface CreateBioLinkRequest {
  title: string;
  url: string;
}

export interface UpdateBioLinkRequest {
  title?: string;
  url?: string;
  is_active?: boolean;
}

export interface ReorderBioLinksRequest {
  ids: string[];
}
