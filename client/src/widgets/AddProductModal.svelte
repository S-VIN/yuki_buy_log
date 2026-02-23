<script lang="ts">
  import { productStore } from '../stores/products.svelte';
  import SelectWidget from './SelectWidget.svelte';
  import TagWidget from './TagWidget.svelte';
  import type { Product, ProductId } from '../models/Product';

  interface Props {
    open: boolean;
    onClose: () => void;
    mode?: 'add' | 'edit';
    productId?: ProductId | null;
    initialName?: string;
    onAdded?: (product: Product) => void;
  }

  let {
    open,
    onClose,
    mode = 'add',
    productId = null,
    initialName = '',
    onAdded,
  }: Props = $props();

  const isEditMode = $derived(mode === 'edit');

  // Local copies so SelectWidget / TagWidget can extend them in-session
  let localVolumeOptions = $state<string[]>([]);
  let localBrandOptions = $state<string[]>([]);
  let localTagOptions = $state<string[]>([]);

  let newName = $state('');
  let newVolume = $state<string | null>(null);
  let newBrand = $state<string | null>(null);
  let newTags = $state<string[]>([]);
  let isSubmitting = $state(false);

  $effect(() => {
    if (open) {
      isSubmitting = false;
      localVolumeOptions = [...productStore.volumes];
      localBrandOptions = [...productStore.brands];
      localTagOptions = [...productStore.tags];

      if (isEditMode && productId != null) {
        const product = productStore.items.find(p => p.id === productId);
        if (product) {
          newName = product.name;
          newVolume = product.volume || null;
          newBrand = product.brand || null;
          newTags = [...product.default_tags];
        } else {
          newName = '';
          newVolume = null;
          newBrand = null;
          newTags = [];
        }
      } else {
        newName = initialName;
        newVolume = null;
        newBrand = null;
        newTags = [];
      }
    }
  });

  function close() {
    if (isSubmitting) return;
    onClose();
  }

  async function confirm() {
    if (!newName.trim() || isSubmitting) return;
    isSubmitting = true;
    try {
      if (isEditMode && productId != null) {
        const existing = productStore.items.find(p => p.id === productId);
        if (!existing) throw new Error('Product not found');
        await productStore.update({
          ...existing,
          name: newName.trim(),
          volume: newVolume ?? '',
          brand: newBrand ?? '',
          default_tags: [...newTags],
        });
      } else {
        const created = await productStore.create({
          name: newName.trim(),
          volume: newVolume ?? '',
          brand: newBrand ?? '',
          default_tags: [...newTags],
        });
        onAdded?.(created);
      }
      onClose();
    } catch (err) {
      console.error(`Failed to ${isEditMode ? 'update' : 'create'} product:`, err);
    } finally {
      isSubmitting = false;
    }
  }

  $effect(() => {
    if (!open) return;
    function onKeydown(e: KeyboardEvent) {
      if (e.key === 'Escape') close();
    }
    document.addEventListener('keydown', onKeydown);
    return () => document.removeEventListener('keydown', onKeydown);
  });
</script>

{#if open}
  <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_noninteractive_element_interactions -->
  <div
    class="modal-overlay"
    role="dialog"
    aria-modal="true"
    tabindex="-1"
    onclick={(e) => { if (e.target === e.currentTarget) close(); }}
  >
    <div class="modal-sheet">
      <div class="modal-body">
        <input
          type="text"
          class="form-input"
          placeholder="Product name…"
          bind:value={newName}
        />

        <div class="form-row">
          <div class="form-col">
            <SelectWidget
              bind:allOptions={localVolumeOptions}
              bind:value={newVolume}
              color="var(--color-green)"
              placeholder="Volume…"
            />
          </div>
          <div class="form-col">
            <SelectWidget
              bind:allOptions={localBrandOptions}
              bind:value={newBrand}
              color="var(--color-yellow)"
              placeholder="Brand…"
            />
          </div>
        </div>

        <TagWidget
          bind:allTags={localTagOptions}
          bind:selectedTags={newTags}
        />
      </div>

      <div class="modal-footer">
        <button
          type="button"
          class="close-btn"
          onclick={close}
        >
          Close
        </button>
        <button
          type="button"
          class="confirm-btn"
          onclick={confirm}
          disabled={!newName.trim() || isSubmitting}
        >
          {isSubmitting ? (isEditMode ? 'Saving…' : 'Adding…') : (isEditMode ? 'Save' : 'Add Product')}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: color-mix(in srgb, var(--color-dark-blue) 50%, transparent);
    z-index: 100;
    display: flex;
    align-items: flex-start;
    justify-content: center;
  }

  .modal-sheet {
    background: var(--color-surface);
    border-radius: 0 0 var(--radius-lg) var(--radius-lg);
    width: 100%;
    box-shadow: var(--shadow-lg);
    padding: var(--space-5);
    padding-top: max(var(--space-5), env(safe-area-inset-top, 0px));
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
  }

  .modal-body {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
  }

  /* Volume + Brand side by side; flex: 1 1 0 forces equal width regardless of content */
  .form-row {
    display: flex;
    gap: var(--space-3);
  }

  .form-col {
    flex: 1 1 0;
    min-width: 0;
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

  .modal-footer {
    display: flex;
    gap: var(--space-3);
    padding-bottom: env(safe-area-inset-bottom, 0px);
  }

  .close-btn {
    flex: 1 1 0;
    height: var(--input-height);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-surface);
    color: var(--color-text);
    font-size: var(--text-base);
    font-weight: 500;
    font-family: inherit;
    cursor: pointer;
    transition: background var(--transition-fast), color var(--transition-fast);
  }

  .close-btn:hover {
    background: var(--color-bg);
    color: var(--color-text);
  }

  .confirm-btn {
    flex: 1 1 0;
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