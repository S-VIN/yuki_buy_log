<script lang="ts">
  import type { Receipt } from '../models/Receipt';

  interface Props {
    receipt: Receipt;
  }

  let { receipt }: Props = $props();
</script>

<div class="card">
  <span class="store">{receipt.store || '—'}</span>
  <div class="tags">
    {#each receipt.common_tags as tag}
      <span class="tag-pill">{tag}</span>
    {/each}
  </div>
  <span class="total">{receipt.total.toFixed(2)}</span>
</div>

<style>
  .card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    padding: var(--card-padding);
    box-shadow: var(--shadow-sm);
    display: flex;
    align-items: center;
    gap: var(--space-4);
    min-height: 44px;
  }

  .store {
    font-size: var(--text-base);
    font-weight: 600;
    color: var(--color-text);
    flex-shrink: 0;
    min-width: 0;
  }

  .tags {
    flex: 1;
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-2);
    justify-content: center;
    min-width: 0;
  }

  .tag-pill {
    display: inline-flex;
    align-items: center;
    background: color-mix(in srgb, var(--color-blue) 10%, transparent);
    color: var(--color-blue);
    border: 1px solid color-mix(in srgb, var(--color-blue) 28%, transparent);
    border-radius: var(--radius-sm);
    padding: 2px 8px;
    font-size: var(--text-xs);
    font-weight: 600;
    white-space: nowrap;
  }

  .total {
    font-size: var(--text-base);
    font-weight: 700;
    color: var(--color-text);
    flex-shrink: 0;
  }
</style>
