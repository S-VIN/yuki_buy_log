import { useState, useEffect } from 'react';
import { Input, Button, Tag, message, Space, Card } from 'antd';
import { fetchGroupMembers, fetchInvites, sendInvite as sendInviteAPI } from '../api';

const Settings = () => {
  const [inviteLogin, setInviteLogin] = useState('');
  const [groupMembers, setGroupMembers] = useState([]);
  const [incomingInvites, setIncomingInvites] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadGroupData();
    loadInvites();
  }, []);

  const loadGroupData = async () => {
    try {
      const data = await fetchGroupMembers();
      setGroupMembers(data.members || []);
    } catch (error) {
      console.error('Failed to fetch group data:', error);
    }
  };

  const loadInvites = async () => {
    try {
      const data = await fetchInvites();
      setIncomingInvites(data.invites || []);
    } catch (error) {
      console.error('Failed to fetch invites:', error);
    }
  };

  const handleSendInvite = async () => {
    if (!inviteLogin.trim()) {
      message.warning('Please enter a username');
      return;
    }

    setLoading(true);
    try {
      const data = await sendInviteAPI(inviteLogin);
      if (data.mutual_invite) {
        message.success('Group created!');
        await loadGroupData();
        await loadInvites();
      } else {
        message.success('Invite sent!');
      }
      setInviteLogin('');
    } catch (error) {
      message.error(error.message || 'Failed to send invite');
    } finally {
      setLoading(false);
    }
  };

  const handleAcceptInvite = async (fromLogin) => {
    setLoading(true);
    try {
      await sendInviteAPI(fromLogin);
      message.success('Invite accepted! Group created/updated');
      await loadGroupData();
      await loadInvites();
    } catch (error) {
      message.error(error.message || 'Failed to accept invite');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: 16 }}>
      <Card style={{ marginTop: 16 }}>
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
                    onClick={() => handleAcceptInvite(invite.from_login)}
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
              onPressEnter={handleSendInvite}
            />
            <Button type="primary" onClick={handleSendInvite} loading={loading}>
              Send
            </Button>
          </Space.Compact>
        </div>
      </Card>
    </div>
  );
};

export default Settings;
