/* eslint-disable react/prop-types */
import { useState, useEffect } from 'react';
import { DatePicker } from 'antd';
import dayjs from 'dayjs';

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
      <div style={{ width: '100%', overflow: 'hidden' }}>
        <input
          type="date"
          value={value || dayjs().format('YYYY-MM-DD')}
          onChange={(e) => {
            const newDate = e.target.value ? dayjs(e.target.value) : dayjs();
            onChange(newDate);
          }}
          style={{
            width: '100%',
            height: '32px',
            padding: '4px 11px',
            fontSize: '14px',
            borderRadius: '6px',
            border: '1px solid #d9d9d9',
            outline: 'none',
            transition: 'all 0.3s',
            boxSizing: 'border-box',
            margin: 0,
            display: 'block',
          }}
          onFocus={(e) => {
            e.target.style.borderColor = '#4096ff';
            e.target.style.boxShadow = '0 0 0 2px rgba(5, 145, 255, 0.1)';
          }}
          onBlur={(e) => {
            e.target.style.borderColor = '#d9d9d9';
            e.target.style.boxShadow = 'none';
          }}
        />
      </div>
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

export default DatepickerCustom;
