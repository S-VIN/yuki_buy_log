import React, { useEffect, useState } from 'react';
import { AutoComplete, message } from 'antd';
import axios from 'axios';
import { ApiUrl } from "../config.jsx";
import { ShopOutlined } from '@ant-design/icons';

interface ShopSelectWidgetProps {
  value: string | null;
  onChange: (value: string) => void;
}

const ShopSelectWidget: React.FC<ShopSelectWidgetProps> = ({ value, onChange }) => {
  const [shops, setShops] = useState<string[]>([]);
  const [filteredShops, setFilteredShops] = useState<string[]>([]);

  useEffect(() => {
    const fetchShops = async () => {
      try {
        const response = await axios.get<string[]>(ApiUrl + '/purchases/shop/unique');
        setShops(response.data);
        setFilteredShops(response.data);
      } catch (error) {
        message.error('Downloading shops error');
        console.error(error);
      }
    };

    fetchShops();
  }, []);

  const handleSearch = (searchValue: string) => {
    setFilteredShops(shops.filter((shop) => shop.toLowerCase().includes(searchValue.toLowerCase())));
  };

  return (
    <AutoComplete
      placeholder="shop"
      style={{ width: '100%' }}
      allowClear
      options={filteredShops.map((shop) => ({ value: shop }))}
      value={value}
      onChange={onChange}
      onSearch={handleSearch}
      filterOption={false}
      suffixIcon={<ShopOutlined style={{color: 'rgba(0, 0, 0, 0.45)' }} />}
    />
  );
};

export default ShopSelectWidget;