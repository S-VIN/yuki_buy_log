н<script lang="ts">
  import { Plus, X } from "lucide-svelte";
  import ProductWidget from "./ProductWidget.svelte";
  import SelectWidget from "./SelectWidget.svelte";
  import TagWidget from "./TagWidget.svelte";
  import type { Product } from "../models/Product";

  interface Props {
    allProducts: Product[];
    value: Product | null;
    id?: string;
  }

  let {
    allProducts = $bindable([]),
    value = $bindable(null),
    id,
  }: Props = $props();

  let inputValue = $state("");
  let isFocused = $state(false);
  let showModal = $state(false);

  // Volume options — static sensible defaults; new ones added via SelectWidget
  let allVolumeOptions = $state([
    "100 мл", "200 мл", "250 мл", "330 мл", "500 мл",
    "750 мл", "1 л", "1.5 л", "2 л",
  ]);

  // Modal form state
  let newName = $state("");
  let newVolume = $state<string | null>(null);
  let newBrand = $state<string | null>(null);
  let newTags = $state<string[]>([]);
  let modalBrandOptions = $state<string[]>([]);
  let modalAllTags = $state<string[]>([]);

  const filteredProducts = $derived(
    inputValue.trim()
      ? allProducts.filter((p) =>
          p.name.toLowerCase().includes(inputValue.toLowerCase())
        )
      : allProducts
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
    newName = inputValue.trim();
    newVolume = null;
    newBrand = null;
    newTags = [];
    modalBrandOptions = Array.from(
      new Set(allProducts.map((p) => p.brand).filter((b) => b.length > 0))
    );
    modalAllTags = Array.from(
      new Set(allProducts.flatMap((p) => p.default_tags))
    );
    showModal = true;
    isFocused = false;
  }

  function closeModal() {
    showModal = false;
  }

  function confirmNewProduct() {
    if (!newName.trim()) return;
    const product: Product = {
      id: crypto.randomUUID(),
      name: newName.trim(),
      volume: newVolume ?? "",
      brand: newBrand ?? "",
      default_tags: [...newTags],
      user_id: "",
    };
    allProducts = [...allProducts, product];
    value = product;
    closeModal();
  }

  function toWidgetProduct(p: Product) {
    return {
      name: p.name,
      volume: p.volume || undefined,
      brand: p.brand || undefined,
      tags: p.default_tags,
    };
  }

  $effect(() => {
    if (!showModal) return;
    function onKeydown(e: KeyboardEvent) {
      if (e.key === "Escape") closeModal();
    }
    document.addEventListener("keydown", onKeydown);
    return () => document.removeEventListener("keydown", onKeydown);
  });
</script>

