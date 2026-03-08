import type {UserId} from "./Invite";

export type ReceiptId = number;

export interface Receipt {
  id: ReceiptId;
  date: Date;
  store: string;
  purchase_ids: bigint[];
  total: number;
  userId: UserId;
}
