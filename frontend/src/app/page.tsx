'use client';

import { useEffect, useState, useRef } from 'react';
import { Bookmark } from '@/types';
import { api } from '@/lib/api';
import { BookmarkCard } from '@/components/bookmarks/BookmarkCard';
import { Button } from '@/components/ui/Button';
import { Pagination } from '@/components/ui/Pagination';
import { UploadModal } from '@/components/bookmarks/UploadModal';
import { SearchAndFilter } from '@/components/bookmarks/SearchAndFilter';
import { SelectionToolbar } from '@/components/bookmarks/SelectionToolbar';
import { Statistics } from '@/components/bookmarks/Statistics';
import { StatisticsRef } from '@/components/bookmarks/Statistics';
import { ThemeToggle } from '@/components/theme/ThemeToggle';

export default function Home() {
  const [bookmarks, setBookmarks] = useState<Bookmark[]>([]);
  const [selectedTag, setSelectedTag] = useState<string>();
  const [searchQuery, setSearchQuery] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(() => {
    if (typeof window !== 'undefined') {
      const savedSize = localStorage.getItem('bookmark-page-size');
      return savedSize ? Number(savedSize) : 20;
    }
    return 20;
  });
  const [totalPages, setTotalPages] = useState(1);
  const [isLoading, setIsLoading] = useState(true);
  const [isUploadModalOpen, setIsUploadModalOpen] = useState(false);
  const [selectedBookmarks, setSelectedBookmarks] = useState<Set<string>>(new Set());
  const [isSelectMode, setIsSelectMode] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [showArchived, setShowArchived] = useState(false);
  const reloadTags = useRef<() => void>(() => {});
  const statisticsRef = useRef<StatisticsRef>(null);

  useEffect(() => {
    loadBookmarks();
  }, [currentPage, selectedTag, searchQuery, showArchived]);

  useEffect(() => {
    console.log('Current archived state:', showArchived);
  }, [showArchived]);

  async function loadBookmarks() {
    setIsLoading(true);
    try {
      const data = await api.getBookmarks({
        tag: selectedTag,
        search: searchQuery,
        page: currentPage,
        limit: pageSize,
        archived: showArchived,
      });

      setBookmarks(data.bookmarks || []);
      setTotalPages(Math.ceil(data.total / pageSize));
    } catch (error) {
      console.error('Failed to load bookmarks:', error);
      setBookmarks([]);
    }
    setIsLoading(false);
  }

  async function handleUpdateTags(id: string, tags: string[]) {
    try {
      await api.updateBookmarkTags(id, tags);
      setBookmarks(currentBookmarks => 
        currentBookmarks.map(bookmark => 
          bookmark.id === id 
            ? { ...bookmark, tags: tags.map((name, index) => ({ 
                id: -index - 1,
                name,
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString()
              })) }
            : bookmark
        )
      );
      reloadTags.current();
      statisticsRef.current?.refresh();
    } catch (error) {
      console.error('Failed to update tags:', error);
    }
  }

  async function handleDeleteTag(tagName: string) {
    try {
      // Update bookmarks state to remove the deleted tag
      setBookmarks(currentBookmarks => 
        currentBookmarks.map(bookmark => ({
          ...bookmark,
          tags: bookmark.tags.filter(tag => tag.name !== tagName)
        }))
      );
      
      // Refresh statistics if needed
      statisticsRef.current?.refresh();
    } catch (error) {
      console.error('Failed to handle tag deletion:', error);
    }
  }

  function handleSearch(query: string) {
    setSearchQuery(query);
  }

  function handleTagSelect(tag: string | undefined) {
    setSelectedTag(tag);
    setCurrentPage(1);
  }

  async function handleDeleteBookmark(id: string) {
    try {
      await api.deleteBookmark(id);
      setBookmarks(currentBookmarks => 
        currentBookmarks.filter(bookmark => bookmark.id !== id)
      );
    } catch (error) {
      console.error('Failed to delete bookmark:', error);
    }
  }

  function toggleSelectMode() {
    setIsSelectMode(!isSelectMode);
    setSelectedBookmarks(new Set());
  }

  function handleToggleSelect(id: string) {
    setSelectedBookmarks(prev => {
      const newSet = new Set(prev);
      if (newSet.has(id)) {
        newSet.delete(id);
      } else {
        newSet.add(id);
      }
      return newSet;
    });
  }

  async function handleDeleteSelected() {
    try {
      const deletePromises = Array.from(selectedBookmarks).map(id => 
        api.deleteBookmark(id)
      );
      await Promise.all(deletePromises);
      
      setBookmarks(currentBookmarks => 
        currentBookmarks.filter(bookmark => !selectedBookmarks.has(bookmark.id))
      );
      setSelectedBookmarks(new Set());
      setIsSelectMode(false);
      setIsDeleteModalOpen(false);
    } catch (error) {
      console.error('Failed to delete bookmarks:', error);
    }
  }

  function handlePageSizeChange(newSize: number) {
    setPageSize(newSize);
    setCurrentPage(1);
  }

  const handleArchive = async (id: string) => {
    try {
      await api.toggleArchiveBookmark(id);
      loadBookmarks();
    } catch (error) {
      console.error('Failed to archive bookmark:', error);
    }
  };

  const filteredBookmarks = bookmarks
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()) ?? [];

  return (
    <main className="px-4 py-8">
      <div className="mb-8 flex items-center justify-between max-w-[2000px] mx-auto">
        <h1 className="text-3xl font-bold">TweetVault</h1>
        <div className="flex items-center gap-2">
          <ThemeToggle />
          <div className="flex items-center gap-2">
            <button
              onClick={() => setShowArchived(!showArchived)}
              className={`px-4 py-2 rounded-lg ${
                showArchived 
                  ? 'bg-purple-600 text-white' 
                  : 'bg-gray-200 text-gray-700'
              }`}
            >
              {showArchived ? 'Show Active' : 'Show Archived'}
            </button>
          </div>
          <Button 
            variant="outline"
            onClick={toggleSelectMode}
          >
            {isSelectMode ? 'Cancel Selection' : 'Select Bookmarks'}
          </Button>
          <Button onClick={() => setIsUploadModalOpen(true)}>
            Upload Bookmarks
          </Button>
        </div>
      </div>

      <div className="max-w-[2000px] mx-auto mb-8">
        <Statistics ref={statisticsRef} />
      </div>

      <div className="max-w-[2000px] mx-auto">
        <SearchAndFilter
          onSearch={handleSearch}
          onTagSelect={handleTagSelect}
          selectedTag={selectedTag}
          onDeleteTag={handleDeleteTag}
          ref={(ref: { loadTags: () => void } | null) => {
            if (ref) {
              reloadTags.current = ref.loadTags;
            }
          }}
        />
      </div>

      {isLoading ? (
        <div className="flex min-h-[400px] items-center justify-center">
          <div className="text-center">
            <div className="mb-2 h-8 w-8 animate-spin rounded-full border-4 border-blue-600 border-t-transparent"></div>
            <p className="text-gray-600">Loading bookmarks...</p>
          </div>
        </div>
      ) : (
        <>
          <div className="columns-2 lg:columns-3 xl:columns-4 gap-4 max-w-[2000px] mx-auto [column-fill:balance]">
            {filteredBookmarks.map((bookmark) => (
              <div key={bookmark.id} className="break-inside-avoid mb-4">
                <BookmarkCard
                  bookmark={bookmark}
                  onUpdateTags={handleUpdateTags}
                  onDelete={handleDeleteBookmark}
                  onArchive={handleArchive}
                  isArchived={bookmark.archived}
                  isSelectable={isSelectMode}
                  isSelected={selectedBookmarks.has(bookmark.id)}
                  onToggleSelect={handleToggleSelect}
                />
              </div>
            ))}
          </div>

          {filteredBookmarks.length === 0 && (
            <div className="mt-8 text-center text-gray-500">
              No bookmarks found
            </div>
          )}

          {filteredBookmarks.length > 0 && (
            <Pagination
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={setCurrentPage}
              pageSize={pageSize}
              onPageSizeChange={handlePageSizeChange}
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

      {isSelectMode && selectedBookmarks.size > 0 && (
        <SelectionToolbar
          selectedCount={selectedBookmarks.size}
          onClearSelection={() => setSelectedBookmarks(new Set())}
          onDeleteSelected={() => setIsDeleteModalOpen(true)}
        />
      )}

      {isDeleteModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[100]">
          <div className="bg-white p-6 rounded-lg max-w-md w-full mx-4">
            <h3 className="text-lg font-semibold mb-2">Delete Bookmarks</h3>
            <p className="mb-4">
              Are you sure you want to delete {selectedBookmarks.size} selected bookmark{selectedBookmarks.size !== 1 ? 's' : ''}?
            </p>
            <div className="flex justify-end gap-2">
              <button
                className="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded"
                onClick={() => setIsDeleteModalOpen(false)}
              >
                Cancel
              </button>
              <button
                className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
                onClick={handleDeleteSelected}
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </main>
  );
}