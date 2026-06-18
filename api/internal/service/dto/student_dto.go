package dto

import "ultrathreads/internal/domain"

// ===== Student DTOs =====

type StudentSignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Verified bool   `json:"verified"`
}

type StudentSignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateStudentRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateStudentRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Verified *bool   `json:"verified"`
	Blocked  *bool   `json:"blocked"`
}

type StudentResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	RegisteredAt int64  `json:"registeredAt"`
	LastVisitAt  int64  `json:"lastVisitAt"`
	Verified     bool   `json:"verified"`
	Blocked      bool   `json:"blocked"`
}

func StudentToResponse(s domain.Student) StudentResponse {
	return StudentResponse{
		ID:           s.ID,
		Name:         s.Name,
		Email:        s.Email,
		RegisteredAt: s.RegisteredAt,
		LastVisitAt:  s.LastVisitAt,
		Verified:     s.Verification.Verified,
		Blocked:      s.Blocked,
	}
}

func StudentsToResponse(students []domain.Student) []StudentResponse {
	res := make([]StudentResponse, len(students))
	for i, s := range students {
		res[i] = StudentToResponse(s)
	}
	return res
}

// ===== User DTOs =====

type UserSignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserSignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	RegisteredAt int64  `json:"registeredAt"`
	LastVisitAt  int64  `json:"lastVisitAt"`
}

func UserToResponse(u domain.User) UserResponse {
	return UserResponse{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Phone:        u.Phone,
		RegisteredAt: u.RegisteredAt,
		LastVisitAt:  u.LastVisitAt,
	}
}

// ===== Auth DTOs =====

type TokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func TokensToResponse(t domain.Tokens) TokensResponse {
	return TokensResponse{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	}
}

type VerifyEmailRequest struct {
	Code string `json:"code" binding:"required"`
}

type CreateSchoolRequest struct {
	Name string `json:"name" binding:"required"`
}
