export interface Webhook {
  id: string;
  url: string;
  events: string[];
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateWebhookRequest {
  url: string;
  events: string[];
}

export interface UpdateWebhookRequest {
  url?: string;
  events?: string[];
  is_active?: boolean;
}

export interface WebhookDelivery {
  id: string;
  webhook_id: string;
  user_id: string;
  event: string;
  request_body: string;
  response_code: number;
  response_body: string;
  error_message?: string;
  success: boolean;
  duration_ms: number;
  created_at: string;
}

export interface DeliveriesResponse {
  deliveries: WebhookDelivery[];
  total: number;
  page: number;
  limit: number;
}

export const WEBHOOK_EVENTS = [
  { value: 'link.created', label: 'Link Created', description: 'Fires when a new link is created' },
  { value: 'link.deleted', label: 'Link Deleted', description: 'Fires when a link is deleted' },
  { value: 'link.clicked', label: 'Link Clicked', description: 'Fires on every redirect click' },
] as const;
