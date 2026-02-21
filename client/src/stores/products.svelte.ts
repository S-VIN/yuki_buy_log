import { fetchProducts, createProduct, updateProduct } from '../lib/api';
import type { Product } from '../models/Product';

let items = $state<Product[]>([]);

const brands = $derived(
  Array.from(new Set(items.map((p) => p.brand).filter((b) => b.length > 0)))
);

const volumes = $derived(
  Array.from(new Set(items.map((p) => p.volume).filter((v) => v.length > 0)))
);

const tags = $derived(
  Array.from(new Set(items.flatMap((p) => p.default_tags)))
);

export const productStore = {
  get items() {
    return items;
  },

  /** Unique brands from all products, reactive. */
  get brands() {
    return brands;
  },

  /** Unique volumes from all products, reactive. */
  get volumes() {
    return volumes;
  },

  /** Unique default_tags from all products, reactive. */
  get tags() {
    return tags;
  },

  async load() {
    const data = await fetchProducts();
    items = ((data as { products: Product[] }).products ?? []) as Product[];
  },

  async create(product: Omit<Product, 'id' | 'user_id'>) {
    const created = (await createProduct(product)) as Product;
    items = [...items, created];
    return created;
  },

  async update(product: Product) {
    const updated = (await updateProduct(product)) as Product;
    items = items.map((p) => (p.id === updated.id ? updated : p));
    return updated;
  },

  clear() {
    items = [];
  },
};