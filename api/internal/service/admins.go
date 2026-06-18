package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ultrathreads/internal/domain"
	"ultrathreads/pkg/auth"
	"ultrathreads/pkg/hash"
)

type AdminsService struct {
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager

	repo        AdminsRepository
	schoolRepo  SchoolsRepository
	studentRepo StudentsRepository

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAdminsService(hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	repo AdminsRepository, schoolRepo SchoolsRepository, studentRepo StudentsRepository,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *AdminsService {
	return &AdminsService{
		hasher:          hasher,
		tokenManager:    tokenManager,
		repo:            repo,
		schoolRepo:      schoolRepo,
		studentRepo:     studentRepo,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AdminsService) SignIn(ctx context.Context, input domain.SchoolSignInInput) (domain.Tokens, error) {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return domain.Tokens{}, err
	}

	admin, err := s.repo.GetByCredentials(ctx, input.SchoolID, input.Email, password)
	if err != nil {
		return domain.Tokens{}, err
	}

	return s.createSession(ctx, admin.ID)
}

func (s *AdminsService) RefreshTokens(ctx context.Context, schoolID uint, refreshToken string) (domain.Tokens, error) {
	admin, err := s.repo.GetByRefreshToken(ctx, schoolID, refreshToken)
	if err != nil {
		return domain.Tokens{}, err
	}

	return s.createSession(ctx, admin.ID)
}

func (s *AdminsService) GetCourses(ctx context.Context, schoolID uint) ([]domain.Course, error) {
	school, err := s.schoolRepo.GetById(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	return school.Courses, nil
}

func (s *AdminsService) GetCourseById(ctx context.Context, schoolID, courseID uint) (domain.Course, error) {
	school, err := s.schoolRepo.GetById(ctx, schoolID)
	if err != nil {
		return domain.Course{}, err
	}

	var searchedCourse domain.Course

	for _, course := range school.Courses {
		if course.ID == courseID {
			searchedCourse = course
		}
	}

	if searchedCourse.ID == 0 {
		return domain.Course{}, errors.New("not found")
	}

	return searchedCourse, nil
}

func (s *AdminsService) CreateStudent(ctx context.Context, inp domain.CreateStudentInput) (domain.Student, error) {
	passwordHash, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return domain.Student{}, err
	}

	student := domain.Student{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     passwordHash,
		RegisteredAt: time.Now().Unix(),
		SchoolID:     inp.SchoolID,
		Verification: domain.Verification{Verified: true},
	}
	err = s.studentRepo.Create(ctx, &student)

	return student, err
}

func (s *AdminsService) UpdateStudent(ctx context.Context, inp domain.UpdateStudentInput) error {
	return s.studentRepo.Update(ctx, inp)
}

func (s *AdminsService) DeleteStudent(ctx context.Context, schoolID, studentID uint) error {
	return s.studentRepo.Delete(ctx, schoolID, studentID)
}

func (s *AdminsService) createSession(ctx context.Context, adminID uint) (domain.Tokens, error) {
	var (
		res domain.Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(fmt.Sprintf("%d", adminID), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL).Unix(),
	}

	err = s.repo.SetSession(ctx, adminID, session)

	return res, err
}
