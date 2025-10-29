import { Card, Typography, Input } from 'antd';
import { observer } from 'mobx-react-lite';
import { useState } from 'react';
import productStore from '../stores/ProductStore.jsx';
import ProductCard from '../widgets/ProductCard.jsx';
import ProductFormModal from '../widgets/ProductFormModal.jsx';

const { Text } = Typography;

const Products = observer(() => {
  const [searchText, setSearchText] = useState('');
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState(null);

  const filteredProducts = productStore.products.filter((product) =>
    product.name.toLowerCase().includes(searchText.toLowerCase())
  );

  const handleProductClick = (product) => {
    setSelectedProduct(product);
    setIsEditModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsEditModalOpen(false);
    setSelectedProduct(null);
  };

  if (productStore.products.length === 0) {
    return (
      <div
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          height: '50vh',
          padding: 16,
        }}
      >
        <Text type="secondary" style={{ fontSize: 16 }}>
          No products added yet
        </Text>
      </div>
    );
  }

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      height: '100vh',
      width: '100%'
    }}>
      {/* Sticky header with search */}
      <div
        style={{
          flexShrink: 0,
          backgroundColor: '#fff',
          padding: '8px 8px 8px 8px',
        }}
      >
        <Card
          bodyStyle={{ padding: '12px 16px' }}
          style={{
            borderRadius: 8,
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
          }}
        >
          <Input
            placeholder="Search by product name..."
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            allowClear
          />
        </Card>
      </div>

      {/* Products list */}
      <div style={{
        flex: 1,
        overflowY: 'auto',
        padding: 8,
        paddingTop: 8,
        minHeight: 0
      }}>
        {filteredProducts.length === 0 ? (
          <div
            style={{
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
              padding: 32,
            }}
          >
            <Text type="secondary">No products match your search</Text>
          </div>
        ) : (
          filteredProducts.map((product) => (
            <ProductCard key={product.id} product={product} onClick={handleProductClick} />
          ))
        )}
      </div>

      <ProductFormModal
        open={isEditModalOpen}
        onClose={handleCloseModal}
        mode="edit"
        product={selectedProduct}
      />
    </div>
  );
});

export default Products;
