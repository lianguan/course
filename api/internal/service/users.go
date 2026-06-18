package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ultrathreads/internal/domain"
	"ultrathreads/pkg/auth"
	"ultrathreads/pkg/hash"
	"ultrathreads/pkg/otp"
)

type UsersService struct {
	repo         UsersRepository
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	otpGenerator otp.Generator

	emailService  Emails
	schoolService Schools

	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
	verificationCodeLength int

	domain string
}

func NewUsersService(repo UsersRepository, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	emailService Emails, schoolsService Schools, accessTTL, refreshTTL time.Duration, otpGenerator otp.Generator,
	verificationCodeLength int, domain string) *UsersService {
	return &UsersService{
		repo:                   repo,
		hasher:                 hasher,
		emailService:           emailService,
		schoolService:          schoolsService,
		tokenManager:           tokenManager,
		accessTokenTTL:         accessTTL,
		refreshTokenTTL:        refreshTTL,
		otpGenerator:           otpGenerator,
		verificationCodeLength: verificationCodeLength,
		domain:                 domain,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input domain.UserSignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	verificationCode := s.otpGenerator.RandomSecret(s.verificationCodeLength)

	user := domain.User{
		Name:         input.Name,
		Password:     passwordHash,
		Phone:        input.Phone,
		Email:        input.Email,
		RegisteredAt: time.Now(),
		LastVisitAt:  time.Now(),
		Verification: domain.Verification{
			Code: verificationCode,
		},
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return err
		}

		return err
	}

	return s.emailService.SendUserVerificationEmail(domain.VerificationEmailInput{
		Email:            user.Email,
		Name:             user.Name,
		VerificationCode: verificationCode,
	})
}

func (s *UsersService) SignIn(ctx context.Context, input domain.UserSignInInput) (domain.Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return domain.Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, input.Email, passwordHash)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.Tokens{}, err
		}

		return domain.Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (domain.Tokens, error) {
	student, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return domain.Tokens{}, err
	}

	return s.createSession(ctx, student.ID)
}

func (s *UsersService) Verify(ctx context.Context, userID uint, hash string) error {
	return s.repo.Verify(ctx, userID, hash)
}

func (s *UsersService) CreateSchool(ctx context.Context, userID uint, schoolName string) (domain.School, error) {
	schoolID, err := s.schoolService.Create(ctx, schoolName)
	if err != nil {
		return domain.School{}, err
	}

	if err := s.repo.AttachSchool(ctx, userID, schoolID); err != nil {
		return domain.School{}, err
	}

	schoolDomain := s.generateSchoolDomain(schoolName)

	if err := s.schoolService.UpdateSettings(ctx, schoolID, domain.UpdateSchoolSettingsInput{
		Domains: []string{schoolDomain},
	}); err != nil {
		return domain.School{}, err
	}

	return domain.School{ID: schoolID, Settings: domain.Settings{Domains: []string{schoolDomain}}}, nil
}

func (s *UsersService) createSession(ctx context.Context, userID uint) (domain.Tokens, error) {
	var (
		res domain.Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(fmt.Sprintf("%d", userID), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, userID, session)

	return res, err
}

func (s *UsersService) generateSchoolDomain(subdomain string) string {
	return fmt.Sprintf("%s.%s", subdomain, s.domain)
}
