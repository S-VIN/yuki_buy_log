export type PurchaseId = string;

export interface Purchase {
  id: PurchaseId;
  product_id: string;
  user_id?: string;
  date: string;
  price: number;
  quantity?: number;
  store?: string;
  receipt_id?: number;
  tags: string[];
}