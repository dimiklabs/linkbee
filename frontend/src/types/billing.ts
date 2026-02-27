export interface UsageCounts {
  links: number;
  api_keys: number;
  webhooks: number;
}

export interface PlanInfo {
  id: string;
  name: string;
  max_links: number;       // -1 = unlimited
  max_api_keys: number;
  max_webhooks: number;
  has_webhooks: boolean;
}

export interface Subscription {
  id: string;
  user_id: string;
  plan_id: string;
  status: string;
  current_period_end?: string;
  cancelled_at?: string;
  created_at: string;
}

export interface SubscriptionWithPlan {
  subscription: Subscription;
  plan: PlanInfo;
}

export type PlanID = 'free' | 'pro' | 'business';

export const PLAN_LABELS: Record<string, string> = {
  free: 'Free',
  pro: 'Pro',
  business: 'Business',
};

export const SUB_STATUS_LABELS: Record<string, string> = {
  active: 'Active',
  cancelled: 'Cancelled',
  expired: 'Expired',
  past_due: 'Past Due',
  paused: 'Paused',
  on_trial: 'Trial',
};
