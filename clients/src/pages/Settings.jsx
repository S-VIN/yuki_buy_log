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
      message.warning('Please enter a username');
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
          message.success('Group created!');
          await fetchGroupData();
          await fetchInvites();
        } else {
          message.success('Invite sent!');
        }
        setInviteLogin('');
      } else {
        const error = await response.text();
        message.error(error || 'Failed to send invite');
      }
    } catch (error) {
      message.error('Failed to send invite');
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
        message.success('Invite accepted! Group created/updated');
        await fetchGroupData();
        await fetchInvites();
      } else {
        const error = await response.text();
        message.error(error || 'Failed to accept invite');
      }
    } catch (error) {
      message.error('Failed to accept invite');
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

      <Card title="Groups and Invites" style={{ marginTop: 16 }}>
        {groupMembers.length > 0 ? (
          <div style={{ marginBottom: 16 }}>
            <div style={{ marginBottom: 8, fontWeight: 'bold' }}>Group Members:</div>
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
            You are not in a group
          </div>
        )}

        {incomingInvites.length > 0 && (
          <div style={{ marginBottom: 16 }}>
            <div style={{ marginBottom: 8, fontWeight: 'bold' }}>Incoming Invites:</div>
            {incomingInvites.map((invite) => (
              <div key={invite.id} style={{ marginBottom: 8 }}>
                <Space>
                  <span>From: <Tag color="green">{invite.from_login}</Tag></span>
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
          <div style={{ marginBottom: 8, fontWeight: 'bold' }}>Send Invite:</div>
          <Space.Compact style={{ width: '100%' }}>
            <Input
              placeholder="Enter username"
              value={inviteLogin}
              onChange={(e) => setInviteLogin(e.target.value)}
              onPressEnter={sendInvite}
            />
            <Button type="primary" onClick={sendInvite} loading={loading}>
              Send
            </Button>
          </Space.Compact>
        </div>
      </Card>
    </div>
  );
};

export default Settings;
