import { useState, useEffect } from 'react';
import { Bookmark, Tag } from '@/types';
import Image from 'next/image';
import { FormattedDate } from './FormattedDate';
import { MediaModal } from './MediaModal';
import { linkify } from '@/lib/utils';
import { api } from '@/lib/api';

interface BookmarkCardProps {
  bookmark: Bookmark;
  onUpdateTags: (id: string, tags: string[]) => void;
  onDelete: (id: string) => void;
  isSelectable: boolean;
  isSelected: boolean;
  onToggleSelect: (id: string) => void;
  onArchive: (id: string) => void;
  isArchived: boolean;
}

export function BookmarkCard({ 
  bookmark, 
  onUpdateTags, 
  onDelete,
  isSelectable,
  isSelected,
  onToggleSelect,
  onArchive,
  isArchived
}: BookmarkCardProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [newTag, setNewTag] = useState('');
  const [selectedMedia, setSelectedMedia] = useState<null | {
    url: string;
    type: 'image' | 'video';
  }>(null);
  const [isExpanded, setIsExpanded] = useState(false);
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [menuPosition, setMenuPosition] = useState({ top: 0, left: 0 });
  const [suggestedTags, setSuggestedTags] = useState<Tag[]>([]);
  const [allTags, setAllTags] = useState<Tag[]>([]);

  useEffect(() => {
    loadTags();
  }, []);

  async function loadTags() {
    try {
      const tags = await api.getTags();
      setAllTags(tags);
    } catch (error) {
      console.error('Failed to load tags:', error);
    }
  }

  const updateSuggestedTags = (input: string) => {
    if (!input.trim()) {
      setSuggestedTags([]);
      return;
    }
    
    const filtered = allTags.filter(tag => 
      tag.name.toLowerCase().includes(input.toLowerCase()) &&
      !bookmark.tags.some(t => t.name === tag.name)
    );
    setSuggestedTags(filtered);
  };

  function handleAddTag(e: React.KeyboardEvent) {
    if (e.key === 'Enter' && newTag.trim()) {
      e.preventDefault();
      const updatedTags = [...bookmark.tags.map(t => t.name), newTag.trim()];
      onUpdateTags(bookmark.id, updatedTags);
      setNewTag('');
    }
  }

  function handleRemoveTag(tagToRemove: string) {
    const updatedTags = bookmark.tags
      .map(t => t.name)
      .filter(t => t !== tagToRemove);
    onUpdateTags(bookmark.id, updatedTags);
  }

  // Function to render text with clickable links
  function renderTextWithLinks(text: string) {
    return text.split(/\s+/).map((word, index) => {
      if (word.match(/^(https?:\/\/[^\s]+)/)) {
        return (
          <a
            key={index}
            href={word}
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 hover:underline break-all"
          >
            {word}
          </a>
        );
      }
      return <span key={index}>{word} </span>;
    });
  }

  async function handleDelete() {
    try {
      await api.deleteBookmark(bookmark.id);
      // You'll need to implement a callback to refresh the bookmark list
      // onDelete(bookmark.id);
    } catch (error) {
      console.error('Failed to delete bookmark:', error);
    }
  }

  const handleDotsClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    const button = e.currentTarget as HTMLElement;
    const rect = button.getBoundingClientRect();
    setMenuPosition({
      top: rect.bottom + window.scrollY,
      left: rect.left + window.scrollX
    });
    setIsMenuOpen(!isMenuOpen);
  };

  // Add this helper function to highlight matching text
  const highlightMatch = (text: string, query: string) => {
    if (!query) return text;
    const parts = text.split(new RegExp(`(${query})`, 'gi'));
    return parts.map((part, i) => 
      part.toLowerCase() === query.toLowerCase() 
        ? <span key={i} className="bg-yellow-200">{part}</span>
        : part
    );
  };

  return (
    <div className="rounded-lg border bg-white p-4 shadow-sm transition-shadow hover:shadow-md h-fit w-full inline-block relative">
      {isSelectable && (
        <div className="absolute top-2 left-2 z-10">
          <input
            type="checkbox"
            checked={isSelected}
            onChange={() => onToggleSelect(bookmark.id)}
            className="h-5 w-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
          />
        </div>
      )}
      <div className={`${isSelectable ? 'pl-8' : ''}`}>
        <div className="flex items-start gap-3 mb-3">
          <Image
            src={bookmark.profile_image_url}
            alt={bookmark.name}
            width={48}
            height={48}
            className="rounded-full"
            quality={100}
            priority
          />
          <div className="flex-1 min-w-0">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2 min-w-0">
                <h3 className="font-semibold truncate">{bookmark.name}</h3>
                <span className="text-sm text-gray-500 truncate">@{bookmark.screen_name}</span>
              </div>
              <div className="flex items-center gap-2">
                <a
                  href={bookmark.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-blue-600 hover:underline"
                >
                  <svg
                    className="h-5 w-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                    />
                  </svg>
                </a>
                <button
                  onClick={handleDotsClick}
                  className="p-1 rounded-full hover:bg-gray-100"
                >
                  <svg className="w-5 h-5 text-gray-500" viewBox="0 0 16 16" fill="currentColor">
                    <path d="M3 9.5a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3zm5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3zm5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3z" />
                  </svg>
                </button>
              </div>
            </div>
            
            <div className="mt-1 break-words">
              {bookmark.full_text.length > 280 ? (
                <div>
                  <p className="text-gray-900 whitespace-pre-wrap">
                    {isExpanded 
                      ? renderTextWithLinks(bookmark.full_text)
                      : renderTextWithLinks(bookmark.full_text.slice(0, 280) + '...')}
                  </p>
                  <button
                    onClick={() => setIsExpanded(!isExpanded)}
                    className="mt-1 text-blue-600 hover:underline text-sm"
                  >
                    {isExpanded ? 'Show less' : 'Read more'}
                  </button>
                </div>
              ) : (
                <p className="text-gray-900 whitespace-pre-wrap">{renderTextWithLinks(bookmark.full_text)}</p>
              )}
            </div>
          </div>
        </div>

        {bookmark.media.length > 0 && (
          <div className="my-4">
            <div className="grid grid-cols-1 gap-2 justify-center">
              {bookmark.media.map((media) => (
                <div 
                  key={media.id} 
                  className="relative aspect-video cursor-pointer overflow-hidden rounded-lg w-full"
                  onClick={() => setSelectedMedia({
                    url: media.original,
                    type: media.type === 'photo' ? 'image' : 'video'
                  })}
                >
                  <Image
                    src={media.thumbnail.replace('name=thumb', 'name=medium')}
                    alt=""
                    fill
                    quality={100}
                    sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
                    className="object-cover hover:opacity-90 transition-opacity"
                    priority
                    unoptimized
                  />
                </div>
              ))}
            </div>
          </div>
        )}

        <div className="space-y-3">
          <div className="flex flex-wrap gap-1.5">
            {bookmark.tags.map((tag) => (
              <span
                key={tag.id}
                className="group flex items-center gap-1 rounded-full bg-blue-100 px-2.5 py-1 text-sm text-blue-700"
              >
                {tag.name}
                <button
                  onClick={() => handleRemoveTag(tag.name)}
                  className="ml-1 hidden rounded-full p-0.5 hover:bg-blue-200 group-hover:block"
                  aria-label={`Remove ${tag.name} tag`}
                >
                  ×
                </button>
              </span>
            ))}
            <div className="relative">
              <input
                type="text"
                placeholder="Add tag..."
                className="w-20 rounded-full bg-blue-50 px-2.5 py-1 text-sm text-blue-700 placeholder:text-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:bg-blue-100 border-none h-[26px] overflow-x-auto whitespace-nowrap scrollbar-none"
                value={newTag}
                onChange={(e) => {
                  setNewTag(e.target.value);
                  updateSuggestedTags(e.target.value);
                }}
                onKeyDown={handleAddTag}
              />
              {suggestedTags.length > 0 && (
                <div className="absolute left-0 top-full mt-1 min-w-[200px] bg-white rounded-lg shadow-lg border z-[100]">
                  {suggestedTags.map(tag => (
                    <button
                      key={tag.id}
                      className="w-full px-4 py-2 text-left text-sm hover:bg-gray-100"
                      onClick={() => {
                        const updatedTags = [...bookmark.tags.map(t => t.name), tag.name];
                        onUpdateTags(bookmark.id, updatedTags);
                        setNewTag('');
                        setSuggestedTags([]);
                      }}
                    >
                      {highlightMatch(tag.name, newTag)}
                    </button>
                  ))}
                </div>
              )}
            </div>
          </div>

          <div className="flex items-center gap-4 text-sm text-gray-500 border-t pt-3">
            <span title="Likes">♥ {bookmark.favorite_count}</span>
            <span title="Retweets">🔄 {bookmark.retweet_count}</span>
            <span title="Views">👁 {bookmark.views_count}</span>
            <span className="ml-auto">
              <FormattedDate date={bookmark.created_at} />
            </span>
          </div>
        </div>

        {selectedMedia && (
          <MediaModal
            media={selectedMedia}
            onClose={() => setSelectedMedia(null)}
          />
        )}

        {isMenuOpen && (
          <>
            <div 
              className="fixed inset-0 z-[90]" 
              onClick={() => setIsMenuOpen(false)}
            />
            <div 
              className="fixed bg-white rounded-lg shadow-lg border z-[91]"
              style={{
                top: `${menuPosition.top}px`,
                left: `${menuPosition.left}px`,
              }}
            >
              <div className="py-1">
                <button
                  className="w-full px-4 py-2 text-left text-sm text-gray-600 hover:bg-gray-100"
                  onClick={() => {
                    setIsMenuOpen(false);
                    onArchive(bookmark.id);
                  }}
                >
                  {isArchived ? 'Unarchive' : 'Archive'}
                </button>
                <button
                  className="w-full px-4 py-2 text-left text-sm text-red-600 hover:bg-gray-100"
                  onClick={() => {
                    setIsMenuOpen(false);
                    setIsDeleteModalOpen(true);
                  }}
                >
                  Delete
                </button>
              </div>
            </div>
          </>
        )}

        {isDeleteModalOpen && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[100]">
            <div className="bg-white p-6 rounded-lg max-w-md w-full mx-4">
              <h3 className="text-lg font-semibold mb-2">Delete Bookmark</h3>
              <p className="mb-4">
                Are you sure you want to delete this bookmark?
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
                  onClick={handleDelete}
                >
                  Delete
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}