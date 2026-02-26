export interface APIKey {
  id: string;
  name: string;
  key_prefix: string;
  last_used_at?: string;
  expires_at?: string;
  created_at: string;
}

export interface CreateAPIKeyRequest {
  name: string;
  expires_at?: string; // ISO-8601 / RFC3339
}

export interface CreateAPIKeyResponse extends APIKey {
  key: string; // full plaintext key — shown once only
}
