import React, { useEffect, useRef, useState } from 'react';
import {Button, Card, message} from 'antd';
import dayjs from 'dayjs';

import NativeDatePicker from "../Widgets/NativeDatePicker";
import SelectProductWidget from "../Widgets/ProductSelectWidget";
import ProductStore from "../Stores/ProductStore";
import Purchase from "../Models/Purchase";
import Product from "../Models/Product";

import ShopSelectWidget from "../Widgets/ShopSelectWidget";

import PriceQuantitySelectWidget from "../Widgets/PriceQuantitySelectWidget";
import TagSelectWidget from "../Widgets/TagSelectWidget";
import ProductCardsWidget from "../Widgets/ProductCardsWidget";
import { UUID } from 'crypto';
import axios from "axios";
import {ApiUrl} from "../config.tsx";


const AddCheckScreen: React.FC = () => {
    const [purchaseList, setPurchaseList] = useState<Purchase[]>([]);
    const [product, setSelectedProduct] = useState<Product | null>(null);
    const [shop, setSelectedShop] = useState<string | null>(null);
    const [tags, setSelectedTags] = useState<string[]>([]);
    const tagSelectWidgetRef = useRef<{ resetTags: () => void }>(null);
    const priceQuantitySelectWidgetRef = useRef<{ reset: () => void }>(null);

    const [date, setSelectedDate] = useState<string | null>(dayjs().format("YYYY-MM-DD"));
    const [price, setPrice] = useState<number | null>(null);
    const [quantity, setQuantity] = useState<number>(1);
    const [messageApi, contextHolder] = message.useMessage();


    useEffect(() => {
    
    }, []);


    // const handleClearForm = () => {
    //     // form.resetFields();
    //     // setProductList([]);
    //     // setSelectedProduct(null);
    //     // setIsCleared(!isCleared); // Переключение состояния для сброса зависимых компонентов
    //     // message.success('Форма успешно очищена.');
    // };

    const handleSelectProduct = (id: UUID | null) => {
        const selectedProduct = id ? ProductStore.getProductById(id) : null;
        setSelectedProduct(selectedProduct || null);
    };

    const handleCloseCheck = async () => {
        if (purchaseList.length < 1)
        {
            return;
        }

        try {
            const transformedPurchases = purchaseList.map(purchase => ({
                product_id: purchase.product.id,
                price: purchase.price,
                quantity: purchase.quantity,
                tags: purchase.tags,
            }));


            // Отправка POST-запроса
            let response = await axios.post(ApiUrl + '/check/', {purchases: transformedPurchases, date: date, shop: shop});
            if (response.status != 200)
            {
                messageApi.error('Error adding product by server');
                return;
            }

            messageApi.success('Add Purchases!');
            setPurchaseList([])
        } catch (error) {
            messageApi.error('Error adding product');
            console.error(error);
        }
    };
    
        const handleAddPurchase = () => {
        console.log('handleAddPurchase', product, price, quantity, tags);

        if (!product || !price || !quantity || !shop) {
            // Handle error: missing required fields
            return;
        }

        const updatedList = [...purchaseList, new Purchase(null, product, price, quantity, tags, null)];
        setPurchaseList(updatedList);
        setSelectedProduct(null);
        tagSelectWidgetRef.current?.resetTags();
        priceQuantitySelectWidgetRef.current?.reset();
    };

    const handleDeletePurchase = (product_id: string) => {
        console.log('handleDeletePurchase', product_id);
        if (!product_id)
        {
            return;
        }
        const updatedList = purchaseList.filter(purchase => purchase.product.id !== product_id);
        setPurchaseList(updatedList);
    };

    const handleDateChange = (date: dayjs.Dayjs) => {
        if (date) {
            setSelectedDate(date.format("YYYY-MM-DD"));
        } else {
            setSelectedDate(null);
        }
    };

    return (

        <div style={{width: '100%'}}>
            {contextHolder}
            <Card
                style={{
                    borderRadius: '8px',
                    boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
                }}

            >
                <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
                    <div style={{ display: 'flex', flexDirection: 'row', gap: '8px' }}>
                        <NativeDatePicker onChange={handleDateChange} value={date ? dayjs(date) : dayjs()} />
                        <ShopSelectWidget value={shop} onChange={setSelectedShop} />
                    </div>

                    <PriceQuantitySelectWidget onPriceChanged={setPrice} onQuantityChanged={setQuantity} ref={priceQuantitySelectWidgetRef}/>

                    <SelectProductWidget onSelect={handleSelectProduct} selectedProductProp={product} />

                    <TagSelectWidget onTagChange={setSelectedTags}  ref={tagSelectWidgetRef}/>

                    <div style={{ display: 'flex', flexDirection: 'row', gap: '8px' }}>
                        <Button
                            type="primary"
                            block
                            style={{height: '32px', fontSize: '16px'}}
                            onClick={handleCloseCheck}
                        >
                            Close check
                        </Button>

                        <Button
                            block
                            style={{height: '32px', fontSize: '16px' }}
                            onClick={handleAddPurchase}
                        >
                            Add purchase
                        </Button>
                    </div>
                </div>
            </Card>

            <ProductCardsWidget productListProp={purchaseList} onDelete={handleDeletePurchase} />
        </div>
    );
};

export default AddCheckScreen;