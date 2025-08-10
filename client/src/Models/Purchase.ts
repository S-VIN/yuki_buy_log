import { UUID } from 'crypto';
import Product from './Product';

class Purchase {
    uuid: UUID | null;
    product: Product;
    price: number;
    quantity: number;
    tags: string[];
    check_id: number | null;

    constructor(uuid: UUID | null, product: Product, price: number, quantity: number, tags: string[], check_id: number | null) {
        this.uuid = uuid;
        this.product = product;
        this.price = price;
        this.quantity = quantity;
        this.tags = tags;
        this.check_id = check_id;
    }
}

export default Purchase;