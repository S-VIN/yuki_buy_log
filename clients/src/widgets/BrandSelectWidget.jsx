import { AutoComplete } from 'antd';
import { observer } from 'mobx-react-lite';
import { useProductStore } from '../stores/DataContext.jsx';

const BrandSelectWidget = observer(({ value, onChange }) => {
  const productStore = useProductStore();

  return (
    <AutoComplete
      placeholder="brand"
      options={productStore.brands.map((b) => ({ value: b }))}
      value={value}
      onChange={onChange}
    />
  );
});

export default BrandSelectWidget;
