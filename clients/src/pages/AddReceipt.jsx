import { useEffect, useRef, useState } from 'react';
import { Button, Card, message } from 'antd';
import dayjs from 'dayjs';
import { useNavigate, useParams } from 'react-router-dom';

import NativeDatePicker from '../widgets/NativeDatePicker.jsx';
import ProductSelectWidget from '../widgets/ProductSelectWidget.jsx';
import ProductStore from '../stores/ProductStore.js';
import Purchase from '../models/Purchase.js';
import ShopSelectWidget from '../widgets/ShopSelectWidget.jsx';
import PriceQuantitySelectWidget from '../widgets/PriceQuantitySelectWidget.jsx';
import TagSelectWidget from '../widgets/TagSelectWidget.jsx';
import ProductCardsWidget from '../widgets/ProductCardsWidget.jsx';

const AddReceipt = () => {
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
  const { id } = useParams();

  useEffect(() => {
    if (id) {
      const receipts = JSON.parse(localStorage.getItem('receipts') || '[]');
      const rec = receipts.find((r) => String(r.id) === id);
      if (rec) {
        setPurchaseList(rec.items);
        setSelectedShop(rec.shop);
        setSelectedDate(rec.date);
      }
    }
  }, [id]);

  const handleSelectProduct = (productId) => {
    const selected = productId ? ProductStore.getProductById(productId) : null;
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
      return;
    }
    const newPurchase = new Purchase(null, product, price, quantity, tags, null);
    setPurchaseList([...purchaseList, newPurchase]);
    setSelectedProduct(null);
    tagSelectWidgetRef.current?.resetTags();
    priceQuantitySelectWidgetRef.current?.reset();
    setSelectedTags([]);
  };

  const handleDeletePurchase = (productId) => {
    setPurchaseList(purchaseList.filter((p) => p.product.id !== productId));
  };

  const handleEditPurchase = (purchase) => {
    setSelectedProduct(purchase.product);
    setPrice(purchase.price);
    setQuantity(purchase.quantity);
    setSelectedTags(purchase.tags || []);
    
    tagSelectWidgetRef.current?.setTags(purchase.tags || []);
    priceQuantitySelectWidgetRef.current?.setValues(purchase.price, purchase.quantity);
    
    setPurchaseList(purchaseList.filter((p) => p.product.id !== purchase.product.id));
  };

  const handleCloseCheck = () => {
    if (purchaseList.length < 1) {
      return;
    }
    const receipts = JSON.parse(localStorage.getItem('receipts') || '[]');
    const receipt = { id: id ? Number(id) : Date.now(), date, shop, items: purchaseList };
    const updated = id
      ? receipts.map((r) => (r.id === receipt.id ? receipt : r))
      : [receipt, ...receipts];
    localStorage.setItem('receipts', JSON.stringify(updated));
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
              {id ? 'Save check' : 'Close check'}
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
};

export default AddReceipt;
