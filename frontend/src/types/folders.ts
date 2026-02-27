export interface FolderResponse {
  id: string;
  name: string;
  color: string;
  click_count: number;
  created_at: string;
  updated_at: string;
}

export interface CreateFolderRequest {
  name: string;
  color?: string;
}

export interface UpdateFolderRequest {
  name?: string;
  color?: string;
}
