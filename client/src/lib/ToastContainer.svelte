<script lang="ts">
  import { fly } from 'svelte/transition';
  import { AlertCircle, X } from 'lucide-svelte';
  import { toastStore } from './toast.svelte';
</script>

<div class="toast-container">
  {#each toastStore.toasts as toast (toast.id)}
    <div class="toast" transition:fly={{ x: 24, duration: 200 }}>
      <div class="toast-content">
        <span class="toast-icon"><AlertCircle size={16} color="var(--color-red)" /></span>
        <span class="toast-message">{toast.message}</span>
        <button class="toast-dismiss" onclick={() => toastStore.dismiss(toast.id)} aria-label="Dismiss">
          <X size={14} />
        </button>
      </div>
      <div class="toast-progress"></div>
    </div>
  {/each}
</div>

<style>
  .toast-container {
    position: fixed;
    top: var(--space-5);
    right: var(--space-5);
    z-index: 9999;
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
    align-items: flex-end;
    pointer-events: none;
  }

  .toast {
    pointer-events: auto;
    background: color-mix(in srgb, var(--color-red) 6%, var(--color-surface));
    border: 1px solid color-mix(in srgb, var(--color-red) 20%, transparent);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-md);
    max-width: 300px;
    min-width: 220px;
    overflow: hidden;
  }

  .toast-content {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-4) var(--space-5);
  }

  .toast-icon {
    flex-shrink: 0;
  }

  .toast-message {
    flex: 1;
    font-size: var(--text-sm);
    color: var(--color-text);
    line-height: 1.4;
  }

  .toast-dismiss {
    all: unset;
    cursor: pointer;
    color: var(--color-disabled);
    flex-shrink: 0;
    transition: color var(--transition-fast);
    display: flex;
    align-items: center;
  }

  .toast-dismiss:hover {
    color: var(--color-text-secondary);
  }

  .toast-progress {
    height: 2px;
    width: 100%;
    background: color-mix(in srgb, var(--color-red) 30%, transparent);
    animation: shrink 3s linear forwards;
  }

  @keyframes shrink {
    from {
      width: 100%;
    }
    to {
      width: 0%;
    }
  }
</style>