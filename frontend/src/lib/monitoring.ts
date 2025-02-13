export const reportWebVitals = (metric: {
    id: string;
    name: string;
    value: number;
}) => {
    // Send to your analytics service
    console.log(metric);
    
    // Example implementation for sending to backend
    if (window.navigator.sendBeacon) {
        const body = JSON.stringify({
            name: metric.name,
            value: metric.value,
            id: metric.id,
        });
        
        window.navigator.sendBeacon('/api/metrics/vitals', body);
    }
};

export const initErrorTracking = () => {
    window.addEventListener('error', (event) => {
        // Send to your error tracking service
        console.error('Global error:', event.error);
    });

    window.addEventListener('unhandledrejection', (event) => {
        // Send to your error tracking service
        console.error('Unhandled promise rejection:', event.reason);
    });
}; 