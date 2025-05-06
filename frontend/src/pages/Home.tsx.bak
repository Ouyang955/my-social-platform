import React from 'react';
import { Layout, Menu, Button } from 'antd';
import { useNavigate } from 'react-router-dom';
import './Home.css';

const { Header, Content } = Layout;

const Home: React.FC = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/');
  };

  return (
    <Layout className="layout">
      <Header className="header">
        <div className="logo">校园智能社交平台</div>
        <div className="header-right">
          <Button type="primary" danger onClick={handleLogout}>
            退出登录
          </Button>
        </div>
      </Header>
      <Content className="content">
        <div className="content-container">
          <h2>欢迎来到校园智能社交平台</h2>
          <p>这里是内容区域，后续会添加帖子列表等功能</p>
        </div>
      </Content>
    </Layout>
  );
};

export default Home; 