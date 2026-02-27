import { defineStore } from 'pinia';
import { ref } from 'vue';
import linksApi from '@/api/links';
import type { CreateLinkRequest, LinkListResponse, LinkResponse, UpdateLinkRequest } from '@/types/links';

export const useLinksStore = defineStore('links', () => {
  const links = ref<LinkResponse[]>([]);
  const total = ref(0);
  const page = ref(1);
  const limit = ref(20);
  const totalPages = ref(0);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const fetchLinks = async (p = 1, l = 20, search = '', folderID = '', starred?: boolean, healthStatus?: string, tags?: string[]) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await linksApi.list(p, l, search, folderID, starred, healthStatus, tags);
      if (response.data) {
        const data = response.data as LinkListResponse;
        links.value = data.links;
        total.value = data.total;
        page.value = data.page;
        limit.value = data.limit;
        totalPages.value = data.total_pages;
      }
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to fetch links';
      error.value = message;
    } finally {
      loading.value = false;
    }
  };

  const createLink = async (data: CreateLinkRequest): Promise<LinkResponse> => {
    loading.value = true;
    error.value = null;
    try {
      const response = await linksApi.create(data);
      if (response.data) {
        links.value.unshift(response.data);
        total.value++;
        return response.data;
      }
      throw new Error('Failed to create link');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to create link';
      error.value = message;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const updateLink = async (id: string, data: UpdateLinkRequest): Promise<LinkResponse> => {
    loading.value = true;
    error.value = null;
    try {
      const response = await linksApi.update(id, data);
      if (response.data) {
        const idx = links.value.findIndex(l => l.id === id);
        if (idx !== -1) links.value[idx] = response.data;
        return response.data;
      }
      throw new Error('Failed to update link');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to update link';
      error.value = message;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const deleteLink = async (id: string) => {
    loading.value = true;
    error.value = null;
    try {
      await linksApi.delete(id);
      links.value = links.value.filter(l => l.id !== id);
      total.value--;
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to delete link';
      error.value = message;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const toggleStar = async (id: string) => {
    try {
      const response = await linksApi.toggleStar(id);
      if (response.data) {
        const idx = links.value.findIndex(l => l.id === id);
        if (idx !== -1) links.value[idx] = response.data!;
      }
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to toggle star';
      error.value = message;
      throw err;
    }
  };

  const checkHealth = async (id: string) => {
    try {
      const response = await linksApi.checkHealth(id);
      if (response.data) {
        const idx = links.value.findIndex(l => l.id === id);
        if (idx !== -1) links.value[idx] = response.data!;
      }
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Health check failed';
      error.value = message;
      throw err;
    }
  };

  return {
    links, total, page, limit, totalPages, loading, error,
    fetchLinks, createLink, updateLink, deleteLink, toggleStar, checkHealth,
  };
});
