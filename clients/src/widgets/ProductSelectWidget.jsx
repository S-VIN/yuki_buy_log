import { useEffect, useState } from 'react';
import { AutoComplete, Input, Tag } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { observer } from 'mobx-react-lite';

import productStore from '../stores/ProductStore.jsx';
import ProductFormModal from './ProductFormModal.jsx';

const ProductSelectWidget = observer(({ onSelect, selectedProductProp }) => {
  const [filteredProducts, setFilteredProducts] = useState([]);
  const [selectedProduct, setSelectedProduct] = useState(selectedProductProp || null);
  const [inputLabel, setInputLabel] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    setFilteredProducts(productStore.products);
    setSelectedProduct(selectedProductProp || null);
    setInputLabel(selectedProductProp ? selectedProductProp.name : '');
  }, [selectedProductProp]);

  const handleSearch = (value) => {
    setInputLabel(value);
    const filtered = productStore.products.filter((p) =>
      p.name.toLowerCase().includes(value.toLowerCase())
    );
    setFilteredProducts(filtered);
  };

  const handleSelect = (value) => {
    const product = productStore.getProductById(value);
    if (product) {
      setInputLabel(product.name);
      setSelectedProduct(product);
      onSelect(product.id);
    } else {
      onSelect(null);
    }
  };

  const handleAddProductClick = () => {
    setOpen(false);
    setIsModalOpen(true);
  };

  const handleProductAdded = () => {
    // After product is added, refresh the filtered products
    setFilteredProducts(productStore.products);
  };

  return (
    <div>
      <AutoComplete
        style={{ width: '100%' }}
        onSearch={handleSearch}
        onSelect={handleSelect}
        placeholder="Name"
        notFoundContent={'No products found'}
        value={inputLabel}
        open={open}
        onOpenChange={setOpen}
        options={[
          {
            value: 'add-new',
            label: (
              <div
                style={{ display: 'flex', gap: 8, color: 'blue', cursor: 'pointer' }}
                onClick={handleAddProductClick}
              >
                <PlusOutlined /> Add new product
              </div>
            ),
          },
          ...filteredProducts.map((product) => ({
            value: product.id,
            label: (
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <span>{product.name}</span>
                <div style={{ display: 'flex', gap: 8 }}>
                  <Tag color="blue">{product.volume}</Tag>
                  <Tag color="green">{product.brand}</Tag>
                </div>
              </div>
            ),
          })),
        ]}
      >
        <Input
          addonAfter={
            selectedProduct ? (
              <div>
                <Tag color="blue">{selectedProduct.volume}</Tag>
                <Tag color="green">{selectedProduct.brand}</Tag>
              </div>
            ) : null
          }
        />
      </AutoComplete>

      <ProductFormModal
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSuccess={handleProductAdded}
        mode="add"
      />
    </div>
  );
});

export default ProductSelectWidget;
