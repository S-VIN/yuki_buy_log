import type {ProductId, UserId} from "./Product";
import type {ReceiptId} from "./Receipt";

export type PurchaseId = bigint;

export interface Purchase {
  id: PurchaseId;
  product_id: ProductId;
  user_id?: UserId;
  date: Date;
  price: number;
  quantity?: number;
  store?: string;
  receipt_id?: ReceiptId;
  tags: string[];
}