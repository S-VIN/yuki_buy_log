import { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { DatePicker } from 'antd';
import dayjs from 'dayjs';
import './DatepickerCustom.css';

const isMobileDevice = () => {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
};

const DatepickerCustom = ({ onChange, value }) => {
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    setIsMobile(isMobileDevice());
  }, []);

  if (isMobile) {
    return (
      <input
        type="date"
        className="date-picker-mobile"
        value={value || dayjs().format('YYYY-MM-DD')}
        onChange={(e) => {
          const newDate = e.target.value ? dayjs(e.target.value) : dayjs();
          onChange(newDate);
        }}
      />
    );
  }

  return (
    <DatePicker
      style={{ width: '100%' }}
      onChange={(d) => onChange(d)}
      value={value ? dayjs(value) : null}
      format="DD-MM-YYYY"
      inputReadOnly
    />
  );
};

DatepickerCustom.propTypes = {
  onChange: PropTypes.func.isRequired,
  value: PropTypes.string,
};

export default DatepickerCustom;
