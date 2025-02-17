const API_BASE_URL = (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080') + '/api';

interface Statistics {
  total_bookmarks: number;
  active_bookmarks: number;
  archived_bookmarks: number;
  total_tags: number;
  top_tags: Array<{
    name: string;
    count: number;
    completed_count: number;
  }>;
}

const handleApiError = (error: any, endpoint: string) => {
  console.error(`API Error (${endpoint}):`, error);
  
  const errorDetails = {
    endpoint,
    message: error.message,
    timestamp: new Date().toISOString(),
    stack: error.stack,
  };

  throw error;
};

export const api = {
  async getBookmarks(params?: {
    tag?: string;
    search?: string;
    page?: number;
    limit?: number;
    archived?: boolean;
  }) {
    try {
      const searchParams = new URLSearchParams();
      if (params?.tag) searchParams.append('tag', params.tag);
      if (params?.search) searchParams.append('search', params.search);
      if (params?.page) searchParams.append('page', params.page.toString());
      if (params?.limit) searchParams.append('limit', params.limit.toString());
      searchParams.append('archived', (params?.archived ?? false).toString());

      const url = `${API_BASE_URL}/bookmarks?${searchParams.toString()}`;
      const res = await fetch(url);

      if (!res.ok) {
        throw new Error(`HTTP error! status: ${res.status}`);
      }
      return res.json();
    } catch (error) {
      return handleApiError(error, 'getBookmarks');
    }
  },

  async getTags() {
    const res = await fetch(`${API_BASE_URL}/tags`);
    return res.json();
  },

  async updateBookmarkTags(id: string, tags: string[]) {
    const res = await fetch(`${API_BASE_URL}/bookmarks/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ tags }),
    });
    return res.json();
  },

  async uploadBookmarks(jsonFile: File, zipFile: File) {
    const formData = new FormData();
    formData.append('jsonFile', jsonFile);
    formData.append('zipFile', zipFile);

    const res = await fetch(`${API_BASE_URL}/upload`, {
      method: 'POST',
      body: formData,
    });
    return res.json();
  },

  async updateTag(id: string, name: string) {
    const res = await fetch(`${API_BASE_URL}/tags/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name }),
    });
    return res.json();
  },

  async deleteTag(id: string) {
    const res = await fetch(`${API_BASE_URL}/tags/${id}`, {
      method: 'DELETE',
    });
    return res.json();
  },

  async getTagBookmarkCount(id: string) {
    const res = await fetch(`${API_BASE_URL}/tags/${id}/count`);
    return res.json();
  },

  async deleteBookmark(id: string) {
    const res = await fetch(`${API_BASE_URL}/bookmarks/${id}`, {
      method: 'DELETE',
    });
    return res.json();
  },

  async getStatistics(): Promise<Statistics> {
    const res = await fetch(`${API_BASE_URL}/statistics`);
    return res.json();
  },

  async toggleArchiveBookmark(id: string) {
    const res = await fetch(`${API_BASE_URL}/bookmarks/${id}/toggle-archive`, {
      method: 'POST',
    });
    return res.json();
  },

  async createTag(name: string) {
    const res = await fetch(`${API_BASE_URL}/tags`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name }),
    });
    return res.json();
  },

  async toggleTagCompletion(bookmarkId: string, tagName: string) {
    const res = await fetch(`${API_BASE_URL}/bookmarks/${bookmarkId}/tags/${tagName}/toggle-completion`, {
      method: 'POST',
    });
    const data = await res.json();
    return data;
  },

  async getBookmark(id: string) {
    const res = await fetch(`${API_BASE_URL}/bookmarks/${id}`);
    return res.json();
  },
};