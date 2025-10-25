/* eslint-disable react/prop-types */
import { useRef, useState } from 'react';
import { Modal, Button } from 'antd';
import TagSelectWidget from './TagSelectWidget.jsx';

const BulkTagsModal = ({ open, onCancel, onAdd }) => {
  const [bulkTags, setBulkTags] = useState([]);
  const bulkTagSelectWidgetRef = useRef(null);

  const handleCancel = () => {
    setBulkTags([]);
    bulkTagSelectWidgetRef.current?.resetTags();
    onCancel();
  };

  const handleAdd = () => {
    onAdd(bulkTags);
    setBulkTags([]);
    bulkTagSelectWidgetRef.current?.resetTags();
  };

  return (
    <Modal
      title="Add Tags to All Purchases"
      open={open}
      onCancel={handleCancel}
      footer={[
        <Button key="cancel" onClick={handleCancel}>
          Cancel
        </Button>,
        <Button key="add" type="primary" onClick={handleAdd}>
          Add
        </Button>,
      ]}
    >
      <div style={{ marginTop: 16, marginBottom: 16 }}>
        <TagSelectWidget onTagChange={setBulkTags} ref={bulkTagSelectWidgetRef} />
      </div>
    </Modal>
  );
};

export default BulkTagsModal;
