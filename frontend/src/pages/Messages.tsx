import React from 'react';
import { List, Avatar, Badge, Tabs, Empty } from 'antd';
import { UserOutlined, BellOutlined, CommentOutlined, LikeOutlined } from '@ant-design/icons';
import './Messages.css';

const { TabPane } = Tabs;

// 模拟消息数据
const notifications = [
  {
    id: 1,
    type: 'like',
    title: '用户1点赞了你的帖子',
    content: '你的校园生活分享帖子获得了新的点赞',
    time: '刚刚',
    unread: true,
  },
  {
    id: 2,
    type: 'comment',
    title: '用户2评论了你的帖子',
    content: '很棒的分享！请问这是在哪个食堂拍的？',
    time: '1小时前',
    unread: true,
  },
  {
    id: 3,
    type: 'follow',
    title: '用户3关注了你',
    content: '你有一个新粉丝',
    time: '昨天',
    unread: false,
  },
];

const Messages: React.FC = () => {
  // 根据消息类型返回对应图标
  const getIconByType = (type: string) => {
    switch (type) {
      case 'like':
        return <LikeOutlined style={{ color: '#ff2442' }} />;
      case 'comment':
        return <CommentOutlined style={{ color: '#1890ff' }} />;
      case 'follow':
        return <UserOutlined style={{ color: '#52c41a' }} />;
      case 'system':
        return <BellOutlined style={{ color: '#faad14' }} />;
      default:
        return <BellOutlined />;
    }
  };

  return (
    <div className="messages-container">
      <div className="messages-header">
        <h2>消息中心</h2>
      </div>

      <Tabs defaultActiveKey="notifications" centered>
        <TabPane tab="通知" key="notifications">
          {notifications.length > 0 ? (
            <List
              itemLayout="horizontal"
              dataSource={notifications}
              renderItem={item => (
                <List.Item className={item.unread ? 'unread-message' : ''}>
                  <List.Item.Meta
                    avatar={
                      <Badge dot={item.unread}>
                        <Avatar icon={getIconByType(item.type)} />
                      </Badge>
                    }
                    title={item.title}
                    description={item.content}
                  />
                  <div className="message-time">{item.time}</div>
                </List.Item>
              )}
            />
          ) : (
            <Empty description="暂无消息" />
          )}
        </TabPane>
        <TabPane tab="私信" key="chats">
          <Empty description="暂无私信" image={Empty.PRESENTED_IMAGE_SIMPLE} />
        </TabPane>
        <TabPane tab="点赞" key="likes">
          <Empty description="暂无点赞" image={Empty.PRESENTED_IMAGE_SIMPLE} />
        </TabPane>
        <TabPane tab="评论" key="comments">
          <Empty description="暂无评论" image={Empty.PRESENTED_IMAGE_SIMPLE} />
        </TabPane>
      </Tabs>
    </div>
  );
};

export default Messages; 