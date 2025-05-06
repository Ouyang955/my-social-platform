import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

// 定义帖子类型接口
interface Post {
  id: number;
  user_id: number;
  content: string;
  images: string;
  tag: string;
  like_count: number;
  comment_count: number;
  created_at: string;
  updated_at: string;
}

// 用户信息接口
interface User {
  username: string;
  id: number;
  avatar?: string;
  nickname?: string;
}

const categories = ['推荐', '穿搭', '美食', '彩妆', '影视', '职场', '情感', '家居', '游戏', '旅行', '健身'];

const Home: React.FC = () => {
  const [activeCategory, setActiveCategory] = useState('推荐');
  const [posts, setPosts] = useState<Post[]>([]);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  // 检查用户是否已登录
  useEffect(() => {
    // 无论是否登录都加载帖子
    fetchPosts();
    
    const token = localStorage.getItem('token');
    if (token) {
      // 已有token，尝试获取用户信息
      fetchUserProfile(token);
    } else {
      setIsAuthenticated(false);
      setUser(null);
    }
  }, []);

  // 获取用户资料
  const fetchUserProfile = async (token: string) => {
    try {
      console.log('获取用户资料...');
      const response = await axios.get('/api/profile', {
        headers: { Authorization: `Bearer ${token}` }
      });
      
      console.log('用户资料:', response.data);
      
      setUser({
        username: response.data.user.username,
        id: response.data.user.id || 0,
        avatar: response.data.user.avatar,
        nickname: response.data.user.nickname
      });
      setIsAuthenticated(true);
    } catch (error) {
      console.error('获取用户资料失败', error);
      localStorage.removeItem('token'); // 清除无效token
      setIsAuthenticated(false);
      setUser(null);
    }
  };

  // 获取所有帖子
  const fetchPosts = async () => {
    try {
      console.log('获取帖子...');
      const response = await axios.get('/api/posts');
      console.log('获取到的帖子数据:', response);
      
      if (response.data && response.data.posts) {
        console.log('设置帖子数据:', response.data.posts);
        setPosts(response.data.posts);
      } else if (response.data && Array.isArray(response.data)) {
        // 如果返回的是数组而不是对象
        console.log('设置帖子数组:', response.data);
        setPosts(response.data);
      } else {
        console.log('帖子数据格式不正确:', response.data);
        setPosts([]);
      }
    } catch (error) {
      console.error('获取帖子失败', error);
      setPosts([]);
    } finally {
      setLoading(false);
    }
  };

  // 处理登录按钮点击
  const handleLoginClick = () => {
    navigate('/login');
  };

  // 处理发布按钮点击
  const handleCreatePost = () => {
    if (isAuthenticated) {
      navigate('/post/create');
    } else {
      // 未登录则跳转到登录页
      navigate('/login');
    }
  };

  // 实现瀑布流布局
  const postsInColumns = () => {
    // 初始化5列的数组
    const columns: Post[][] = [[], [], [], [], []];
    // 根据窗口宽度确定列数
    let numColumns = 5;
    if (window.innerWidth < 1200) numColumns = 4;
    if (window.innerWidth < 900) numColumns = 3;
    if (window.innerWidth < 600) numColumns = 2;
    if (window.innerWidth < 400) numColumns = 1;
    
    // 按顺序分配帖子到各列中
    posts.forEach((post, index) => {
      const columnIndex = index % numColumns;
      columns[columnIndex].push(post);
    });

    return columns.slice(0, numColumns);
  };

  // 处理图片URL
  const getImageUrl = (post: Post) => {
    console.log('处理图片URL:', post.images);
    
    // 如果没有图片
    if (!post.images) {
      console.log('没有图片，使用默认图片');
      return 'https://via.placeholder.com/400/300?text=无图片';
    }
    
    try {
      // 检查是否是完整URL
      if (post.images.startsWith('http')) {
        console.log('使用完整URL:', post.images);
        return post.images;
      }
      
      // 检查是否是上传图片路径
      if (post.images.startsWith('/uploads/')) {
        const baseUrl = axios.defaults.baseURL || 'http://localhost:8080';
        const fullUrl = `${baseUrl}${post.images}`;
        console.log('拼接上传路径:', fullUrl);
        return fullUrl;
      }
      
      // 尝试解析JSON数组
      try {
        const images = JSON.parse(post.images);
        console.log('解析JSON成功:', images);
        
        if (Array.isArray(images) && images.length > 0) {
          let imageUrl = images[0];
          
          // 检查解析后的URL是否需要添加前缀
          if (imageUrl.startsWith('/uploads/')) {
            const baseUrl = axios.defaults.baseURL || 'http://localhost:8080';
            imageUrl = `${baseUrl}${imageUrl}`;
          }
          
          console.log('使用JSON数组中的第一张图片:', imageUrl);
          return imageUrl;
        }
      } catch (jsonError) {
        console.log('不是有效的JSON格式');
      }
      
      // 如果所有尝试都失败，直接返回原始值
      console.log('使用原始图片路径:', post.images);
      return post.images;
    } catch (e) {
      console.error('处理图片URL出错:', e);
      return 'https://via.placeholder.com/400/300?text=图片错误';
    }
  };

  // 窗口大小变化时重新计算列数
  useEffect(() => {
    const handleResize = () => {
      setPosts([...posts]); // 触发重新渲染
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [posts]);

  const columns = postsInColumns();

  // 在帖子列表的渲染部分，添加以下函数
  const getUserNickname = (userId: number) => {
    // 如果是当前用户
    if (isAuthenticated && user && user.id === userId) {
      return user.nickname || user.username;
    }
    // 默认显示用户ID
    return `用户 ${userId}`;
  };

  // 用户头像组件
  const UserAvatar: React.FC<{ user: User | null }> = ({ user }) => {
    if (!user) return null;
    
    if (user.avatar) {
      return (
        <img 
          src={user.avatar} 
          alt={user?.nickname || user?.username} 
          className="w-full h-full object-cover"
        />
      );
    }
    
    return (
      <div className="w-full h-full bg-[#1890ff] flex items-center justify-center text-white font-bold">
        {(user?.nickname || user?.username)?.substring(0, 1).toUpperCase()}
      </div>
    );
  };

  // 帖子用户头像组件
  const PostUserAvatar: React.FC<{ userId: number }> = ({ userId }) => {
    // 创建一个简单的随机颜色生成器，根据用户ID生成稳定的颜色
    const getColorByUserId = (id: number) => {
      const colors = ['#1890ff', '#52c41a', '#faad14', '#f5222d', '#722ed1'];
      return colors[id % colors.length];
    };

    return (
      <div 
        className="w-6 h-6 rounded-full flex items-center justify-center text-xs text-white"
        style={{ backgroundColor: getColorByUserId(userId) }}
      >
        {userId}
      </div>
    );
  };

  if (loading) {
    return <div className="flex justify-center items-center min-h-screen">加载中...</div>;
  }

  // 如果没有渲染内容，显示一个简单的界面
  if (!document.body.classList.contains('home-rendered')) {
    document.body.classList.add('home-rendered');
    console.log("Home组件渲染中...");
  }

  return (
    <div className="flex min-h-screen bg-gray-50">
      {/* 左侧导航栏（大屏显示） */}
      <aside className="hidden md:flex flex-col w-60 bg-white shadow-lg py-8 px-4 sticky top-0 h-screen">
        <div className="text-3xl font-bold text-[#1890ff] mb-8">校园社区</div>
        <button className="flex items-center w-full mb-4 py-2 px-4 rounded-lg hover:bg-gray-100 text-gray-700">
          <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          发现
        </button>
        <button 
          className="flex items-center w-full mb-4 py-2 px-4 rounded-lg hover:bg-gray-100 text-gray-700"
          onClick={handleCreatePost}
        >
          <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          发布
        </button>
        <button className="flex items-center w-full mb-8 py-2 px-4 rounded-lg hover:bg-gray-100 text-gray-700">
          <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
          通知
        </button>
        {/* 登录/用户信息区 */}
        <div className="bg-[#e6f7ff] rounded-lg p-4 text-center">
          {isAuthenticated ? (
            <div>
              <div className="flex items-center justify-center mb-2">
                <div 
                  className="w-10 h-10 rounded-full overflow-hidden cursor-pointer hover:opacity-80 transition"
                  onClick={() => navigate('/profile')}
                >
                  <UserAvatar user={user} />
                </div>
                <span className="ml-2 font-medium text-[#003c88]">{user?.nickname || user?.username}</span>
              </div>
              <button 
                className="w-full py-1 bg-[#1890ff] text-white rounded hover:bg-[#40a9ff] transition"
                onClick={() => {
                  localStorage.removeItem('token');
                  setIsAuthenticated(false);
                  setUser(null);
                }}
              >
                退出登录
              </button>
            </div>
          ) : (
            <div>
              <div className="font-bold text-[#003c88] mb-2">未登录</div>
              <button 
                className="w-full py-1 bg-[#1890ff] text-white rounded hover:bg-[#40a9ff] transition"
                onClick={handleLoginClick}
              >
                登录
              </button>
            </div>
          )}
        </div>
      </aside>

      {/* 右侧主内容区 */}
      <div className="flex-1 flex flex-col">
        {/* 顶部导航栏 */}
        <nav className="sticky top-[60px] z-40 bg-white shadow-sm flex items-center justify-between px-4 md:px-8 h-16">
          <div className="flex-1 max-w-xl mx-auto relative">
            <input
              className="w-full px-10 py-2 rounded-full bg-gray-100 focus:outline-none focus:ring-2 focus:ring-[#1890ff] text-sm"
              placeholder="搜索你感兴趣的内容"
            />
            <svg className="w-5 h-5 absolute left-3 top-2.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
        </nav>
        
        {/* 分类标签 */}
        <div className="bg-white shadow-sm px-4 md:px-8 py-3 flex space-x-4 overflow-x-auto">
          {categories.map(cat => (
            <button
              key={cat}
              className={`px-4 py-1.5 rounded-full whitespace-nowrap transition text-sm font-medium ${
                activeCategory === cat
                  ? 'bg-[#1890ff] text-white'
                  : 'text-gray-700 hover:bg-gray-100'
              }`}
              onClick={() => setActiveCategory(cat)}
            >
              {cat}
            </button>
          ))}
        </div>
        
        {/* 瀑布流内容区 */}
        <main className="flex-1 p-4 md:p-8">
          {posts.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              还没有帖子，赶紧发布一个吧！
            </div>
          ) : (
            <div className="flex flex-wrap -mx-2">
              {columns.map((column, colIndex) => (
                <div key={colIndex} className="px-2 w-full sm:w-1/2 lg:w-1/3 xl:w-1/4 2xl:w-1/5">
                  {column.map(post => (
                    <div key={post.id} className="bg-white rounded-xl mb-4 shadow-sm hover:shadow-md transition overflow-hidden cursor-pointer">
                      <div className="relative pb-[100%] md:pb-[75%]">
                        <img 
                          src={getImageUrl(post)} 
                          alt={post.content.substring(0, 20)} 
                          className="absolute inset-0 w-full h-full object-cover"
                          onError={(e) => {
                            console.log('图片加载失败，替换为默认图片');
                            (e.target as HTMLImageElement).src = 'https://via.placeholder.com/400/300?text=加载失败';
                          }}
                        />
                      </div>
                      <div className="p-3">
                        <h3 className="text-sm font-medium mb-2 line-clamp-2 leading-snug">{post.content}</h3>
                        <div className="flex items-center">
                          <PostUserAvatar userId={post.user_id} />
                          <span className="text-xs text-gray-600 ml-2">{getUserNickname(post.user_id)}</span>
                          <div className="flex ml-auto items-center text-gray-400 text-xs">
                            <svg className="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                              <path d="M2 10.5a1.5 1.5 0 113 0v6a1.5 1.5 0 01-3 0v-6zM6 10.333v5.43a2 2 0 001.106 1.79l.05.025A4 4 0 008.943 18h5.416a2 2 0 001.962-1.608l1.2-6A2 2 0 0015.56 8H12V4a2 2 0 00-2-2 1 1 0 00-1 1v.667a4 4 0 01-.8 2.4L6.8 7.933a4 4 0 00-.8 2.4z" />
                            </svg>
                            {post.like_count || 0}
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              ))}
            </div>
          )}
        </main>
      </div>
    </div>
  );
};

export default Home; 