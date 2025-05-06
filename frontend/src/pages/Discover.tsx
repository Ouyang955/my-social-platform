import React, { useState, useEffect } from 'react';
import { Card, Avatar, Input, Tabs, List, Spin, Empty } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import axios from 'axios';
import './Discover.css';

const { TabPane } = Tabs;

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

const Discover: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [posts, setPosts] = useState<Post[]>([]);
  const [activeTab, setActiveTab] = useState('推荐');

  const categories = ['推荐', '校园', '美食', '旅行', '时尚', '宠物', '运动'];

  useEffect(() => {
    fetchPosts();
  }, [activeTab]);

  const fetchPosts = async () => {
    try {
      setLoading(true);
      const response = await axios.get('/api/posts');
      setPosts(response.data.posts || []);
      setLoading(false);
    } catch (error) {
      console.error('获取帖子失败', error);
      setLoading(false);
    }
  };

  return (
    <div className="discover-container">
      <div className="discover-header">
        <Input
          prefix={<SearchOutlined />}
          placeholder="搜索感兴趣的内容"
          className="search-input"
        />
      </div>

      <Tabs 
        activeKey={activeTab} 
        onChange={setActiveTab}
        centered
        className="category-tabs"
      >
        {categories.map(category => (
          <TabPane tab={category} key={category}>
            {loading ? (
              <div className="loading-container">
                <Spin />
              </div>
            ) : posts.length > 0 ? (
              <List
                grid={{ gutter: 16, column: 2 }}
                dataSource={posts}
                renderItem={post => (
                  <List.Item>
                    <Card
                      cover={post.images ? <img alt="帖子图片" src={post.images.split(',')[0]} /> : null}
                      hoverable
                      className="post-card"
                    >
                      <Card.Meta
                        title={post.tag}
                        description={post.content.length > 30 ? post.content.substring(0, 30) + '...' : post.content}
                        avatar={<Avatar>{post.user_id}</Avatar>}
                      />
                      <div className="post-info">
                        <span>❤️ {post.like_count}</span>
                        <span>💬 {post.comment_count}</span>
                      </div>
                    </Card>
                  </List.Item>
                )}
              />
            ) : (
              <Empty description="暂无内容" />
            )}
          </TabPane>
        ))}
      </Tabs>
    </div>
  );
};

export default Discover; 