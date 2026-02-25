<script lang="ts">
    import type { Product } from '../models/Product';
    import TagWidget from "./TagWidget.svelte";

    interface Props {
        product: Product;
    }

    let { product }: Props = $props();
</script>

<div class="product-widget">

    <span class="product-name">{product.name}</span>

    <div class="tags-row">
        {#if product.volume}
            <TagWidget text={product.volume} color="--color-green"/>
        {/if}
        {#if product.brand}
            <TagWidget text={product.brand} color="--color-yellow"/>
        {/if}
    </div>

</div>

<style>
    .product-widget {
        display: flex;
        justify-content: space-between;
        align-items: center;
        flex-direction: row;
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