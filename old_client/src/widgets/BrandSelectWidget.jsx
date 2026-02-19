import { AutoComplete } from 'antd';
import { observer } from 'mobx-react-lite';
import productStore from '../stores/ProductStore.jsx';

const BrandSelectWidget = observer(({ value, onChange }) => {

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
