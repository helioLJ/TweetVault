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
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/80"
      onClick={onClose}
    >
      <div className="relative max-h-[90vh] max-w-[90vw]" onClick={e => e.stopPropagation()}>
        {media.type === 'image' ? (
          <Image
            src={media.url}
            alt=""
            className="rounded-lg"
            width={1200}
            height={800}
            style={{ objectFit: 'contain' }}
          />
        ) : (
          <video
            src={media.url}
            controls
            className="rounded-lg"
            style={{ maxHeight: '90vh', maxWidth: '90vw' }}
          />
        )}
        <button
          onClick={onClose}
          className="absolute -top-4 -right-4 rounded-full bg-white p-2 shadow-lg hover:bg-gray-100"
          aria-label="Close modal"
        >
          <svg
            className="h-6 w-6"
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
      </div>
    </div>
  );
} 