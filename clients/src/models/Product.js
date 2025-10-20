class Product {
  constructor(id, name, volume, brand, default_tags) {
    this.id = id;
    this.name = name;
    this.volume = volume;
    this.brand = brand;
    this.default_tags = Array.isArray(default_tags) ? default_tags : [];
  }

  getDescription() {
    return `Product [ID: ${this.id}, Name: ${this.name}, Volume: ${this.volume}, Brand: ${this.brand}, Default tags: ${this.default_tags.join(', ')}]`;
  }
}

export default Product;
