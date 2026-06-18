package domain

import "time"

type Survey struct {
	Title     string           `json:"title"`     // 问卷标题
	Questions []SurveyQuestion `json:"questions"` // 问题列表
	Required  bool             `json:"required"`  // 是否必填
}

type SurveyQuestion struct {
	ID            uint     `json:"id"`            // 问题ID
	Question      string   `json:"question"`      // 问题内容
	AnswerType    string   `json:"answerType"`    // 答案类型
	AnswerOptions []string `json:"answerOptions"` // 答案选项
}

type SurveyResult struct {
	ID          uint             `gorm:"primaryKey;autoIncrement" json:"id"` // 结果ID
	Student     StudentInfoShort `gorm:"serializer:json" json:"student"`     // 学生信息
	ModuleID    uint             `gorm:"not null;index" json:"moduleId"`     // 所属模块ID
	SubmittedAt time.Time        `gorm:"not null" json:"submittedAt"`        // 提交时间
	Answers     []SurveyAnswer   `gorm:"serializer:json" json:"answers"`     // 答案列表
}

type SurveyAnswer struct {
	QuestionID uint   `json:"questionId"` // 问题ID
	Answer     string `json:"answer"`     // 答案内容
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
