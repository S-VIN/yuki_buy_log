<script lang="ts">
  import { Search, X, Plus } from 'lucide-svelte';
  import { productStore } from '../stores/products.svelte';
  import ProductWidget from '../widgets/ProductWidget.svelte';
  import AddProductModal from '../widgets/AddProductModal.svelte';
  import type { ProductId } from '../models/Product';

  let searchQuery = $state('');
  let showModal = $state(false);
  let modalMode = $state<'add' | 'edit'>('add');
  let editingProductId = $state<ProductId | null>(null);

  const filteredProducts = $derived(
    searchQuery.trim()
      ? productStore.items.filter(p =>
          p.name.toLowerCase().includes(searchQuery.toLowerCase())
        )
      : productStore.items
  );

  function clearSearch() {
    searchQuery = '';
  }

  function openAdd() {
    modalMode = 'add';
    editingProductId = null;
    showModal = true;
  }

  function openEdit(id: ProductId) {
    modalMode = 'edit';
    editingProductId = id;
    showModal = true;
  }
</script>

<div class="page">
  <div class="search-header">
    <div class="search-row">
      <span class="search-icon-wrap"><Search size={16} /></span>
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Search products…"
        class="search-input"
      />
      {#if searchQuery}
        <button
          type="button"
          class="clear-btn"
          onclick={clearSearch}
          aria-label="Clear search"
        >
          <X size={14} />
        </button>
      {/if}
    </div>

    <button
      type="button"
      class="add-btn"
      onclick={openAdd}
      aria-label="Add product"
    >
      <Plus size={18} />
    </button>
  </div>

  <div class="list-area">
    {#if filteredProducts.length === 0}
      <div class="empty-state">
        {searchQuery ? 'No products match your search' : 'No products yet'}
      </div>
    {:else}
      <div class="product-list">
        {#each filteredProducts as product (product.id)}
          <button type="button" class="product-item" onclick={() => openEdit(product.id)}>
            <ProductWidget {product} needTags={true} />
          </button>
        {/each}
      </div>
    {/if}
  </div>
</div>

<AddProductModal
  open={showModal}
  onClose={() => (showModal = false)}
  mode={modalMode}
  productId={editingProductId}
/>

<style>
  .page {
    height: 100%;
    display: flex;
    flex-direction: column;
    background: var(--color-bg);
    overflow: hidden;
  }

  /* Sticky search header — visually elevated above the scrollable list */
  .search-header {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    gap: var(--space-3);
    background: var(--color-surface);
    border-bottom: 1px solid var(--color-border);
    box-shadow: var(--shadow-sm);
    padding: var(--space-4);
    z-index: 1;
  }

  .search-row {
    flex: 1;
    display: flex;
    align-items: center;
    gap: var(--space-3);
    height: var(--input-height);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    padding: 0 var(--space-4);
    background: var(--color-bg);
    transition: border-color var(--transition-base), box-shadow var(--transition-base);
  }

  .search-row:focus-within {
    border-color: var(--color-blue);
    box-shadow: var(--focus-ring);
  }

  .search-icon-wrap {
    display: inline-flex;
    align-items: center;
    color: var(--color-disabled);
    flex-shrink: 0;
  }

  .search-input {
    flex: 1;
    min-width: 0;
    border: none;
    outline: none;
    font-size: var(--text-base);
    background: transparent;
    color: var(--color-text);
  }

  .search-input::placeholder {
    color: var(--color-disabled);
  }

  .clear-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    border-radius: var(--radius-xs);
    background: none;
    border: none;
    color: var(--color-text-secondary);
    cursor: pointer;
    padding: 0;
    opacity: 0.6;
    transition: opacity var(--transition-fast), background var(--transition-fast), color var(--transition-fast);
    flex-shrink: 0;
  }

  .clear-btn:hover {
    opacity: 1;
    background: color-mix(in srgb, var(--color-red) 10%, transparent);
    color: var(--color-red);
  }

  .add-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: var(--input-height);
    height: var(--input-height);
    flex-shrink: 0;
    border-radius: var(--radius-md);
    border: 1px solid color-mix(in srgb, var(--color-blue) 30%, transparent);
    background: color-mix(in srgb, var(--color-blue) 8%, transparent);
    color: var(--color-blue);
    cursor: pointer;
    transition: background var(--transition-fast), border-color var(--transition-fast);
  }

  .add-btn:hover {
    background: color-mix(in srgb, var(--color-blue) 15%, transparent);
    border-color: color-mix(in srgb, var(--color-blue) 50%, transparent);
  }

  .add-btn:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
  }

  /* Scrollable list — scrollbar hidden per SPA convention */
  .list-area {
    flex: 1;
    overflow-y: auto;
    scrollbar-width: none;
    padding: var(--space-4);
  }

  .list-area::-webkit-scrollbar {
    display: none;
  }

  .product-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }

  .product-item {
    display: block;
    width: 100%;
    text-align: left;
    font: inherit;
    color: inherit;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    padding: var(--card-padding);
    box-shadow: var(--shadow-sm);
    transition: box-shadow var(--transition-fast), background var(--transition-fast);
    cursor: pointer;
  }

  .product-item:hover {
    box-shadow: var(--shadow-md);
    background: var(--color-bg);
  }

  .empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 120px;
    font-size: var(--text-base);
    color: var(--color-disabled);
  }
</style>