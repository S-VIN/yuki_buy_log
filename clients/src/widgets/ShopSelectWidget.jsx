import { AutoComplete } from 'antd';
import { observer } from 'mobx-react-lite';
import { usePurchaseStore } from '../stores/DataContext.jsx';

const ShopSelectWidget = observer(({ value, onChange }) => {
  const purchaseStore = usePurchaseStore();

  return (
    <AutoComplete
      style={{ width: '100%' }}
      value={value}
      options={purchaseStore.shops.map((s) => ({ value: s }))}
      onChange={onChange}
      placeholder="shop"
    />
  );
});

export default ShopSelectWidget;
