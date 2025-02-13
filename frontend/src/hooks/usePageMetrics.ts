import { useEffect } from 'react';

export function usePageMetrics(pageName: string) {
  useEffect(() => {
    // Record page load timing
    const pageLoadTime = performance.now();

    // Report initial load
    if (navigator.sendBeacon) {
      navigator.sendBeacon('/api/metrics/pageview', JSON.stringify({
        page: pageName,
        loadTime: pageLoadTime,
        timestamp: new Date().toISOString(),
      }));
    }

    // Setup performance observer
    const observer = new PerformanceObserver((list) => {
      list.getEntries().forEach((entry) => {
        if (navigator.sendBeacon) {
          navigator.sendBeacon('/api/metrics/performance', JSON.stringify({
            page: pageName,
            metric: entry.name,
            value: entry.startTime,
            timestamp: new Date().toISOString(),
          }));
        }
      });
    });

    observer.observe({ entryTypes: ['largest-contentful-paint', 'first-input', 'layout-shift'] });

    return () => {
      observer.disconnect();
    };
  }, [pageName]);
} 