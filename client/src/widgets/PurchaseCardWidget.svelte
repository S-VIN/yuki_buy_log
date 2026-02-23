<script lang="ts">
  import { Trash2, Pencil } from 'lucide-svelte';
  import ProductWidget from './ProductWidget.svelte';
  import type { Product } from '../models/Product';

  interface Props {
    product: Product;
    price: number;
    quantity: number;
    tags: string[];
    onEdit: () => void;
    onDelete: () => void;
  }

  let { product, price, quantity, tags, onEdit, onDelete }: Props = $props();
</script>

<div class="card">
  <div class="card-product">
    <ProductWidget {product} />
  </div>
  <div class="card-footer">
    <div class="card-info">
      <span class="price-qty">
        {price.toFixed(2)} × {quantity}
        {#if quantity > 1}
          <span class="total">= {(price * quantity).toFixed(2)}</span>
        {/if}
      </span>
      {#if tags.length > 0}
        <div class="card-tags">
          {#each tags as tag}
            <span class="tag-pill">{tag}</span>
          {/each}
        </div>
      {/if}
    </div>
    <div class="card-actions">
      <button
        type="button"
        class="action-btn edit-btn"
        onclick={onEdit}
        aria-label="Edit"
      >
        <Pencil size={14} />
      </button>
      <button
        type="button"
        class="action-btn delete-btn"
        onclick={onDelete}
        aria-label="Delete"
      >
        <Trash2 size={14} />
      </button>
    </div>
  </div>
</div>

<style>
  .card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    padding: var(--card-padding);
    box-shadow: var(--shadow-sm);
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
  }

  .card-footer {
    display: flex;
    align-items: flex-start;
    gap: var(--space-3);
  }

  .card-info {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .price-qty {
    font-size: var(--text-sm);
    color: var(--color-text-secondary);
    font-weight: 500;
  }

  .total {
    color: var(--color-text);
    font-weight: 600;
  }

  .card-tags {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-2);
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

  .card-actions {
    display: flex;
    gap: var(--space-2);
    flex-shrink: 0;
  }

  .action-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border);
    background: none;
    cursor: pointer;
    transition: background var(--transition-fast), border-color var(--transition-fast),
      color var(--transition-fast);
  }

  .edit-btn {
    color: var(--color-text-secondary);
  }

  .edit-btn:hover {
    background: color-mix(in srgb, var(--color-blue) 10%, transparent);
    border-color: color-mix(in srgb, var(--color-blue) 30%, transparent);
    color: var(--color-blue);
  }

  .delete-btn {
    color: var(--color-text-secondary);
  }

  .delete-btn:hover {
    background: color-mix(in srgb, var(--color-red) 10%, transparent);
    border-color: color-mix(in srgb, var(--color-red) 30%, transparent);
    color: var(--color-red);
  }
</style>