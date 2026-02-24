<script lang="ts">
  import { auth } from '../lib/auth.svelte';
  import { fetchGroupMembers, leaveGroup as apiLeaveGroup } from '../lib/api';
  import { inviteStore } from '../stores/invites.svelte';

  interface GroupMember {
    group_id: number;
    user_id: number;
    login: string;
    member_number: number;
  }

  const MEMBER_COLORS = [
    '--color-blue',
    '--color-green',
    '--color-yellow',
    '--color-red',
    '--color-dark-blue',
  ] as const;

  function getMemberColor(memberNumber: number): string {
    return `var(${MEMBER_COLORS[(memberNumber - 1) % MEMBER_COLORS.length]})`;
  }

  let members = $state<GroupMember[]>([]);
  let inviteLogin = $state('');
  let loading = $state(false);
  let error = $state<string | null>(null);
  let successMsg = $state<string | null>(null);

  async function loadData() {
    try {
      const [groupData] = await Promise.all([
        fetchGroupMembers(),
        inviteStore.load(),
      ]);
      members = groupData.members ?? [];
    } catch (e) {
      console.error('Failed to load profile data:', e);
    }
  }

  $effect(() => {
    loadData();
  });

  function clearMessages() {
    error = null;
    successMsg = null;
  }

  async function handleSendInvite() {
    if (!inviteLogin.trim()) return;
    loading = true;
    clearMessages();
    try {
      const data = await inviteStore.send(inviteLogin.trim());
      if (data.message === 'group created') {
        successMsg = 'Group created!';
        await loadData();
      } else {
        successMsg = 'Invite sent!';
      }
      inviteLogin = '';
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Failed to send invite';
    } finally {
      loading = false;
    }
  }

  async function handleAcceptInvite(fromLogin: string) {
    loading = true;
    clearMessages();
    try {
      await inviteStore.accept(fromLogin);
      successMsg = 'Invite accepted! Group created.';
      await loadData();
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Failed to accept invite';
    } finally {
      loading = false;
    }
  }

  async function handleLeaveGroup() {
    loading = true;
    clearMessages();
    try {
      await apiLeaveGroup();
      successMsg = 'Left the group';
      await loadData();
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Failed to leave group';
    } finally {
      loading = false;
    }
  }

  function handleLogout() {
    auth.logout();
  }
</script>

