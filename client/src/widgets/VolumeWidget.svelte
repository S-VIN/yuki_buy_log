<script lang="ts">
  import { untrack } from 'svelte';
  import { Package } from 'lucide-svelte';
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
    placeholder = 'Volume…',
  }: Props = $props();

  let allOptions = $state<string[]>([]);

  $effect(() => {
    const storeVolumes = productStore.volumes;
    const current = untrack(() => allOptions);
    const newOnes = storeVolumes.filter((v) => !current.includes(v));
    if (newOnes.length > 0) allOptions = [...current, ...newOnes];
  });
</script>

<SelectWidget
  {id}
  {placeholder}
  bind:allOptions
  bind:value
  color="var(--color-green)"
>
  {#snippet icon()}
    <Package size={14} />
  {/snippet}
</SelectWidget>