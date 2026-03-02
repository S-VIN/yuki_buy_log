<script lang="ts">
  import type { Product } from '../models/Product';
  import TagWidget from "./TagWidget.svelte";

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
      <TagWidget color="--color-green" text={product.volume}/>
    {/if}
    {#if product.brand}
      <TagWidget color="--color-yellow" text={product.brand}/>
    {/if}

    {#if needTags}
      {#each product.default_tags ?? [] as tag}
        <TagWidget color="--color-blue" text={tag}/>
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
</style>
