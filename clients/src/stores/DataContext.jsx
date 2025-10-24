import { createContext, useContext, useState } from 'react';
import PropTypes from 'prop-types';
import { fetchProducts, createProduct, fetchPurchases, createPurchase } from '../api.js';
import Product from '../models/Product.js';

const DataContext = createContext({
  products: [],
  purchases: [],
  loadData: async () => {},
  addProduct: async () => {},
  addPurchase: async () => {},
  refreshProducts: async () => {},
  refreshPurchases: async () => {},
});

export const DataProvider = ({ children }) => {
  const [products, setProducts] = useState([]);
  const [purchases, setPurchases] = useState([]);
  const [loading, setLoading] = useState(false);

  const refreshProducts = async () => {
    try {
      const response = await fetchProducts();
      const productsData = response.products || [];
      const productObjects = productsData.map(
        (p) => new Product(String(p.id), p.name, p.volume, p.brand, p.default_tags || [])
      );
      setProducts(productObjects);
      return productObjects;
    } catch (error) {
      console.error('Failed to fetch products:', error);
      throw error;
    }
  };

  const refreshPurchases = async () => {
    try {
      const response = await fetchPurchases();
      const purchasesData = response.purchases || [];
      setPurchases(purchasesData);
      return purchasesData;
    } catch (error) {
      console.error('Failed to fetch purchases:', error);
      throw error;
    }
  };

  const loadData = async () => {
    setLoading(true);
    try {
      await Promise.all([refreshProducts(), refreshPurchases()]);
    } catch (error) {
      console.error('Failed to load data:', error);
    } finally {
      setLoading(false);
    }
  };

  const addProduct = async (productData) => {
    try {
      const response = await createProduct(productData);
      const newProduct = new Product(
        String(response.id),
        response.name,
        response.volume,
        response.brand,
        response.default_tags || []
      );
      setProducts([...products, newProduct]);
      return newProduct;
    } catch (error) {
      console.error('Failed to create product:', error);
      throw error;
    }
  };

  const addPurchase = async (purchaseData) => {
    try {
      const response = await createPurchase(purchaseData);
      setPurchases([...purchases, response]);
      return response;
    } catch (error) {
      console.error('Failed to create purchase:', error);
      throw error;
    }
  };

  return (
    <DataContext.Provider
      value={{
        products,
        purchases,
        loading,
        loadData,
        addProduct,
        addPurchase,
        refreshProducts,
        refreshPurchases,
      }}
    >
      {children}
    </DataContext.Provider>
  );
};

DataProvider.propTypes = {
  children: PropTypes.node.isRequired,
};

export const useData = () => useContext(DataContext);

export default DataContext;
