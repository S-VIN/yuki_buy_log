<script lang="ts">
  import { Plus, X, Tags } from "lucide-svelte";

  interface Props {
    allTags: string[];
    selectedTags: string[];
    id?: string;
  }

  let { allTags = $bindable(), selectedTags = $bindable(), id }: Props = $props();

  let inputValue = $state('');
  let isFocused = $state(false);

  function addTag(raw: string) {
    const tag = raw.trim();
    if (!tag) return;
    if (!allTags.includes(tag)) {
      allTags = [...allTags, tag];
    }
    if (!selectedTags.includes(tag)) {
      selectedTags = [...selectedTags, tag];
    }
    inputValue = "";
  }

  const filteredTags = $derived(
    allTags.filter(
      (t) =>
        !selectedTags.includes(t) &&
        t.toLowerCase().includes(inputValue.toLowerCase())
    )
  );

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      const raw = inputValue.trim();
      if (!raw) return;
      e.preventDefault();
      addTag(raw);
    }
  }

  function removeTag(tag: string) {
    selectedTags = selectedTags.filter((t) => t !== tag);
  }

  function clearAll() {
    selectedTags = [];
    inputValue = "";
  }
</script>

<div class="tag-widget">
  <div class="input-row">
    <div class="tags-and-input">
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
        {id}
        bind:value={inputValue}
        placeholder={selectedTags.length === 0 ? "Add tags…" : ""}
        enterkeyhint="done"
        onkeydown={handleKeydown}
        onfocus={() => (isFocused = true)}
        onblur={() => (isFocused = false)}
        class="tag-input"
      />
    </div>
    {#if selectedTags.length > 0 || inputValue.trim()}
      <button
        type="button"
        class="clear-btn"
        onclick={clearAll}
        aria-label="Clear all tags"
      >
        <X size={14} />
      </button>
    {:else}
      <span class="input-icon"><Tags size={14} /></span>
    {/if}
  </div>

  {#if isFocused && inputValue.trim()}
    <div class="dropdown" role="listbox" tabindex="-1" onmousedown={(e) => e.preventDefault()}>
      {#if inputValue.trim() && !allTags.includes(inputValue.trim())}
        <button type="button" class="option new-tag" onclick={() => addTag(inputValue)}>
          <span class="new-tag-icon"><Plus size={12} /></span>
          <span class="new-tag-label">Create</span>
          <span class="pill">{inputValue.trim()}</span>
        </button>
        {#if filteredTags.length > 0}
          <div class="divider"></div>
        {/if}
      {/if}
      {#each filteredTags as tag}
        <button type="button" class="option" onclick={() => addTag(tag)}>
          {tag}
        </button>
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
    flex-wrap: nowrap;
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

  .tags-and-input {
    display: flex;
    flex-wrap: wrap;
    flex: 1;
    min-width: 0;
    gap: var(--space-2);
    align-items: center;
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
    transition: opacity var(--transition-fast), background var(--transition-fast), color var(--transition-fast);
    flex-shrink: 0;
  }

  .clear-btn:hover {
    opacity: 1;
    background: color-mix(in srgb, var(--color-red) 10%, transparent);
    color: var(--color-red);
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

  .input-icon {
    display: inline-flex;
    align-items: center;
    color: var(--color-disabled);
    flex-shrink: 0;
    transition: color var(--transition-fast);
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
    width: 100%;
    background: none;
    border: none;
    font: inherit;
    text-align: left;
    padding: var(--space-4) var(--space-5);
    cursor: pointer;
    font-size: var(--text-base);
    color: var(--color-text);
    border-radius: var(--radius-sm);
    transition: background var(--transition-fast);
    user-select: none;
  }

  .option:hover {
    background: var(--color-bg);
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