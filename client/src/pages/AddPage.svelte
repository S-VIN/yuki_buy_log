<script lang="ts">
  import { untrack } from 'svelte';
  import ProductSelectWidget from '../widgets/ProductSelectWidget.svelte';
  import TagWidget from '../widgets/TagWidget.svelte';
  import SelectWidget from '../widgets/SelectWidget.svelte';
  import PurchaseCardWidget from '../widgets/PurchaseCardWidget.svelte';
  import { productStore } from '../stores/products.svelte';
  import { purchaseStore } from '../stores/purchases.svelte';
  import type { Product } from '../models/Product';
  import type { Purchase } from '../models/Purchase';

  interface PendingPurchase {
    uuid: string;
    product: Product;
    price: number;
    quantity: number;
    tags: string[];
  }

  // ─── Receipt-level state ──────────────────────────────────────
  let selectedDate = $state(new Date().toISOString().slice(0, 10));
  let selectedShop = $state<string | null>(null);
  let allShops = $state<string[]>([]);

  // ─── Item-level state ─────────────────────────────────────────
  let selectedProduct = $state<Product | null>(null);
  let price = $state('');
  let quantity = $state('1');
  let selectedTags = $state<string[]>([]);

  // ─── Tags pool (merged from store + user additions) ───────────
  let allTags = $state<string[]>([]);
  $effect(() => {
    const storeTags = productStore.tags;
    const current = untrack(() => allTags);
    const newOnes = storeTags.filter((t) => !current.includes(t));
    if (newOnes.length > 0) allTags = [...current, ...newOnes];
  });

  // ─── Auto-fill tags when product changes ─────────────────────
  // _prevProductId is plain JS (not $state) — not a reactive dependency
  let _prevProductId = '';
  $effect(() => {
    const id = selectedProduct?.id ?? '';
    if (id !== _prevProductId) {
      _prevProductId = id;
      if (selectedProduct) {
        selectedTags = [...(selectedProduct.default_tags ?? [])];
      }
    }
  });

  // ─── Check cache ──────────────────────────────────────────────
  let pendingPurchases = $state<PendingPurchase[]>([]);

  const canAdd = $derived(!!selectedProduct && !!price && parseFloat(price) > 0);
  const canClose = $derived(pendingPurchases.length > 0);

  function handleAddPurchase() {
    if (!canAdd || !selectedProduct) return;
    const priceVal = parseFloat(price);
    const quantityVal = parseInt(quantity, 10) || 1;

    pendingPurchases = [
      ...pendingPurchases,
      {
        uuid: crypto.randomUUID(),
        product: selectedProduct,
        price: priceVal,
        quantity: quantityVal,
        tags: [...selectedTags],
      },
    ];

    selectedProduct = null;
    price = '';
    quantity = '1';
    selectedTags = [];
  }

  let isSubmitting = $state(false);

  async function handleCloseCheck() {
    if (!canClose || isSubmitting) return;
    isSubmitting = true;
    try {
      const receiptId = Math.floor(Date.now() / 1000);
      for (const pending of pendingPurchases) {
        await purchaseStore.create({
          product_id: pending.product.id,
          price: pending.price,
          quantity: pending.quantity,
          date: new Date(selectedDate).toISOString(),
          store: selectedShop ?? '',
          tags: pending.tags,
          receipt_id: receiptId,
        } as Omit<Purchase, 'id'>);
      }
      pendingPurchases = [];
      selectedDate = new Date().toISOString().slice(0, 10);
      selectedShop = null;
    } finally {
      isSubmitting = false;
    }
  }

  function deletePurchase(uuid: string) {
    pendingPurchases = pendingPurchases.filter((p) => p.uuid !== uuid);
  }

  function editPurchase(pending: PendingPurchase) {
    _prevProductId = pending.product.id;
    selectedProduct = pending.product;
    price = String(pending.price);
    quantity = String(pending.quantity);
    selectedTags = [...pending.tags];
    deletePurchase(pending.uuid);
  }
</script>

