<script lang="ts">
  import { Trash2, Pencil } from 'lucide-svelte';
  import type { Product } from '../models/Product';
  import ProductOneLineWidget from "./ProductOneLineWidget.svelte";
  import TagWidget from "./TagWidget.svelte";

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

    <div class="card-row">
      <ProductOneLineWidget {product}/>

      <div class="actions-row">
        <button type="button" class="action-btn edit-btn" onclick={onEdit} aria-label="Edit">
          <Pencil size={14} />
        </button>
        <button type="button" class="action-btn delete-btn" onclick={onDelete} aria-label="Delete">
          <Trash2 size={14} />
        </button>
      </div>
    </div>

    <div class="card-row">
      <div class="tags-row">
        {#each tags as tag (tag)}
          <TagWidget text={tag} color="--color-blue"/>
        {/each}
      </div>

      <span class="price-label">{price} x {quantity} = {price * quantity}₽</span>
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

  .card-product {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .card-row {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
  }

  .actions-row {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: var(--space-1);
    flex-shrink: 0;
  }

  .tags-row {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    align-items: center;
    gap: var(--space-2);
  }

  .price-label {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-secondary);
    white-space: nowrap;
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