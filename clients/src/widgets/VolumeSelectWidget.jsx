import PropTypes from 'prop-types';
import { AutoComplete } from 'antd';

const VolumeSelectWidget = ({ value, onChange, volumes = [] }) => (
  <AutoComplete
    placeholder="volume"
    options={volumes.map((v) => ({ value: v }))}
    value={value}
    onChange={onChange}
  />
);

VolumeSelectWidget.propTypes = {
  value: PropTypes.string,
  onChange: PropTypes.func.isRequired,
  volumes: PropTypes.arrayOf(PropTypes.string),
};

export default VolumeSelectWidget;
