<script lang="ts">
  import CardWidget from "../widgets/CardWidget.svelte";
  import ProductWidget from "../widgets/ProductWidget.svelte";

  const demoProduct = {
    name: "Oat Milk Barista Edition",
    volume: "1 L",
    brand: "Oatly",
    tags: ["dairy-free", "vegan", "coffee"],
  };
</script>

<div class="page">
  <h2>Home</h2>
  <p>Home page placeholder</p>

  <div class="demo-section">
    <CardWidget>
      <ProductWidget product={demoProduct} />
    </CardWidget>
  </div>
</div>

<style>
  .page {
    padding: var(--page-padding);
  }

  .demo-section {
    margin-top: var(--space-6);
  }
</style>
