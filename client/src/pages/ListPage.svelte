<script lang="ts">
  import { receiptStore } from '../stores/receipts.svelte';
  import ReceiptCardWidget from '../widgets/ReceiptCardWidget.svelte';
  import type { Receipt } from '../models/Receipt';

  type ListItem =
    | { type: 'date'; date: string; dayTotal: number }
    | { type: 'receipt'; receipt: Receipt };

  function formatDate(dateStr: string): string {
    const d = new Date(dateStr.slice(0, 10) + 'T12:00:00');
    return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' });
  }

  const listItems = $derived.by((): ListItem[] => {
    const receipts = receiptStore.items;

    // Sort by date descending (newest first)
    const sorted = [...receipts].sort((a, b) => b.date.localeCompare(a.date));

    // Group by YYYY-MM-DD, preserving sorted order
    const dateOrder: string[] = [];
    const dateMap = new Map<string, Receipt[]>();
    for (const r of sorted) {
      const d = r.date.slice(0, 10);
      if (!dateMap.has(d)) {
        dateOrder.push(d);
        dateMap.set(d, []);
      }
      dateMap.get(d)!.push(r);
    }

    // Build flat list: date header + receipts for that date
    const items: ListItem[] = [];
    for (const date of dateOrder) {
      const dateReceipts = dateMap.get(date)!;
      const dayTotal = dateReceipts.reduce((sum, r) => sum + r.total, 0);
      items.push({ type: 'date', date, dayTotal });
      for (const r of dateReceipts) {
        items.push({ type: 'receipt', receipt: r });
      }
    }

    return items;
  });
</script>

<div class="page">
  {#if listItems.length === 0}
    <div class="empty">No receipts yet</div>
  {:else}
    <div class="list">
      {#each listItems as item (item.type === 'date' ? 'date-' + item.date : 'receipt-' + item.receipt.id)}
        {#if item.type === 'date'}
          <div class="date-header">
            <span class="date-label">{formatDate(item.date)}</span>
            <span class="date-total">{item.dayTotal.toFixed(2)}</span>
          </div>
        {:else}
          <ReceiptCardWidget receipt={item.receipt} />
        {/if}
      {/each}
    </div>
  {/if}
</div>

<style>
  .page {
    height: 100%;
    overflow-y: auto;
    background: var(--color-bg);
    padding: var(--page-padding);
    scrollbar-width: none;
  }

  .page::-webkit-scrollbar {
    display: none;
  }

  .list {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
    padding: var(--space-4);
  }

  .date-header {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    padding: var(--space-4) var(--space-2) var(--space-2);
    margin-top: var(--space-4);
  }

  .date-header:first-child {
    margin-top: 0;
  }

  .date-label {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-secondary);
    text-transform: lowercase;
  }

  .date-total {
    font-size: var(--text-sm);
    font-weight: 700;
    color: var(--color-text-secondary);
  }

  .empty {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 120px;
    font-size: var(--text-base);
    color: var(--color-disabled);
  }
</style>
