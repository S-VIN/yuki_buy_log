<script lang="ts">
  import { Plus, X, Search } from "lucide-svelte";
  import ProductWidget from "./ProductWidget.svelte";
  import AddProductModal from "./AddProductModal.svelte";
  import { productStore } from "../stores/products.svelte";
  import type { Product } from "../models/Product";

  interface Props {
    value: Product | null;
    id?: string;
  }

  let {
    value = $bindable(null),
    id,
  }: Props = $props();

  let inputValue = $state("");
  let isFocused = $state(false);
  let showModal = $state(false);

  const filteredProducts = $derived(
    inputValue.trim()
      ? productStore.items.filter((p) =>
          p.name.toLowerCase().includes(inputValue.toLowerCase())
        )
      : productStore.items
  );

  function selectProduct(product: Product) {
    value = product;
    inputValue = "";
    isFocused = false;
  }

  function clearSelection() {
    value = null;
    inputValue = "";
  }

  function openModal() {
    showModal = true;
    isFocused = false;
  }

  function handleProductAdded(product: Product) {
    value = product;
  }
</script>

<div class="product-select-widget">
  {#if value}
    <div class="input-row selected">
      <div class="product-preview">
        <ProductWidget product={value} needTags={false} />
      </div>
      <button
        type="button"
        class="clear-btn"
        onclick={clearSelection}
        aria-label="Clear selection"
      >
        <X size={14} />
      </button>
    </div>
  {:else}
    <div class="input-row" class:focused={isFocused}>
      <input
        {id}
        type="text"
        class="select-input"
        placeholder="Select product…"
        autocomplete="off"
        bind:value={inputValue}
        onfocus={() => (isFocused = true)}
        onblur={() => (isFocused = false)}
      />
      {#if inputValue.trim()}
        <button
          type="button"
          class="clear-btn"
          onclick={() => { inputValue = ''; }}
          aria-label="Clear"
        >
          <X size={14} />
        </button>
      {:else}
        <span class="input-icon"><Search size={14} /></span>
      {/if}
    </div>

    {#if isFocused && inputValue.trim()}
      <div
        class="dropdown"
        role="listbox"
        tabindex="-1"
        onmousedown={(e) => e.preventDefault()}
      >
        <button type="button" class="option add-option" onclick={openModal}>
          <span class="add-icon"><Plus size={12} /></span>
          <span class="add-label">Add new product</span>
          {#if inputValue.trim()}
            <span class="add-name-hint">"{inputValue.trim()}"</span>
          {/if}
        </button>

        {#if filteredProducts.length > 0}
          <div class="divider"></div>
          {#each filteredProducts as product (product.id)}
            <button
              type="button"
              class="option product-option"
              onclick={() => selectProduct(product)}
            >
              <ProductWidget product={product} needTags={false} />
            </button>
          {/each}
        {/if}
      </div>
    {/if}
  {/if}
</div>

<AddProductModal
  open={showModal}
  onClose={() => (showModal = false)}
  initialName={inputValue}
  onAdded={handleProductAdded}
/>

<style>
  .product-select-widget {
    position: relative;
  }

  /* ── Input row ─────────────────────────────────────────────── */
  .input-row {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    min-height: var(--input-height);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    padding: var(--space-2) var(--space-4);
    background: var(--color-surface);
    transition: border-color var(--transition-base), box-shadow var(--transition-base);
    cursor: text;
  }

  .input-row.focused {
    border-color: var(--color-blue);
    box-shadow: var(--focus-ring);
  }

  .input-row.selected {
    cursor: default;
    align-items: flex-start;
    padding: var(--space-3) var(--space-4);
    min-height: auto;
  }

  .product-preview {
    flex: 1;
    min-width: 0;
  }

  /* ── Clear button ──────────────────────────────────────────── */
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
    transition: opacity var(--transition-fast), background var(--transition-fast),
      color var(--transition-fast);
    flex-shrink: 0;
    align-self: center;
  }

  .clear-btn:hover {
    opacity: 1;
    background: color-mix(in srgb, var(--color-red) 10%, transparent);
    color: var(--color-red);
  }

  /* ── Text input ────────────────────────────────────────────── */
  .select-input {
    flex: 1;
    min-width: 60px;
    border: none;
    outline: none;
    font-size: var(--text-base);
    background: transparent;
    color: var(--color-text);
    padding: var(--space-1) 0;
  }

  .select-input::placeholder {
    color: var(--color-disabled);
  }

  .input-icon {
    display: inline-flex;
    align-items: center;
    color: var(--color-disabled);
    flex-shrink: 0;
    transition: color var(--transition-fast);
  }

  /* ── Dropdown ──────────────────────────────────────────────── */
  .dropdown {
    position: absolute;
    left: 0;
    width: 100%;
    z-index: 10;
    max-height: 280px;
    overflow-y: auto;
    scrollbar-width: none;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-lg);
    margin-top: var(--space-3);
    padding: var(--space-2);
    box-shadow: var(--shadow-md);
  }

  .dropdown::-webkit-scrollbar {
    display: none;
  }

  .divider {
    height: 1px;
    background: var(--color-border);
    margin: var(--space-2);
  }

  .option {
    display: flex;
    align-items: center;
    padding: var(--space-4) var(--space-5);
    cursor: pointer;
    font-size: var(--text-base);
    color: var(--color-text);
    border-radius: var(--radius-sm);
    transition: background var(--transition-fast);
    user-select: none;
    width: 100%;
    background: none;
    border: none;
    font: inherit;
    text-align: left;
  }

  .option:hover {
    background: var(--color-bg);
  }

  .option:focus {
    outline: none;
    background: var(--color-bg);
  }

  .product-option {
    flex-direction: column;
    align-items: stretch;
  }

  .add-option {
    gap: 6px;
    color: var(--color-text-secondary);
  }

  .add-option:hover .add-label,
  .add-option:focus .add-label {
    color: var(--color-text);
  }

  .add-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 18px;
    height: 18px;
    border-radius: var(--radius-xs);
    background: color-mix(in srgb, var(--color-blue) 10%, transparent);
    color: var(--color-blue);
    flex-shrink: 0;
  }

  .add-label {
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--color-text-secondary);
    transition: color var(--transition-fast);
  }

  .add-name-hint {
    font-size: var(--text-sm);
    color: var(--color-text);
    font-style: italic;
  }
</style>