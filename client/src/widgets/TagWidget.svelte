<script lang="ts">
  let { color = '', text = '' }: { color?: string; text?: string } = $props();
</script>

<span class="pill" style:--color={`var(${color})`}>{text}</span>


<style>
    .pill {
        display: inline-flex;
        align-items: center;
        border-radius: var(--radius-sm);
        padding: 2px 8px;
        font-size: var(--text-sm);
        font-weight: 600;
        letter-spacing: 0.01em;
        white-space: nowrap;
        line-height: 1.5;

        background: color-mix(in srgb, var(--color) 15%, transparent);
        color: color-mix(in srgb, var(--color) 80%, var(--color-text));
        border: 1px solid color-mix(in srgb, var(--color) 40%, transparent);
    }
</style>