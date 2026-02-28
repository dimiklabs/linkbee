import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // Landing page (public)
    {
      path: '/',
      name: 'landing',
      component: () => import('@/pages/LandingPage.vue'),
      meta: { guest: false }, // accessible by all
    },

    // Auth pages (guest only)
    {
      path: '/login',
      name: 'login',
      component: () => import('@/pages/auth/LoginPage.vue'),
      meta: { guest: true },
    },
    {
      path: '/signup',
      name: 'signup',
      component: () => import('@/pages/auth/SignupPage.vue'),
      meta: { guest: true },
    },
    {
      path: '/forgot-password',
      name: 'forgot-password',
      component: () => import('@/pages/auth/ForgotPasswordPage.vue'),
      meta: { guest: true },
    },
    {
      path: '/reset-password',
      name: 'reset-password',
      component: () => import('@/pages/auth/ResetPasswordPage.vue'),
      meta: { guest: true },
    },
    {
      path: '/verify-email',
      name: 'verify-email',
      component: () => import('@/pages/auth/VerifyEmailPage.vue'),
      meta: { guest: true },
    },
    {
      path: '/auth/callback',
      name: 'oauth-callback',
      component: () => import('@/pages/auth/OAuthCallbackPage.vue'),
    },

    // Dashboard (auth required)
    {
      path: '/dashboard',
      component: () => import('@/layouts/DashboardLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/dashboard/overview',
        },
        {
          path: 'overview',
          name: 'overview',
          component: () => import('@/pages/dashboard/OverviewPage.vue'),
        },
        {
          path: 'links',
          name: 'links',
          component: () => import('@/pages/dashboard/LinksPage.vue'),
        },
        {
          path: 'analytics',
          name: 'global-analytics',
          component: () => import('@/pages/dashboard/GlobalAnalyticsPage.vue'),
        },
        {
          path: 'links/:id',
          name: 'link-analytics',
          component: () => import('@/pages/dashboard/LinkAnalyticsPage.vue'),
        },
        {
          path: 'comparison',
          name: 'link-comparison',
          component: () => import('@/pages/dashboard/LinkComparisonPage.vue'),
        },
        {
          path: 'reports',
          name: 'reports',
          component: () => import('@/pages/dashboard/ReportsPage.vue'),
        },
        {
          path: 'api-keys',
          name: 'api-keys',
          component: () => import('@/pages/dashboard/ApiKeysPage.vue'),
        },
        {
          path: 'webhooks',
          name: 'webhooks',
          component: () => import('@/pages/dashboard/WebhooksPage.vue'),
        },
        {
          path: 'bio',
          name: 'bio',
          component: () => import('@/pages/dashboard/BioPage.vue'),
        },
        {
          path: 'domains',
          name: 'domains',
          component: () => import('@/pages/dashboard/DomainsPage.vue'),
        },
        {
          path: 'audit-logs',
          name: 'audit-logs',
          component: () => import('@/pages/dashboard/AuditLogsPage.vue'),
        },
        {
          path: 'billing',
          name: 'billing',
          component: () => import('@/pages/dashboard/BillingPage.vue'),
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/pages/dashboard/SettingsPage.vue'),
        },
        {
          path: 'security',
          name: 'security',
          component: () => import('@/pages/dashboard/SecurityPage.vue'),
        },
        {
          path: 'teams',
          name: 'teams',
          component: () => import('@/pages/dashboard/TeamsPage.vue'),
        },
      ],
    },

    // Admin (admin role required)
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/pages/AdminPage.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
    },

    // Password-protected link unlock page
    {
      path: '/unlock/:slug',
      name: 'unlock',
      component: () => import('@/pages/UnlockPage.vue'),
    },

    // Public bio page
    {
      path: '/bio/:username',
      name: 'public-bio',
      component: () => import('@/pages/PublicBioPage.vue'),
    },

    // Catch all
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/pages/NotFoundPage.vue'),
    },
  ],
});

let initAttempted = false;

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore();

  if (authStore.isAuthenticated && !authStore.profile && !initAttempted) {
    initAttempted = true;
    await authStore.init();
  }

  if (!authStore.isAuthenticated) {
    initAttempted = false;
  }

  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);
  const isGuestRoute = to.matched.some((record) => record.meta.guest);
  const requiresAdmin = to.matched.some((record) => record.meta.requiresAdmin);

  if (requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } });
  } else if (isGuestRoute && authStore.isAuthenticated) {
    next({ name: 'overview' });
  } else if (requiresAdmin && authStore.profile?.role !== 'admin') {
    next({ name: 'overview' });
  } else {
    next();
  }
});

export default router;
