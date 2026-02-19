class Purchase {
  constructor(uuid, product, price, quantity, tags = [], checkId = null) {
    this.uuid = uuid;
    this.product = product; // instance of Product
    this.price = price;
    this.quantity = quantity;
    this.tags = tags;
    this.check_id = checkId;
  }
}

export default Purchase;
