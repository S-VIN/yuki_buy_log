<script lang="ts">
  import { receiptStore } from '../stores/receipts.svelte';
  import { purchaseStore } from '../stores/purchases.svelte';
  import { productStore } from '../stores/products.svelte';
  import { navigation } from '../lib/navigation.svelte';
  import ReceiptCardWidget from '../widgets/ReceiptCardWidget.svelte';
  import ChecksHeaderWidget from "../widgets/ChecksHeaderWidget.svelte";

  async function handleReceiptClick(receiptId: number) {
    const receipt = receiptStore.items.find((r) => r.id === receiptId);
    if (!receipt) return;

    const purchases = purchaseStore.items.filter((p) => receipt.purchase_ids.includes(p.id));

    const items = purchases
      .map((p) => {
        const product = productStore.items.find((pr) => pr.id === p.product_id);
        if (!product) return null;
        return {
          uuid: crypto.randomUUID(),
          product,
          price: p.price,
          quantity: p.quantity ?? 1,
          tags: [...p.tags],
        };
      })
      .filter((x): x is NonNullable<typeof x> => x !== null);

    for (const purchase of purchases) {
      await purchaseStore.delete(purchase.id);
    }

    navigation.go('add', {
      date: receipt.date.toISOString().slice(0, 10),
      shop: receipt.store || null,
      items,
    });
  }
</script>

<div class="page">
  {#if receiptStore.groupedByDate.length === 0}
    <div class="empty">No receipts yet</div>
  {:else}
    <div class="list">
      {#each receiptStore.groupedByDate as group (group.date.toDateString())}
        <ChecksHeaderWidget date={group.date} total={group.total} />
        {#each group.receipts as receipt (receipt.id)}
          <ReceiptCardWidget {receipt} onclick={() => handleReceiptClick(receipt.id)} />
        {/each}
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

  .empty {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 120px;
    font-size: var(--text-base);
    color: var(--color-disabled);
  }
</style>
