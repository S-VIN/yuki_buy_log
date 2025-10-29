import { useRef, useState, useEffect } from 'react';
import { Button, Card, message } from 'antd';
import { AppstoreAddOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import { useNavigate, useLocation } from 'react-router-dom';
import { observer } from 'mobx-react-lite';

import DatepickerCustom from '../widgets/DatepickerCustom.jsx';
import ProductSelectWidget from '../widgets/ProductSelectWidget.jsx';
import productStore from '../stores/ProductStore.jsx';
import purchaseStore from '../stores/PurchaseStore.jsx';
import checkCache from '../stores/CheckCache.jsx';
import ShopSelectWidget from '../widgets/ShopSelectWidget.jsx';
import PriceQuantitySelectWidget from '../widgets/PriceQuantitySelectWidget.jsx';
import TagSelectWidget from '../widgets/TagSelectWidget.jsx';
import ProductCardsWidget from '../widgets/ProductCardsWidget.jsx';
import BulkTagsModal from '../widgets/BulkTagsModal.jsx';

const AddReceipt = observer(() => {
  const location = useLocation();
  const [product, setSelectedProduct] = useState(null);
  const [shop, setSelectedShop] = useState(null);
  const tagSelectWidgetRef = useRef(null);
  const priceQuantitySelectWidgetRef = useRef(null);
  const [date, setSelectedDate] = useState(dayjs().format('YYYY-MM-DD'));
  const [price, setPrice] = useState(null);
  const [quantity, setQuantity] = useState(1);
  const [tags, setSelectedTags] = useState([]);
  const [isBulkTagModalOpen, setIsBulkTagModalOpen] = useState(false);
  const [messageApi, contextHolder] = message.useMessage();
  const navigate = useNavigate();

  useEffect(() => {
    checkCache.clear();

    if (location.state?.receipt) {
      const { receipt } = location.state;
      setSelectedDate(dayjs(receipt.date).format('YYYY-MM-DD'));
      setSelectedShop(receipt.shop);

      receipt.items.forEach((item) => {
        const product = productStore.getProductById(String(item.product_id));
        checkCache.addPurchase(product, item.price, item.quantity, item.tags || []);
      });
    }
  }, [location.state]);

  const handleSelectProduct = (productId) => {
    const selected = productId ? productStore.getProductById(productId) : null;
    setSelectedProduct(selected);
    
    if (selected && selected.default_tags && selected.default_tags.length > 0) {
      setSelectedTags(selected.default_tags);
      tagSelectWidgetRef.current?.setTags(selected.default_tags);
    } else {
      setSelectedTags([]);
      tagSelectWidgetRef.current?.resetTags();
    }
  };

  const handleAddPurchase = () => {
    if (!product || !price || !quantity || !shop) {
      messageApi.warning('Please fill all fields');
      return;
    }

    try {
      checkCache.addPurchase(product, price, quantity, tags);

      setSelectedProduct(null);
      tagSelectWidgetRef.current?.resetTags();
      priceQuantitySelectWidgetRef.current?.reset();
      setSelectedTags([]);

      messageApi.success('Purchase added!');
    } catch (error) {
      messageApi.error(`Failed to add purchase: ${error.message}`);
      console.error('Add purchase error:', error);
    }
  };

  const handleDeletePurchase = (purchase) => {
    try {
      checkCache.removePurchase(purchase.uuid);
      messageApi.success('Purchase deleted!');
    } catch (error) {
      messageApi.error(`Failed to delete purchase: ${error.message}`);
      console.error('Delete purchase error:', error);
    }
  };

  const handleEditPurchase = (purchase) => {
    setSelectedProduct(purchase.product);
    setPrice(purchase.price);
    setQuantity(purchase.quantity);
    setSelectedTags(purchase.tags || []);

    tagSelectWidgetRef.current?.setTags(purchase.tags || []);
    priceQuantitySelectWidgetRef.current?.setValues(purchase.price, purchase.quantity);

    checkCache.removePurchase(purchase.uuid);
  };

  const handleCloseCheck = async () => {
    if (checkCache.isEmpty) {
      messageApi.warning('Add at least one purchase');
      return;
    }

    try {
      const receiptId = Math.floor(Date.now() / 1000);
      const purchases = checkCache.getPurchases();

      for (const purchase of purchases) {
        const purchaseData = {
          product_id: Number(purchase.product.id),
          quantity: purchase.quantity,
          price: purchase.price,
          date: dayjs(date).toISOString(),
          store: shop,
          tags: purchase.tags,
          receipt_id: receiptId,
        };

        await purchaseStore.addPurchase(purchaseData);
      }

      checkCache.clear();
      await purchaseStore.loadPurchases();
      messageApi.success('Receipt saved!');
      navigate('/receipts');
    } catch (error) {
      messageApi.error(`Failed to save receipt: ${error.message}`);
      console.error('Close check error:', error);
    }
  };

  const handleDateChange = (d) => {
    setSelectedDate(d ? d.format('YYYY-MM-DD') : dayjs().format('YYYY-MM-DD'));
  };

  const handleOpenBulkTagModal = () => {
    if (checkCache.isEmpty) {
      messageApi.warning('Add at least one purchase first');
      return;
    }
    setIsBulkTagModalOpen(true);
  };

  const handleCloseBulkTagModal = () => {
    setIsBulkTagModalOpen(false);
  };

  const handleAddBulkTags = (bulkTags) => {
    if (bulkTags.length === 0) {
      messageApi.warning('Please select at least one tag');
      return;
    }

    try {
      checkCache.addTagsToAllPurchases(bulkTags);
      messageApi.success(`Added ${bulkTags.length} tag(s) to all purchases`);
      setIsBulkTagModalOpen(false);
    } catch (error) {
      messageApi.error(`Failed to add tags: ${error.message}`);
      console.error('Add bulk tags error:', error);
    }
  };

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      height: '100vh',
      width: '100%',
      padding: 8
    }}>
      {contextHolder}
      <Card style={{ borderRadius: 8, boxShadow: '0 2px 8px rgba(0,0,0,0.1)', flexShrink: 0 }}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8 }}>
            <DatepickerCustom onChange={handleDateChange} value={date} />
            <ShopSelectWidget value={shop} onChange={setSelectedShop} />
          </div>
          <PriceQuantitySelectWidget
            onPriceChanged={setPrice}
            onQuantityChanged={setQuantity}
            ref={priceQuantitySelectWidgetRef}
          />
          <ProductSelectWidget onSelect={handleSelectProduct} selectedProductProp={product} />
          <div style={{ display: 'flex', flexDirection: 'row', gap: 8, alignItems: 'center' }}>
            <div style={{ flex: 1 }}>
              <TagSelectWidget onTagChange={setSelectedTags} ref={tagSelectWidgetRef} options={productStore.tags} />
            </div>
            <Button
              icon={<AppstoreAddOutlined />}
              onClick={handleOpenBulkTagModal}
              style={{ height: 32, flexShrink: 0 }}
              title="Add tags to all purchases"
            />
          </div>
          <div style={{ display: 'flex', flexDirection: 'row', gap: 8 }}>
            <Button type="primary" block style={{ height: 32, fontSize: 16 }} onClick={handleCloseCheck}>
              Close check
            </Button>
            <Button block style={{ height: 32, fontSize: 16 }} onClick={handleAddPurchase}>
              Add purchase
            </Button>
          </div>
        </div>
      </Card>
      <div style={{ flex: 1, minHeight: 0, marginTop: 8 }}>
        <ProductCardsWidget productListProp={checkCache.purchases} onDelete={handleDeletePurchase} onEdit={handleEditPurchase} />
      </div>

      <BulkTagsModal
        open={isBulkTagModalOpen}
        onCancel={handleCloseBulkTagModal}
        onAdd={handleAddBulkTags}
        options={productStore.tags}
      />
    </div>
  );
});

export default AddReceipt;
