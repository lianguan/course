package domain

// Session 会话值对象
type Session struct {
	RefreshToken string // 刷新令牌
	ExpiresAt    int64  // 过期时间（Unix 时间戳）
}
