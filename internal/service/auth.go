package service

import (
	"errors"
	"my-social-platform/internal/dto"
	"my-social-platform/internal/model"
	"my-social-platform/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword - 哈希密码
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword - 校验密码
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// User到UserDTO的转换函数
func ToUserDTO(user *model.User) *dto.UserDTO {
	return &dto.UserDTO{
		ID:       user.ID,
		Username: user.Username,
	}
}

// Register - 注册用户
// Register - 用户注册方法
// 类似于Java中的Service层方法:
// @Service
//
//	public class AuthService {
//	    @Autowired
//	    private UserRepository userRepository;
//
//	    public User register(String username, String password) throws Exception {
//	        // 1. 密码加密
//	        String hashedPassword = passwordEncoder.encode(password);
//
//	        // 2. 创建用户对象
//	        User user = new User();
//	        user.setUsername(username);
//	        user.setPassword(hashedPassword);
//
//	        // 3. 保存到数据库
//	        return userRepository.save(user);
//	    }
//	}
func Register(username, password string) (*dto.UserDTO, error) {
	// 1. 对密码进行加密,类似于Spring Security的passwordEncoder
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 2. 创建用户对象,类似于Java中的User实体类
	user := &model.User{
		Username: username,
		Password: hashedPassword,
	}

	// 3. 保存到数据库,类似于JPA的save方法
	// repository.DB.Create相当于userRepository.save()
	if err := repository.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return ToUserDTO(user), nil
}

// Login - 用户登录
// Login - 用户登录方法
// 类似于Java中的Service层登录方法:
// @Service
//
//	public class AuthService {
//	    @Autowired
//	    private UserRepository userRepository;
//
//	    public User login(String username, String password) throws Exception {
//	        // 1. 根据用户名查找用户
//	        User user = userRepository.findByUsername(username)
//	            .orElseThrow(() -> new Exception("User not found"));
//
//	        // 2. 验证密码
//	        if(!passwordEncoder.matches(password, user.getPassword())) {
//	            throw new Exception("Invalid password");
//	        }
//
//	        // 3. 返回用户信息
//	        return user;
//	    }
//	}
func Login(username, password string) (*dto.UserDTO, error) {
	var user model.User
	// 1. 根据用户名查找用户
	// repository.DB.Where相当于JPA的findByUsername方法
	// First方法相当于findOne或getOne方法
	if err := repository.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// 2. 验证密码
	// VerifyPassword方法类似于Spring Security的passwordEncoder.matches方法
	if !VerifyPassword(user.Password, password) {
		return nil, errors.New("invalid password")
	}

	// 3. 验证通过,返回用户信息
	return ToUserDTO(&user), nil
}

// GetUserByUsername - 根据用户名查找用户（用于JWT生成）
func GetUserByUsername(username string, user *model.User) error {
	// First函数用于查询数据库并返回第一条匹配的记录
	// 它会将查询结果填充到传入的user指针中
	// 这里user是一个空的model.User指针，First会用查询到的数据填充它
	// 如果找不到记录，会返回gorm.ErrRecordNotFound错误
	return repository.DB.Where("username = ?", username).First(user).Error
}
