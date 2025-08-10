import React, { useEffect, useState, ChangeEvent } from "react";
import dayjs from "dayjs";
import "antd/dist/reset.css";
import { Input } from "antd";

interface NativeDatePickerProps {
  value?: dayjs.Dayjs;
  onChange?: (date: dayjs.Dayjs) => void;
  style?: React.CSSProperties;
}

const NativeDatePicker: React.FC<NativeDatePickerProps> = ({ value, onChange, style, }) => {
  const [_, setSelectedDate] = useState<string>();

  useEffect(() => {

  }, []);

  const handleDateChange = (event: ChangeEvent<HTMLInputElement>): void => {
    const date = event.target.value; // Получаем значение из input
    const formattedDate = dayjs(date, "YYYY-MM-DD").format("DD-MM-YYYY");
    setSelectedDate(formattedDate);
    onChange && onChange(dayjs(date)); // Передаем оригинальную дату
  };

  return (
    <Input
      type="date"
      value={value?.format("YYYY-MM-DD")}
      onChange={handleDateChange}
      data-date-format="DD MMMM YYYY"
      style={{ height: "32px", ...style }}
    />
  );
};

export default NativeDatePicker;