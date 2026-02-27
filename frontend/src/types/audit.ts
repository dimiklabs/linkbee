export interface AuditLog {
  id: string
  user_id: string
  action: string
  resource_type: string
  resource_id: string
  resource_name: string
  ip_address: string
  user_agent: string
  created_at: string
}

export interface AuditLogsResponse {
  logs: AuditLog[]
  total: number
  page: number
  limit: number
}

// Human-readable labels for each audit action
export const ACTION_LABELS: Record<string, string> = {
  user_signup: 'Account created',
  user_login: 'Login',
  user_logout: 'Logout',
  password_changed: 'Password changed',
  account_deleted: 'Account deleted',
  link_created: 'Link created',
  link_updated: 'Link updated',
  link_deleted: 'Link deleted',
  links_imported: 'Links imported',
  domain_added: 'Domain added',
  domain_verified: 'Domain verified',
  domain_deleted: 'Domain deleted',
  api_key_created: 'API key created',
  api_key_revoked: 'API key revoked',
}

export const RESOURCE_LABELS: Record<string, string> = {
  user: 'Account',
  link: 'Link',
  domain: 'Domain',
  api_key: 'API Key',
}

export const ACTION_BADGE: Record<string, string> = {
  user_signup: 'text-bg-success',
  user_login: 'text-bg-success',
  user_logout: 'text-bg-secondary',
  password_changed: 'text-bg-warning',
  account_deleted: 'text-bg-danger',
  link_created: 'text-bg-primary',
  link_updated: 'text-bg-info',
  link_deleted: 'text-bg-danger',
  links_imported: 'text-bg-primary',
  domain_added: 'text-bg-primary',
  domain_verified: 'text-bg-success',
  domain_deleted: 'text-bg-danger',
  api_key_created: 'text-bg-primary',
  api_key_revoked: 'text-bg-danger',
}
