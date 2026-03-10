type TabId = 'add' | 'list' | 'products' | 'profile';

let _go: ((tab: TabId) => void) | null = null;

export const navigation = {
  register(fn: (tab: TabId) => void) {
    _go = fn;
  },
  go(tab: TabId) {
    _go?.(tab);
  },
};
