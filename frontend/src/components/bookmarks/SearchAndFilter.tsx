import { useState, useEffect, forwardRef, useImperativeHandle } from 'react';
import { api } from '@/lib/api';
import { Tag } from '@/types';
import { TagMenu } from './TagMenu';

interface SearchAndFilterProps {
  onSearch: (query: string) => void;
  onTagSelect: (tag: string | undefined) => void;
  selectedTag?: string;
  onDeleteTag?: (tagName: string) => void;
}

export interface SearchAndFilterRef {
  loadTags: () => void;
}

export const SearchAndFilter = forwardRef<SearchAndFilterRef, SearchAndFilterProps>(
  ({ onSearch, onTagSelect, selectedTag, onDeleteTag }, ref) => {
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
      
      // Use a case-insensitive search
      const searchTerm = input.toLowerCase();
      const filtered = tags.filter(tag => 
        tag.name.toLowerCase().includes(searchTerm)
      );
      setSuggestedTags(filtered);
    };

    const handleAddTag = async (tagName: string) => {
      if (!tagName.trim()) return;
      
      try {
        // First select the tag if it exists
        if (tags.some(t => t.name.toLowerCase() === tagName.toLowerCase())) {
          onTagSelect(tagName);
          setNewTagInput('');
          setSuggestedTags([]);
          return;
        }

        // Create new tag via API without selecting it
        await api.createTag(tagName);
        await loadTags(); // Reload tags after creating
        
        setNewTagInput('');
        setSuggestedTags([]);
      } catch (error) {
        console.error('Failed to add tag:', error);
      }
    };

    const handleDeleteTag = async (tagName: string) => {
      if (onDeleteTag) {
        onDeleteTag(tagName);
      }
      // If the deleted tag was selected, reset to show all bookmarks
      if (selectedTag === tagName) {
        onTagSelect(undefined);
      }
      await loadTags();
    };

    return (
      <div className="mb-6 flex flex-col gap-4">
        <div className="relative w-full">
          <input
            type="search"
            placeholder="Search bookmarks..."
            className="w-full rounded-lg border border-gray-300 dark:border-gray-700 px-4 py-2 pl-10 bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder:text-gray-500 dark:placeholder:text-gray-400"
            value={searchQuery}
            onChange={handleSearch}
          />
          <svg
            className="absolute left-3 top-2.5 h-5 w-5 text-gray-400 dark:text-gray-500"
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
                : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 hover:bg-gray-200 dark:hover:bg-gray-600'
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
                    : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 hover:bg-gray-200 dark:hover:bg-gray-600'
                }`}
              >
                <span>{tag.name}</span>
                <TagMenu 
                  tag={tag} 
                  onSuccess={loadTags} 
                  selectedTag={selectedTag}
                  onDeleteTag={handleDeleteTag}
                />
              </div>
            </div>
          ))}
          <div className="flex items-center rounded-full overflow-hidden">
            <div className="px-3 py-1 text-sm whitespace-nowrap flex items-center gap-1 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-600 hover:border-gray-300 dark:hover:border-gray-500">
              <svg 
                className="w-3.5 h-3.5 text-gray-500 dark:text-gray-400" 
                fill="none" 
                viewBox="0 0 24 24" 
                stroke="currentColor"
              >
                <path 
                  strokeLinecap="round" 
                  strokeLinejoin="round" 
                  strokeWidth={2} 
                  d="M12 4v16m8-8H4" 
                />
              </svg>
              <input
                type="text"
                placeholder="Add tag..."
                className="bg-transparent text-gray-700 dark:text-gray-200 placeholder:text-gray-500 dark:placeholder:text-gray-400 focus:outline-none w-20"
                value={newTagInput}
                onChange={(e) => {
                  const value = e.target.value;
                  setNewTagInput(value);
                  updateSuggestedTags(value);
                }}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && newTagInput.trim()) {
                    handleAddTag(newTagInput.trim());
                  }
                }}
              />
            </div>
          </div>
        </div>
      </div>
    );
  }
);