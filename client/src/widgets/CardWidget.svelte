<script lang="ts">
</script>

<div class="card">
    <slot/>
</div>

<style>
    .card {
        background: var(--color-surface);
        border-radius: var(--radius-lg);
        padding: var(--card-padding);
        border: 1px solid var(--color-border);
        box-shadow: var(--shadow-sm);
    }
</style>