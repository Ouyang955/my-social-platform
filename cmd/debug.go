package debug

import (
	"fmt"
	"log"
	"my-social-platform/internal/model"
	"my-social-platform/internal/repository"
)

// GenerateTestData 生成测试数据
func GenerateTestData() {
	// 初始化数据库连接
	repository.InitDB()
	defer repository.CloseDB()

	// 检查帖子数量
	var count int64
	repository.DB.Model(&model.Post{}).Count(&count)
	fmt.Printf("数据库中有 %d 条帖子\n", count)

	// 如果没有帖子，创建一些测试数据
	if count == 0 {
		fmt.Println("没有帖子数据，创建测试数据...")

		// 先检查是否有用户
		var userCount int64
		repository.DB.Model(&model.User{}).Count(&userCount)
		if userCount == 0 {
			// 创建测试用户
			testUser := model.User{
				Username: "testuser",
				Password: "password", // 实际应用中应该加密存储
				Nickname: "测试用户",
				Bio:      "这是一个测试账号",
			}
			if err := repository.DB.Create(&testUser).Error; err != nil {
				log.Fatal("创建测试用户失败:", err)
			}
			fmt.Printf("创建测试用户成功，ID: %d\n", testUser.ID)
		}

		// 获取一个用户ID用于创建帖子
		var firstUser model.User
		repository.DB.First(&firstUser)

		// 创建测试帖子
		testPosts := []model.Post{
			{
				UserID:  firstUser.ID,
				Content: "这是第一条测试帖子内容，欢迎来到校园社区！",
				Images:  "https://picsum.photos/400/300",
				Tag:     "测试",
			},
			{
				UserID:  firstUser.ID,
				Content: "这是第二条测试帖子，分享一下校园美食！",
				Images:  "https://picsum.photos/400/300",
				Tag:     "美食",
			},
			{
				UserID:  firstUser.ID,
				Content: "周末去看了电影，推荐给大家！",
				Images:  "https://picsum.photos/400/300",
				Tag:     "影视",
			},
		}

		for _, post := range testPosts {
			if err := repository.DB.Create(&post).Error; err != nil {
				log.Printf("创建帖子失败: %v\n", err)
			} else {
				fmt.Printf("创建帖子成功，ID: %d\n", post.ID)
			}
		}

		// 再次检查帖子数量
		repository.DB.Model(&model.Post{}).Count(&count)
		fmt.Printf("现在数据库中有 %d 条帖子\n", count)
	} else {
		// 如果有帖子，显示部分数据
		var posts []model.Post
		repository.DB.Limit(5).Find(&posts)

		fmt.Println("帖子数据示例:")
		for _, post := range posts {
			fmt.Printf("ID: %d, 用户ID: %d, 内容: %s, 图片: %s\n",
				post.ID, post.UserID, post.Content, post.Images)
		}
	}
}
