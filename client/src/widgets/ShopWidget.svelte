<script lang="ts">
  import { untrack } from 'svelte';
  import { ShoppingBag } from 'lucide-svelte';
  import { purchaseStore } from '../stores/purchases.svelte';
  import SelectWidget from './SelectWidget.svelte';

  interface Props {
    value: string | null;
    id?: string;
    placeholder?: string;
  }

  let {
    value = $bindable(null),
    id,
    placeholder = 'Shop…',
  }: Props = $props();

  let allOptions = $state<string[]>([]);

  $effect(() => {
    const storeShops = purchaseStore.shops;
    const current = untrack(() => allOptions);
    const newOnes = storeShops.filter((s) => !current.includes(s));
    if (newOnes.length > 0) allOptions = [...current, ...newOnes];
  });
</script>

<SelectWidget
  {id}
  {placeholder}
  bind:allOptions
  bind:value
  color="var(--color-yellow)"
>
  {#snippet icon()}
    <ShoppingBag size={14} />
  {/snippet}
</SelectWidget>