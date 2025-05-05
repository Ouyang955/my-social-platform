# 启动流程指南（MySQL + 后端 + 前端）

本项目启动顺序建议：**先启动 MySQL 数据库，再启动后端服务，最后启动前端服务。**

---

## 1. 启动 MySQL 数据库

### Windows 下常用方法：

#### 方法一：使用服务管理器
1. 按 `Win + R` 输入 `services.msc` 回车。
2. 找到 `MySQL` 或 `MySQL80` 服务。
3. 右键点击，选择"启动"。

#### 方法二：命令行启动
```powershell
# 打开 PowerShell 或 CMD，输入：
net start mysql
# 或
net start mysql80
```

#### 方法三：手动启动（如果是 zip 版安装）
```powershell
# 进入 MySQL 安装目录的 bin 目录
cd C:\path\to\mysql\bin
mysqld
```

#### 检查 MySQL 是否启动
```powershell
# 登录测试
mysql -u root -p
# 如果无密码直接回车
```

---

## 2. 启动后端服务（Go 项目）

1. 打开命令行，进入项目根目录：
```powershell
cd D:\vsworkspace\my-social-platform
```
2. 启动后端服务：
```powershell
go run cmd/main.go
```
3. 如果看到 `Server running at :8080` 或无报错即启动成功。

---

## 3. 启动前端服务（React 项目）

1. 打开新命令行窗口，进入前端目录：
```powershell
cd D:\vsworkspace\my-social-platform\frontend
```
2. 启动前端开发服务器：
```powershell
npm start
```
3. 浏览器访问：
```
http://localhost:3000
```

---

## 常见问题与解决
- **端口被占用**：检查 3306（MySQL）、8080（后端）、3000（前端）端口是否被其他程序占用。
- **数据库连接失败**：确认 MySQL 已启动，用户名/密码/数据库名配置正确。
- **依赖未安装**：
  - 后端：`go mod tidy`
  - 前端：`npm install`
- **日志文件找不到**：确认后端 logs 目录有写入权限。

---

## 一键重启建议
可以将常用命令写入批处理（.bat）或 shell 脚本，提升效率。

---

如有疑问，欢迎随时查阅本指南或联系开发者！ 