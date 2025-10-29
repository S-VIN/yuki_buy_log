import { makeAutoObservable, runInAction } from 'mobx';
import { fetchProducts, createProduct, updateProduct } from '../api.js';
import Product from '../models/Product.js';

class ProductStore {
  products = [];
  loading = false;
  error = null;

  constructor() {
    makeAutoObservable(this);
  }

  async loadProducts() {
    this.loading = true;
    this.error = null;

    try {
      const response = await fetchProducts();
      const productsData = response.products || [];

      runInAction(() => {
        this.products = productsData.map(
          (p) => new Product(String(p.id), p.name, p.volume, p.brand, p.default_tags || [])
        );
        this.loading = false;
      });
    } catch (error) {
      runInAction(() => {
        this.error = error.message;
        this.loading = false;
      });
      throw error;
    }
  }

  async addProduct(productData) {
    try {
      const response = await createProduct(productData);
      const newProduct = new Product(
        String(response.id),
        response.name,
        response.volume,
        response.brand,
        response.default_tags || []
      );

      runInAction(() => {
        this.products.push(newProduct);
      });

      return newProduct;
    } catch (error) {
      runInAction(() => {
        this.error = error.message;
      });
      throw error;
    }
  }

  async updateProduct(productData) {
    try {
      const response = await updateProduct(productData);
      const updatedProduct = new Product(
        String(response.id),
        response.name,
        response.volume,
        response.brand,
        response.default_tags || []
      );

      runInAction(() => {
        const index = this.products.findIndex((p) => p.id === updatedProduct.id);
        if (index !== -1) {
          this.products[index] = updatedProduct;
        }
      });

      return updatedProduct;
    } catch (error) {
      runInAction(() => {
        this.error = error.message;
      });
      throw error;
    }
  }

  getProductById(id) {
    return this.products.find((p) => p.id === id);
  }

  get brands() {
    const brandSet = new Set();
    this.products.forEach((p) => p.brand && brandSet.add(p.brand));
    return Array.from(brandSet).sort();
  }

  get volumes() {
    const volumeSet = new Set();
    this.products.forEach((p) => p.volume && volumeSet.add(p.volume));
    return Array.from(volumeSet).sort();
  }

  get tags() {
    const tagSet = new Set();
    this.products.forEach((p) => {
      if (Array.isArray(p.default_tags)) {
        p.default_tags.forEach((tag) => tagSet.add(tag));
      }
    });
    return Array.from(tagSet).sort();
  }
}

export default new ProductStore();