<div class="page">
  <div class="card">
    {#if members.length > 0}
      <div class="section">
        <span class="section-title">Group Members</span>
        <div class="tags-row">
          {#each members as member}
            <span
              class="member-tag"
              style="--member-color: {getMemberColor(member.member_number)}"
            >
              {member.login}
            </span>
          {/each}
        </div>
      </div>
    {:else}
      <span class="empty-group">You are not in a group</span>
    {/if}

    {#if inviteStore.items.length > 0}
      <div class="section">
        <span class="section-title">Incoming Invites</span>
        <div class="invites-list">
          {#each inviteStore.items as invite}
            <div class="invite-row">
              <span class="invite-label">
                From: <span class="invite-login">{invite.from_login}</span>
              </span>
              <button
                class="btn-accept"
                onclick={() => handleAcceptInvite(invite.from_login)}
                disabled={loading}
              >
                Accept
              </button>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <div class="section">
      <span class="section-title">Send Invite</span>
      <div class="input-row">
        <input
          class="invite-input"
          type="text"
          placeholder="Enter username"
          bind:value={inviteLogin}
          onkeydown={(e) => e.key === 'Enter' && handleSendInvite()}
          disabled={loading}
        />
        <button
          class="btn-send"
          onclick={handleSendInvite}
          disabled={loading || !inviteLogin.trim()}
        >
          Send
        </button>
      </div>
    </div>

    {#if successMsg}
      <div class="status success">{successMsg}</div>
    {/if}
    {#if error}
      <div class="status error">{error}</div>
    {/if}

    {#if members.length > 0}
      <button class="btn-leave" onclick={handleLeaveGroup} disabled={loading}>
        Leave Group
      </button>
    {/if}
  </div>

  <button class="btn-logout" onclick={handleLogout}>
    Logout
  </button>
</div>

<style>
  .page {
    padding: var(--page-padding);
    display: flex;
    flex-direction: column;
    gap: var(--section-gap);
    height: 100%;
    overflow-y: auto;
    scrollbar-width: none;
  }
  .page::-webkit-scrollbar {
    display: none;
  }

  .card {
    background: var(--color-surface);
    border-radius: var(--radius-lg);
    padding: var(--card-padding);
    border: 1px solid var(--color-border);
    box-shadow: var(--shadow-sm);
    display: flex;
    flex-direction: column;
    gap: var(--section-gap);
  }

  .section {
    display: flex;
    flex-direction: column;
    gap: var(--label-gap);
  }

  .section-title {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text);
  }

  .empty-group {
    font-size: var(--text-sm);
    color: var(--color-text-secondary);
  }

  .tags-row {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-2);
  }

  .member-tag {
    display: inline-block;
    padding: var(--space-1) var(--space-4);
    border-radius: var(--radius-sm);
    font-size: var(--text-sm);
    color: var(--member-color);
    background: color-mix(in srgb, var(--member-color) 12%, transparent);
    border: 1px solid color-mix(in srgb, var(--member-color) 25%, transparent);
  }

  .invites-list {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .invite-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-4);
  }

  .invite-label {
    font-size: var(--text-sm);
    color: var(--color-text-secondary);
  }

  .invite-login {
    font-weight: 500;
    color: var(--color-text);
  }

  .btn-accept {
    all: unset;
    padding: var(--space-2) var(--space-5);
    border-radius: var(--radius-md);
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--color-blue);
    background: color-mix(in srgb, var(--color-blue) 10%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-blue) 25%, transparent);
    cursor: pointer;
    transition: background var(--transition-fast);
    white-space: nowrap;
  }
  .btn-accept:hover:not(:disabled) {
    background: color-mix(in srgb, var(--color-blue) 18%, transparent);
  }
  .btn-accept:disabled {
    color: var(--color-disabled);
    border-color: var(--color-border);
    background: transparent;
    cursor: not-allowed;
  }

  .input-row {
    display: flex;
    gap: var(--form-gap);
  }

  .invite-input {
    flex: 1;
    height: var(--input-height);
    padding: var(--input-padding);
    border-radius: var(--radius-md);
    border: 1px solid var(--color-border);
    background: var(--color-bg);
    color: var(--color-text);
    font-size: var(--text-base);
    outline: none;
    transition:
      border-color var(--transition-fast),
      box-shadow var(--transition-base);
  }
  .invite-input:focus {
    border-color: var(--color-blue);
    box-shadow: var(--focus-ring);
  }
  .invite-input:disabled {
    color: var(--color-disabled);
    border-color: var(--color-border);
  }

  .btn-send {
    all: unset;
    height: var(--input-height);
    padding: 0 var(--space-6);
    border-radius: var(--radius-md);
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--color-surface);
    background: var(--color-blue);
    cursor: pointer;
    transition: background var(--transition-fast);
    white-space: nowrap;
  }
  .btn-send:hover:not(:disabled) {
    background: color-mix(in srgb, var(--color-blue) 85%, black);
  }
  .btn-send:disabled {
    background: var(--color-disabled);
    cursor: not-allowed;
  }

  .status {
    font-size: var(--text-sm);
    padding: var(--space-2) var(--space-4);
    border-radius: var(--radius-sm);
  }
  .status.success {
    color: var(--color-green);
    background: color-mix(in srgb, var(--color-green) 10%, transparent);
  }
  .status.error {
    color: var(--color-red);
    background: color-mix(in srgb, var(--color-red) 10%, transparent);
  }

  .btn-leave {
    all: unset;
    width: 100%;
    box-sizing: border-box;
    height: var(--input-height);
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-md);
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--color-red);
    background: color-mix(in srgb, var(--color-red) 10%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-red) 25%, transparent);
    cursor: pointer;
    transition: background var(--transition-fast);
  }
  .btn-leave:hover:not(:disabled) {
    background: color-mix(in srgb, var(--color-red) 18%, transparent);
  }
  .btn-leave:disabled {
    color: var(--color-disabled);
    border-color: var(--color-border);
    background: transparent;
    cursor: not-allowed;
  }

  .btn-logout {
    all: unset;
    width: 100%;
    box-sizing: border-box;
    height: var(--input-height);
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-md);
    font-size: var(--text-base);
    font-weight: 600;
    color: var(--color-surface);
    background: var(--color-red);
    cursor: pointer;
    transition: background var(--transition-fast);
  }
  .btn-logout:hover {
    background: color-mix(in srgb, var(--color-red) 85%, black);
  }
  .btn-logout:active {
    background: color-mix(in srgb, var(--color-red) 75%, black);
  }
</style>