<div class="page">
  <div class="controls-header">
    <div class="form">
      <div class="two-col">
        <div class="col">
          <input
            id="date-input"
            type="date"
            class="native-input"
            bind:value={selectedDate}
          />
        </div>
        <div class="col">
          <SelectWidget
            id="shop-select"
            bind:allOptions={allShops}
            bind:value={selectedShop}
            color="var(--color-yellow)"
            placeholder="Shop…"
          />
        </div>
      </div>

      <div class="two-col">
        <div class="col">
          <input
            id="price-input"
            type="number"
            inputmode="decimal"
            min="0"
            step="0.01"
            class="native-input"
            placeholder="Price"
            bind:value={price}
          />
        </div>
        <div class="col">
          <input
            id="qty-input"
            type="number"
            inputmode="numeric"
            min="1"
            step="1"
            class="native-input"
            placeholder="Qty"
            bind:value={quantity}
          />
        </div>
      </div>

      <ProductSelectWidget id="product-select" bind:value={selectedProduct} />

      <TagWidget id="tags-input" bind:allTags bind:selectedTags />

      <div class="actions-row">
        <button
          type="button"
          class="btn btn-primary"
          onclick={handleCloseCheck}
          disabled={!canClose || isSubmitting}
        >
          {isSubmitting ? 'Saving…' : 'Close check'}
        </button>
        <button
          type="button"
          class="btn btn-secondary"
          onclick={handleAddPurchase}
          disabled={!canAdd}
        >
          Add purchase
        </button>
      </div>
    </div>
  </div>

  <div class="list-area">
    {#if pendingPurchases.length === 0}
      <div class="empty-state">No items yet</div>
    {:else}
      <div class="purchase-list">
        {#each pendingPurchases as pending (pending.uuid)}
          <PurchaseCardWidget
            product={pending.product}
            price={pending.price}
            quantity={pending.quantity}
            tags={pending.tags}
            onEdit={() => editPurchase(pending)}
            onDelete={() => deletePurchase(pending.uuid)}
          />
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .page {
    height: 100%;
    display: flex;
    flex-direction: column;
    background: var(--color-bg);
    overflow: hidden;
  }

  /* ─── Controls header ────────────────────────────────────────── */
  .controls-header {
    flex-shrink: 0;
    background: var(--color-surface);
    border-bottom: 1px solid var(--color-border);
    box-shadow: var(--shadow-sm);
    padding: var(--space-3);
    z-index: 1;
  }

  .form {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
  }

  .two-col {
    display: flex;
    gap: var(--space-3);
  }

  .col {
    flex: 1;
    min-width: 0;
  }

  /* ─── Native inputs ──────────────────────────────────────────── */
  .native-input {
    height: var(--input-height);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    padding: var(--input-padding);
    font-size: var(--text-base);
    color: var(--color-text);
    background: var(--color-surface);
    outline: none;
    width: 100%;
    box-sizing: border-box;
    transition: border-color var(--transition-base), box-shadow var(--transition-base);
  }

  .native-input:focus {
    border-color: var(--color-blue);
    box-shadow: var(--focus-ring);
  }

  .native-input::placeholder {
    color: var(--color-disabled);
  }

  .native-input[type='number']::-webkit-inner-spin-button,
  .native-input[type='number']::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }

  .native-input[type='number'] {
    -moz-appearance: textfield;
    appearance: textfield;
  }

  /* ─── Action buttons ─────────────────────────────────────────── */
  .actions-row {
    display: flex;
    gap: var(--space-3);
  }

  .btn {
    flex: 1;
    height: var(--input-height);
    border-radius: var(--radius-md);
    font-size: var(--text-base);
    font-weight: 500;
    cursor: pointer;
    border: 1px solid;
    transition: background var(--transition-fast), border-color var(--transition-fast),
      opacity var(--transition-fast);
  }

  .btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .btn-primary {
    background: var(--color-blue);
    border-color: var(--color-blue);
    color: #fff;
  }

  .btn-primary:not(:disabled):hover {
    background: color-mix(in srgb, var(--color-blue) 85%, black);
    border-color: color-mix(in srgb, var(--color-blue) 85%, black);
  }

  .btn-secondary {
    background: color-mix(in srgb, var(--color-blue) 8%, transparent);
    border-color: color-mix(in srgb, var(--color-blue) 30%, transparent);
    color: var(--color-blue);
  }

  .btn-secondary:not(:disabled):hover {
    background: color-mix(in srgb, var(--color-blue) 15%, transparent);
    border-color: color-mix(in srgb, var(--color-blue) 50%, transparent);
  }

  /* ─── List area ──────────────────────────────────────────────── */
  .list-area {
    flex: 1;
    overflow-y: auto;
    scrollbar-width: none;
    padding: var(--space-4);
  }

  .list-area::-webkit-scrollbar {
    display: none;
  }

  .purchase-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
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