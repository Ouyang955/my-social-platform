import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, message } from 'antd';
import axios from 'axios';
import './Login.css';
// @ts-ignore
import logoImage from '../assets/scut-logo.png';

interface LoginForm {
  username: string;
  password: string;
}

const Login: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState('');

  const onFinish = async (values: LoginForm) => {
    setLoading(true);
    try {
      const response = await axios.post('/login', values);
      if (response.data.token) {
        localStorage.setItem('token', response.data.token);
        message.success('登录成功！');
        navigate('/');
      }
    } catch (error) {
      message.error('登录失败，请检查用户名和密码');
      setErrorMsg('用户名或密码错误');
      console.error('登录错误:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-container">
      <div className="login-box">
        <div className="logo-container">
          <img src={logoImage} alt="华南理工大学校徽" className="school-logo" />
        </div>
        <h1>校园智能社交平台</h1>
        <Form
          name="login"
          onFinish={onFinish}
          autoComplete="off"
          layout="vertical"
        >
          <Form.Item
            label="用户名"
            name="username"
            rules={[{ required: true, message: '请输入用户名！' }]}
          >
            <Input size="large" placeholder="请输入用户名" />
          </Form.Item>

          <Form.Item
            label="密码"
            name="password"
            rules={[{ required: true, message: '请输入密码！' }]}
          >
            <Input.Password size="large" placeholder="请输入密码" />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              size="large"
              block
              loading={loading}
            >
              登录
            </Button>
          </Form.Item>

          {errorMsg && (
            <div className="error-message" style={{ color: 'red', marginBottom: '10px', textAlign: 'center' }}>
              {errorMsg}
            </div>
          )}

          <div className="register-link">
            还没有账号？<a onClick={() => navigate('/register')}>立即注册</a>
          </div>
        </Form>
      </div>
    </div>
  );
};

export default Login; 