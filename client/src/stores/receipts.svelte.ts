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

    const total = ps.reduce((sum, p) => sum + p.price * (p.quantity ?? 1), 0);

    receipts.push({
      id: rid,
      date: first.date,
      store: first.store ?? '',
      purchase_ids: ps.map((p) => p.id),
      total,
      userId: first.user_id!,
    });
  }

  return receipts;
});

const groupedByDate = $derived.by(() => {
  const sorted = [...items].sort((a, b) => b.date.getTime() - a.date.getTime());

  const groups: { date: Date; total: number; receipts: Receipt[] }[] = [];
  for (const receipt of sorted) {
    const last = groups.at(-1);
    if (last && last.date.toDateString() === receipt.date.toDateString()) {
      last.receipts.push(receipt);
      last.total += receipt.total;
    } else {
      groups.push({ date: receipt.date, total: receipt.total, receipts: [receipt] });
    }
  }
  return groups;
});

export const receiptStore = {
  get items() {
    return items;
  },
  get groupedByDate() {
    return groupedByDate;
  },
};
