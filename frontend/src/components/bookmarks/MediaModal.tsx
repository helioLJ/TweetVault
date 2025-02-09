import { useEffect } from 'react';
import Image from 'next/image';

interface MediaModalProps {
  media: {
    url: string;
    type: 'image' | 'video';
  };
  onClose: () => void;
}

export function MediaModal({ media, onClose }: MediaModalProps) {
  // Close on escape key
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, [onClose]);

  return (
    <>
      {/* Prevent scroll on body when modal is open */}
      <style jsx global>{`
        body {
          overflow: hidden;
        }
      `}</style>
      
      <div
        className="fixed inset-0 z-50 flex items-center justify-center bg-black/90"
        onClick={onClose}
      >
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-white hover:text-gray-300 z-50"
          aria-label="Close modal"
        >
          <svg
            className="h-8 w-8"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>

        <div 
          className="relative max-h-[90vh] max-w-[90vw]"
          onClick={e => e.stopPropagation()}
        >
          {media.type === 'image' ? (
            <Image
              src={media.url}
              alt=""
              className="rounded-lg"
              width={1200}
              height={800}
              style={{ 
                maxHeight: '90vh',
                width: 'auto',
                objectFit: 'contain'
              }}
            />
          ) : (
            <video
              src={media.url}
              controls
              className="rounded-lg"
              style={{ maxHeight: '90vh', maxWidth: '90vw' }}
            />
          )}
        </div>
      </div>
    </>
  );
} 