import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { TabBar } from 'antd-mobile';
import { 
  HomeOutlined, 
  CompassOutlined, 
  PlusCircleOutlined, 
  MessageOutlined, 
  UserOutlined 
} from '@ant-design/icons';
import './Navbar.css';

const Navbar: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { pathname } = location;

  const tabs = [
    {
      key: '/home',
      title: '首页',
      icon: <HomeOutlined />,
    },
    {
      key: '/discover',
      title: '发现',
      icon: <CompassOutlined />,
    },
    {
      key: '/post/create',
      title: '发布',
      icon: <PlusCircleOutlined style={{ fontSize: '32px' }} />,
    },
    {
      key: '/messages',
      title: '消息',
      icon: <MessageOutlined />,
    },
    {
      key: '/profile',
      title: '我',
      icon: <UserOutlined />,
    },
  ];

  const setRouteActive = (value: string) => {
    navigate(value);
  };

  // 不在特定页面时不显示底部导航
  const hideNavbarPaths = ['/login', '/register'];
  if (hideNavbarPaths.includes(pathname)) {
    return null;
  }

  return (
    <div className="navbar-container">
      <TabBar activeKey={pathname} onChange={(value: string) => setRouteActive(value)}>
        {tabs.map(item => (
          <TabBar.Item key={item.key} icon={item.icon} title={item.title} />
        ))}
      </TabBar>
    </div>
  );
};

export default Navbar; 