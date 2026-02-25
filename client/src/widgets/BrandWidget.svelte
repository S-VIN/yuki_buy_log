<script lang="ts">
  import { untrack } from 'svelte';
  import { Tag } from 'lucide-svelte';
  import { productStore } from '../stores/products.svelte';
  import SelectWidget from './SelectWidget.svelte';

  interface Props {
    value: string | null;
    id?: string;
    placeholder?: string;
  }

  let {
    value = $bindable(null),
    id,
    placeholder = 'Brand…',
  }: Props = $props();

  let allOptions = $state<string[]>([]);

  $effect(() => {
    const storeBrands = productStore.brands;
    const current = untrack(() => allOptions);
    const newOnes = storeBrands.filter((b) => !current.includes(b));
    if (newOnes.length > 0) allOptions = [...current, ...newOnes];
  });
</script>

<SelectWidget
  {id}
  {placeholder}
  bind:allOptions
  bind:value
  color="--color-yellow"
>
  {#snippet icon()}
    <Tag size={14} />
  {/snippet}
</SelectWidget>