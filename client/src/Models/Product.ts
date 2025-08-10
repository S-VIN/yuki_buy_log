import { UUID } from "crypto";

class Product {
    id: UUID | null;
    name: string;
    volume: string;
    category: string;
    brand: string;

    /**
     * Конструктор для создания объекта покупки
     * @param {string} id - Уникальный идентификатор покупки
     * @param {string} name - Название покупки
     * @param {string} volume - Объём покупки
     * @param {Array<string>} category - Массив тегов, связанных с покупкой
     * @param {string} brand - Бренд покупки
     */
    constructor(id: UUID | null, name: string, volume: string, category: string, brand: string) {
        this.id = id; // Уникальный идентификатор
        this.name = name; // Название
        this.volume = volume; // Объём
        this.category = category;
        this.brand = brand; // Бренд
    }

    /**
     * Метод для получения описания объекта в текстовом формате
     * @returns {string} Строка с описанием объекта
     */
    getDescription(): string {
        return `Product [ID: ${this.id}, Name: ${this.name}, Volume: ${this.volume}, category: ${this.category}, Brand: ${this.brand}]`;
    }


}

export default Product;
