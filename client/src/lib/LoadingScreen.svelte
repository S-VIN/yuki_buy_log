<script lang="ts">
</script>

<div class="screen">
  <div class="center">
    <div class="card">
      <img src="/icon.svg" alt="Yuki Buy Log" class="logo" />
    </div>

    <div class="brand">
      <span class="name">yuki</span>
      <span class="tagline">buy log</span>
    </div>

    <div class="spinner-row">
      <div class="spinner"></div>
    </div>
  </div>
</div>

<style>
  .screen {
    position: fixed;
    inset: 0;
    background: var(--color-bg);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 50;
  }

  .center {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--space-8);
  }

  /* ── Portrait card ── */

  .card {
    width: 110px;
    height: 110px;
    background: var(--color-surface);
    border: 1.5px solid var(--color-border);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-lg);
    overflow: hidden;
    display: flex;
    align-items: flex-end;
    justify-content: center;
    animation: rise 0.5s ease both;
    /* subtle gloss line at the top */
    background-image: linear-gradient(
      to bottom,
      rgba(255, 255, 255, 0.7) 0px,
      rgba(255, 255, 255, 0) 32px
    );
  }

  .logo {
    width: 100px;
    height: 100px;
    object-fit: contain;
    display: block;
    flex-shrink: 0;
  }

  /* ── Brand text ── */

  .brand {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--space-2);
    animation: rise 0.5s ease 0.1s both;
  }

  .name {
    font-size: 1.75rem;
    font-weight: 700;
    color: var(--color-dark-blue);
    letter-spacing: -0.03em;
    line-height: 1;
  }

  .tagline {
    font-size: var(--text-sm);
    color: var(--color-text-secondary);
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  /* ── Spinner ── */

  .spinner-row {
    animation: rise 0.5s ease 0.2s both;
  }

  .spinner {
    width: 24px;
    height: 24px;
    border-radius: var(--radius-full);
    border: 2px solid color-mix(in srgb, var(--color-blue) 22%, transparent);
    border-top-color: var(--color-blue);
    animation: spin 0.8s linear infinite;
  }

  /* ── Keyframes ── */

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  @keyframes rise {
    from {
      opacity: 0;
      transform: translateY(16px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>