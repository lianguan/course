package domain

import "context"

// UsersRepository 用户仓储接口
type UsersRepository interface {
	Create(ctx context.Context, user User) error
	GetByCredentials(ctx context.Context, email, password string) (User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (User, error)
	Verify(ctx context.Context, userID uint, code string) error
	SetSession(ctx context.Context, userID uint, session Session) error
	AttachSchool(ctx context.Context, userID, schoolID uint) error
}

// SchoolsRepository 学校仓储接口
type SchoolsRepository interface {
	Create(ctx context.Context, name string) (uint, error)
	GetByDomain(ctx context.Context, domainName string) (School, error)
	GetById(ctx context.Context, id uint) (School, error)
	UpdateSettings(ctx context.Context, id uint, inp UpdateSchoolSettingsInput) error
	SetFondyCredentials(ctx context.Context, id uint, fondy Fondy) error
}

// StudentsRepository 学生仓储接口
type StudentsRepository interface {
	Create(ctx context.Context, student *Student) error
	Update(ctx context.Context, inp UpdateStudentInput) error
	Delete(ctx context.Context, schoolID, studentID uint) error
	GetByCredentials(ctx context.Context, schoolID uint, email, password string) (Student, error)
	GetByRefreshToken(ctx context.Context, schoolID uint, refreshToken string) (Student, error)
	GetById(ctx context.Context, schoolID, id uint) (Student, error)
	GetBySchool(ctx context.Context, schoolID uint, query GetStudentsQuery) ([]Student, int64, error)
	SetSession(ctx context.Context, studentID uint, session Session) error
	GiveAccessToModule(ctx context.Context, studentID, moduleID uint) error
	AttachOffer(ctx context.Context, studentID, offerID uint, moduleIDs []uint) error
	DetachOffer(ctx context.Context, studentID, offerID uint, moduleIDs []uint) error
	Verify(ctx context.Context, code string) (Student, error)
}

// StudentLessonsRepository 学生课时仓储接口
type StudentLessonsRepository interface {
	AddFinished(ctx context.Context, studentID, lessonID uint) error
	SetLastOpened(ctx context.Context, studentID, lessonID uint) error
}

// AdminsRepository 管理员仓储接口
type AdminsRepository interface {
	GetByCredentials(ctx context.Context, schoolID uint, email, password string) (Admin, error)
	GetByRefreshToken(ctx context.Context, schoolID uint, refreshToken string) (Admin, error)
	SetSession(ctx context.Context, id uint, session Session) error
	GetById(ctx context.Context, id uint) (Admin, error)
}

// CoursesRepository 课程仓储接口
type CoursesRepository interface {
	Create(ctx context.Context, schoolID uint, course Course) (uint, error)
	Update(ctx context.Context, inp UpdateCourseInput) error
	Delete(ctx context.Context, schoolID, courseID uint) error
}

// ModulesRepository 模块仓储接口
type ModulesRepository interface {
	Create(ctx context.Context, module Module) (uint, error)
	GetPublishedByCourseID(ctx context.Context, courseID uint) ([]Module, error)
	GetByCourseID(ctx context.Context, courseID uint) ([]Module, error)
	GetPublishedByID(ctx context.Context, moduleID uint) (Module, error)
	GetByID(ctx context.Context, moduleID uint) (Module, error)
	GetByPackages(ctx context.Context, packageIDs []uint) ([]Module, error)
	Update(ctx context.Context, inp UpdateModuleInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	DeleteByCourse(ctx context.Context, schoolID, courseID uint) error
	AddLesson(ctx context.Context, schoolID, id uint, lesson Lesson) error
	GetByLesson(ctx context.Context, lessonID uint) (Module, error)
	UpdateLesson(ctx context.Context, inp UpdateLessonInput) error
	DeleteLesson(ctx context.Context, schoolID, id uint) error
	DetachPackageFromAll(ctx context.Context, schoolID, packageID uint) error
	AttachPackage(ctx context.Context, schoolID, packageID uint, modules []uint) error
	AttachSurvey(ctx context.Context, schoolID, id uint, survey Survey) error
	DetachSurvey(ctx context.Context, schoolID, id uint) error
}

// LessonContentRepository 课时内容仓储接口
type LessonContentRepository interface {
	GetByLessons(ctx context.Context, lessonIDs []uint) ([]LessonContent, error)
	GetByLesson(ctx context.Context, lessonID uint) (LessonContent, error)
	Update(ctx context.Context, schoolID, lessonID uint, content string) error
	DeleteContent(ctx context.Context, schoolID uint, lessonIDs []uint) error
}

// PackagesRepository 套餐仓储接口
type PackagesRepository interface {
	Create(ctx context.Context, pkg Package) (uint, error)
	Update(ctx context.Context, inp UpdatePackageInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	GetByCourse(ctx context.Context, courseID uint) ([]Package, error)
	GetByID(ctx context.Context, id uint) (Package, error)
	GetByIDs(ctx context.Context, ids []uint) ([]Package, error)
}

// OffersRepository 优惠仓储接口
type OffersRepository interface {
	Create(ctx context.Context, offer Offer) (uint, error)
	Update(ctx context.Context, inp UpdateOfferInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	GetBySchool(ctx context.Context, schoolID uint) ([]Offer, error)
	GetByID(ctx context.Context, id uint) (Offer, error)
	GetByPackages(ctx context.Context, packageIDs []uint) ([]Offer, error)
	GetByIDs(ctx context.Context, ids []uint) ([]Offer, error)
}

// PromoCodesRepository 优惠码仓储接口
type PromoCodesRepository interface {
	Create(ctx context.Context, promocode PromoCode) (uint, error)
	Update(ctx context.Context, inp UpdatePromoCodeInput) error
	Delete(ctx context.Context, schoolID, id uint) error
	GetByCode(ctx context.Context, schoolID uint, code string) (PromoCode, error)
	GetByID(ctx context.Context, schoolID, id uint) (PromoCode, error)
	GetBySchool(ctx context.Context, schoolID uint) ([]PromoCode, error)
}

// OrdersRepository 订单仓储接口
type OrdersRepository interface {
	Create(ctx context.Context, order Order) error
	AddTransaction(ctx context.Context, id uint, transaction Transaction) (Order, error)
	GetBySchool(ctx context.Context, schoolID uint, pagination GetOrdersQuery) ([]Order, int64, error)
	GetByID(ctx context.Context, id uint) (Order, error)
	SetStatus(ctx context.Context, id uint, status string) error
}

// FilesRepository 文件仓储接口
type FilesRepository interface {
	Create(ctx context.Context, file File) (uint, error)
	UpdateStatus(ctx context.Context, fileName string, status FileStatus) error
	GetForUploading(ctx context.Context) (File, error)
	UpdateStatusAndSetURL(ctx context.Context, id uint, url string) error
	GetByID(ctx context.Context, id, schoolID uint) (File, error)
}

// SurveyResultsRepository 问卷结果仓储接口
type SurveyResultsRepository interface {
	Save(ctx context.Context, results SurveyResult) error
	GetAllByModule(ctx context.Context, moduleID uint, pagination *PaginationQuery) ([]SurveyResult, int64, error)
	GetByStudent(ctx context.Context, moduleID, studentID uint) (SurveyResult, error)
}
