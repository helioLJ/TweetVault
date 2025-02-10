import { useState } from 'react';
import { Tag } from '@/types';
import { api } from '@/lib/api';

interface TagMenuProps {
  tag: Tag;
  onSuccess: () => void;
  selectedTag?: string;
}

export function TagMenu({ tag, onSuccess, selectedTag }: TagMenuProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [newName, setNewName] = useState(tag.name);
  const [isDeleting, setIsDeleting] = useState(false);
  const [bookmarkCount, setBookmarkCount] = useState<number | null>(null);
  const [menuPosition, setMenuPosition] = useState({ top: 0, left: 0 });

  async function handleEdit() {
    try {
      await api.updateTag(tag.id.toString(), newName);
      setIsEditing(false);
      setIsMenuOpen(false);
      onSuccess();
    } catch (error) {
      console.error('Failed to update tag:', error);
    }
  }

  async function handleDelete() {
    try {
      await api.deleteTag(tag.id.toString());
      setIsDeleting(false);
      setIsMenuOpen(false);
      onSuccess();
    } catch (error) {
      console.error('Failed to delete tag:', error);
    }
  }

  async function showDeleteConfirmation() {
    try {
      const response = await api.getTagBookmarkCount(tag.id.toString());
      setBookmarkCount(response.count);
      setIsDeleting(true);
    } catch (error) {
      console.error('Failed to get bookmark count:', error);
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

  if (isDeleting) {
    return (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[100]" style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0 }}>
        <div className="bg-white p-6 rounded-lg max-w-md w-full mx-4">
          <h3 className="text-lg font-semibold mb-2">Delete Tag</h3>
          <p className="mb-4">
            Are you sure you want to delete the tag "{tag.name}"? 
            {bookmarkCount !== null && bookmarkCount > 0 && (
              <span className="block mt-2 text-gray-600">
                This tag is used in {bookmarkCount} bookmark{bookmarkCount !== 1 ? 's' : ''}. 
                The bookmarks won't be deleted.
              </span>
            )}
          </p>
          <div className="flex justify-end gap-2">
            <button
              className="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded"
              onClick={() => setIsDeleting(false)}
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
    );
  }

  return (
    <span className="inline-flex items-center ml-1">
      <div 
        onClick={handleDotsClick}
        className={`p-1 rounded-full cursor-pointer ${
          selectedTag === tag.name ? 'text-white hover:bg-blue-700' : 'text-gray-500 hover:bg-gray-200'
        }`}
      >
        <svg className="w-4 h-4" viewBox="0 0 16 16" fill="currentColor">
          <path d="M3 9.5a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3zm5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3zm5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3z" />
        </svg>
      </div>
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
            {isEditing ? (
              <div className="p-2 min-w-[200px]">
                <input
                  type="text"
                  value={newName}
                  onChange={(e) => setNewName(e.target.value)}
                  className="w-full px-2 py-1 border rounded text-gray-900"
                  autoFocus
                />
                <div className="flex justify-end gap-2 mt-2">
                  <button
                    className="px-2 py-1 text-sm text-gray-600 hover:bg-gray-100 rounded"
                    onClick={() => {
                      setIsEditing(false);
                      setIsMenuOpen(false);
                    }}
                  >
                    Cancel
                  </button>
                  <button
                    className="px-2 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
                    onClick={handleEdit}
                  >
                    Save
                  </button>
                </div>
              </div>
            ) : (
              <div className="py-1 min-w-[120px] flex flex-col text-gray-900">
                <button
                  className="w-full px-4 py-2 text-left text-sm hover:bg-gray-100"
                  onClick={() => setIsEditing(true)}
                >
                  Rename
                </button>
                <button
                  className="w-full px-4 py-2 text-left text-sm text-red-600 hover:bg-gray-100"
                  onClick={showDeleteConfirmation}
                >
                  Delete
                </button>
              </div>
            )}
          </div>
        </>
      )}
    </span>
  );
}