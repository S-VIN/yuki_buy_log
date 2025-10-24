import { makeAutoObservable, runInAction } from 'mobx';
import { fetchPurchases, createPurchase, deletePurchase } from '../api.js';

class PurchaseStore {
  purchases = [];
  loading = false;
  error = null;

  constructor() {
    makeAutoObservable(this);
  }

  async loadPurchases() {
    this.loading = true;
    this.error = null;

    try {
      const response = await fetchPurchases();
      const purchasesData = response.purchases || [];

      runInAction(() => {
        this.purchases = purchasesData;
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

  async addPurchase(purchaseData) {
    try {
      const response = await createPurchase(purchaseData);

      runInAction(() => {
        this.purchases.push(response);
      });

      return response;
    } catch (error) {
      runInAction(() => {
        this.error = error.message;
      });
      throw error;
    }
  }

  async removePurchase(purchaseId) {
    try {
      await deletePurchase(purchaseId);

      runInAction(() => {
        this.purchases = this.purchases.filter((p) => p.id !== purchaseId);
      });

      return true;
    } catch (error) {
      runInAction(() => {
        this.error = error.message;
      });
      throw error;
    }
  }

  getPurchasesByReceiptId(receiptId) {
    return this.purchases.filter((p) => String(p.receipt_id) === String(receiptId));
  }

  get receipts() {
    const receiptMap = {};

    this.purchases.forEach((purchase) => {
      const receiptId = purchase.receipt_id;
      if (!receiptMap[receiptId]) {
        receiptMap[receiptId] = {
          id: receiptId,
          date: purchase.date,
          store: purchase.store,
          items: [],
        };
      }
      receiptMap[receiptId].items.push(purchase);
    });

    return Object.values(receiptMap);
  }

  get stores() {
    const storeSet = new Set();
    this.purchases.forEach((p) => p.store && storeSet.add(p.store));
    return Array.from(storeSet).sort();
  }
}

export default new PurchaseStore();
