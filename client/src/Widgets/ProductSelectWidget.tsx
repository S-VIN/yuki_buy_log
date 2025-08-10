import React, { useEffect, useState } from 'react';
import { AutoComplete, Button, Form, Input, message, Modal, Tag } from 'antd';
import { PlusOutlined } from '@ant-design/icons';

import ProductStore from '../Stores/ProductStore';
import VolumeSelectWidget from "./VolumeSelectWidget";
import BrandSelectWidget from "./BrandSelectWidget";
import CategorySelectWidget from "./CategorySelectWidget";
import Product from '../Models/Product';
import { UUID } from 'crypto';


interface SelectProductProps {
  style?: React.CSSProperties;
  onSelect: (id: UUID | null) => void;
  selectedProductProp?: Product | null;
}

const SelectProduct: React.FC<SelectProductProps> = ({ style, onSelect, selectedProductProp }) => {
  const [filteredProducts, setFilteredProducts] = useState<Product[]>([]);
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(selectedProductProp || null);
  const [inputLabel, setInputLabel] = useState<string>('');
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [form] = Form.useForm();
  const [isDropdownOpen, setIsDropdownOpen] = useState<boolean>(false);

  useEffect(() => {
    let products = ProductStore.getProducts();
    setFilteredProducts(products);
    setSelectedProduct(selectedProductProp || null);
    setInputLabel(selectedProductProp ? selectedProductProp.name : '');
  }, [selectedProductProp]);

  const handleSearch = (value: string) => {
    setInputLabel(value);

    const filtered = ProductStore.getProducts().filter((product) =>
      product.name.toLowerCase().includes(value.toLowerCase())
    );
    setFilteredProducts(filtered);
  };

  const handleSelect = (value: string) => {
    const id = value as UUID;
    onSelect(id);
    const product = id ? ProductStore.getProductById(id) : null;
    if (product) {
      setInputLabel(product.name);
      setSelectedProduct(product);
    }
  };

  const handleAddProductClick = () => {
    setIsDropdownOpen(false); // Закрываем выпадающий список
    setIsModalOpen(true);
  };

  const handleAddProduct = async (values: { name: string; volume: string; brand: string; category: string }) => {
    try {
      let productFromServer = await ProductStore.addProduct(new Product(null, values.name, values.volume, values.category, values.brand));
      if (productFromServer) {
        handleSelect(productFromServer.id as UUID)
        message.success('Product added successfully!');
      }
      setIsModalOpen(false);
      form.resetFields();
    } catch (error) {
      message.error('Failed to add product.');
    }
  };

  return (
    <div style={{ ...style }}>
      <AutoComplete
        style={{ width: '100%' }}
        onSearch={handleSearch}
        onSelect={(value) => handleSelect(value)}
        placeholder="Name"
        notFoundContent={'No products found'}
        value={inputLabel}
        open={isDropdownOpen}
        onOpenChange={setIsDropdownOpen}
        options={[
          {
            value: 'add-new-product',
            label: (
              <div
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '8px',
                  color: 'blue',
                  cursor: 'pointer',
                }}
                onClick={handleAddProductClick}
              >
                <PlusOutlined />
                Add new product
              </div>
            ),
          },
          ...filteredProducts.map((product) => ({
            value: product.id,
            label: (
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <span>{product.name}</span>
                <div style={{ display: 'flex', gap: '8px' }}>
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

      <Modal
        title="Add New Product"
        open={isModalOpen}
        style={{ top: 8, padding: '0px' }}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
      >
        <Form form={form} layout="horizontal" onFinish={handleAddProduct}>
          {/* Название */}
          <Form.Item
            label="name"
            name="name"
            rules={[{ required: true, message: 'name required' }]}
          >
            <Input placeholder="name" />
          </Form.Item>

          {/* Объём */}
          <Form.Item
            label="volume"
            name="volume"
            rules={[{ required: true, message: 'volume required' }]}
          >
            <VolumeSelectWidget 
              value={form.getFieldValue('volume')}
              onChange={(value) => form.setFieldsValue({ volume: value })}
            />
          </Form.Item>

          {/* Бренд */}
          <Form.Item
            label="brand"
            name="brand"
            rules={[{ required: true, message: 'brand required' }]}
          >
            <BrandSelectWidget 
              value={form.getFieldValue('brand')}
              onChange={(value) => form.setFieldsValue({ brand: value })}
            />
          </Form.Item>

          {/* Категория */}
          <Form.Item
            label="category"
            name="category"
            rules={[{ required: true, message: 'category required' }]}
          >
            <CategorySelectWidget 
              value={form.getFieldValue('category')}
              onChange={(value) => form.setFieldsValue({ category: value })}
            />
          </Form.Item>

          {/* Кнопка подтверждения */}
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              add
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default SelectProduct;
