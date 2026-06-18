package service

import (
	"context"
	"time"

	"ultrathreads/internal/config"
	"ultrathreads/internal/domain"
	"ultrathreads/pkg/auth"
	"ultrathreads/pkg/cache"
	"ultrathreads/pkg/email"
	"ultrathreads/pkg/hash"
	"ultrathreads/pkg/otp"
	"ultrathreads/pkg/storage"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Users interface {
	SignUp(ctx context.Context, input domain.UserSignUpInput) error
	SignIn(ctx context.Context, input domain.UserSignInInput) (domain.Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (domain.Tokens, error)
	Verify(ctx context.Context, userID uint, hash string) error
	CreateSchool(ctx context.Context, userID uint, schoolName string) (domain.School, error)
}

type Schools interface {
	Create(ctx context.Context, name string) (uint, error)
	GetByDomain(ctx context.Context, domainName string) (domain.School, error)
	GetById(ctx context.Context, id uint) (domain.School, error)
	UpdateSettings(ctx context.Context, schoolID uint, input domain.UpdateSchoolSettingsInput) error
	ConnectFondy(ctx context.Context, input domain.ConnectFondyInput) error
	ConnectSendPulse(ctx context.Context, input domain.ConnectSendPulseInput) error
}

type Students interface {
	SignUp(ctx context.Context, input domain.StudentSignUpInput) error
	SignIn(ctx context.Context, input domain.SchoolSignInInput) (domain.Tokens, error)
	RefreshTokens(ctx context.Context, schoolID uint, refreshToken string) (domain.Tokens, error)
	Verify(ctx context.Context, hash string) error
	GetModuleContent(ctx context.Context, schoolID, studentID, moduleID uint) (domain.ModuleContent, error)
	GetLesson(ctx context.Context, studentID, lessonID uint) (domain.Lesson, error)
	SetLessonFinished(ctx context.Context, studentID, lessonID uint) error
	GiveAccessToOffer(ctx context.Context, studentID uint, offer domain.Offer) error
	RemoveAccessToOffer(ctx context.Context, studentID uint, offer domain.Offer) error
	GetById(ctx context.Context, schoolID, id uint) (domain.Student, error)
	GetBySchool(ctx context.Context, schoolID uint, query domain.GetStudentsQuery) ([]domain.Student, int64, error)
}

type StudentLessons interface {
	AddFinished(ctx context.Context, studentID, lessonID uint) error
	SetLastOpened(ctx context.Context, studentID, lessonID uint) error
}

type Admins interface {
	SignIn(ctx context.Context, input domain.SchoolSignInInput) (domain.Tokens, error)
	RefreshTokens(ctx context.Context, schoolID uint, refreshToken string) (domain.Tokens, error)
	GetCourses(ctx context.Context, schoolID uint) ([]domain.Course, error)
	GetCourseById(ctx context.Context, schoolID, courseID uint) (domain.Course, error)
	CreateStudent(ctx context.Context, inp domain.CreateStudentInput) (domain.Student, error)
	UpdateStudent(ctx context.Context, inp domain.UpdateStudentInput) error
	DeleteStudent(ctx context.Context, schoolID, studentID uint) error
}

type Files interface {
	UploadAndSaveFile(ctx context.Context, file domain.File) (string, error)
	Save(ctx context.Context, file domain.File) (uint, error)
	UpdateStatus(ctx context.Context, fileName string, status domain.FileStatus) error
	GetByID(ctx context.Context, id, schoolID uint) (domain.File, error)
	InitStorageUploaderWorkers(ctx context.Context)
}

type Emails interface {
	SendStudentVerificationEmail(domain.VerificationEmailInput) error
	SendUserVerificationEmail(domain.VerificationEmailInput) error
	SendStudentPurchaseSuccessfulEmail(domain.StudentPurchaseSuccessfulEmailInput) error
	AddStudentToList(ctx context.Context, email, name string, schoolID uint) error
}

type Courses interface {
	Create(ctx context.Context, schoolID uint, name string) (uint, error)
	Update(ctx context.Context, inp domain.UpdateCourseInput) error
	Delete(ctx context.Context, schoolID, courseID uint) error
}

type PromoCodes interface {
	Create(ctx context.Context, inp domain.CreatePromoCodeInput) (uint, error)
	Update(ctx context.Context, inp domain.UpdatePromoCodeInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	GetByCode(ctx context.Context, schoolID uint, code string) (domain.PromoCode, error)
	GetById(ctx context.Context, schoolID, id uint) (domain.PromoCode, error)
	GetBySchool(ctx context.Context, schoolID uint) ([]domain.PromoCode, error)
}

type Offers interface {
	Create(ctx context.Context, inp domain.CreateOfferInput) (uint, error)
	Update(ctx context.Context, inp domain.UpdateOfferInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	GetById(ctx context.Context, id uint) (domain.Offer, error)
	GetByModule(ctx context.Context, schoolID, moduleID uint) ([]domain.Offer, error)
	GetByCourse(ctx context.Context, courseID uint) ([]domain.Offer, error)
	GetAll(ctx context.Context, schoolID uint) ([]domain.Offer, error)
	GetByIds(ctx context.Context, ids []uint) ([]domain.Offer, error)
}

type Modules interface {
	Create(ctx context.Context, inp domain.CreateModuleInput) (uint, error)
	Update(ctx context.Context, inp domain.UpdateModuleInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	DeleteByCourse(ctx context.Context, schoolID, courseID uint) error
	GetPublishedByCourseId(ctx context.Context, courseID uint) ([]domain.Module, error)
	GetByCourseId(ctx context.Context, courseID uint) ([]domain.Module, error)
	GetById(ctx context.Context, moduleID uint) (domain.Module, error)
	GetByPackages(ctx context.Context, packageIDs []uint) ([]domain.Module, error)
	GetWithContent(ctx context.Context, moduleID uint) (domain.Module, error)
	GetByLesson(ctx context.Context, lessonID uint) (domain.Module, error)
}

type Lessons interface {
	Create(ctx context.Context, inp domain.AddLessonInput) (uint, error)
	GetById(ctx context.Context, lessonID uint) (domain.Lesson, error)
	Update(ctx context.Context, inp domain.UpdateLessonInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	DeleteContent(ctx context.Context, schoolID uint, lessonIDs []uint) error
}

type Packages interface {
	Create(ctx context.Context, inp domain.CreatePackageInput) (uint, error)
	Update(ctx context.Context, inp domain.UpdatePackageInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	GetByCourse(ctx context.Context, courseID uint) ([]domain.Package, error)
	GetById(ctx context.Context, id uint) (domain.Package, error)
	GetByIds(ctx context.Context, ids []uint) ([]domain.Package, error)
}

type Orders interface {
	Create(ctx context.Context, studentID, offerID, promocodeID uint) (uint, error)
	AddTransaction(ctx context.Context, id uint, transaction domain.Transaction) (domain.Order, error)
	GetBySchool(ctx context.Context, schoolID uint, query domain.GetOrdersQuery) ([]domain.Order, int64, error)
	GetById(ctx context.Context, id uint) (domain.Order, error)
	SetStatus(ctx context.Context, id uint, status string) error
}

type Payments interface {
	GeneratePaymentLink(ctx context.Context, orderID uint) (string, error)
	ProcessTransaction(ctx context.Context, callback interface{}) error
}

type Surveys interface {
	Create(ctx context.Context, inp domain.CreateSurveyInput) error
	Delete(ctx context.Context, schoolID, moduleID uint) error
	SaveStudentAnswers(ctx context.Context, inp domain.SaveStudentAnswersInput) error
	GetResultsByModule(ctx context.Context, moduleID uint,
		pagination *domain.PaginationQuery) ([]domain.SurveyResult, int64, error)
	GetStudentResults(ctx context.Context, moduleID, studentID uint) (domain.SurveyResult, error)
}

// Repository interfaces for dependency inversion
type SchoolsRepository interface {
	Create(ctx context.Context, name string) (uint, error)
	GetByDomain(ctx context.Context, domainName string) (domain.School, error)
	GetById(ctx context.Context, id uint) (domain.School, error)
	UpdateSettings(ctx context.Context, id uint, inp domain.UpdateSchoolSettingsInput) error
	SetFondyCredentials(ctx context.Context, id uint, fondy domain.Fondy) error
}

type StudentsRepository interface {
	Create(ctx context.Context, student *domain.Student) error
	Update(ctx context.Context, inp domain.UpdateStudentInput) error
	Delete(ctx context.Context, schoolID, studentID uint) error
	GetByCredentials(ctx context.Context, schoolID uint, email, password string) (domain.Student, error)
	GetByRefreshToken(ctx context.Context, schoolID uint, refreshToken string) (domain.Student, error)
	GetById(ctx context.Context, schoolID, id uint) (domain.Student, error)
	GetBySchool(ctx context.Context, schoolID uint, query domain.GetStudentsQuery) ([]domain.Student, int64, error)
	SetSession(ctx context.Context, studentID uint, session domain.Session) error
	GiveAccessToModule(ctx context.Context, studentID, moduleID uint) error
	AttachOffer(ctx context.Context, studentID, offerID uint, moduleIDs []uint) error
	DetachOffer(ctx context.Context, studentID, offerID uint, moduleIDs []uint) error
	Verify(ctx context.Context, code string) (domain.Student, error)
}

type StudentLessonsRepository interface {
	AddFinished(ctx context.Context, studentID, lessonID uint) error
	SetLastOpened(ctx context.Context, studentID, lessonID uint) error
}

type CoursesRepository interface {
	Create(ctx context.Context, schoolID uint, course domain.Course) (uint, error)
	Update(ctx context.Context, inp domain.UpdateCourseInput) error
	Delete(ctx context.Context, schoolID, courseID uint) error
}

type ModulesRepository interface {
	Create(ctx context.Context, module domain.Module) (uint, error)
	GetPublishedByCourseID(ctx context.Context, courseID uint) ([]domain.Module, error)
	GetByCourseID(ctx context.Context, courseID uint) ([]domain.Module, error)
	GetPublishedByID(ctx context.Context, moduleID uint) (domain.Module, error)
	GetByID(ctx context.Context, moduleID uint) (domain.Module, error)
	GetByPackages(ctx context.Context, packageIDs []uint) ([]domain.Module, error)
	Update(ctx context.Context, inp domain.UpdateModuleInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	DeleteByCourse(ctx context.Context, schoolID, courseID uint) error
	AddLesson(ctx context.Context, schoolID, id uint, lesson domain.Lesson) error
	GetByLesson(ctx context.Context, lessonID uint) (domain.Module, error)
	UpdateLesson(ctx context.Context, inp domain.UpdateLessonInput) error
	DeleteLesson(ctx context.Context, schoolID, id uint) error
	DetachPackageFromAll(ctx context.Context, schoolID, packageID uint) error
	AttachPackage(ctx context.Context, schoolID, packageID uint, modules []uint) error
	AttachSurvey(ctx context.Context, schoolID, id uint, survey domain.Survey) error
	DetachSurvey(ctx context.Context, schoolID, id uint) error
}

type LessonContentRepository interface {
	GetByLessons(ctx context.Context, lessonIDs []uint) ([]domain.LessonContent, error)
	GetByLesson(ctx context.Context, lessonID uint) (domain.LessonContent, error)
	Update(ctx context.Context, schoolID, lessonID uint, content string) error
	DeleteContent(ctx context.Context, schoolID uint, lessonIDs []uint) error
}

type PackagesRepository interface {
	Create(ctx context.Context, pkg domain.Package) (uint, error)
	Update(ctx context.Context, inp domain.UpdatePackageInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	GetByCourse(ctx context.Context, courseID uint) ([]domain.Package, error)
	GetByID(ctx context.Context, id uint) (domain.Package, error)
	GetByIDs(ctx context.Context, ids []uint) ([]domain.Package, error)
}

type OffersRepository interface {
	Create(ctx context.Context, offer domain.Offer) (uint, error)
	Update(ctx context.Context, inp domain.UpdateOfferInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	GetBySchool(ctx context.Context, schoolID uint) ([]domain.Offer, error)
	GetByID(ctx context.Context, id uint) (domain.Offer, error)
	GetByPackages(ctx context.Context, packageIDs []uint) ([]domain.Offer, error)
	GetByIDs(ctx context.Context, ids []uint) ([]domain.Offer, error)
}

type PromoCodesRepository interface {
	Create(ctx context.Context, promocode domain.PromoCode) (uint, error)
	Update(ctx context.Context, inp domain.UpdatePromoCodeInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	GetByCode(ctx context.Context, schoolID uint, code string) (domain.PromoCode, error)
	GetByID(ctx context.Context, schoolID, id uint) (domain.PromoCode, error)
	GetBySchool(ctx context.Context, schoolID uint) ([]domain.PromoCode, error)
}

type OrdersRepository interface {
	Create(ctx context.Context, order domain.Order) error
	AddTransaction(ctx context.Context, id uint, transaction domain.Transaction) (domain.Order, error)
	GetBySchool(ctx context.Context, schoolID uint, pagination domain.GetOrdersQuery) ([]domain.Order, int64, error)
	GetByID(ctx context.Context, id uint) (domain.Order, error)
	SetStatus(ctx context.Context, id uint, status string) error
}

type AdminsRepository interface {
	GetByCredentials(ctx context.Context, schoolID uint, email, password string) (domain.Admin, error)
	GetByRefreshToken(ctx context.Context, schoolID uint, refreshToken string) (domain.Admin, error)
	SetSession(ctx context.Context, id uint, session domain.Session) error
	GetById(ctx context.Context, id uint) (domain.Admin, error)
}

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	Verify(ctx context.Context, userID uint, code string) error
	SetSession(ctx context.Context, userID uint, session domain.Session) error
	AttachSchool(ctx context.Context, userID, schoolID uint) error
}

type FilesRepository interface {
	Create(ctx context.Context, file domain.File) (uint, error)
	UpdateStatus(ctx context.Context, fileName string, status domain.FileStatus) error
	GetForUploading(ctx context.Context) (domain.File, error)
	UpdateStatusAndSetURL(ctx context.Context, id uint, url string) error
	GetByID(ctx context.Context, id, schoolID uint) (domain.File, error)
}

type SurveyResultsRepository interface {
	Save(ctx context.Context, results domain.SurveyResult) error
	GetAllByModule(ctx context.Context, moduleID uint, pagination *domain.PaginationQuery) ([]domain.SurveyResult, int64, error)
	GetByStudent(ctx context.Context, moduleID, studentID uint) (domain.SurveyResult, error)
}

type Services struct {
	Schools        Schools
	Students       Students
	StudentLessons StudentLessons
	Courses        Courses
	PromoCodes     PromoCodes
	Offers         Offers
	Packages       Packages
	Modules        Modules
	Lessons        Lessons
	Payments       Payments
	Orders         Orders
	Admins         Admins
	Files          Files
	Users          Users
	Surveys        Surveys
	Emails         Emails
}

type Deps struct {
	SchoolsRepo            SchoolsRepository
	StudentsRepo           StudentsRepository
	StudentLessonsRepo     StudentLessonsRepository
	CoursesRepo            CoursesRepository
	ModulesRepo            ModulesRepository
	LessonContentRepo      LessonContentRepository
	PackagesRepo           PackagesRepository
	OffersRepo             OffersRepository
	PromoCodesRepo         PromoCodesRepository
	OrdersRepo             OrdersRepository
	AdminsRepo             AdminsRepository
	UsersRepo              UsersRepository
	FilesRepo              FilesRepository
	SurveyResultsRepo      SurveyResultsRepository
	Cache                  cache.Cache
	Hasher                 hash.PasswordHasher
	TokenManager           auth.TokenManager
	EmailSender            email.Sender
	EmailConfig            config.EmailConfig
	StorageProvider        storage.Provider
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	FondyCallbackURL       string
	CacheTTL               int64
	OtpGenerator           otp.Generator
	VerificationCodeLength int
	Environment            string
	Domain                 string
}

func NewServices(deps Deps) *Services {
	schoolsService := NewSchoolsService(deps.SchoolsRepo, deps.Cache, deps.CacheTTL)
	emailsService := NewEmailsService(deps.EmailSender, deps.EmailConfig, *schoolsService, deps.Cache)
	modulesService := NewModulesService(deps.ModulesRepo, deps.LessonContentRepo)
	coursesService := NewCoursesService(deps.CoursesRepo, modulesService)
	packagesService := NewPackagesService(deps.PackagesRepo, deps.ModulesRepo)
	offersService := NewOffersService(deps.OffersRepo, modulesService, packagesService)
	promoCodesService := NewPromoCodeService(deps.PromoCodesRepo)
	lessonsService := NewLessonsService(deps.ModulesRepo, deps.LessonContentRepo)
	studentLessonsService := NewStudentLessonsService(deps.StudentLessonsRepo)
	studentsService := NewStudentsService(deps.StudentsRepo, modulesService, offersService, lessonsService, deps.Hasher,
		deps.TokenManager, emailsService, studentLessonsService, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.OtpGenerator, deps.VerificationCodeLength)
	ordersService := NewOrdersService(deps.OrdersRepo, offersService, promoCodesService, studentsService)
	usersService := NewUsersService(deps.UsersRepo, deps.Hasher, deps.TokenManager, emailsService, schoolsService,
		deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.OtpGenerator, deps.VerificationCodeLength, deps.Domain)

	return &Services{
		Schools:        schoolsService,
		Students:       studentsService,
		StudentLessons: studentLessonsService,
		Courses:        coursesService,
		PromoCodes:     promoCodesService,
		Offers:         offersService,
		Modules:        modulesService,
		Payments: NewPaymentsService(ordersService, offersService, studentsService, emailsService, schoolsService,
			deps.FondyCallbackURL),
		Orders: ordersService,
		Admins: NewAdminsService(deps.Hasher, deps.TokenManager, deps.AdminsRepo, deps.SchoolsRepo, deps.StudentsRepo,
			deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Packages: packagesService,
		Lessons:  lessonsService,
		Files:    NewFilesService(deps.FilesRepo, deps.StorageProvider, deps.Environment),
		Users:    usersService,
		Surveys:  NewSurveysService(deps.ModulesRepo, deps.SurveyResultsRepo, deps.StudentsRepo),
		Emails:   emailsService,
	}
}
