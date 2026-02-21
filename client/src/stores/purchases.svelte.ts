import { fetchPurchases, createPurchase, deletePurchase } from '../lib/api';
import type { Purchase } from '../models/Purchase';

let items = $state<Purchase[]>([]);

export const purchaseStore = {
  get items() {
    return items;
  },

  async load() {
    const data = await fetchPurchases();
    items = ((data as { purchases: Purchase[] }).purchases ?? []) as Purchase[];
  },

  async create(purchase: Omit<Purchase, 'id'>) {
    const created = (await createPurchase(purchase)) as Purchase;
    items = [...items, created];
    return created;
  },

  async delete(id: string) {
    await deletePurchase(id);
    items = items.filter((p) => p.id !== id);
  },

  clear() {
    items = [];
  },
};