import axios from 'axios';
import {Button, Form, Input, message} from 'antd';
import CategorySelectWidget from "../Widgets/CategorySelectWidget.js";
import BrandSelectWidget from "../Widgets/BrandSelectWidget.js";
import VolumeSelectWidget from "../Widgets/VolumeSelectWidget.js";
import ProductStore from "../Stores/ProductStore.js";
import Product from "../Models/Product.js";
import {ApiUrl} from "../config.jsx";


const AddProductScreen = () => {
    const [form] = Form.useForm();

    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();
            // Преобразование заглавных букв в строчные
            const formattedValues = {
                name: values.name.toLowerCase(),
                volume: values.volume.toLowerCase(),
                brand: values.brand.toLowerCase(),
                category: values.category.toLowerCase(),
            };

            // Отправка POST-запроса
            let response = await axios.post(ApiUrl + '/products/', formattedValues);

            await ProductStore.addProduct(new Product(response.data.id, values.name, values.volume, values.brand, values.category))
            message.success('Add Purchase!');
            form.resetFields();
        } catch (error) {
            message.error('Error adding product');
            console.error(error);
        }
    };

    return (<div
            style={{
                maxWidth: '400px',
                margin: '0px',
                padding: '20px',
                border: '1px solid #ddd',
                borderRadius: '8px',
                backgroundColor: '#fff',
            }}
        >
            <h2 style={{textAlign: 'center', marginBottom: '20px'}}>Добавить товар</h2>

            <Form
                form={form}
                layout="vertical"
                onFinish={handleSubmit}
            >
                {/* Название */}
                <Form.Item
                    label="name"
                    name="name"
                    rules={[{required: true, message: 'Name required'}]}
                >
                    <Input placeholder="Введите название"/>
                </Form.Item>

                {/* Объём */}
                <Form.Item
                    label="volume"
                    name="volume"
                    rules={[{required: true, message: 'Volume required'}]}
                >
                    <VolumeSelectWidget/>
                </Form.Item>

                {/* Бренд */}
                <Form.Item
                    label="Brand"
                    name="brand"
                    rules={[{required: true, message: 'Brand required'}]}
                >
                    <BrandSelectWidget/>
                </Form.Item>

                {/* Категория */}
                <Form.Item
                    label="Category"
                    name="category"
                    rules={[{required: true, message: 'Category required'}]}
                >
                    <CategorySelectWidget/>
                </Form.Item>

                {/* Кнопка подтверждения */}
                <Form.Item>
                    <Button type="primary" htmlType="submit" block>
                        Add Product
                    </Button>
                </Form.Item>
            </Form>
        </div>);
};

export default AddProductScreen;
