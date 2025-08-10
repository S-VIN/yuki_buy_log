/* eslint-disable react/prop-types */
import { AutoComplete } from 'antd';

const volumes = ['1L', '500g', '250g'];

const VolumeSelectWidget = ({ value, onChange }) => (
  <AutoComplete
    placeholder="volume"
    options={volumes.map((v) => ({ value: v }))}
    value={value}
    onChange={onChange}
  />
);

export default VolumeSelectWidget;
