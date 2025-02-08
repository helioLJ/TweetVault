import { useState } from 'react';
import { Bookmark } from '@/types';
import Image from 'next/image';
import { FormattedDate } from './FormattedDate';
import { MediaModal } from './MediaModal';

interface BookmarkCardProps {
  bookmark: Bookmark;
  onUpdateTags: (id: string, tags: string[]) => void;
}

export function BookmarkCard({ bookmark, onUpdateTags }: BookmarkCardProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [newTag, setNewTag] = useState('');
  const [selectedMedia, setSelectedMedia] = useState<null | {
    url: string;
    type: 'image' | 'video';
  }>(null);

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

  return (
    <div className="rounded-lg border bg-white p-4 shadow-sm transition-shadow hover:shadow-md">
      <div className="flex items-start gap-3">
        <Image
          src={bookmark.profile_image_url}
          alt={bookmark.name}
          width={48}
          height={48}
          className="rounded-full"
        />
        <div className="flex-1">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <h3 className="font-semibold">{bookmark.name}</h3>
              <span className="text-sm text-gray-500">@{bookmark.screen_name}</span>
            </div>
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
          </div>
          <p className="mt-1 text-gray-900">{bookmark.full_text}</p>
          
          {bookmark.media.length > 0 && (
            <div className="mt-3 grid grid-cols-2 gap-2">
              {bookmark.media.map((media) => (
                <div 
                  key={media.id} 
                  className="relative aspect-video cursor-pointer"
                  onClick={() => setSelectedMedia({
                    url: media.original,
                    type: media.type === 'photo' ? 'image' : 'video'
                  })}
                >
                  <Image
                    src={media.thumbnail}
                    alt=""
                    fill
                    sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
                    className="rounded-lg object-cover hover:opacity-90 transition-opacity"
                  />
                </div>
              ))}
            </div>
          )}
          
          <div className="mt-3">
            <div className="flex flex-wrap gap-1">
              {bookmark.tags.map((tag) => (
                <span
                  key={tag.id}
                  className="group flex items-center gap-1 rounded-full bg-blue-100 px-2 py-0.5 text-sm text-blue-700"
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
                  className="inline-flex h-6 rounded-full bg-gray-100 px-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  value={newTag}
                  onChange={(e) => setNewTag(e.target.value)}
                  onKeyDown={handleAddTag}
                />
              </div>
            </div>
          </div>
          
          <div className="mt-3 flex items-center gap-4 text-sm text-gray-500">
            <span title="Likes">♥ {bookmark.favorite_count}</span>
            <span title="Retweets">🔄 {bookmark.retweet_count}</span>
            <span title="Views">👁 {bookmark.views_count}</span>
            <span className="ml-auto">
              <FormattedDate date={bookmark.created_at} />
            </span>
          </div>
        </div>
      </div>

      {/* Media Modal */}
      {selectedMedia && (
        <MediaModal
          media={selectedMedia}
          onClose={() => setSelectedMedia(null)}
        />
      )}
    </div>
  );
}