import apiClient from './client';
import type { SubscriptionWithPlan } from '@/types/billing';
import type { ApiResponse } from '@/types/auth';

export default {
  getSubscription(): Promise<{ data: ApiResponse<SubscriptionWithPlan> }> {
    return apiClient.get('/billing/subscription');
  },

  getCheckoutURL(plan: string): Promise<{ data: ApiResponse<{ checkout_url: string }> }> {
    return apiClient.get(`/billing/checkout/${plan}`);
  },
};
