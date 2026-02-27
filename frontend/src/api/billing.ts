import apiClient from './client';
import type { SubscriptionWithPlan, UsageCounts } from '@/types/billing';
import type { ApiResponse } from '@/types/auth';

export default {
  getSubscription(): Promise<{ data: ApiResponse<SubscriptionWithPlan> }> {
    return apiClient.get('/billing/subscription');
  },

  getUsage(): Promise<{ data: ApiResponse<UsageCounts> }> {
    return apiClient.get('/billing/usage');
  },

  getCheckoutURL(plan: string): Promise<{ data: ApiResponse<{ checkout_url: string }> }> {
    return apiClient.get(`/billing/checkout/${plan}`);
  },
};
