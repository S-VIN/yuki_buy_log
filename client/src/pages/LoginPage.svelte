<script lang="ts">
  import { auth } from '../lib/auth.svelte';
  import { apiLogin, apiRegister } from '../lib/api';

  import CardWidget from "../widgets/CardWidget.svelte";

  let activeTab = $state<'login' | 'register'>('login');
  let login = $state('');
  let password = $state('');
  let error = $state('');
  let loading = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    error = '';
    loading = true;

    try {
      const token = activeTab === 'login'
        ? await apiLogin(login, password)
        : await apiRegister(login, password);
      auth.login(token);
    } catch {
      error = activeTab === 'login' ? 'Login failed' : 'Registration failed';
    } finally {
      loading = false;
    }
  }

  function switchTab(tab: 'login' | 'register') {
    activeTab = tab;
    error = '';
  }
</script>

<div class="login-container">
  <CardWidget>
    <div class="tabs">
      <button
        class="tab"
        class:active={activeTab === 'login'}
        onclick={() => switchTab('login')}
      >
        Sign In
      </button>
      <button
        class="tab"
        class:active={activeTab === 'register'}
        onclick={() => switchTab('register')}
      >
        Sign Up
      </button>
    </div>

    <form onsubmit={handleSubmit}>
      <input
        type="text"
        placeholder="Login"
        bind:value={login}
        required
        autocomplete="username"
      />
      <input
        type="password"
        placeholder="Password"
        bind:value={password}
        required
        autocomplete={activeTab === 'login' ? 'current-password' : 'new-password'}
      />

      {#if error}
        <p class="error">{error}</p>
      {/if}

      <button type="submit" class="submit-btn" disabled={loading}>
        {#if loading}
          ...
        {:else}
          {activeTab === 'login' ? 'Sign In' : 'Sign Up'}
        {/if}
      </button>
    </form>
  </CardWidget>
</div>

<style>
  .login-container {
    width: 100%;
    height: 100%;
    padding: var(--page-padding);
    display: flex;
    flex-direction: column;
    align-items: stretch;
    justify-content: start;
  }

  .tabs {
    display: flex;
    margin-bottom: var(--space-6);
    border-bottom: 1px solid var(--color-border);
  }

  .tab {
    all: unset;
    flex: 1;
    text-align: center;
    padding: var(--space-4) 0;
    cursor: pointer;
    color: var(--color-text-secondary);
    font-size: var(--text-base);
    font-weight: 500;
    transition: color var(--transition-fast), border-color var(--transition-fast);
    border-bottom: 2px solid transparent;
    margin-bottom: -1px;
  }

  .tab.active {
    color: var(--color-blue);
    border-bottom-color: var(--color-blue);
  }

  form {
    display: flex;
    flex-direction: column;
    gap: var(--form-gap);
  }

  input {
    width: 100%;
    padding: var(--input-padding);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    font-size: var(--text-base);
    background: var(--color-surface);
    color: var(--color-text);
    outline: none;
    transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
  }

  input:focus {
    border-color: var(--color-blue);
    box-shadow: var(--focus-ring);
  }

  input::placeholder {
    color: var(--color-disabled);
  }

  .submit-btn {
    width: 100%;
    height: var(--input-height);
    padding: 0 var(--space-5);
    border: none;
    border-radius: var(--radius-md);
    background: var(--color-blue);
    color: #fff;
    font-size: var(--text-base);
    font-weight: 600;
    cursor: pointer;
    transition: opacity var(--transition-fast);
    margin-top: var(--space-2);
  }

  .submit-btn:hover:not(:disabled) {
    opacity: 0.88;
  }

  .submit-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .error {
    color: var(--color-red);
    font-size: var(--text-sm);
    text-align: center;
    margin: 0;
  }
</style>