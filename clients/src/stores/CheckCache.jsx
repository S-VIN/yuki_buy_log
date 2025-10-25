import { makeAutoObservable } from 'mobx';
import Purchase from '../models/Purchase.js';

class CheckCache {
  purchases = [];

  constructor() {
    makeAutoObservable(this);
  }

  addPurchase(product, price, quantity, tags = []) {
    const uuid = `temp_${Date.now()}_${Math.random()}`;
    const newPurchase = new Purchase(uuid, product, price, quantity, tags, null);
    this.purchases.push(newPurchase);
    return newPurchase;
  }

  removePurchase(uuid) {
    this.purchases = this.purchases.filter((p) => p.uuid !== uuid);
  }

  updatePurchase(uuid, product, price, quantity, tags) {
    const index = this.purchases.findIndex((p) => p.uuid === uuid);
    if (index !== -1) {
      this.purchases[index] = new Purchase(uuid, product, price, quantity, tags, null);
    }
  }

  getPurchases() {
    return this.purchases;
  }

  clear() {
    this.purchases = [];
  }

  get isEmpty() {
    return this.purchases.length === 0;
  }

  addTagsToAllPurchases(newTags) {
    this.purchases.forEach((purchase) => {
      const existingTags = purchase.tags || [];
      const combinedTags = [...new Set([...existingTags, ...newTags])];
      purchase.tags = combinedTags;
    });
  }
}

export default new CheckCache();
