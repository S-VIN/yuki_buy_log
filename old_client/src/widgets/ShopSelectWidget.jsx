import { AutoComplete } from 'antd';
import { observer } from 'mobx-react-lite';
import purchaseStore from '../stores/PurchaseStore.jsx';

const ShopSelectWidget = observer(({ value, onChange }) => {

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
