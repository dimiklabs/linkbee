export interface DomainResponse {
  id: string
  user_id: string
  domain: string
  status: 'pending' | 'verified' | 'failed'
  verify_token: string
  created_at: string
  updated_at: string
}
