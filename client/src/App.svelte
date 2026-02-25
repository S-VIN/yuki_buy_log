<script lang="ts">
  import { fade } from 'svelte/transition';
  import { ShoppingCart, Receipt, Package, User } from 'lucide-svelte';
  import { auth } from './lib/auth.svelte';
  import { productStore } from './stores/products.svelte';
  import { purchaseStore } from './stores/purchases.svelte';
  import LoadingScreen from './lib/LoadingScreen.svelte';
  import LoginPage from './pages/LoginPage.svelte';
  import AddPage from './pages/AddPage.svelte';
  import ListPage from './pages/ListPage.svelte';
  import ProductListPage from './pages/ProductListPage.svelte';
  import ProfilePage from './pages/ProfilePage.svelte';

  type TabId = 'add' | 'list' | 'products' | 'profile';

  let activeTab = $state<TabId>('add');

  const menuItems = [
    { id: 'add', icon: ShoppingCart },
    { id: 'list', icon: Receipt },
    { id: 'products', icon: Package },
    { id: 'profile', icon: User },
  ] as const;

  export function navigateTo(tab: TabId) {
    activeTab = tab;
  }

  let isLoading = $state(auth.isAuthenticated);

  $effect(() => {
    if (auth.isAuthenticated) {
      isLoading = true;
      Promise.all([
        productStore.load(),
        purchaseStore.load()
      ])
        .catch(console.error)
        .finally(() => {
          isLoading = false;
        });
    }
  });
</script>

{#if auth.isAuthenticated}
  <div class="app-shell">
    <main class="content">
      <div class="tab-panel" hidden={activeTab !== 'add'}>
        <AddPage />
      </div>
      <div class="tab-panel" hidden={activeTab !== 'list'}>
        <ListPage />
      </div>
      <div class="tab-panel" hidden={activeTab !== 'products'}>
        <ProductListPage />
      </div>
      <div class="tab-panel" hidden={activeTab !== 'profile'}>
        <ProfilePage />
      </div>
    </main>

    <nav class="bottom-nav">
      {#each menuItems as item}
        <button
          class="nav-item"
          role="tab"
          aria-selected={activeTab === item.id}
          data-active={activeTab === item.id || undefined}
          onclick={() => (activeTab = item.id)}
        >
          <item.icon size={24} strokeWidth={1.5} />
        </button>
      {/each}
    </nav>
  </div>
{:else}
  <LoginPage />
{/if}

{#if isLoading}
  <div transition:fade={{ duration: 280 }}>
    <LoadingScreen />
  </div>
{/if}

<style>
  .app-shell {
    display: flex;
    flex-direction: column;
    /* --app-height зафиксирована при загрузке в main.ts — клавиатура
       не вызывает пересчёт высоты для position: fixed элемента */
    height: var(--app-height, 100%);
    width: 100%;
    overflow: hidden;
    position: fixed;
    top: 0;
    left: 0;
  }

  .content {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    -webkit-overflow-scrolling: touch;
  }

  .content > .tab-panel {
    height: 100%;
  }

  .content > [hidden] {
    display: none;
  }

  .bottom-nav {
    display: flex;
    align-items: flex-start;
    justify-content: space-around;
    height: calc(var(--nav-height) + env(safe-area-inset-bottom, 0px));
    min-height: calc(var(--nav-height) + env(safe-area-inset-bottom, 0px));
    background: var(--color-surface);
    border-top: 1px solid var(--color-border);
    padding-bottom: env(safe-area-inset-bottom, 0px);
    user-select: none;
    -webkit-user-select: none;
    -webkit-tap-highlight-color: transparent;
  }

  .nav-item {
    all: unset;
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    flex: 1;
    height: var(--nav-height);
    cursor: pointer;
    color: var(--color-disabled);
    transition: color var(--transition-fast);
    -webkit-tap-highlight-color: transparent;
  }

  .nav-item[data-active] {
    color: var(--color-blue);
  }

  .nav-item[data-active]::after {
    content: '';
    position: absolute;
    top: 0;
    left: 50%;
    transform: translateX(-50%);
    width: 24px;
    height: 2px;
    background: var(--color-blue);
    border-radius: var(--radius-full);
  }

  .nav-item:focus-visible {
    outline: none;
    box-shadow: var(--focus-ring);
    border-radius: var(--radius-xs);
  }
</style>