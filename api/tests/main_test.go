package tests

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"ultrathreads/internal/config"
	v1 "ultrathreads/internal/delivery/http/v1"
	"ultrathreads/internal/repository"
	"ultrathreads/internal/service"
	"ultrathreads/pkg/auth"
	"ultrathreads/pkg/cache"
	"ultrathreads/pkg/database/mysql"
	emailmock "ultrathreads/pkg/email/mock"
	"ultrathreads/pkg/hash"
	"ultrathreads/pkg/otp"
	"gorm.io/gorm"
)

var dbURI string

func init() {
	dbURI = os.Getenv("TEST_DB_DSN")
}

type APITestSuite struct {
	suite.Suite

	db       *gorm.DB
	handler  *v1.Handler
	services *service.Services
	repos    *repository.Repositories

	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	mocks        *mocks
}

type mocks struct {
	emailSender  *emailmock.EmailSender
	otpGenerator *otp.MockGenerator
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	db, err := mysql.NewClient(dbURI)
	if err != nil {
		s.FailNow("Failed to connect to MySQL", err)
	}
	s.db = db

	s.initMocks()
	s.initDeps()

	if err := s.populateDB(); err != nil {
		s.FailNow("Failed to populate DB", err)
	}
}

func (s *APITestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	sqlDB.Close()
}

func (s *APITestSuite) initDeps() {
	// Init domain deps
	repos := repository.NewRepositories(s.db)
	memCache := cache.NewMemoryCache()
	hasher := hash.NewSHA1Hasher("salt")

	tokenManager, err := auth.NewManager("signing_key")
	if err != nil {
		s.FailNow("Failed to initialize token manager", err)
	}

	services := service.NewServices(service.Deps{
		SchoolsRepo:          repos.Schools,
		StudentsRepo:         repos.Students,
		StudentLessonsRepo:   repos.StudentLessons,
		CoursesRepo:          repos.Courses,
		ModulesRepo:          repos.Modules,
		LessonContentRepo:    repos.LessonContent,
		PackagesRepo:         repos.Packages,
		OffersRepo:           repos.Offers,
		PromoCodesRepo:       repos.PromoCodes,
		OrdersRepo:           repos.Orders,
		AdminsRepo:           repos.Admins,
		UsersRepo:            repos.Users,
		FilesRepo:            repos.Files,
		SurveyResultsRepo:    repos.SurveyResults,
		Cache:                memCache,
		Hasher:               hasher,
		TokenManager:         tokenManager,
		EmailSender:          s.mocks.emailSender,
		EmailConfig: config.EmailConfig{
			Templates: config.EmailTemplates{
				Verification:       "../templates/verification_email.html",
				PurchaseSuccessful: "../templates/purchase_successful.html",
			},
			Subjects: config.EmailSubjects{
				Verification:       "Спасибо за регистрацию, %s!",
				PurchaseSuccessful: "Покупка прошла успешно!",
			},
		},
		AccessTokenTTL:         time.Minute * 15,
		RefreshTokenTTL:        time.Minute * 15,
		CacheTTL:               int64(time.Minute.Seconds()),
		OtpGenerator:           s.mocks.otpGenerator,
		VerificationCodeLength: 8,
	})

	s.repos = repos
	s.services = services
	s.handler = v1.NewHandler(services, tokenManager)
	s.hasher = hasher
	s.tokenManager = tokenManager
}

func (s *APITestSuite) initMocks() {
	s.mocks = &mocks{
		emailSender:  new(emailmock.EmailSender),
		otpGenerator: new(otp.MockGenerator),
	}
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

func (s *APITestSuite) populateDB() error {
	// Create school
	if err := s.db.Create(&school).Error; err != nil {
		return err
	}

	// Create packages
	for _, pkg := range packages {
		if err := s.db.Create(pkg).Error; err != nil {
			return err
		}
	}

	// Create offers
	for _, offer := range offers {
		if err := s.db.Create(offer).Error; err != nil {
			return err
		}
	}

	// Create modules
	for _, module := range modules {
		if err := s.db.Create(module).Error; err != nil {
			return err
		}
	}

	// Create promocodes
	for _, promocode := range promocodes {
		if err := s.db.Create(promocode).Error; err != nil {
			return err
		}
	}

	return nil
}
