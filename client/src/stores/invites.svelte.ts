import { fetchInvites, sendInvite } from '../lib/api';
import type { Invite } from '../models/Invite';

let items = $state<Invite[]>([]);

export const inviteStore = {
  get items() {
    return items;
  },

  async load() {
    const data = await fetchInvites();
    items = ((data as { invites: Invite[] }).invites ?? []) as Invite[];
  },

  async send(login: string) {
    return (await sendInvite(login)) as { message: string };
  },

  async accept(fromLogin: string) {
    return (await sendInvite(fromLogin)) as { message: string };
  },

  clear() {
    items = [];
  },
};
