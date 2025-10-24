import { useRef, useState, useMemo } from 'react';
import { Button, Card, message } from 'antd';
import dayjs from 'dayjs';
import { useNavigate } from 'react-router-dom';
import { observer } from 'mobx-react-lite';

import NativeDatePicker from '../widgets/NativeDatePicker.jsx';
import ProductSelectWidget from '../widgets/ProductSelectWidget.jsx';
import { useProductStore, usePurchaseStore } from '../stores/DataContext.jsx';
import Purchase from '../models/Purchase.js';
import ShopSelectWidget from '../widgets/ShopSelectWidget.jsx';
import PriceQuantitySelectWidget from '../widgets/PriceQuantitySelectWidget.jsx';
import TagSelectWidget from '../widgets/TagSelectWidget.jsx';
import ProductCardsWidget from '../widgets/ProductCardsWidget.jsx';

const AddReceipt = observer(() => {
  const productStore = useProductStore();
  const purchaseStore = usePurchaseStore();
  const [purchaseList, setPurchaseList] = useState([]);
  const [product, setSelectedProduct] = useState(null);
  const [shop, setSelectedShop] = useState(null);
  const tagSelectWidgetRef = useRef(null);
  const priceQuantitySelectWidgetRef = useRef(null);
  const [date, setSelectedDate] = useState(dayjs().format('YYYY-MM-DD'));
  const [price, setPrice] = useState(null);
  const [quantity, setQuantity] = useState(1);
  const [tags, setSelectedTags] = useState([]);
  const [messageApi, contextHolder] = message.useMessage();
  const navigate = useNavigate();

  const receiptId = useMemo(() => Date.now(), []);

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

  const handleAddPurchase = async () => {
    if (!product || !price || !quantity || !shop) {
      messageApi.warning('Please fill all fields');
      return;
    }

    try {
      const purchaseData = {
        product_id: Number(product.id),
        quantity: quantity,
        price: price,
        date: dayjs(date).toISOString(),
        store: shop,
        tags: tags,
        receipt_id: receiptId,
      };

      const serverPurchase = await purchaseStore.addPurchase(purchaseData);

      const newPurchase = new Purchase(serverPurchase.id, product, price, quantity, tags, receiptId);
      setPurchaseList([...purchaseList, newPurchase]);

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

  const handleDeletePurchase = async (purchase) => {
    try {
      if (purchase.uuid) {
        await purchaseStore.removePurchase(purchase.uuid);
        messageApi.success('Purchase deleted!');
      }
      setPurchaseList(purchaseList.filter((p) => p.uuid !== purchase.uuid));
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

    setPurchaseList(purchaseList.filter((p) => p.uuid !== purchase.uuid));
  };

  const handleCloseCheck = async () => {
    if (purchaseList.length < 1) {
      messageApi.warning('Add at least one purchase');
      return;
    }

    await purchaseStore.loadPurchases();
    messageApi.success('Receipt saved!');
    setPurchaseList([]);
    navigate('/receipts');
  };

  const handleDateChange = (d) => {
    setSelectedDate(d ? d.format('YYYY-MM-DD') : dayjs().format('YYYY-MM-DD'));
  };

  return (
    <div style={{ width: '100%', padding: 8 }}>
      {contextHolder}
      <Card style={{ borderRadius: 8, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8 }}>
            <NativeDatePicker onChange={handleDateChange} value={date} />
            <ShopSelectWidget value={shop} onChange={setSelectedShop} />
          </div>
          <PriceQuantitySelectWidget
            onPriceChanged={setPrice}
            onQuantityChanged={setQuantity}
            ref={priceQuantitySelectWidgetRef}
          />
          <ProductSelectWidget onSelect={handleSelectProduct} selectedProductProp={product} />
          <TagSelectWidget onTagChange={setSelectedTags} ref={tagSelectWidgetRef} />
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
      <ProductCardsWidget productListProp={purchaseList} onDelete={handleDeletePurchase} onEdit={handleEditPurchase} />
    </div>
  );
});

export default AddReceipt;
