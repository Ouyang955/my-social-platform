import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Menu, Avatar, Dropdown } from 'antd';
import { 
  HomeOutlined, 
  CompassOutlined, 
  MessageOutlined, 
  UserOutlined,
  SettingOutlined,
  LogoutOutlined
} from '@ant-design/icons';
import axios from 'axios';
import './Header.css';

interface User {
  id: number;
  username: string;
  nickname?: string;
  avatar?: string;
}

const Header: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    // 检查用户是否已登录
    const token = localStorage.getItem('token');
    if (token) {
      fetchUserProfile(token);
    }
  }, []);

  const fetchUserProfile = async (token: string) => {
    try {
      const response = await axios.get('/api/profile', {
        headers: { Authorization: `Bearer ${token}` }
      });
      
      setUser({
        id: response.data.user.id,
        username: response.data.user.username,
        nickname: response.data.user.nickname,
        avatar: response.data.user.avatar
      });
      setIsAuthenticated(true);
    } catch (error) {
      localStorage.removeItem('token');
      setIsAuthenticated(false);
      setUser(null);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    setIsAuthenticated(false);
    setUser(null);
    navigate('/login');
  };

  const handleMenuClick = (path: string) => {
    navigate(path);
  };

  // 用户下拉菜单选项
  const userMenuItems = [
    {
      key: 'profile',
      label: '个人主页',
      icon: <UserOutlined />,
      onClick: () => navigate('/profile')
    },
    {
      key: 'settings',
      label: '编辑资料',
      icon: <SettingOutlined />,
      onClick: () => navigate('/profile/edit')
    },
    {
      key: 'logout',
      label: '退出登录',
      icon: <LogoutOutlined />,
      onClick: handleLogout
    }
  ];

  return (
    <header className="app-header">
      <div className="header-content">
        <div className="logo" onClick={() => navigate('/')}>
          校园社区
        </div>
        
        <Menu 
          mode="horizontal" 
          selectedKeys={[location.pathname]}
          className="nav-menu"
        >
          <Menu.Item key="/" icon={<HomeOutlined />} onClick={() => handleMenuClick('/')}>
            首页
          </Menu.Item>
          <Menu.Item key="/discover" icon={<CompassOutlined />} onClick={() => handleMenuClick('/discover')}>
            发现
          </Menu.Item>
          {isAuthenticated && (
            <Menu.Item key="/messages" icon={<MessageOutlined />} onClick={() => handleMenuClick('/messages')}>
              消息
            </Menu.Item>
          )}
        </Menu>
        
        <div className="user-area">
          {isAuthenticated ? (
            <Dropdown 
              menu={{ items: userMenuItems }}
              placement="bottomRight"
              arrow
            >
              <div className="user-avatar">
                <Avatar 
                  src={user?.avatar} 
                  icon={!user?.avatar && <UserOutlined />}
                  size="large"
                >
                  {!user?.avatar && (user?.nickname || user?.username)?.substring(0, 1).toUpperCase()}
                </Avatar>
              </div>
            </Dropdown>
          ) : (
            <button 
              className="login-btn"
              onClick={() => navigate('/login')}
            >
              登录
            </button>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header; 