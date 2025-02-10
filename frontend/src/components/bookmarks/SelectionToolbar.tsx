interface SelectionToolbarProps {
  selectedCount: number;
  onClearSelection: () => void;
  onDeleteSelected: () => void;
}

export function SelectionToolbar({ 
  selectedCount, 
  onClearSelection, 
  onDeleteSelected 
}: SelectionToolbarProps) {
  return (
    <div className="fixed bottom-0 left-0 right-0 bg-white border-t shadow-lg z-50 py-4 px-6">
      <div className="max-w-[2000px] mx-auto flex items-center justify-between">
        <div className="flex items-center gap-4">
          <span className="font-medium">
            {selectedCount} bookmark{selectedCount !== 1 ? 's' : ''} selected
          </span>
          <button
            onClick={onClearSelection}
            className="text-gray-600 hover:text-gray-800"
          >
            Clear selection
          </button>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={onDeleteSelected}
            className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
          >
            Delete Selected
          </button>
        </div>
      </div>
    </div>
  );
} 