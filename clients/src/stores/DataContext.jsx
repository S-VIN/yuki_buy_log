import { createContext, useContext } from 'react';
import PropTypes from 'prop-types';
import productStore from './ProductStore.jsx';
import purchaseStore from './PurchaseStore.jsx';

const DataContext = createContext({
  productStore,
  purchaseStore,
  loadData: async () => {},
});

export const DataProvider = ({ children }) => {
  const loadData = async () => {
    try {
      await Promise.all([
        productStore.loadProducts(),
        purchaseStore.loadPurchases(),
      ]);
    } catch (error) {
      console.error('Failed to load data:', error);
    }
  };

  return (
    <DataContext.Provider
      value={{
        productStore,
        purchaseStore,
        loadData,
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

// Convenience hooks for direct store access
export const useProductStore = () => {
  const { productStore } = useContext(DataContext);
  return productStore;
};

export const usePurchaseStore = () => {
  const { purchaseStore } = useContext(DataContext);
  return purchaseStore;
};

export default DataContext;
