import { useState, useEffect, forwardRef, useImperativeHandle } from 'react';
import { api } from '@/lib/api';

export interface StatisticsRef {
  refresh: () => void;
}

export const Statistics = forwardRef<StatisticsRef>((_, ref) => {
  const [stats, setStats] = useState<{
    total_bookmarks: number;
    active_bookmarks: number;
    archived_bookmarks: number;
    total_tags: number;
    top_tags: Array<{ name: string; count: number }>;
  } | null>(null);

  const [isLoading, setIsLoading] = useState(true);

  async function loadStatistics() {
    try {
      const data = await api.getStatistics();
      setStats(data);
    } catch (error) {
      console.error('Failed to load statistics:', error);
    } finally {
      setIsLoading(false);
    }
  }

  useImperativeHandle(ref, () => ({
    refresh: loadStatistics
  }));

  useEffect(() => {
    loadStatistics();
  }, []);

  if (isLoading) {
    return (
      <div className="bg-white rounded-lg border p-6">
        <div className="animate-pulse">
          <div className="h-6 bg-gray-200 rounded w-1/4 mb-8"></div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <div className="h-5 bg-gray-200 rounded w-1/3 mb-6"></div>
              <div className="space-y-4">
                <div className="h-4 bg-gray-200 rounded w-full"></div>
                <div className="h-4 bg-gray-200 rounded w-full"></div>
                <div className="h-4 bg-gray-200 rounded w-full"></div>
              </div>
            </div>
            <div>
              <div className="h-5 bg-gray-200 rounded w-1/3 mb-6"></div>
              <div className="space-y-4">
                {[...Array(5)].map((_, i) => (
                  <div key={i} className="h-4 bg-gray-200 rounded w-full"></div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!stats) return null;

  return (
    <div className="bg-white rounded-lg border p-6">
      <h2 className="text-xl font-semibold mb-8 text-gray-900">Dashboard Overview</h2>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div className="space-y-6">
          <h3 className="font-medium text-gray-900 text-lg">Summary</h3>
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-blue-50 rounded-lg p-4">
              <div className="text-blue-600 text-sm font-medium mb-1">Active Bookmarks</div>
              <div className="text-2xl font-bold text-gray-900">{stats.active_bookmarks}</div>
            </div>
            <div className="bg-orange-50 rounded-lg p-4">
              <div className="text-orange-600 text-sm font-medium mb-1">Archived</div>
              <div className="text-2xl font-bold text-gray-900">{stats.archived_bookmarks}</div>
            </div>
            <div className="bg-green-50 rounded-lg p-4">
              <div className="text-green-600 text-sm font-medium mb-1">Total Bookmarks</div>
              <div className="text-2xl font-bold text-gray-900">{stats.total_bookmarks}</div>
            </div>
            <div className="bg-purple-50 rounded-lg p-4">
              <div className="text-purple-600 text-sm font-medium mb-1">Tags</div>
              <div className="text-2xl font-bold text-gray-900">{stats.total_tags}</div>
            </div>
          </div>
        </div>

        <div className="space-y-6">
          <h3 className="font-medium text-gray-900 text-lg">Popular Tags</h3>
          <div className="space-y-3">
            {stats.top_tags.map((tag, index) => (
              <div 
                key={tag.name} 
                className="flex items-center justify-between p-3 rounded-lg bg-gray-50 hover:bg-gray-100 transition-colors"
              >
                <div className="flex items-center gap-3">
                  <span className="text-gray-500 text-sm">#{index + 1}</span>
                  <span className="text-gray-900 font-medium">{tag.name}</span>
                </div>
                <span className="bg-blue-100 text-blue-800 text-sm font-medium px-2.5 py-0.5 rounded-full">
                  {tag.count}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
});