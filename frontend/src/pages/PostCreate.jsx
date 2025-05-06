import React, { useState } from "react";
import { Form, Input, Button, Upload, Select, message, Card, Spin } from "antd";
import { PlusOutlined, LoadingOutlined } from "@ant-design/icons";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import './PostCreate.css'; // 将添加样式文件

// 标签选项
const TAG_OPTIONS = [
  { label: "校园生活", value: "校园生活" },
  { label: "升学就业", value: "升学就业" },
  { label: "考试资料", value: "考试资料" },
  { label: "情感交流", value: "情感交流" },
  { label: "兴趣爱好", value: "兴趣爱好" },
  { label: "美食分享", value: "美食分享" },
];

const PostCreate = () => {
  const [fileList, setFileList] = useState([]);
  const [uploading, setUploading] = useState(false);
  const [uploadedUrls, setUploadedUrls] = useState([]); // 存储已上传图片的URL
  const [uploadingImage, setUploadingImage] = useState(false); // 单张图片上传状态
  const navigate = useNavigate();

  // 图片上传前校验
  const beforeUpload = (file) => {
    const isImage = ["image/jpeg", "image/png", "image/webp"].includes(file.type);
    if (!isImage) {
      message.error("仅支持 jpg/png/webp 格式图片");
      return Upload.LIST_IGNORE;
    }
    if (file.size / 1024 / 1024 > 32) {
      message.error("单张图片不能超过32MB");
      return Upload.LIST_IGNORE;
    }
    if (fileList.length >= 9) { // 小红书风格，最多9张
      message.error("最多上传9张图片");
      return Upload.LIST_IGNORE;
    }
    return true;
  };

  // 自定义上传方法，调用真实的图片上传接口
  const customRequest = async ({ file, onSuccess, onError }) => {
    setUploadingImage(true);
    const formData = new FormData();
    formData.append('image', file);
    
    try {
      const token = localStorage.getItem("token");
      const response = await axios.post('/api/upload/image', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          'Authorization': `Bearer ${token}`
        }
      });
      
      // 上传成功，添加到URL列表
      const imageUrl = response.data.url;
      setUploadedUrls([...uploadedUrls, imageUrl]);
      
      onSuccess(response, file);
      message.success(`${file.name} 上传成功`);
    } catch (error) {
      onError(error);
      message.error(`${file.name} 上传失败: ${error.response?.data?.error || '未知错误'}`);
    } finally {
      setUploadingImage(false);
    }
  };

  // 图片上传变更
  const handleChange = ({ fileList: newFileList }) => setFileList(newFileList);

  // 提交表单
  const onFinish = async (values) => {
    if (fileList.length === 0) {
      message.error("请至少上传一张图片");
      return;
    }
    setUploading(true);

    try {
      // 收集所有上传图片的URL
      const uploadedImageUrls = fileList
        .filter(file => file.status === 'done')
        .map(file => file.response?.url || file.response?.data?.url || '');
      
      const token = localStorage.getItem("token");
      const res = await axios.post(
        "/api/posts",
        {
          content: values.content,
          images: JSON.stringify(uploadedImageUrls), // 将图片URL数组转为JSON字符串
          tag: values.tag,
        },
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );
      message.success("发布成功！");
      // 跳转到首页
      navigate('/');
    } catch (err) {
      message.error("发布失败：" + (err.response?.data?.error || "网络错误"));
    } finally {
      setUploading(false);
    }
  };

  // 上传按钮
  const uploadButton = (
    <div>
      {uploadingImage ? <LoadingOutlined /> : <PlusOutlined />}
      <div style={{ marginTop: 8 }}>上传图片</div>
    </div>
  );

  return (
    <div className="post-create-container">
      <Card
        title={<div className="post-create-title">发布笔记</div>}
        className="post-create-card"
        bordered={false}
      >
        <Form layout="vertical" onFinish={onFinish}>
          <Form.Item 
            label={<span className="upload-label">图片 <span className="required">*</span></span>}
            required
          >
            <Upload
              listType="picture-card"
              fileList={fileList}
              onChange={handleChange}
              beforeUpload={beforeUpload}
              customRequest={customRequest}
              multiple
              maxCount={9}
              className="image-uploader"
            >
              {fileList.length >= 9 ? null : uploadButton}
            </Upload>
            <div className="upload-hint">
              支持jpg/png/webp格式，单张不超32MB，最多9张
            </div>
          </Form.Item>
          <Form.Item
            label="配文"
            name="content"
            rules={[{ required: true, message: "请输入内容" }]}
          >
            <Input.TextArea 
              rows={4} 
              maxLength={1000} 
              showCount 
              placeholder="分享你的故事..." 
              className="content-textarea"
            />
          </Form.Item>
          <Form.Item
            label="标签"
            name="tag"
            rules={[{ required: true, message: "请选择标签" }]}
          >
            <Select 
              options={TAG_OPTIONS} 
              placeholder="请选择标签" 
              className="tag-select"
            />
          </Form.Item>
          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={uploading}
              className="publish-button"
            >
              {uploading ? "发布中..." : "发布笔记"}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default PostCreate; 