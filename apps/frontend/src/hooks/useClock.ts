import { useEffect, useState } from 'react';
export function useClock() {
  const [time, setTime] = useState(Date.now);
  useEffect(() => {
    const interval = window.setInterval(() => setTime(Date.now()), 1000);
    return () => window.clearInterval(interval);
  }, []);
  return time;
}
