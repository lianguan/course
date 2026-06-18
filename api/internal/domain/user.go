package domain

// User 用户实体
type User struct {
	ID           uint         // 用户ID
	Name         string       // 用户名
	Email        string       // 邮箱
	Phone        string       // 电话
	Password     string       // 密码
	RegisteredAt int64        // 注册时间（Unix 时间戳）
	LastVisitAt  int64        // 最后登录时间（Unix 时间戳）
	Verification Verification // 邮箱验证信息
	Schools      []uint       // 关联学校ID列表
}

// UserSignUpInput 用户注册输入（Service 层使用）
type UserSignUpInput struct {
	Name     string
	Email    string
	Phone    string
	Password string
}

// UserSignInInput 用户登录输入（Service 层使用）
type UserSignInInput struct {
	Email    string
	Password string
}

// Tokens 认证令牌（Service 层使用）
type Tokens struct {
	AccessToken  string
	RefreshToken string
}

// VerificationEmailInput 验证邮件输入（Service 层使用）
type VerificationEmailInput struct {
	Email            string
	Name             string
	VerificationCode string
	Domain           string
}

// StudentPurchaseSuccessfulEmailInput 学生购买成功邮件输入（Service 层使用）
type StudentPurchaseSuccessfulEmailInput struct {
	Email      string
	Name       string
	CourseName string
}
