import React, { forwardRef, useEffect, useImperativeHandle, useState } from 'react';
import { Select } from 'antd';
import axios from 'axios';
import { ApiUrl } from "../config";

interface TagSelectWidgetProps {
  onTagChange: (value: string[]) => void;
  style?: React.CSSProperties;
}

interface TagSelectWidgetRef {
  resetTags: () => void;
}

const TagSelectWidget = forwardRef<TagSelectWidgetRef, TagSelectWidgetProps>(({ onTagChange, style }, ref) => {
  const [options, setOptions] = useState<string[]>([]);
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const fetchTags = async () => {
    setIsLoading(true);
    try {
      const response = await axios.get<string[]>(ApiUrl + '/purchases/tag/unique');
      setOptions(response.data);
    } catch (error) {
      console.error('Tag downloading error', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchTags();
  }, []);

  useImperativeHandle(ref, () => ({
    resetTags: () => {
      setSelectedTags([]);
    },
  }));

  const handleTagSelect = (value: string[]) => {
    console.log('handle tag select', value);
    setSelectedTags(value);
    onTagChange(value);
  };

  return (
    <Select
      mode="tags"
      placeholder="tags"
      style={{ width: '100%', ...style }}
      onChange={handleTagSelect}
      tokenSeparators={[',']}
      loading={isLoading}
      value={selectedTags}
      options={options.map((item) => ({
        value: item,
        label: item,
      }))}
    />
  );
});

export default TagSelectWidget;