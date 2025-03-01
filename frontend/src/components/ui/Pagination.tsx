interface PaginationProps {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
    pageSize: number;
    onPageSizeChange: (size: number) => void;
  }
  
  const PAGE_SIZE_OPTIONS = [12, 20, 32, 48];
  const PAGE_SIZE_KEY = 'bookmark-page-size';
  
  export function Pagination({ 
    currentPage, 
    totalPages, 
    onPageChange, 
    pageSize,
    onPageSizeChange 
  }: PaginationProps) {
    // Calculate visible page numbers
    let pages = Array.from({ length: totalPages }, (_, i) => i + 1);
    if (totalPages > 7) {
      if (currentPage <= 4) {
        pages = [...pages.slice(0, 5), -1, totalPages];
      } else if (currentPage >= totalPages - 3) {
        pages = [1, -1, ...pages.slice(totalPages - 5)];
      } else {
        pages = [1, -1, currentPage - 1, currentPage, currentPage + 1, -1, totalPages];
      }
    }
  
    const handlePageSizeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
      const newSize = Number(e.target.value);
      localStorage.setItem(PAGE_SIZE_KEY, newSize.toString());
      onPageSizeChange(newSize);
    };
  
    // Helper to render a button for a page number
    const renderPageButton = (page: number, key: string | number) => {
      return (
        <button
          key={key}
          type="button"
          className={`px-3 py-1 rounded ${
            page === currentPage ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-700'
          }`}
          onClick={(e) => {
            // Blur the button so it doesn't retain focus after click.
            e.currentTarget.blur();
            onPageChange(page);
          }}
        >
          {page}
        </button>
      );
    };
  
    return (
      <div className="mt-8 flex flex-col items-center gap-4">
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={(e) => {
              e.currentTarget.blur();
              onPageChange(currentPage - 1);
            }}
            disabled={currentPage === 1}
            className="rounded-md border dark:border-gray-700 px-3 py-1 text-sm disabled:opacity-50 hover:bg-gray-50 dark:hover:bg-gray-800 dark:text-gray-300 disabled:dark:hover:bg-transparent"
          >
            Previous
          </button>
  
          {pages.map((page, i) => (
            page === -1 ? (
              <span key={`ellipsis-${i}`} className="px-2 dark:text-gray-400">...</span>
            ) : (
              renderPageButton(page, page)
            )
          ))}
  
          <button
            type="button"
            onClick={(e) => {
              e.currentTarget.blur();
              onPageChange(currentPage + 1);
            }}
            disabled={currentPage === totalPages}
            className="rounded-md border dark:border-gray-700 px-3 py-1 text-sm disabled:opacity-50 hover:bg-gray-50 dark:hover:bg-gray-800 dark:text-gray-300 disabled:dark:hover:bg-transparent"
          >
            Next
          </button>
        </div>
  
        <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
          <span>Show</span>
          <select
            value={pageSize}
            onChange={handlePageSizeChange}
            className="rounded border dark:border-gray-700 px-2 py-1 dark:bg-gray-800 dark:text-gray-300"
          >
            {PAGE_SIZE_OPTIONS.map(size => (
              <option key={size} value={size}>
                {size}
              </option>
            ))}
          </select>
          <span>items per page</span>
        </div>
      </div>
    );
  }