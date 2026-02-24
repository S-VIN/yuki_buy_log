<script lang="ts">
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
    </div>
    <div class="col">
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

  .input {
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

  .input:focus {
    border-color: var(--color-blue);
    box-shadow: var(--focus-ring);
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
