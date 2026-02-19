import { useEffect } from 'react';
import { Button, Form, Input, Modal, message } from 'antd';
import { observer } from 'mobx-react-lite';
import PropTypes from 'prop-types';

import productStore from '../stores/ProductStore.jsx';
import VolumeSelectWidget from './VolumeSelectWidget.jsx';
import BrandSelectWidget from './BrandSelectWidget.jsx';
import DefaultTagsWidget from './DefaultTagsWidget.jsx';

const ProductFormModal = observer(({ open, onClose, onSuccess, mode = 'add', product = null }) => {
  const [form] = Form.useForm();

  useEffect(() => {
    if (open) {
      if (mode === 'edit' && product) {
        // Pre-fill form with product data for editing
        form.setFieldsValue({
          name: product.name,
          volume: product.volume,
          brand: product.brand,
          default_tags: product.default_tags || [],
        });
      } else {
        // Reset form for adding new product
        form.resetFields();
      }
    }
  }, [open, mode, product, form]);

  const handleSubmit = async (values) => {
    try {
      if (mode === 'edit') {
        // Update existing product
        const productData = {
          id: parseInt(product.id, 10),
          ...values,
        };
        await productStore.updateProduct(productData);
        message.success('Product updated successfully!');
      } else {
        // Add new product
        await productStore.addProduct(values);
        message.success('Product added successfully!');
      }

      form.resetFields();
      onClose();
      if (onSuccess) {
        onSuccess();
      }
    } catch {
      message.error(`Failed to ${mode === 'edit' ? 'update' : 'add'} product.`);
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onClose();
  };

  return (
    <Modal
      title={mode === 'edit' ? 'Edit Product' : 'Add New Product'}
      open={open}
      style={{ top: 8 }}
      bodyStyle={{ padding: 8 }}
      onCancel={handleCancel}
      footer={null}
    >
      <Form form={form} layout="horizontal" onFinish={handleSubmit}>
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
            {mode === 'edit' ? 'save' : 'add'}
          </Button>
        </Form.Item>
      </Form>
    </Modal>
  );
});

ProductFormModal.propTypes = {
  open: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  onSuccess: PropTypes.func,
  mode: PropTypes.oneOf(['add', 'edit']),
  product: PropTypes.shape({
    id: PropTypes.string,
    name: PropTypes.string,
    volume: PropTypes.string,
    brand: PropTypes.string,
    default_tags: PropTypes.arrayOf(PropTypes.string),
  }),
};

export default ProductFormModal;
