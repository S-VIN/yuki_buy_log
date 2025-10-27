// Цвета для номеров участников группы (1-5)
// Используем цвета, которые хорошо работают с Ant Design Tag компонентом
const MEMBER_COLORS = {
  1: 'green',
  2: 'red',
  3: 'gold',
  4: 'blue',
  5: 'purple',
};

// Получить цвет для номера участника
export const getMemberColor = (memberNumber) => {
  return MEMBER_COLORS[memberNumber] || 'default';
};

// Получить hex цвет для номера участника (для использования в CSS)
const MEMBER_HEX_COLORS = {
  1: '#52c41a', // green
  2: '#ff4d4f', // red
  3: '#faad14', // gold
  4: '#1677ff', // blue
  5: '#722ed1', // purple
};

export const getMemberHexColor = (memberNumber) => {
  return MEMBER_HEX_COLORS[memberNumber] || '#d9d9d9';
};
