/* eslint-disable react/prop-types */
import { useEffect, useState } from 'react';
import { AutoComplete, Button, Form, Input, Modal, Tag, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';

import ProductStore from '../stores/ProductStore.js';
import VolumeSelectWidget from './VolumeSelectWidget.jsx';
import BrandSelectWidget from './BrandSelectWidget.jsx';
import DefaultTagsWidget from './DefaultTagsWidget.jsx';

const ProductSelectWidget = ({ onSelect, selectedProductProp }) => {
  const [filteredProducts, setFilteredProducts] = useState([]);
  const [selectedProduct, setSelectedProduct] = useState(selectedProductProp || null);
  const [inputLabel, setInputLabel] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form] = Form.useForm();
  const [open, setOpen] = useState(false);

  useEffect(() => {
    const products = ProductStore.getProducts();
    setFilteredProducts(products);
    setSelectedProduct(selectedProductProp || null);
    setInputLabel(selectedProductProp ? selectedProductProp.name : '');
  }, [selectedProductProp]);

  const handleSearch = (value) => {
    setInputLabel(value);
    const filtered = ProductStore.getProducts().filter((p) =>
      p.name.toLowerCase().includes(value.toLowerCase())
    );
    setFilteredProducts(filtered);
  };

  const handleSelect = (value) => {
    const product = ProductStore.getProductById(value);
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

  const handleAddProduct = async (values) => {
    const product = await ProductStore.addProduct(values);
    if (product) {
      message.success('Product added successfully!');
      setIsModalOpen(false);
      form.resetFields();
      handleSelect(product.id);
    } else {
      message.error('Failed to add product.');
    }
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

      <Modal
        title="Add New Product"
        open={isModalOpen}
        style={{ top: 8, padding: 0 }}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
      >
        <Form form={form} layout="horizontal" onFinish={handleAddProduct}>
          <Form.Item label="name" name="name" rules={[{ required: true, message: 'name required' }]}> 
            <Input placeholder="name" />
          </Form.Item>
          <Form.Item label="volume" name="volume" rules={[{ required: true, message: 'volume required' }]}> 
            <VolumeSelectWidget value={form.getFieldValue('volume')} onChange={(v) => form.setFieldsValue({ volume: v })} />
          </Form.Item>
          <Form.Item label="brand" name="brand" rules={[{ required: true, message: 'brand required' }]}> 
            <BrandSelectWidget value={form.getFieldValue('brand')} onChange={(v) => form.setFieldsValue({ brand: v })} />
          </Form.Item>
          <Form.Item label="default tags" name="default_tags"> 
            <DefaultTagsWidget value={form.getFieldValue('default_tags')} onChange={(v) => form.setFieldsValue({ default_tags: v })} />
          </Form.Item>
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

export default ProductSelectWidget;
