import type { Product } from '../models/Product';

export interface PendingReceiptItem {
  uuid: string;
  product: Product;
  price: number;
  quantity: number;
  tags: string[];
}

export interface ReceiptEdit {
  date: string; // YYYY-MM-DD
  shop: string | null;
  items: PendingReceiptItem[];
}
