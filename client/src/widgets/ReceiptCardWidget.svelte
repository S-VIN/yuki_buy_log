<script lang="ts">
  import type { Receipt } from '../models/Receipt';

  interface Props {
    receipt: Receipt;
    onclick?: () => void;
  }

  let { receipt, onclick }: Props = $props();
</script>

<div class="card" role={onclick ? 'button' : undefined} tabindex={onclick ? 0 : undefined}
  {onclick} onkeydown={onclick ? (e) => { if (e.key === 'Enter') onclick(); } : undefined}>
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

  .card[role='button'] {
    cursor: pointer;
    transition: background var(--transition-fast), border-color var(--transition-fast);
    -webkit-tap-highlight-color: transparent;
  }

  .card[role='button']:hover {
    background: color-mix(in srgb, var(--color-surface) 85%, var(--color-blue));
    border-color: color-mix(in srgb, var(--color-border) 60%, var(--color-blue));
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
