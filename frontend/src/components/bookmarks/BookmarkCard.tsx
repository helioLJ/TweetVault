import { useState } from 'react';
import { Bookmark } from '@/types';
import Image from 'next/image';
import { FormattedDate } from './FormattedDate';
import { MediaModal } from './MediaModal';
import { linkify } from '@/lib/utils';

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
  const [isExpanded, setIsExpanded] = useState(false);

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

  return (
    <div className="rounded-lg border bg-white p-4 shadow-sm transition-shadow hover:shadow-md h-fit w-full inline-block">
      <div className="flex items-start gap-3 mb-3">
        <Image
          src={bookmark.profile_image_url}
          alt={bookmark.name}
          width={48}
          height={48}
          className="rounded-full"
          quality={95}
          priority
        />
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 min-w-0">
              <h3 className="font-semibold truncate">{bookmark.name}</h3>
              <span className="text-sm text-gray-500 truncate">@{bookmark.screen_name}</span>
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
                  src={media.thumbnail}
                  alt=""
                  fill
                  quality={100}
                  sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
                  className="object-cover hover:opacity-90 transition-opacity"
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
          <input
            type="text"
            placeholder="Add tag..."
            className="w-20 rounded-full bg-blue-50 px-2.5 py-1 text-sm text-blue-700 placeholder:text-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:bg-blue-100 border-none h-[26px] overflow-x-auto whitespace-nowrap scrollbar-none"
            value={newTag}
            onChange={(e) => setNewTag(e.target.value)}
            onKeyDown={handleAddTag}
          />
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
    </div>
  );
}