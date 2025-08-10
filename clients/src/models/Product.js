class Product {
  constructor(id, name, volume, category, brand) {
    this.id = id;
    this.name = name;
    this.volume = volume;
    this.category = category;
    this.brand = brand;
  }

  getDescription() {
    return `Product [ID: ${this.id}, Name: ${this.name}, Volume: ${this.volume}, category: ${this.category}, Brand: ${this.brand}]`;
  }
}

export default Product;
