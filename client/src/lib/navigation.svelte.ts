type TabId = 'add' | 'list' | 'products' | 'profile';

let _go: ((tab: TabId, data?: unknown) => void) | null = null;

export const navigation = {
  register(fn: (tab: TabId, data?: unknown) => void) {
    _go = fn;
  },
  go(tab: TabId, data?: unknown) {
    _go?.(tab, data);
  },
};
