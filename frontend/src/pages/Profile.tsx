import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { Avatar, Card, Row, Col, Typography, Tabs, Statistic, Button, Spin, Empty, message } from 'antd';
import { EditOutlined, HeartOutlined, TeamOutlined, UserOutlined } from '@ant-design/icons';
import './Profile.css';

const { Title, Paragraph } = Typography;
const { TabPane } = Tabs;

// 用户资料接口
interface UserProfile {
  id: number;
  username: string;
  nickname: string;
  avatar: string;
  bio: string;
  follow_count: number;
  fans_count: number;
  like_count: number;
}

// 帖子接口
interface Post {
  id: number;
  user_id: number;
  content: string;
  images: string;
  tag: string;
  like_count: number;
  comment_count: number;
  created_at: string;
}

const Profile: React.FC = () => {
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    fetchProfile();
  }, []);

  // 获取用户资料和帖子
  const fetchProfile = async () => {
    const token = localStorage.getItem('token');
    if (!token) {
      message.error('您未登录，请先登录');
      navigate('/login');
      return;
    }

    try {
      setLoading(true);
      const response = await axios.get('/api/profile', {
        headers: { Authorization: `Bearer ${token}` }
      });

      setProfile(response.data.user);
      setPosts(response.data.posts || []);
      setLoading(false);
    } catch (error) {
      console.error('获取个人资料失败', error);
      message.error('获取个人资料失败');
      setLoading(false);
      
      // 如果是401错误，可能是token过期，重定向到登录页
      if (axios.isAxiosError(error) && error.response?.status === 401) {
        localStorage.removeItem('token');
        navigate('/login');
      }
    }
  };

  // 渲染帖子列表
  const renderPosts = () => {
    if (posts.length === 0) {
      return <Empty description="还没有发布任何帖子" />;
    }

    return (
      <Row gutter={[16, 16]}>
        {posts.map(post => (
          <Col xs={24} sm={12} md={8} key={post.id}>
            <Card 
              hoverable 
              cover={post.images ? <img alt="帖子图片" src={post.images.split(',')[0]} /> : null}
              onClick={() => navigate(`/post/${post.id}`)}
            >
              <Card.Meta
                title={post.tag}
                description={
                  <Paragraph ellipsis={{ rows: 2 }}>
                    {post.content}
                  </Paragraph>
                }
              />
              <div className="post-stats">
                <span><HeartOutlined /> {post.like_count}</span>
                <span>评论 {post.comment_count}</span>
              </div>
            </Card>
          </Col>
        ))}
      </Row>
    );
  };

  if (loading) {
    return (
      <div className="profile-loading">
        <Spin size="large" />
      </div>
    );
  }

  if (!profile) {
    return <Empty description="用户资料不存在" />;
  }

  return (
    <div className="profile-container">
      <div className="profile-header">
        <Row gutter={[24, 24]} align="middle">
          <Col xs={24} sm={6} className="profile-avatar">
            <Avatar 
              size={100} 
              src={profile.avatar || undefined}
              icon={!profile.avatar && <UserOutlined />}
            />
          </Col>
          <Col xs={24} sm={18}>
            <Title level={4}>{profile.nickname || profile.username}</Title>
            <Paragraph type="secondary">ID: {profile.id}</Paragraph>
            <Paragraph>{profile.bio || '这个人很懒，什么都没写~'}</Paragraph>
            
            <Row gutter={16} className="profile-stats">
              <Col span={8}>
                <Statistic 
                  title="获赞" 
                  value={profile.like_count} 
                  prefix={<HeartOutlined />} 
                />
              </Col>
              <Col span={8}>
                <Statistic 
                  title="关注" 
                  value={profile.follow_count} 
                  prefix={<TeamOutlined />} 
                />
              </Col>
              <Col span={8}>
                <Statistic 
                  title="粉丝" 
                  value={profile.fans_count} 
                  prefix={<TeamOutlined />} 
                />
              </Col>
            </Row>
            
            <Button 
              type="primary" 
              icon={<EditOutlined />}
              onClick={() => navigate('/profile/edit')}
            >
              编辑资料
            </Button>
          </Col>
        </Row>
      </div>

      <div className="profile-content">
        <Tabs defaultActiveKey="posts">
          <TabPane tab="我的帖子" key="posts">
            {renderPosts()}
          </TabPane>
          <TabPane tab="我的喜欢" key="likes">
            <Empty description="暂无收藏内容" />
          </TabPane>
        </Tabs>
      </div>
    </div>
  );
};

export default Profile; 