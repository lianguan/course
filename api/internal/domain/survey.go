package domain

// Survey 问卷值对象
type Survey struct {
	Title     string           // 问卷标题
	Questions []SurveyQuestion // 问题列表
	Required  bool             // 是否必填
}

// SurveyQuestion 问卷问题值对象
type SurveyQuestion struct {
	ID            uint     // 问题ID
	Question      string   // 问题内容
	AnswerType    string   // 答案类型
	AnswerOptions []string // 答案选项
}

// SurveyResult 问卷结果实体
type SurveyResult struct {
	ID          uint             // 结果ID
	Student     StudentInfoShort // 学生信息
	ModuleID    uint             // 所属模块ID
	SubmittedAt int64            // 提交时间（Unix 时间戳）
	Answers     []SurveyAnswer   // 答案列表
}

// SurveyAnswer 问卷答案值对象
type SurveyAnswer struct {
	QuestionID uint   // 问题ID
	Answer     string // 答案内容
}

// CreateSurveyInput 问卷创建输入（Service 层使用）
type CreateSurveyInput struct {
	ModuleID uint
	SchoolID uint
	Survey   Survey
}

// SaveStudentAnswersInput 学生答案保存输入（Service 层使用）
type SaveStudentAnswersInput struct {
	ModuleID  uint
	StudentID uint
	SchoolID  uint
	Answers   []SurveyAnswer
}
