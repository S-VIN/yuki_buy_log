<script lang="ts">
  import { Combobox } from "melt/builders";
  import { Plus } from "lucide-svelte";

  interface Props {
    allTags: string[];
    selectedTags: string[];
    id?: string;
  }

  let { allTags = $bindable(), selectedTags = $bindable(), id }: Props = $props();

  function addNewTag(raw: string) {
    const tag = raw.trim();
    if (!tag) return;
    if (!allTags.includes(tag)) {
      allTags = [...allTags, tag];
    }
    if (!selectedTags.includes(tag)) {
      selectedTags = [...selectedTags, tag];
    }
    combobox.inputValue = "";
  }

  const combobox = new Combobox<string, true>({
    multiple: true,
    value: () => selectedTags,
    onValueChange: (v) => {
      selectedTags = [...v];
      combobox.inputValue = "";
    },
  });

  const filteredTags = $derived(
    allTags.filter(
      (t) =>
        !selectedTags.includes(t) &&
        t.toLowerCase().includes(combobox.inputValue.toLowerCase())
    )
  );

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      const raw = combobox.inputValue.trim();
      if (!raw || allTags.includes(raw)) return;
      e.preventDefault();
      addNewTag(raw);
    }
  }

  function removeTag(tag: string) {
    selectedTags = selectedTags.filter((t) => t !== tag);
  }

  let isFocused = $state(false);
</script>

<div class="tag-widget">
  <div class="input-row">
    {#each selectedTags as tag}
      <span class="pill">
        {tag}
        <button type="button" class="pill-remove" onclick={() => removeTag(tag)} aria-label="Remove {tag}">
          <svg width="10" height="10" viewBox="0 0 10 10" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M1 1L9 9M9 1L1 9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
      </span>
    {/each}
    <input
      {...combobox.input}
      {id}
      placeholder={selectedTags.length === 0 ? "Add tags…" : ""}
      enterkeyhint="done"
      onkeydown={(e) => { handleKeydown(e); if (!e.defaultPrevented) combobox.input.onkeydown?.(e); }}
      onfocus={() => (isFocused = true)}
      onblur={() => (isFocused = false)}
      class="tag-input"
    />
  </div>

  {#if isFocused && combobox.open && (filteredTags.length > 0 || combobox.inputValue.trim())}
    <div class="dropdown" role="listbox" tabindex="-1" onmousedown={(e) => e.preventDefault()}>
      {#if combobox.inputValue.trim() && !allTags.includes(combobox.inputValue.trim())}
        <button type="button" class="option new-tag" onclick={() => addNewTag(combobox.inputValue)}>
          <span class="new-tag-icon"><Plus size={12} /></span>
          <span class="new-tag-label">Create</span>
          <span class="pill">{combobox.inputValue.trim()}</span>
        </button>
        {#if filteredTags.length > 0}
          <div class="divider"></div>
        {/if}
      {/if}
      {#each filteredTags as tag}
        <div {...combobox.getOption(tag)} class="option">
          {tag}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .tag-widget {
    position: relative;
  }

  .input-row {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-2);
    align-items: center;
    min-height: var(--input-height);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    padding: var(--space-2) var(--space-4);
    background: var(--color-surface);
    transition: border-color var(--transition-base), box-shadow var(--transition-base);
    cursor: text;
  }

  .input-row:focus-within {
    border-color: var(--color-blue);
    box-shadow: var(--focus-ring);
  }

  .pill {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    background: color-mix(in srgb, var(--color-blue) 10%, transparent);
    color: var(--color-blue);
    border: 1px solid color-mix(in srgb, var(--color-blue) 28%, transparent);
    border-radius: var(--radius-sm);
    padding: 2px 5px 2px 8px;
    font-size: var(--text-sm);
    font-weight: 600;
    letter-spacing: 0.01em;
    white-space: nowrap;
    line-height: 1.5;
  }

  .pill-remove {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 16px;
    height: 16px;
    border-radius: var(--radius-xs);
    background: none;
    border: none;
    color: var(--color-blue);
    cursor: pointer;
    padding: 0;
    opacity: 0.5;
    transition: opacity var(--transition-fast), background var(--transition-fast);
    flex-shrink: 0;
  }

  .pill-remove:hover {
    opacity: 1;
    background: color-mix(in srgb, var(--color-blue) 18%, transparent);
  }

  .tag-input {
    flex: 1;
    min-width: 60px;
    border: none;
    outline: none;
    font-size: var(--text-base);
    background: transparent;
    color: var(--color-text);
    padding: var(--space-1) 0;
  }

  .tag-input::placeholder {
    color: var(--color-disabled);
  }

  .dropdown {
    position: absolute;
    left: 0;
    width: 100%;
    z-index: 10;
    max-height: 220px;
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
  }

  .option[data-highlighted] {
    background: var(--color-bg);
  }

  .option[aria-selected="true"] {
    font-weight: 600;
    color: var(--color-blue);
  }

  .option.new-tag {
    gap: 6px;
    color: var(--color-text-secondary);
    background: none;
    border: none;
    font: inherit;
    text-align: left;
    width: 100%;
    cursor: pointer;
    font-size: var(--text-base);
  }

  .new-tag-icon {
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

  .new-tag-label {
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--color-text-secondary);
  }

  .option.new-tag .pill {
    padding: 1px 8px;
  }

  .option.new-tag:hover,
  .option.new-tag:focus {
    background: var(--color-bg);
    color: var(--color-text);
    outline: none;
  }

  .option.new-tag:hover .new-tag-label,
  .option.new-tag:focus .new-tag-label {
    color: var(--color-text);
  }
</style>