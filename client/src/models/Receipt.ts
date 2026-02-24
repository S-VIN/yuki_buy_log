export type ReceiptId = number;

export interface Receipt {
  id: ReceiptId;
  date: string;
  store: string;
  common_tags: string[];
  purchase_ids: string[];
  total: number;
}
