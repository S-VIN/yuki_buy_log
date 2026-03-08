export type InviteId = bigint;
export type UserId = bigint;

export interface Invite {
  id: InviteId;
  from_user_id: UserId;
  to_user_id: UserId;
  from_login: string;
  to_login: string;
  created_at: Date;
}
