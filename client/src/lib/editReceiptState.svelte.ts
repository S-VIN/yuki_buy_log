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

let data = $state<ReceiptEdit | null>(null);
let version = $state(0);

export const editReceiptState = {
  get data() {
    return data;
  },
  get version() {
    return version;
  },
  set(d: ReceiptEdit) {
    data = d;
    version++;
  },
  clear() {
    data = null;
  },
};
