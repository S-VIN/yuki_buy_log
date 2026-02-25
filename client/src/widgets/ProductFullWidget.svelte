<script lang="ts">
  import type { Product } from '../models/Product';

  interface Props {
    product: Product;
    needTags: boolean;
  }

  let { product, needTags }: Props = $props();
</script>

<div class="product-widget">
  <span class="product-name">{product.name}</span>
  <div class="tags-row">
    {#if product.volume}
      <span class="pill pill-green">{product.volume}</span>
    {/if}
    {#if product.brand}
      <span class="pill pill-yellow">{product.brand}</span>
    {/if}

    {#if needTags}
      {#each product.default_tags ?? [] as tag}
        <span class="pill pill-blue">{tag}</span>
      {/each}
    {/if}

  </div>
</div>

<style>
  .product-widget {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
  }

  .product-name {
    font-size: var(--text-base);
    font-weight: 700;
    color: var(--color-text);
    line-height: 1.4;
  }

  .tags-row {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-2);
    align-items: center;
  }

  .pill {
    display: inline-flex;
    align-items: center;
    border-radius: var(--radius-sm);
    padding: 2px 8px;
    font-size: var(--text-sm);
    font-weight: 600;
    letter-spacing: 0.01em;
    white-space: nowrap;
    line-height: 1.5;
  }

  .pill-blue {
    background: color-mix(in srgb, var(--color-blue) 10%, transparent);
    color: var(--color-blue);
    border: 1px solid color-mix(in srgb, var(--color-blue) 28%, transparent);
  }

  .pill-green {
    background: color-mix(in srgb, var(--color-green) 10%, transparent);
    color: var(--color-green);
    border: 1px solid color-mix(in srgb, var(--color-green) 28%, transparent);
  }

  .pill-yellow {
    background: color-mix(in srgb, var(--color-yellow) 15%, transparent);
    color: color-mix(in srgb, var(--color-yellow) 80%, var(--color-text));
    border: 1px solid color-mix(in srgb, var(--color-yellow) 40%, transparent);
  }
</style>
