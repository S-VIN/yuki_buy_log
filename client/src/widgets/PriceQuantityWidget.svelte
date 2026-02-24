<script lang="ts">
  import { X, DollarSign, Hash } from 'lucide-svelte';

  interface Props {
    price: string;
    quantity: string;
    priceId?: string;
    qtyId?: string;
  }

  let { price = $bindable(), quantity = $bindable(), priceId, qtyId }: Props = $props();

  const formula = $derived.by(() => {
    if (!price || !quantity) return '';
    const p = parseFloat(price);
    const q = parseFloat(quantity);
    if (isNaN(p) || isNaN(q) || p <= 0 || q <= 0) return '';
    const total = p * q;
    const totalStr = total % 1 === 0 ? String(total) : total.toFixed(2);
    return `${price} × ${quantity} = ${totalStr}`;
  });
</script>

<div class="widget">
  <div class="inputs-row">
    <div class="col">
      <div class="input-wrap">
        <input
          id={priceId}
          type="number"
          inputmode="decimal"
          min="0"
          step="0.01"
          class="input"
          placeholder="Price"
          bind:value={price}
        />
        {#if price}
          <button
            type="button"
            class="clear-btn"
            onclick={() => { price = ''; }}
            aria-label="Clear price"
          >
            <X size={14} />
          </button>
        {:else}
          <span class="input-icon"><DollarSign size={14} /></span>
        {/if}
      </div>
    </div>
    <div class="col">
      <div class="input-wrap">
        <input
          id={qtyId}
          type="number"
          inputmode="numeric"
          min="1"
          step="1"
          class="input"
          placeholder="Qty"
          bind:value={quantity}
        />
        {#if quantity}
          <button
            type="button"
            class="clear-btn"
            onclick={() => { quantity = ''; }}
            aria-label="Clear quantity"
          >
            <X size={14} />
          </button>
        {:else}
          <span class="input-icon"><Hash size={14} /></span>
        {/if}
      </div>
    </div>
  </div>
  <div class="formula-row">
    <span class="formula">{formula}</span>
  </div>
</div>

<style>
  .widget {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .inputs-row {
    display: flex;
    gap: var(--space-3);
  }

  .col {
    flex: 1;
    min-width: 0;
  }

  .input-wrap {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    height: var(--input-height);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    padding: 0 var(--space-4);
    background: var(--color-surface);
    transition: border-color var(--transition-base), box-shadow var(--transition-base);
  }

  .input-wrap:focus-within {
    border-color: var(--color-blue);
    box-shadow: var(--focus-ring);
  }

  .input {
    flex: 1;
    min-width: 0;
    border: none;
    outline: none;
    padding: 0;
    font-size: var(--text-base);
    color: var(--color-text);
    background: transparent;
  }

  .input::placeholder {
    color: var(--color-disabled);
  }

  .input[type='number']::-webkit-inner-spin-button,
  .input[type='number']::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }

  .input[type='number'] {
    -moz-appearance: textfield;
    appearance: textfield;
  }

  .input-icon {
    display: inline-flex;
    align-items: center;
    color: var(--color-disabled);
    flex-shrink: 0;
    transition: color var(--transition-fast);
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
    transition: opacity var(--transition-fast), background var(--transition-fast),
      color var(--transition-fast);
    flex-shrink: 0;
  }

  .clear-btn:hover {
    opacity: 1;
    background: color-mix(in srgb, var(--color-red) 10%, transparent);
    color: var(--color-red);
  }

  /* Formula row always reserves space to avoid layout shifts */
  .formula-row {
    min-height: calc(var(--text-xs) * 1.5);
    display: flex;
    align-items: center;
    padding: 0 var(--space-2);
  }

  .formula {
    font-size: var(--text-xs);
    color: var(--color-disabled);
    line-height: 1;
  }
</style>
