import { purchaseStore } from './purchases.svelte';
import type { Receipt } from '../models/Receipt';

const items = $derived.by((): Receipt[] => {
  const purchases = purchaseStore.items;

  if (purchases.length === 0) return [];

  // Group purchases by receipt_id
  const receiptMap = new Map<number, typeof purchases>();
  for (const p of purchases) {
    const rid = p.receipt_id ?? 0;
    if (!receiptMap.has(rid)) receiptMap.set(rid, []);
    receiptMap.get(rid)!.push(p);
  }

  const receipts: Receipt[] = [];
  for (const [rid, ps] of receiptMap) {
    const first = ps[0];

    // Common tags: tags present in every purchase of this receipt
    const common_tags = first.tags.filter((t) => ps.every((p) => p.tags.includes(t)));

    // Total: sum of price × quantity for all purchases
    const total = ps.reduce((sum, p) => sum + p.price * (p.quantity ?? 1), 0);

    receipts.push({
      id: rid,
      date: first.date,
      store: first.store ?? '',
      common_tags,
      purchase_ids: ps.map((p) => p.id),
      total,
    });
  }

  return receipts;
});

export const receiptStore = {
  get items() {
    return items;
  },
};
