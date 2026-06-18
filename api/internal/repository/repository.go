package repository

import (
	"context"

	"ultrathreads/internal/domain"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Users interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	Verify(ctx context.Context, userID uint, code string) error
	SetSession(ctx context.Context, userID uint, session domain.Session) error
	AttachSchool(ctx context.Context, userID, schoolID uint) error
}

type Schools interface {
	Create(ctx context.Context, name string) (uint, error)
	GetByDomain(ctx context.Context, domainName string) (domain.School, error)
	GetById(ctx context.Context, id uint) (domain.School, error)
	UpdateSettings(ctx context.Context, id uint, inp domain.UpdateSchoolSettingsInput) error
	SetFondyCredentials(ctx context.Context, id uint, fondy domain.Fondy) error
}

type Students interface {
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

type StudentLessons interface {
	AddFinished(ctx context.Context, studentID, lessonID uint) error
	SetLastOpened(ctx context.Context, studentID, lessonID uint) error
}

type Admins interface {
	GetByCredentials(ctx context.Context, schoolID uint, email, password string) (domain.Admin, error)
	GetByRefreshToken(ctx context.Context, schoolID uint, refreshToken string) (domain.Admin, error)
	SetSession(ctx context.Context, id uint, session domain.Session) error
	GetById(ctx context.Context, id uint) (domain.Admin, error)
}

type Courses interface {
	Create(ctx context.Context, schoolID uint, course domain.Course) (uint, error)
	Update(ctx context.Context, inp domain.UpdateCourseInput) error
	Delete(ctx context.Context, schoolID, courseID uint) error
}

type Modules interface {
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

type LessonContent interface {
	GetByLessons(ctx context.Context, lessonIDs []uint) ([]domain.LessonContent, error)
	GetByLesson(ctx context.Context, lessonID uint) (domain.LessonContent, error)
	Update(ctx context.Context, schoolID, lessonID uint, content string) error
	DeleteContent(ctx context.Context, schoolID uint, lessonIDs []uint) error
}

type Packages interface {
	Create(ctx context.Context, pkg domain.Package) (uint, error)
	Update(ctx context.Context, inp domain.UpdatePackageInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	GetByCourse(ctx context.Context, courseID uint) ([]domain.Package, error)
	GetByID(ctx context.Context, id uint) (domain.Package, error)
	GetByIDs(ctx context.Context, ids []uint) ([]domain.Package, error)
}

type Offers interface {
	Create(ctx context.Context, offer domain.Offer) (uint, error)
	Update(ctx context.Context, inp domain.UpdateOfferInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	GetBySchool(ctx context.Context, schoolID uint) ([]domain.Offer, error)
	GetByID(ctx context.Context, id uint) (domain.Offer, error)
	GetByPackages(ctx context.Context, packageIDs []uint) ([]domain.Offer, error)
	GetByIDs(ctx context.Context, ids []uint) ([]domain.Offer, error)
}

type PromoCodes interface {
	Create(ctx context.Context, promocode domain.PromoCode) (uint, error)
	Update(ctx context.Context, inp domain.UpdatePromoCodeInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	GetByCode(ctx context.Context, schoolID uint, code string) (domain.PromoCode, error)
	GetByID(ctx context.Context, schoolID, id uint) (domain.PromoCode, error)
	GetBySchool(ctx context.Context, schoolID uint) ([]domain.PromoCode, error)
}

type Orders interface {
	Create(ctx context.Context, order domain.Order) error
	AddTransaction(ctx context.Context, id uint, transaction domain.Transaction) (domain.Order, error)
	GetBySchool(ctx context.Context, schoolID uint, pagination domain.GetOrdersQuery) ([]domain.Order, int64, error)
	GetByID(ctx context.Context, id uint) (domain.Order, error)
	SetStatus(ctx context.Context, id uint, status string) error
}

type Files interface {
	Create(ctx context.Context, file domain.File) (uint, error)
	UpdateStatus(ctx context.Context, fileName string, status domain.FileStatus) error
	GetForUploading(ctx context.Context) (domain.File, error)
	UpdateStatusAndSetURL(ctx context.Context, id uint, url string) error
	GetByID(ctx context.Context, id, schoolID uint) (domain.File, error)
}

type SurveyResults interface {
	Save(ctx context.Context, results domain.SurveyResult) error
	GetAllByModule(ctx context.Context, moduleID uint, pagination *domain.PaginationQuery) ([]domain.SurveyResult, int64, error)
	GetByStudent(ctx context.Context, moduleID, studentID uint) (domain.SurveyResult, error)
}

type Repositories struct {
	Schools        Schools
	Students       Students
	StudentLessons StudentLessons
	Courses        Courses
	Modules        Modules
	Packages       Packages
	LessonContent  LessonContent
	Offers         Offers
	PromoCodes     PromoCodes
	Orders         Orders
	Admins         Admins
	Users          Users
	Files          Files
	SurveyResults  SurveyResults
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Schools:        NewSchoolsRepo(db),
		Students:       NewStudentsRepo(db),
		StudentLessons: NewStudentLessonsRepo(db),
		Courses:        NewCoursesRepo(db),
		Modules:        NewModulesRepo(db),
		LessonContent:  NewLessonContentRepo(db),
		Offers:         NewOffersRepo(db),
		PromoCodes:     NewPromocodeRepo(db),
		Orders:         NewOrdersRepo(db),
		Admins:         NewAdminsRepo(db),
		Packages:       NewPackagesRepo(db),
		Users:          NewUsersRepo(db),
		Files:          NewFilesRepo(db),
		SurveyResults:  NewSurveyResultsRepo(db),
	}
}
