<script lang="ts">
  import CardWidget from "../widgets/CardWidget.svelte";
  import ProductWidget from "../widgets/ProductWidget.svelte";
  import SelectWidget from "../widgets/SelectWidget.svelte";
  import ProductSelectWidget from "../widgets/ProductSelectWidget.svelte";
  import type { Product } from "../models/Product";

  const demoProduct = {
    id: "demo",
    name: "Oat Milk Barista Edition",
    volume: "1 L",
    brand: "Oatly",
    default_tags: ["dairy-free", "vegan", "coffee"],
    user_id: "",
  };

  let selectedCategory = $state<string | null>(null);
  let categoryOptions = $state(["Молочное", "Напитки", "Снеки", "Заморозка", "Бакалея"]);

  let selectedBrand = $state<string | null>(null);
  let brandOptions = $state(["Oatly", "Alpro", "ВкусВилл", "Valio"]);

  let selectedProduct = $state<Product | null>(null);
</script>

<div class="page">
  <h2>Home</h2>
  <p>Home page placeholder</p>

  <div class="demo-section">
    <CardWidget>
      <ProductWidget product={demoProduct} />
    </CardWidget>
  </div>

  <div class="demo-section">
    <h3 class="demo-title">ProductSelectWidget demo</h3>
    <div class="demo-field">
      <label class="demo-label" for="sel-product">Product</label>
      <ProductSelectWidget
        id="sel-product"
        bind:value={selectedProduct}
      />
    </div>
  </div>

  <div class="demo-section">
    <h3 class="demo-title">SelectWidget demo</h3>
    <div class="demo-row">
      <div class="demo-field">
        <label class="demo-label" for="sel-category">Категория</label>
        <SelectWidget
          id="sel-category"
          allOptions={categoryOptions}
          bind:value={selectedCategory}
          color="var(--color-blue)"
          placeholder="Выбрать категорию…"
        />
      </div>
      <div class="demo-field">
        <label class="demo-label" for="sel-brand">Бренд</label>
        <SelectWidget
          id="sel-brand"
          allOptions={brandOptions}
          bind:value={selectedBrand}
          color="var(--color-green)"
          placeholder="Выбрать бренд…"
        />
      </div>
    </div>
    <p class="demo-state">
      Категория: <strong>{selectedCategory ?? "—"}</strong> &nbsp;|&nbsp;
      Бренд: <strong>{selectedBrand ?? "—"}</strong>
    </p>
  </div>
</div>

<style>
  .page {
    padding: var(--page-padding);
  }

  .demo-section {
    margin-top: var(--space-6);
  }

  .demo-title {
    font-size: var(--text-base);
    font-weight: 600;
    color: var(--color-text-secondary);
    margin-bottom: var(--space-4);
  }

  .demo-row {
    display: flex;
    gap: var(--space-4);
    flex-wrap: wrap;
  }

  .demo-field {
    flex: 1;
    min-width: 200px;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .demo-label {
    font-size: var(--text-sm);
    font-weight: 500;
    color: var(--color-text-secondary);
  }

  .demo-state {
    margin-top: var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-text-secondary);
  }
</style>
