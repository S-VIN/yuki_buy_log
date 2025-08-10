import { Card, Tabs, Form, Input, Button, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../stores/AuthContext.jsx';
import API_URL from '../api.js';

const Login = () => {
  const [loginForm] = Form.useForm();
  const [registerForm] = Form.useForm();
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleLogin = async (values) => {
    try {
      const res = await fetch(`${API_URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(values),
      });
      if (!res.ok) throw new Error('login failed');
      const data = await res.json();
      login(data.token);
      navigate('/');
    } catch {
      message.error('Login failed');
    }
  };

  const handleRegister = async (values) => {
    try {
      const res = await fetch(`${API_URL}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(values),
      });
      if (!res.ok) throw new Error('register failed');
      const data = await res.json();
      login(data.token);
      navigate('/');
    } catch {
      message.error('Registration failed');
    }
  };

  const items = [
    {
      key: 'login',
      label: 'Login',
      children: (
        <Form form={loginForm} layout="vertical" onFinish={handleLogin}>
          <Form.Item name="login" rules={[{ required: true, message: 'Login' }]}>
            <Input placeholder="Login" />
          </Form.Item>
          <Form.Item name="password" rules={[{ required: true, message: 'Password' }]}>
            <Input.Password placeholder="Password" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              Sign In
            </Button>
          </Form.Item>
        </Form>
      ),
    },
    {
      key: 'register',
      label: 'Register',
      children: (
        <Form form={registerForm} layout="vertical" onFinish={handleRegister}>
          <Form.Item name="login" rules={[{ required: true, message: 'Login' }]}>
            <Input placeholder="Login" />
          </Form.Item>
          <Form.Item name="password" rules={[{ required: true, message: 'Password' }]}>
            <Input.Password placeholder="Password" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              Sign Up
            </Button>
          </Form.Item>
        </Form>
      ),
    },
  ];

  return (
    <div style={{ padding: 16 }}>
      <Card>
        <Tabs defaultActiveKey="login" items={items} />
      </Card>
    </div>
  );
};

export default Login;
