package student

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/service"
	mock_service "ultrathreads/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func TestHandler_userSignUp(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUsers, input domain.UserSignUpInput)

	tests := []struct {
		name         string
		body         string
		input        domain.UserSignUpInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name: "ok",
			body: `{"name":"John Doe","email":"john@example.com","phone":"+1234567890","password":"password123"}`,
			input: domain.UserSignUpInput{
				Name:     "John Doe",
				Email:    "john@example.com",
				Phone:    "+1234567890",
				Password: "password123",
			},
			mockBehavior: func(r *mock_service.MockUsers, input domain.UserSignUpInput) {
				r.EXPECT().SignUp(context.Background(), input).Return(nil)
			},
			statusCode:   201,
			responseBody: "",
		},
		{
			name:         "invalid input body",
			body:         `{wrong}`,
			mockBehavior: func(r *mock_service.MockUsers, input domain.UserSignUpInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name: "user already exists",
			body: `{"name":"John Doe","email":"john@example.com","phone":"+1234567890","password":"password123"}`,
			input: domain.UserSignUpInput{
				Name:     "John Doe",
				Email:    "john@example.com",
				Phone:    "+1234567890",
				Password: "password123",
			},
			mockBehavior: func(r *mock_service.MockUsers, input domain.UserSignUpInput) {
				r.EXPECT().SignUp(context.Background(), input).Return(domain.ErrUserAlreadyExists)
			},
			statusCode:   400,
			responseBody: `{"message":"user with such email already exists"}`,
		},
		{
			name: "service error",
			body: `{"name":"John Doe","email":"john@example.com","phone":"+1234567890","password":"password123"}`,
			input: domain.UserSignUpInput{
				Name:     "John Doe",
				Email:    "john@example.com",
				Phone:    "+1234567890",
				Password: "password123",
			},
			mockBehavior: func(r *mock_service.MockUsers, input domain.UserSignUpInput) {
				r.EXPECT().SignUp(context.Background(), input).Return(errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			users := mock_service.NewMockUsers(c)
			tt.mockBehavior(users, tt.input)

			services := &service.Services{Users: users}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/users/sign-up", handler.userSignUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users/sign-up", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_userSignIn(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUsers, input domain.UserSignInInput)

	tests := []struct {
		name         string
		body         string
		input        domain.UserSignInInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name: "ok",
			body: `{"email":"john@example.com","password":"password123"}`,
			input: domain.UserSignInInput{
				Email:    "john@example.com",
				Password: "password123",
			},
			mockBehavior: func(r *mock_service.MockUsers, input domain.UserSignInInput) {
				r.EXPECT().SignIn(context.Background(), input).Return(domain.Tokens{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"accessToken":"access-token","refreshToken":"refresh-token"}`,
		},
		{
			name:         "invalid input body",
			body:         `{wrong}`,
			mockBehavior: func(r *mock_service.MockUsers, input domain.UserSignInInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name: "user not found",
			body: `{"email":"john@example.com","password":"password123"}`,
			input: domain.UserSignInInput{
				Email:    "john@example.com",
				Password: "password123",
			},
			mockBehavior: func(r *mock_service.MockUsers, input domain.UserSignInInput) {
				r.EXPECT().SignIn(context.Background(), input).Return(domain.Tokens{}, domain.ErrUserNotFound)
			},
			statusCode:   400,
			responseBody: `{"message":"user doesn't exists"}`,
		},
		{
			name: "service error",
			body: `{"email":"john@example.com","password":"password123"}`,
			input: domain.UserSignInInput{
				Email:    "john@example.com",
				Password: "password123",
			},
			mockBehavior: func(r *mock_service.MockUsers, input domain.UserSignInInput) {
				r.EXPECT().SignIn(context.Background(), input).Return(domain.Tokens{}, errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			users := mock_service.NewMockUsers(c)
			tt.mockBehavior(users, tt.input)

			services := &service.Services{Users: users}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/users/sign-in", handler.userSignIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users/sign-in", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_userRefresh(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUsers, token string)

	tests := []struct {
		name         string
		body         string
		token        string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:  "ok",
			body:  `{"token":"refresh-token"}`,
			token: "refresh-token",
			mockBehavior: func(r *mock_service.MockUsers, token string) {
				r.EXPECT().RefreshTokens(context.Background(), token).Return(domain.Tokens{
					AccessToken:  "new-access-token",
					RefreshToken: "new-refresh-token",
				}, nil)
			},
			statusCode:   200,
			responseBody: `{"accessToken":"new-access-token","refreshToken":"new-refresh-token"}`,
		},
		{
			name:         "invalid input body",
			body:         `{wrong}`,
			mockBehavior: func(r *mock_service.MockUsers, token string) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:  "service error",
			body:  `{"token":"refresh-token"}`,
			token: "refresh-token",
			mockBehavior: func(r *mock_service.MockUsers, token string) {
				r.EXPECT().RefreshTokens(context.Background(), token).Return(domain.Tokens{}, errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			users := mock_service.NewMockUsers(c)
			tt.mockBehavior(users, tt.token)

			services := &service.Services{Users: users}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/users/auth/refresh", handler.userRefresh)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users/auth/refresh", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_userVerify(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUsers, userId uint, code string)

	tests := []struct {
		name         string
		code         string
		userId       uint
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:   "ok",
			code:   "123456",
			userId: 1,
			mockBehavior: func(r *mock_service.MockUsers, userId uint, code string) {
				r.EXPECT().Verify(context.Background(), userId, code).Return(nil)
			},
			statusCode:   200,
			responseBody: `{"message":"success"}`,
		},
		{
			name:         "empty code",
			code:         "",
			userId:       1,
			mockBehavior: func(r *mock_service.MockUsers, userId uint, code string) {},
			statusCode:   404,
			responseBody: `404 page not found`,
		},
		{
			name:   "invalid verification code",
			code:   "123456",
			userId: 1,
			mockBehavior: func(r *mock_service.MockUsers, userId uint, code string) {
				r.EXPECT().Verify(context.Background(), userId, code).Return(domain.ErrVerificationCodeInvalid)
			},
			statusCode:   400,
			responseBody: `{"message":"verification code is invalid"}`,
		},
		{
			name:   "service error",
			code:   "123456",
			userId: 1,
			mockBehavior: func(r *mock_service.MockUsers, userId uint, code string) {
				r.EXPECT().Verify(context.Background(), userId, code).Return(errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			users := mock_service.NewMockUsers(c)
			tt.mockBehavior(users, tt.userId, tt.code)

			services := &service.Services{Users: users}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/users/verify/:code", func(c *gin.Context) {
				c.Set(userCtx, fmt.Sprintf("%d", tt.userId))
			}, handler.userVerify)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users/verify/"+tt.code, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestHandler_userCreateSchool(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUsers, userId uint, name string)

	tests := []struct {
		name         string
		body         string
		userId       uint
		schoolName   string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:       "ok",
			body:       `{"name":"Test School"}`,
			userId:     1,
			schoolName: "Test School",
			mockBehavior: func(r *mock_service.MockUsers, userId uint, name string) {
				r.EXPECT().CreateSchool(context.Background(), userId, name).Return(domain.School{
					ID:   1,
					Name: name,
				}, nil)
			},
			statusCode:   201,
			responseBody: `{"id":1,"name":"Test School","subtitle":"","description":"","registeredAt":"0001-01-01T00:00:00Z","settings":{"color":"","domains":null,"contactInfo":{"businessName":"","registrationNumber":"","address":"","email":"","phone":""},"pages":{"confidential":"","serviceAgreement":"","newsletterConsent":""},"showPaymentImages":false,"logo":"","googleAnalyticsCode":"","fondy":{"merchantId":"","merchantPassword":"","connected":false},"sendpulse":{"id":"","secret":"","listId":"","connected":false},"disableRegistration":false}}`,
		},
		{
			name:         "invalid input body",
			body:         `{wrong}`,
			userId:       1,
			mockBehavior: func(r *mock_service.MockUsers, userId uint, name string) {},
			statusCode:   400,
		},
		{
			name:       "service error",
			body:       `{"name":"Test School"}`,
			userId:     1,
			schoolName: "Test School",
			mockBehavior: func(r *mock_service.MockUsers, userId uint, name string) {
				r.EXPECT().CreateSchool(context.Background(), userId, name).Return(domain.School{}, errors.New("service error"))
			},
			statusCode:   500,
			responseBody: `{"message":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			users := mock_service.NewMockUsers(c)
			tt.mockBehavior(users, tt.userId, tt.schoolName)

			services := &service.Services{Users: users}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/users/schools", func(c *gin.Context) {
				c.Set(userCtx, fmt.Sprintf("%d", tt.userId))
			}, handler.userCreateSchool)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users/schools", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			if tt.responseBody != "" {
				assert.Equal(t, tt.responseBody, w.Body.String())
			}
		})
	}
}
