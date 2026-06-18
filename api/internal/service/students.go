package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ultrathreads/internal/domain"
	"ultrathreads/pkg/auth"
	"ultrathreads/pkg/hash"
	"ultrathreads/pkg/logger"
	"ultrathreads/pkg/otp"
)

type StudentsService struct {
	repo         StudentsRepository
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	otpGenerator otp.Generator

	modulesService        Modules
	offersService         Offers
	emailService          Emails
	lessonsService        Lessons
	studentLessonsService StudentLessons

	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
	verificationCodeLength int
}

func NewStudentsService(repo StudentsRepository, modulesService Modules, offersService Offers, lessonsService Lessons, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	emailService Emails, studentLessonsService StudentLessons, accessTTL, refreshTTL time.Duration, otpGenerator otp.Generator, verificationCodeLength int) *StudentsService {
	return &StudentsService{
		repo:                   repo,
		modulesService:         modulesService,
		offersService:          offersService,
		hasher:                 hasher,
		emailService:           emailService,
		lessonsService:         lessonsService,
		studentLessonsService:  studentLessonsService,
		tokenManager:           tokenManager,
		accessTokenTTL:         accessTTL,
		refreshTokenTTL:        refreshTTL,
		otpGenerator:           otpGenerator,
		verificationCodeLength: verificationCodeLength,
	}
}

func (s *StudentsService) SignUp(ctx context.Context, input domain.StudentSignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	student := domain.NewStudent(input.Name, input.Email, passwordHash, input.SchoolID)

	if input.Verified {
		student.MarkAsVerified()

		go s.addStudentToList(context.Background(), *student)

		return s.repo.Create(ctx, student)
	}

	verificationCode := s.otpGenerator.RandomSecret(s.verificationCodeLength)
	student.SetVerificationCode(verificationCode)

	if err := s.repo.Create(ctx, student); err != nil {
		return err
	}

	return s.emailService.SendStudentVerificationEmail(domain.VerificationEmailInput{
		Email:            student.Email,
		Name:             student.Name,
		VerificationCode: verificationCode,
		Domain:           input.SchoolDomain,
	})
}

func (s *StudentsService) SignIn(ctx context.Context, input domain.SchoolSignInInput) (domain.Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return domain.Tokens{}, err
	}

	student, err := s.repo.GetByCredentials(ctx, input.SchoolID, input.Email, passwordHash)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.Tokens{}, err
		}

		return domain.Tokens{}, err
	}

	if student.IsBlocked() {
		return domain.Tokens{}, domain.ErrStudentBlocked
	}

	return s.createSession(ctx, student.ID)
}

func (s *StudentsService) RefreshTokens(ctx context.Context, schoolID uint, refreshToken string) (domain.Tokens, error) {
	student, err := s.repo.GetByRefreshToken(ctx, schoolID, refreshToken)
	if err != nil {
		return domain.Tokens{}, err
	}

	if student.IsBlocked() {
		return domain.Tokens{}, domain.ErrStudentBlocked
	}

	return s.createSession(ctx, student.ID)
}

func (s *StudentsService) Verify(ctx context.Context, hash string) error {
	student, err := s.repo.Verify(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrVerificationCodeInvalid) {
			return domain.ErrVerificationCodeInvalid
		}

		return err
	}

	logger.Info(student)

	go s.addStudentToList(context.Background(), student)

	return nil
}

