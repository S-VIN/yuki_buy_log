/* eslint-disable react/prop-types */
import { AutoComplete } from 'antd';

const VolumeSelectWidget = ({ value, onChange, volumes = [] }) => (
  <AutoComplete
    placeholder="volume"
    options={volumes.map((v) => ({ value: v }))}
    value={value}
    onChange={onChange}
  />
);

export default VolumeSelectWidget;
