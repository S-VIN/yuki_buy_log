import Product from '../models/Product.js';

const products = [
  new Product('1', 'Milk', '1L', 'Dairy', 'BrandA'),
  new Product('2', 'Bread', '500g', 'Bakery', 'BrandB'),
  new Product('3', 'Coffee', '250g', 'Beverages', 'BrandC')
];

const getProducts = () => products;

const getProductById = (id) => products.find((p) => p.id === id);

const addProduct = async (productData) => {
  const id = crypto.randomUUID();
  const product = new Product(id, productData.name, productData.volume, productData.category, productData.brand);
  products.push(product);
  return product;
};

export default { getProducts, getProductById, addProduct };
