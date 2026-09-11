import { useEffect } from 'react';
import { useNavigation } from '../stores/navigation';
export function useDirtyForm(dirty: boolean) {
  useEffect(() => {
    useNavigation.getState().setDirty(dirty);
    return () => useNavigation.getState().setDirty(false);
  }, [dirty]);
}
