import { useEffect, useState } from 'react';
import { List, Input, Button, message } from 'antd';
import { authFetch } from '../api.js';

const Settings = () => {
  const [members, setMembers] = useState([]);
  const [invitations, setInvitations] = useState([]);
  const [login, setLogin] = useState('');

  const load = async () => {
    try {
      const mRes = await authFetch('/family/members');
      if (mRes.ok) {
        const data = await mRes.json();
        setMembers(data.members);
      }
      const iRes = await authFetch('/family/invitations');
      if (iRes.ok) {
        const data = await iRes.json();
        setInvitations(data.invitations);
      }
    } catch (e) {
      // ignore
    }
  };

  useEffect(() => {
    load();
  }, []);

  const sendInvite = async () => {
    if (!login) return;
    try {
      const res = await authFetch('/family/invite', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login }),
      });
      if (!res.ok) throw new Error();
      message.success('Invitation sent');
      setLogin('');
    } catch {
      message.error('Failed to invite');
    }
  };

  const respond = async (inviter, accept) => {
    try {
      const res = await authFetch('/family/respond', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: inviter, accept }),
      });
      if (!res.ok) throw new Error();
      message.success(accept ? 'Accepted' : 'Declined');
      load();
    } catch {
      message.error('Failed to respond');
    }
  };

  return (
    <div style={{ padding: 16 }}>
      <List
        header="Family members"
        dataSource={members}
        locale={{ emptyText: 'No family' }}
        renderItem={(item) => <List.Item>{item}</List.Item>}
      />
      <div style={{ marginTop: 16, display: 'flex', gap: 8 }}>
        <Input
          placeholder="Login"
          value={login}
          onChange={(e) => setLogin(e.target.value)}
        />
        <Button type="primary" onClick={sendInvite}>
          Invite
        </Button>
      </div>
      {invitations.length > 0 && (
        <List
          style={{ marginTop: 16 }}
          header="Invitations"
          dataSource={invitations}
          renderItem={(item) => (
            <List.Item
              actions={[
                <Button type="link" onClick={() => respond(item, true)} key="a">
                  Accept
                </Button>,
                <Button type="link" onClick={() => respond(item, false)} key="d">
                  Decline
                </Button>,
              ]}
            >
              {item}
            </List.Item>
          )}
        />
      )}
    </div>
  );
};

export default Settings;
