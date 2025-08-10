import {makeAutoObservable} from 'mobx';
import axios from 'axios';
import Product from '../Models/Product';
import {ApiUrl} from "../config.jsx";
import {UUID} from "crypto";

class ProductStore {
    products: Product[] = [];

    constructor() {
        makeAutoObservable(this);
        this.loadProducts(); // Инициализация при создании
    }

    // Метод для загрузки продуктов с сервера
    async loadProducts() {
        try {
            const response = await axios.get<Product[]>(ApiUrl + '/products/unique');

            for (const item of response.data) {
                this.products.push(new Product(item.id, item.name, item.volume, item.category, item.brand));
            }
        } catch (error) {
            console.error('Ошибка загрузки продуктов:', error);
        }
    }

    // Получение всех продуктов
    getProducts() {
        return this.products;
    }

    // Получение продукта по ID
    getProductById(id: string) {
        return this.products.find((product) => product.id === id);
    }

    // Добавление нового продукта
    async addProduct(productData: Product) {
        type ProductResponse = {
          product_id: UUID;
        };
        try {
            // Создаем новый продукт и отправляем его на сервер
            const response = await axios.post<ProductResponse>(ApiUrl + '/products/', productData);
            const item = response.data;
            const newProduct = new Product(item.product_id, productData.name, productData.volume, productData.category, productData.brand);

            // Добавляем продукт в store
            this.products.push(newProduct);
            return newProduct;
        } catch (error) {
            return null;
        }
    }
}

const productStore = new ProductStore();
export default productStore;
