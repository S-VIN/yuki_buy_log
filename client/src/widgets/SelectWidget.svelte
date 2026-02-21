<script lang="ts">
  import { Combobox } from "melt/builders";
  import { Plus, X } from "lucide-svelte";

  interface Props {
    allOptions: string[];
    value: string | null;
    color?: string;
    id?: string;
    placeholder?: string;
  }

  let {
    allOptions = $bindable(),
    value = $bindable(null),
    color = "var(--color-blue)",
    id,
    placeholder = "Select…",
  }: Props = $props();

  function selectOption(raw: string) {
    const opt = raw.trim();
    if (!opt) return;
    if (!allOptions.includes(opt)) {
      allOptions = [...allOptions, opt];
    }
    value = opt;
    combobox.inputValue = "";
  }

  function clearSelection() {
    value = null;
    combobox.inputValue = "";
  }

  const combobox = new Combobox<string, false>({
    multiple: false,
    value: () => value ?? undefined,
    onValueChange: (v) => {
      if (v !== undefined) {
        selectOption(v);
      }
    },
  });

  const filteredOptions = $derived(
    allOptions.filter((o) =>
      o.toLowerCase().includes(combobox.inputValue.toLowerCase())
    )
  );

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      const raw = combobox.inputValue.trim();
      if (!raw) return;
      e.preventDefault();
      selectOption(raw);
    }
  }

  let isFocused = $state(false);
</script>

<div class="select-widget">
  <div class="input-row">
    {#if value}
      <div class="value-area">
        <span class="pill" style="--pill-color: {color}">{value}</span>
      </div>
      <button
        type="button"
        class="clear-btn"
        onclick={clearSelection}
        aria-label="Clear selection"
      >
        <X size={14} />
      </button>
    {:else}
      <input
        {...combobox.input}
        {id}
        {placeholder}
        enterkeyhint="done"
        onkeydown={(e) => {
          handleKeydown(e);
          if (!e.defaultPrevented) combobox.input.onkeydown?.(e);
        }}
        onfocus={() => (isFocused = true)}
        onblur={() => (isFocused = false)}
        class="select-input"
      />
    {/if}
  </div>

  {#if isFocused && !value && combobox.open && combobox.inputValue.trim()}
    <div
      class="dropdown"
      role="listbox"
      tabindex="-1"
      onmousedown={(e) => e.preventDefault()}
    >
      {#if combobox.inputValue.trim() && !allOptions.includes(combobox.inputValue.trim())}
        <button
          type="button"
          class="option new-option"
          onclick={() => selectOption(combobox.inputValue)}
        >
          <span class="new-option-icon"><Plus size={12} /></span>
          <span class="new-option-label">Create</span>
          <span class="pill" style="--pill-color: {color}">{combobox.inputValue.trim()}</span>
        </button>
        {#if filteredOptions.length > 0}
          <div class="divider"></div>
        {/if}
      {/if}
      {#each filteredOptions as option}
        <div {...combobox.getOption(option)} class="option">
          {option}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .select-widget {
    position: relative;
  }

  .input-row {
    display: flex;
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

  .input-row:focus-within {
    border-color: var(--color-blue);
    box-shadow: var(--focus-ring);
  }

  .pill {
    display: inline-flex;
    align-items: center;
    background: color-mix(in srgb, var(--pill-color) 10%, transparent);
    color: var(--pill-color);
    border: 1px solid color-mix(in srgb, var(--pill-color) 28%, transparent);
    border-radius: var(--radius-sm);
    padding: 2px 8px;
    font-size: var(--text-sm);
    font-weight: 600;
    letter-spacing: 0.01em;
    white-space: nowrap;
    line-height: 1.5;
  }

  .value-area {
    flex: 1;
    min-width: 0;
    display: flex;
    align-items: center;
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

  .select-input {
    flex: 1;
    min-width: 60px;
    border: none;
    outline: none;
    font-size: var(--text-base);
    background: transparent;
    color: var(--color-text);
    padding: var(--space-1) 0;
  }

  .select-input::placeholder {
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

  .option.new-option {
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

  .new-option-icon {
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

  .new-option-label {
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--color-text-secondary);
  }

  .option.new-option .pill {
    padding: 1px 8px;
  }

  .option.new-option:hover,
  .option.new-option:focus {
    background: var(--color-bg);
    color: var(--color-text);
    outline: none;
  }

  .option.new-option:hover .new-option-label,
  .option.new-option:focus .new-option-label {
    color: var(--color-text);
  }
</style>
