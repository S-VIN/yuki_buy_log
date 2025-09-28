/* eslint-disable react/prop-types */
import { DatePicker } from 'antd';
import dayjs from 'dayjs';

const NativeDatePicker = ({ onChange, value }) => (
  <DatePicker
    style={{ width: '100%' }}
    onChange={(d) => onChange(d)}
    value={value ? dayjs(value) : null}
    format="DD-MM-YYYY"
    inputReadOnly
  />
);

export default NativeDatePicker;
