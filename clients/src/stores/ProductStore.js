import Product from '../models/Product.js';

const products = [
  new Product('1', 'Milk', '1L', 'BrandA', ['healthy', 'drink']),
  new Product('2', 'Bread', '500g', 'BrandB', ['food', 'carbs']),
  new Product('3', 'Coffee', '250g', 'BrandC', ['energy', 'drink'])
];

const getProducts = () => products;

const getProductById = (id) => products.find((p) => p.id === id);

const addProduct = async (productData) => {
  const id = crypto.randomUUID();
  const product = new Product(id, productData.name, productData.volume, productData.brand, productData.default_tags || []);
  products.push(product);
  return product;
};

export default { getProducts, getProductById, addProduct };
