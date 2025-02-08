'use client';

interface FormattedDateProps {
  date: string;
}

export function FormattedDate({ date }: FormattedDateProps) {
  return (
    <span className="text-xs text-gray-400">
      {new Date(date).toLocaleDateString()}
    </span>
  );
} 