func (s *StudentsService) GetModuleContent(ctx context.Context, schoolID, studentID, moduleID uint) (domain.ModuleContent, error) {
	module, err := s.modulesService.GetWithContent(ctx, moduleID)
	if err != nil {
		return domain.ModuleContent{}, err
	}

	student, err := s.repo.GetById(ctx, schoolID, studentID)
	if err != nil {
		return domain.ModuleContent{}, err
	}

	if student.IsModuleAvailable(module) {
		return domain.ModuleContent{
			Lessons: module.Lessons,
			Survey:  module.Survey,
		}, nil
	}

	offers, err := s.offersService.GetByModule(ctx, schoolID, module.ID)
	if err != nil {
		return domain.ModuleContent{}, err
	}

	if len(offers) != 0 {
		return domain.ModuleContent{}, domain.ErrModuleIsNotAvailable
	}

	if err := s.repo.GiveAccessToModule(ctx, studentID, moduleID); err != nil {
		return domain.ModuleContent{}, err
	}

	return domain.ModuleContent{
		Lessons: module.Lessons,
		Survey:  module.Survey,
	}, nil
}

func (s *StudentsService) GetLesson(ctx context.Context, studentID, lessonID uint) (domain.Lesson, error) {
	if err := s.isLessonAvailable(ctx, studentID, lessonID); err != nil {
		return domain.Lesson{}, err
	}

	lesson, err := s.lessonsService.GetById(ctx, lessonID)
	if err != nil {
		return domain.Lesson{}, err
	}

	if err := s.studentLessonsService.SetLastOpened(ctx, studentID, lessonID); err != nil {
		return domain.Lesson{}, err
	}

	return lesson, nil
}

func (s *StudentsService) SetLessonFinished(ctx context.Context, studentID, lessonID uint) error {
	if err := s.isLessonAvailable(ctx, studentID, lessonID); err != nil {
		return err
	}

	return s.studentLessonsService.AddFinished(ctx, studentID, lessonID)
}

func (s *StudentsService) GiveAccessToOffer(ctx context.Context, studentID uint, offer domain.Offer) error {
	modules, err := s.modulesService.GetByPackages(ctx, offer.PackageIDs)
	if err != nil {
		return err
	}

	moduleIDs := make([]uint, len(modules))
	for i := range modules {
		moduleIDs[i] = modules[i].ID
	}

	return s.repo.AttachOffer(ctx, studentID, offer.ID, moduleIDs)
}

func (s *StudentsService) RemoveAccessToOffer(ctx context.Context, studentID uint, offer domain.Offer) error {
	modules, err := s.modulesService.GetByPackages(ctx, offer.PackageIDs)
	if err != nil {
		return err
	}

	moduleIDs := make([]uint, len(modules))
	for i := range modules {
		moduleIDs[i] = modules[i].ID
	}

	return s.repo.DetachOffer(ctx, studentID, offer.ID, moduleIDs)
}

func (s *StudentsService) GetById(ctx context.Context, schoolID, id uint) (domain.Student, error) {
	return s.repo.GetById(ctx, schoolID, id)
}

func (s *StudentsService) GetBySchool(ctx context.Context, schoolID uint, query domain.GetStudentsQuery) ([]domain.Student, int64, error) {
	return s.repo.GetBySchool(ctx, schoolID, query)
}

func (s *StudentsService) createSession(ctx context.Context, studentID uint) (domain.Tokens, error) {
	var (
		res domain.Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(fmt.Sprintf("%d", studentID), s.accessTokenTTL)
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

	err = s.repo.SetSession(ctx, studentID, session)

	return res, err
}

func (s *StudentsService) isLessonAvailable(ctx context.Context, studentID, lessonID uint) error {
	module, err := s.modulesService.GetByLesson(ctx, lessonID)
	if err != nil {
		return err
	}

	student, err := s.GetById(ctx, module.SchoolID, studentID)
	if err != nil {
		return err
	}

	if !student.IsModuleAvailable(module) {
		return domain.ErrModuleIsNotAvailable
	}

	return nil
}

func (s *StudentsService) addStudentToList(ctx context.Context, student domain.Student) {
	if err := s.emailService.AddStudentToList(ctx, student.Email, student.Name, student.SchoolID); err != nil {
		if err == domain.ErrSendPulseIsNotConnected {
			return
		}

		logger.Errorf("[SENDPULSE] failed to add email to the list: %s", err.Error())
	}
}
