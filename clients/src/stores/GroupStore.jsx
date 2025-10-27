import { makeAutoObservable, runInAction } from 'mobx';
import { fetchGroupMembers } from '../api.js';

// Цвета для номеров участников группы (1-5)
// Используем цвета, которые хорошо работают с Ant Design Tag компонентом
const MEMBER_COLORS = {
  1: 'green',
  2: 'red',
  3: 'yellow',
  4: 'blue',
  5: 'purple',
};

// Hex цвета для номеров участников (для использования в CSS)
const MEMBER_HEX_COLORS = {
  1: '#52c41a', // green
  2: '#ff4d4f', // red
  3: '#fadb14', // yellow
  4: '#1677ff', // blue
  5: '#722ed1', // purple
};

class GroupStore {
  members = [];
  loading = false;
  error = null;
  currentUserId = null;

  constructor() {
    makeAutoObservable(this);
  }

  async loadGroupMembers() {
    this.loading = true;
    this.error = null;

    try {
      const response = await fetchGroupMembers();
      const membersData = response.members || [];
      const currentUserId = response.current_user_id;

      runInAction(() => {
        this.members = membersData;
        this.currentUserId = currentUserId;
        this.loading = false;
      });
    } catch (error) {
      runInAction(() => {
        this.error = error.message;
        this.loading = false;
        // В случае ошибки (например, пользователь не в группе), очищаем участников
        this.members = [];
      });
    }
  }

  // Проверяет, находится ли пользователь в мультиюзерной группе
  get isInMultiUserGroup() {
    return this.members.length > 1;
  }

  // Получить участника по user_id
  getMemberByUserId(userId) {
    return this.members.find((m) => m.user_id === userId);
  }

  // Получить номер участника по user_id
  getMemberNumberByUserId(userId) {
    const member = this.getMemberByUserId(userId);
    return member ? member.member_number : null;
  }

  // Проверяет, является ли чек текущим пользователем
  isCurrentUserPurchase(userId) {
    return this.currentUserId === userId;
  }

  // Получить информацию об участнике для отображения (логин и цвет)
  // Возвращает null если не в мультиюзерной группе или участник не найден
  getMemberInfo(userId) {
    if (!this.isInMultiUserGroup || !userId) {
      return null;
    }

    const member = this.getMemberByUserId(userId);
    if (!member) {
      return null;
    }

    return {
      login: member.login,
      memberNumber: member.member_number,
    };
  }

  // Получить цвет для номера участника (для Ant Design Tag)
  getMemberColor(memberNumber) {
    return MEMBER_COLORS[memberNumber] || 'default';
  }

  // Получить hex цвет для номера участника (для CSS)
  getMemberHexColor(memberNumber) {
    return MEMBER_HEX_COLORS[memberNumber] || '#d9d9d9';
  }
}

export default new GroupStore();