<div class="product-select-widget">
  {#if value}
    <div class="input-row selected">
      <div class="product-preview">
        <ProductWidget product={toWidgetProduct(value)} />
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
    </div>

    {#if isFocused}
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
              <ProductWidget product={toWidgetProduct(product)} />
            </button>
          {/each}
        {/if}
      </div>
    {/if}
  {/if}
</div>

{#if showModal}
  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
  <div
    class="modal-overlay"
    role="dialog"
    aria-modal="true"
    tabindex="-1"
    onmousedown={closeModal}
  >
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="modal-sheet" onmousedown={(e) => e.stopPropagation()}>
      <div class="modal-header">
        <h3 class="modal-title">New Product</h3>
        <button
          type="button"
          class="modal-close"
          onclick={closeModal}
          aria-label="Close"
        >
          <X size={18} />
        </button>
      </div>

      <div class="modal-body">
        <div class="form-field">
          <label class="form-label" for="new-product-name">Name</label>
          <input
            id="new-product-name"
            type="text"
            class="form-input"
            placeholder="Product name…"
            bind:value={newName}
          />
        </div>

        <div class="form-field">
          <label class="form-label" for="new-product-volume">Volume</label>
          <SelectWidget
            id="new-product-volume"
            bind:allOptions={allVolumeOptions}
            bind:value={newVolume}
            color="var(--color-green)"
            placeholder="Select volume…"
          />
        </div>

        <div class="form-field">
          <label class="form-label" for="new-product-brand">Brand</label>
          <SelectWidget
            id="new-product-brand"
            bind:allOptions={modalBrandOptions}
            bind:value={newBrand}
            color="var(--color-yellow)"
            placeholder="Select brand…"
          />
        </div>

        <div class="form-field">
          <label class="form-label" for="new-product-tags">Default Tags</label>
          <TagWidget
            id="new-product-tags"
            bind:allTags={modalAllTags}
            bind:selectedTags={newTags}
          />
        </div>
      </div>

      <div class="modal-footer">
        <button
          type="button"
          class="confirm-btn"
          onclick={confirmNewProduct}
          disabled={!newName.trim()}
        >
          Add Product
        </button>
      </div>
    </div>
  </div>
{/if}

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

  /* ── Dropdown ──────────────────────────────────────────────── */
  .dropdown {
    position: absolute;
    left: 0;
    width: 100%;
    z-index: 10;
    max-height: 280px;
    overflow-y: auto;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-lg);
    margin-top: var(--space-3);
    padding: var(--space-2);
    box-shadow: var(--shadow-md);
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

  /* ── Modal overlay ─────────────────────────────────────────── */
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: color-mix(in srgb, var(--color-dark-blue) 50%, transparent);
    z-index: 100;
    display: flex;
    align-items: flex-start;
    justify-content: center;
  }

  /* ── Modal sheet (drops from top) ──────────────────────────── */
  .modal-sheet {
    background: var(--color-surface);
    border-radius: 0 0 var(--radius-lg) var(--radius-lg);
    width: 100%;
    box-shadow: var(--shadow-lg);
    padding: var(--space-6);
    padding-top: max(var(--space-6), env(safe-area-inset-top, 0px));
    display: flex;
    flex-direction: column;
    gap: var(--space-6);
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .modal-title {
    font-size: var(--text-lg);
    font-weight: 600;
    color: var(--color-text);
    margin: 0;
  }

  .modal-close {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: var(--radius-xs);
    background: none;
    border: none;
    color: var(--color-text-secondary);
    cursor: pointer;
    padding: 0;
    transition: background var(--transition-fast), color var(--transition-fast);
  }

  .modal-close:hover {
    background: var(--color-bg);
    color: var(--color-text);
  }

  /* ── Modal form ────────────────────────────────────────────── */
  .modal-body {
    display: flex;
    flex-direction: column;
    gap: var(--form-gap);
  }

  .form-field {
    display: flex;
    flex-direction: column;
    gap: var(--label-gap);
  }

  .form-label {
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--color-text-secondary);
  }

  .form-input {
    height: var(--input-height);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    padding: var(--input-padding);
    background: var(--color-surface);
    color: var(--color-text);
    font-size: var(--text-base);
    font-family: inherit;
    width: 100%;
    box-sizing: border-box;
    outline: none;
    transition: border-color var(--transition-base), box-shadow var(--transition-base);
  }

  .form-input::placeholder {
    color: var(--color-disabled);
  }

  .form-input:focus {
    border-color: var(--color-blue);
    box-shadow: var(--focus-ring);
  }

  /* ── Modal footer ──────────────────────────────────────────── */
  .modal-footer {
    padding-bottom: env(safe-area-inset-bottom, 0px);
  }

  .confirm-btn {
    width: 100%;
    height: var(--input-height);
    border: none;
    border-radius: var(--radius-md);
    background: var(--color-blue);
    color: #fff;
    font-size: var(--text-base);
    font-weight: 600;
    font-family: inherit;
    cursor: pointer;
    transition: opacity var(--transition-fast);
  }

  .confirm-btn:hover:not(:disabled) {
    opacity: 0.88;
  }

  .confirm-btn:disabled {
    opacity: 0.38;
    cursor: not-allowed;
  }
</style>