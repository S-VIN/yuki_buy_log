export type InviteId = number;
export type UserId = number;

export interface Invite {
  id: InviteId;
  from_user_id: UserId;
  to_user_id: UserId;
  from_login: string;
  to_login: string;
  created_at: string;
}
