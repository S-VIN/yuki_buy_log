import { useEffect, useState } from 'react';
import { AutoComplete, Button, Form, Input, Modal, Tag, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { observer } from 'mobx-react-lite';

import productStore from '../stores/ProductStore.jsx';
import VolumeSelectWidget from './VolumeSelectWidget.jsx';
import BrandSelectWidget from './BrandSelectWidget.jsx';
import DefaultTagsWidget from './DefaultTagsWidget.jsx';

const ProductSelectWidget = observer(({ onSelect, selectedProductProp }) => {
  const [filteredProducts, setFilteredProducts] = useState([]);
  const [selectedProduct, setSelectedProduct] = useState(selectedProductProp || null);
  const [inputLabel, setInputLabel] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form] = Form.useForm();
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

  const handleAddProduct = async (values) => {
    try {
      const product = await productStore.addProduct(values);
      message.success('Product added successfully!');
      setIsModalOpen(false);
      form.resetFields();
      handleSelect(product.id);
    } catch {
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
        style={{ top: 8 }}
        bodyStyle={{ padding: 8 }}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
      >
        <Form form={form} layout="horizontal" onFinish={handleAddProduct}>
          <Form.Item name="name" rules={[{ required: true, message: 'name required' }]} style={{ marginBottom: 8 }}>
            <Input placeholder="name" />
          </Form.Item>
          <Form.Item name="volume" rules={[{ required: true, message: 'volume required' }]} style={{ marginBottom: 8 }}>
            <VolumeSelectWidget value={form.getFieldValue('volume')} onChange={(v) => form.setFieldsValue({ volume: v })} volumes={productStore.volumes} />
          </Form.Item>
          <Form.Item name="brand" rules={[{ required: true, message: 'brand required' }]} style={{ marginBottom: 8 }}>
            <BrandSelectWidget value={form.getFieldValue('brand')} onChange={(v) => form.setFieldsValue({ brand: v })} />
          </Form.Item>
          <Form.Item name="default_tags" style={{ marginBottom: 8 }}>
            <DefaultTagsWidget value={form.getFieldValue('default_tags')} onChange={(v) => form.setFieldsValue({ default_tags: v })} options={productStore.tags} placeholder="default tags (optional)" />
          </Form.Item>
          <Form.Item style={{ marginBottom: 0 }}>
            <Button type="primary" htmlType="submit" block>
              add
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
});

export default ProductSelectWidget;
