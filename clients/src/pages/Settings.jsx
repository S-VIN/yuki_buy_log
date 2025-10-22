import { useState, useEffect } from 'react';
import { List, Switch, Input, Button, Tag, message, Space, Card, Divider } from 'antd';
import { authFetch } from '../api';

const Settings = () => {
  const [inviteLogin, setInviteLogin] = useState('');
  const [groupMembers, setGroupMembers] = useState([]);
  const [incomingInvites, setIncomingInvites] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchGroupData();
    fetchInvites();
  }, []);

  const fetchGroupData = async () => {
    try {
      const response = await authFetch('/group');
      if (response.ok) {
        const data = await response.json();
        setGroupMembers(data.members || []);
      }
    } catch (error) {
      console.error('Failed to fetch group data:', error);
    }
  };

  const fetchInvites = async () => {
    try {
      const response = await authFetch('/invite');
      if (response.ok) {
        const data = await response.json();
        setIncomingInvites(data.invites || []);
      }
    } catch (error) {
      console.error('Failed to fetch invites:', error);
    }
  };

  const sendInvite = async () => {
    if (!inviteLogin.trim()) {
      message.warning('Введите логин пользователя');
      return;
    }

    setLoading(true);
    try {
      const response = await authFetch('/invite', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: inviteLogin }),
      });

      if (response.ok) {
        const data = await response.json();
        if (data.mutual_invite) {
          message.success('Группа создана!');
          await fetchGroupData();
          await fetchInvites();
        } else {
          message.success('Инвайт отправлен!');
        }
        setInviteLogin('');
      } else {
        const error = await response.text();
        message.error(error || 'Ошибка при отправке инвайта');
      }
    } catch (error) {
      message.error('Ошибка при отправке инвайта');
    } finally {
      setLoading(false);
    }
  };

  const acceptInvite = async (fromLogin) => {
    setLoading(true);
    try {
      const response = await authFetch('/invite', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: fromLogin }),
      });

      if (response.ok) {
        message.success('Инвайт принят! Группа создана/обновлена');
        await fetchGroupData();
        await fetchInvites();
      } else {
        const error = await response.text();
        message.error(error || 'Ошибка при принятии инвайта');
      }
    } catch (error) {
      message.error('Ошибка при принятии инвайта');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: 16 }}>
      <List>
        <List.Item actions={[<Switch key="notify" />]}>Notifications</List.Item>
        <List.Item actions={[<Switch key="dark" />]}>Dark mode</List.Item>
      </List>

      <Divider />

      <Card title="Группа и инвайты" style={{ marginTop: 16 }}>
        {groupMembers.length > 0 ? (
          <div style={{ marginBottom: 16 }}>
            <div style={{ marginBottom: 8, fontWeight: 'bold' }}>Участники группы:</div>
            <Space wrap>
              {groupMembers.map((member) => (
                <Tag color="blue" key={member.user_id}>
                  {member.login}
                </Tag>
              ))}
            </Space>
          </div>
        ) : (
          <div style={{ marginBottom: 16, color: '#888' }}>
            Вы не состоите в группе
          </div>
        )}

        {incomingInvites.length > 0 && (
          <div style={{ marginBottom: 16 }}>
            <div style={{ marginBottom: 8, fontWeight: 'bold' }}>Входящие инвайты:</div>
            {incomingInvites.map((invite) => (
              <div key={invite.id} style={{ marginBottom: 8 }}>
                <Space>
                  <span>От пользователя: <Tag color="green">{invite.from_login}</Tag></span>
                  <Button
                    type="primary"
                    size="small"
                    onClick={() => acceptInvite(invite.from_login)}
                    loading={loading}
                  >
                    Accept
                  </Button>
                </Space>
              </div>
            ))}
          </div>
        )}

        <div>
          <div style={{ marginBottom: 8, fontWeight: 'bold' }}>Отправить инвайт:</div>
          <Space.Compact style={{ width: '100%' }}>
            <Input
              placeholder="Введите логин пользователя"
              value={inviteLogin}
              onChange={(e) => setInviteLogin(e.target.value)}
              onPressEnter={sendInvite}
            />
            <Button type="primary" onClick={sendInvite} loading={loading}>
              Отправить
            </Button>
          </Space.Compact>
        </div>
      </Card>
    </div>
  );
};

export default Settings;
