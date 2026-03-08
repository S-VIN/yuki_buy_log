<script lang="ts">
  import type { Receipt } from '../models/Receipt';

  interface Props {
    receipt: Receipt;
  }

  let { receipt }: Props = $props();
</script>

<div class="card">
  <span class="store">{receipt.store || '—'}</span>
  <span class="total">{receipt.total.toFixed(2)}₽</span>
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
    flex: 1;
    font-size: var(--text-base);
    font-weight: 600;
    color: var(--color-text);
    min-width: 0;
  }

  .total {
    font-size: var(--text-base);
    font-weight: 700;
    color: var(--color-text);
    flex-shrink: 0;
  }
</style>
