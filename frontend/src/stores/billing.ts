import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import billingApi from '@/api/billing';

export const useBillingStore = defineStore('billing', () => {
  const planId = ref<string | null>(null);
  const loaded = ref(false);

  const isPaidPlan = computed(() => planId.value === 'pro' || planId.value === 'growth');

  async function fetchPlan() {
    if (loaded.value) return;
    try {
      const res = await billingApi.getSubscription();
      planId.value = res.data?.data?.plan?.id ?? 'free';
    } catch {
      planId.value = 'free';
    } finally {
      loaded.value = true;
    }
  }

  function reset() {
    planId.value = null;
    loaded.value = false;
  }

  return { planId, loaded, isPaidPlan, fetchPlan, reset };
});
