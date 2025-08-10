import { List, Switch } from 'antd';

const Settings = () => (
  <div style={{ padding: 16 }}>
    <List>
      <List.Item actions={[<Switch key="notify" />]}>Notifications</List.Item>
      <List.Item actions={[<Switch key="dark" />]}>Dark mode</List.Item>
    </List>
  </div>
);

export default Settings;
