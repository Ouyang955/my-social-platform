import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Home from './pages/Home';
import Register from './pages/Register';
import PostCreate from './pages/PostCreate';
import Profile from './pages/Profile';
import ProfileEdit from './pages/ProfileEdit';
import Discover from './pages/Discover';
import Messages from './pages/Messages';
import Header from './components/Header';
import './App.css';

const App: React.FC = () => {
  return (
    <div className="app-container">
      <Header />
      <div className="main-content">
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/home" element={<Home />} />
          <Route path="/discover" element={<Discover />} />
          <Route path="/post/create" element={<PostCreate />} />
          <Route path="/messages" element={<Messages />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/profile/edit" element={<ProfileEdit />} />
          <Route path="/" element={<Home />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </div>
    </div>
  );
};

export default App;
 