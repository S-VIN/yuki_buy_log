<script lang="ts">
  import { Tabs } from 'melt/builders';
  import { Home, Search, PlusCircle, List, User } from 'lucide-svelte';
  import { auth } from './lib/auth.svelte';
  import LoginPage from './pages/LoginPage.svelte';
  import HomePage from './pages/HomePage.svelte';
  import SearchPage from './pages/SearchPage.svelte';
  import AddPage from './pages/AddPage.svelte';
  import ListPage from './pages/ListPage.svelte';
  import ProfilePage from './pages/ProfilePage.svelte';

  type TabId = 'home' | 'search' | 'add' | 'list' | 'profile';

  const tabs = new Tabs<TabId>({
    value: 'home',
    orientation: 'horizontal',
    selectWhenFocused: false,
    loop: true,
  });

  const menuItems = [
    { id: 'home', icon: Home },
    { id: 'search', icon: Search },
    { id: 'add', icon: PlusCircle },
    { id: 'list', icon: List },
    { id: 'profile', icon: User },
  ] as const;

  export function navigateTo(tab: TabId) {
    tabs.value = tab;
  }
</script>

{#if auth.isAuthenticated}
  <div class="app-shell">
    <main class="content">
      <div {...tabs.getContent('home')}>
        <HomePage />
      </div>
      <div {...tabs.getContent('search')}>
        <SearchPage />
      </div>
      <div {...tabs.getContent('add')}>
        <AddPage />
      </div>
      <div {...tabs.getContent('list')}>
        <ListPage />
      </div>
      <div {...tabs.getContent('profile')}>
        <ProfilePage />
      </div>
    </main>

    <nav class="bottom-nav" {...tabs.triggerList}>
      {#each menuItems as item}
        <button class="nav-item" {...tabs.getTrigger(item.id)}>
          <item.icon size={24} strokeWidth={1.5} />
        </button>
      {/each}
    </nav>
  </div>
{:else}
  <LoginPage />
{/if}

<style>
  .app-shell {
    display: flex;
    flex-direction: column;
    height: 100%;
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

  .content > [data-melt-tabs-content] {
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