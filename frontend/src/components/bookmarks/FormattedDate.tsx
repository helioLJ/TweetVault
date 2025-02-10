'use client';

interface FormattedDateProps {
  date: string;
}

export function FormattedDate({ date }: FormattedDateProps) {
  return (
    <span className="text-xs text-gray-400">
      {new Date(date).toLocaleDateString(undefined, {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      })}
    </span>
  );
} 