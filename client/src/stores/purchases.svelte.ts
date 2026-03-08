import { fetchPurchases, createPurchase, deletePurchase } from '../lib/api';
import type { Purchase } from '../models/Purchase';

let items = $state<Purchase[]>([]);

const shops = $derived(
  Array.from(new Set(items.map((p) => p.store).filter((s): s is string => !!s && s.length > 0)))
);

export const purchaseStore = {
  get items() {
    return items;
  },

  /** Unique store names from all purchases, reactive. */
  get shops() {
    return shops;
  },

  async load() {
    const data = await fetchPurchases();
    items = ((data as { purchases: Purchase[] }).purchases ?? []).map((p) => ({
      ...p,
      date: new Date(p.date as unknown as string),
    }));
  },

  async create(purchase: Omit<Purchase, 'id'>) {
    const created = (await createPurchase(purchase)) as Purchase;
    const withDate = { ...created, date: new Date(created.date as unknown as string) };
    items = [...items, withDate];
    return withDate;
  },

  async delete(id: bigint) {
    await deletePurchase(id);
    items = items.filter((p) => p.id !== id);
  },

  clear() {
    items = [];
  },
};