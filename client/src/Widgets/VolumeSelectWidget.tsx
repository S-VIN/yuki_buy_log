import React, { useEffect, useState } from 'react';
import { AutoComplete, message } from 'antd';
import axios from 'axios';
import { ApiUrl } from "../config";

interface VolumeSelectWidgetProps {
  value: string;
  onChange: (value: string) => void;
}

const VolumeSelectWidget: React.FC<VolumeSelectWidgetProps> = ({ value, onChange }) => {
  const [volumes, setVolumes] = useState<string[]>([]);
  const [filteredVolumes, setFilteredVolumes] = useState<string[]>([]);

  useEffect(() => {
    const fetchVolumes = async () => {
      try {
        const response = await axios.get<string[]>(ApiUrl + '/products/volume/unique');
        setVolumes(response.data);
        setFilteredVolumes(response.data);
      } catch (error) {
        message.error('Volumes download error');
        console.error(error);
      }
    };

    fetchVolumes();
  }, []);

  const handleSearch = (searchValue: string) => {
    setFilteredVolumes(
      volumes.filter((volume) =>
        volume.toLowerCase().includes(searchValue.toLowerCase())
      )
    );
  };

  return (
    <AutoComplete
      placeholder="Volume"
      options={filteredVolumes.map((volume) => ({ value: volume }))}
      value={value}
      onChange={onChange}
      onSearch={handleSearch}
      filterOption={false}
    />
  );
};

export default VolumeSelectWidget;
export type { VolumeSelectWidgetProps };
