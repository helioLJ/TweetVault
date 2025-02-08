const API_BASE_URL = 'http://localhost:8080/api';

export const api = {
  async getBookmarks(params?: {
    tag?: string;
    search?: string;
    page?: number;
    limit?: number;
  }) {
    const searchParams = new URLSearchParams();
    if (params?.tag) searchParams.append('tag', params.tag);
    if (params?.search) searchParams.append('search', params.search);
    if (params?.page) searchParams.append('page', params.page.toString());
    if (params?.limit) searchParams.append('limit', params.limit.toString());

    const url = `${API_BASE_URL}/bookmarks?${searchParams.toString()}`;
    const res = await fetch(url);
    return res.json();
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
};