export type ProductId = string;
export type UserId = string;

export interface Product {
  id: ProductId;
  name: string;
  volume: string;
  brand: string;
  default_tags: string[];
  user_id: UserId;
}