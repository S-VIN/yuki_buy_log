<script lang="ts">
  import TagWidget from "../widgets/TagWidget.svelte";

  let selectedTags: string[] = $state([]);
  const allTags = [
    "продукты",
    "одежда",
    "электроника",
    "еда",
    "напитки",
    "бытовая химия",
  ];
</script>

<div class="page">
  <h2>Add</h2>

  <section>
    <label class="field-label" for="tags-input">Теги</label>
    <TagWidget {allTags} bind:selectedTags id="tags-input" />
  </section>

  {#if selectedTags.length > 0}
    <p class="debug">Выбрано: {selectedTags.join(", ")}</p>
  {/if}
</div>

<style>
  .page {
    padding: var(--page-padding);
  }

  section {
    margin-top: var(--section-gap);
  }

  .field-label {
    display: block;
    font-size: var(--text-sm);
    color: var(--color-text-secondary);
    margin-bottom: var(--label-gap);
  }

  .debug {
    margin-top: var(--section-gap);
    font-size: var(--text-sm);
    color: var(--color-text-secondary);
  }
</style>
