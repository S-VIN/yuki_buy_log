<script lang="ts">
    const { date, total }: { date: Date; total: number } = $props();

    function formatDate(date: Date): string {
        return date.toLocaleDateString('en-En', { day: 'numeric', month: 'long', year: 'numeric' });
    }

</script>

<div class="date-header">
    <span class="date-label">{formatDate(date)}</span>
    <span class="date-total">{total.toFixed(2)}₽</span>
</div>


<style>
    .date-header {
        display: flex;
        align-items: baseline;
        justify-content: space-between;
        padding: var(--space-4) var(--space-2) var(--space-2);
        margin-top: var(--space-4);
    }

    .date-header:first-child {
        margin-top: 0;
    }

    .date-label {
        font-size: var(--text-sm);
        font-weight: 600;
        color: var(--color-text-secondary);
        text-transform: lowercase;
    }

    .date-total {
        font-size: var(--text-sm);
        font-weight: 700;
        color: var(--color-text-secondary);
    }
</style>