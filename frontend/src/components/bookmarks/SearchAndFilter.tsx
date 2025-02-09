import { useState, useEffect, forwardRef, useImperativeHandle } from 'react';
import { api } from '@/lib/api';
import { Tag } from '@/types';

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
            <button
              key={tag.id}
              onClick={() => onTagSelect(tag.name)}
              className={`rounded-full px-3 py-1 text-sm whitespace-nowrap ${
                selectedTag === tag.name
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              }`}
            >
              {tag.name}
            </button>
          ))}
        </div>
      </div>
    );
  }
);