import { useState, useEffect, forwardRef, useImperativeHandle } from 'react';
import { api } from '@/lib/api';
import { Tag } from '@/types';
import { TagMenu } from './TagMenu';

interface SearchAndFilterProps {
  onSearch: (query: string) => void;
  onTagSelect: (tag: string | undefined) => void;
  selectedTag?: string;
}

export interface SearchAndFilterRef {
  loadTags: () => void;
}

export const SearchAndFilter = forwardRef<SearchAndFilterRef, SearchAndFilterProps>(
  ({ onSearch, onTagSelect, selectedTag }, ref) => {
    const [tags, setTags] = useState<Tag[]>([]);
    const [searchQuery, setSearchQuery] = useState('');
    const [newTagInput, setNewTagInput] = useState('');
    const [suggestedTags, setSuggestedTags] = useState<Tag[]>([]);

    useImperativeHandle(ref, () => ({
      loadTags
    }));

    useEffect(() => {
      loadTags();
    }, []);

    async function loadTags() {
      try {
        const response = await api.getTags();
        setTags(response);
      } catch (error) {
        console.error('Failed to load tags:', error);
      }
    }

    const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
      const value = e.target.value;
      setSearchQuery(value);
      onSearch(value);
    };

    const updateSuggestedTags = (input: string) => {
      if (!input.trim()) {
        setSuggestedTags([]);
        return;
      }
      
      const filtered = tags.filter(tag => 
        tag.name.toLowerCase().includes(input.toLowerCase())
      );
      setSuggestedTags(filtered);
    };

    const handleAddTag = (tagName: string) => {
      onTagSelect(tagName);
      setNewTagInput('');
      setSuggestedTags([]);
    };

    return (
      <div className="mb-6 flex flex-col gap-4">
        <div className="relative w-full">
          <input
            type="search"
            placeholder="Search bookmarks..."
            className="w-full rounded-lg border border-gray-300 px-4 py-2 pl-10"
            value={searchQuery}
            onChange={handleSearch}
          />
          <svg
            className="absolute left-3 top-2.5 h-5 w-5 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        </div>
        <div className="flex flex-wrap gap-2 overflow-x-auto pb-2">
          <button
            onClick={() => onTagSelect(undefined)}
            className={`rounded-full px-3 py-1 text-sm whitespace-nowrap ${
              !selectedTag
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            All
          </button>
          {tags.map((tag) => (
            <div key={tag.id} className="flex items-center rounded-full overflow-hidden">
              <div
                onClick={() => onTagSelect(tag.name)}
                className={`px-3 py-1 text-sm whitespace-nowrap flex items-center cursor-pointer ${
                  selectedTag === tag.name
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                <span>{tag.name}</span>
                <TagMenu 
                  tag={tag} 
                  onSuccess={loadTags} 
                  selectedTag={selectedTag}
                />
              </div>
            </div>
          ))}
          <div className="relative inline-block">
            <input
              type="text"
              placeholder="Add tag..."
              className="rounded-full px-3 py-1 text-sm whitespace-nowrap bg-gray-100 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={newTagInput}
              onChange={(e) => setNewTagInput(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter' && newTagInput.trim()) {
                  handleAddTag(newTagInput.trim());
                }
              }}
            />
          </div>
        </div>
      </div>
    );
  }
);