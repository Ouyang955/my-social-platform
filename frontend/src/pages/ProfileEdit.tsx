import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { 
  Form, 
  Input, 
  Button, 
  Upload, 
  Card, 
  message, 
  Spin, 
  Avatar 
} from 'antd';
import { 
  LoadingOutlined, 
  PlusOutlined, 
  UserOutlined,
  ArrowLeftOutlined 
} from '@ant-design/icons';
import './ProfileEdit.css';

interface UserProfile {
  id: number;
  username: string;
  nickname: string;
  avatar: string;
  bio: string;
}

const ProfileEdit: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(true);
  const [uploadLoading, setUploadLoading] = useState(false);
  const [imageUrl, setImageUrl] = useState<string>();
  const navigate = useNavigate();

  useEffect(() => {
    fetchUserProfile();
  }, []);

  // 获取用户资料
  const fetchUserProfile = async () => {
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

      const user = response.data.user;
      setImageUrl(user.avatar);
      
      // 设置表单初始值
      form.setFieldsValue({
        nickname: user.nickname,
        bio: user.bio
      });
      
      setLoading(false);
    } catch (error) {
      console.error('获取用户资料失败', error);
      message.error('获取用户资料失败');
      setLoading(false);
      
      if (axios.isAxiosError(error) && error.response?.status === 401) {
        localStorage.removeItem('token');
        navigate('/login');
      }
    }
  };

  // 提交表单
  const onFinish = async (values: any) => {
    const token = localStorage.getItem('token');
    if (!token) {
      message.error('您未登录，请先登录');
      navigate('/login');
      return;
    }

    try {
      const { nickname, bio } = values;
      await axios.put('/api/profile', 
        { nickname, bio, avatar: imageUrl },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      
      message.success('个人资料更新成功');
      navigate('/profile');
    } catch (error) {
      console.error('更新个人资料失败', error);
      message.error('更新个人资料失败');
    }
  };

  // 上传头像前的检查
  const beforeUpload = (file: File) => {
    const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
    if (!isJpgOrPng) {
      message.error('只能上传JPG/PNG格式的图片!');
      return false;
    }
    
    const isLt2M = file.size / 1024 / 1024 < 2;
    if (!isLt2M) {
      message.error('图片大小不能超过2MB!');
      return false;
    }
    
    return true;
  };

  // 头像上传状态变化
  const handleChange = (info: any) => {
    if (info.file.status === 'uploading') {
      setUploadLoading(true);
      return;
    }
    
    if (info.file.status === 'done') {
      setUploadLoading(false);
      setImageUrl(info.file.response.url);
      message.success('头像上传成功');
    }
  };

  // 自定义头像上传
  const uploadButton = (
    <div>
      {uploadLoading ? <LoadingOutlined /> : <PlusOutlined />}
      <div style={{ marginTop: 8 }}>上传头像</div>
    </div>
  );

  if (loading) {
    return (
      <div className="profile-edit-loading">
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div className="profile-edit-container">
      <Card 
        title={
          <div className="profile-edit-header">
            <Button 
              type="text" 
              icon={<ArrowLeftOutlined />} 
              onClick={() => navigate('/profile')}
            />
            <span>编辑个人资料</span>
          </div>
        }
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={onFinish}
        >
          <div className="avatar-upload-container">
            <Avatar 
              size={100} 
              src={imageUrl} 
              icon={!imageUrl && <UserOutlined />}
            />
            <Upload
              name="avatar"
              action="/api/upload/image"
              headers={{ Authorization: `Bearer ${localStorage.getItem('token') || ''}` }}
              showUploadList={false}
              beforeUpload={beforeUpload}
              onChange={handleChange}
              className="avatar-uploader"
            >
              <Button type="primary">更换头像</Button>
            </Upload>
          </div>

          <Form.Item
            name="nickname"
            label="昵称"
            rules={[{ max: 20, message: '昵称不能超过20个字符' }]}
          >
            <Input placeholder="填写你的昵称" />
          </Form.Item>

          <Form.Item
            name="bio"
            label="个性签名"
            rules={[{ max: 100, message: '个性签名不能超过100个字符' }]}
          >
            <Input.TextArea 
              placeholder="介绍一下自己吧" 
              rows={4} 
              maxLength={100} 
              showCount 
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              保存修改
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default ProfileEdit; 