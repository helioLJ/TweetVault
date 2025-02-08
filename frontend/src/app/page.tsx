'use client';

import { useEffect, useState } from 'react';
import { Bookmark } from '@/types';
import { api } from '@/lib/api';
import { BookmarkCard } from '@/components/bookmarks/BookmarkCard';
import { Button } from '@/components/ui/Button';
import { Pagination } from '@/components/ui/Pagination';
import { UploadModal } from '@/components/bookmarks/UploadModal';
import { SearchAndFilter } from '@/components/bookmarks/SearchAndFilter';

const ITEMS_PER_PAGE = 12;

export default function Home() {
  const [bookmarks, setBookmarks] = useState<Bookmark[]>([]);
  const [selectedTag, setSelectedTag] = useState<string>();
  const [searchQuery, setSearchQuery] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [isLoading, setIsLoading] = useState(true);
  const [isUploadModalOpen, setIsUploadModalOpen] = useState(false);

  useEffect(() => {
    loadBookmarks();
  }, [selectedTag, searchQuery, currentPage]);

  async function loadBookmarks() {
    setIsLoading(true);
    try {
      const data = await api.getBookmarks({
        tag: selectedTag,
        search: searchQuery,
        page: currentPage,
        limit: ITEMS_PER_PAGE,
      });
      setBookmarks(data.bookmarks);
      setTotalPages(Math.ceil(data.total / ITEMS_PER_PAGE));
    } catch (error) {
      console.error('Failed to load bookmarks:', error);
    }
    setIsLoading(false);
  }

  async function handleUpdateTags(id: string, tags: string[]) {
    try {
      await api.updateBookmarkTags(id, tags);
      loadBookmarks();
    } catch (error) {
      console.error('Failed to update tags:', error);
    }
  }

  function handleSearch(query: string) {
    setSearchQuery(query);
    setCurrentPage(1);
  }

  function handleTagSelect(tag: string | undefined) {
    setSelectedTag(tag);
    setCurrentPage(1);
  }

  return (
    <main className="container mx-auto px-4 py-8">
      <div className="mb-8 flex items-center justify-between">
        <h1 className="text-3xl font-bold">TweetVault</h1>
        <Button onClick={() => setIsUploadModalOpen(true)}>
          Upload Bookmarks
        </Button>
      </div>

      <SearchAndFilter
        onSearch={handleSearch}
        onTagSelect={handleTagSelect}
        selectedTag={selectedTag}
      />

      {isLoading ? (
        <div className="flex min-h-[400px] items-center justify-center">
          <div className="text-center">
            <div className="mb-2 h-8 w-8 animate-spin rounded-full border-4 border-blue-600 border-t-transparent"></div>
            <p className="text-gray-600">Loading bookmarks...</p>
          </div>
        </div>
      ) : (
        <>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {bookmarks.map((bookmark) => (
              <BookmarkCard
                key={bookmark.id}
                bookmark={bookmark}
                onUpdateTags={handleUpdateTags}
              />
            ))}
          </div>

          {bookmarks.length === 0 && (
            <div className="mt-8 text-center text-gray-500">
              No bookmarks found
            </div>
          )}

          {bookmarks.length > 0 && (
            <Pagination
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={setCurrentPage}
            />
          )}
        </>
      )}

      <UploadModal
        isOpen={isUploadModalOpen}
        onClose={() => setIsUploadModalOpen(false)}
        onSuccess={() => {
          setCurrentPage(1);
          loadBookmarks();
        }}
      />
    </main>
  );
}