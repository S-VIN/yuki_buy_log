export type ProductId = bigint;
export type UserId = bigint;

export interface Product {
  id: ProductId;
  name: string;
  volume: string;
  brand: string;
  default_tags: string[];
  user_id: UserId;
}