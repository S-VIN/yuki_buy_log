<script lang="ts">
  import { receiptStore } from '../stores/receipts.svelte';
  import ReceiptCardWidget from '../widgets/ReceiptCardWidget.svelte';
  import ChecksHeaderWidget from "../widgets/ChecksHeaderWidget.svelte";
</script>

<div class="page">
  {#if receiptStore.groupedByDate.length === 0}
    <div class="empty">No receipts yet</div>
  {:else}
    <div class="list">
      {#each receiptStore.groupedByDate as group (group.date.toDateString())}
        <ChecksHeaderWidget date={group.date} total={group.total} />
        {#each group.receipts as receipt (receipt.id)}
          <ReceiptCardWidget {receipt} />
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
