package dto

import "ultrathreads/internal/domain"

// ===== School DTOs =====

type SchoolResponse struct {
	ID           uint             `json:"id"`
	Name         string           `json:"name"`
	Subtitle     string           `json:"subtitle"`
	Description  string           `json:"description"`
	RegisteredAt int64            `json:"registeredAt"`
	Settings     SettingsResponse `json:"settings"`
}

type SettingsResponse struct {
	Color               string             `json:"color"`
	Domains             []string           `json:"domains"`
	ContactInfo         ContactInfoResp    `json:"contactInfo"`
	Pages               PagesResp          `json:"pages"`
	ShowPaymentImages   bool               `json:"showPaymentImages"`
	Logo                string             `json:"logo"`
	GoogleAnalyticsCode string             `json:"googleAnalyticsCode"`
	DisableRegistration bool               `json:"disableRegistration"`
}

type ContactInfoResp struct {
	BusinessName       string `json:"businessName"`
	RegistrationNumber string `json:"registrationNumber"`
	Address            string `json:"address"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
}

type PagesResp struct {
	Confidential      string `json:"confidential"`
	ServiceAgreement  string `json:"serviceAgreement"`
	NewsletterConsent string `json:"newsletterConsent"`
}

func SchoolToResponse(s domain.School) SchoolResponse {
	return SchoolResponse{
		ID:           s.ID,
		Name:         s.Name,
		Subtitle:     s.Subtitle,
		Description:  s.Description,
		RegisteredAt: s.RegisteredAt,
		Settings: SettingsResponse{
			Color:   s.Settings.Color,
			Domains: s.Settings.Domains,
			ContactInfo: ContactInfoResp{
				BusinessName:       s.Settings.ContactInfo.BusinessName,
				RegistrationNumber: s.Settings.ContactInfo.RegistrationNumber,
				Address:            s.Settings.ContactInfo.Address,
				Email:              s.Settings.ContactInfo.Email,
				Phone:              s.Settings.ContactInfo.Phone,
			},
			Pages: PagesResp{
				Confidential:      s.Settings.Pages.Confidential,
				ServiceAgreement:  s.Settings.Pages.ServiceAgreement,
				NewsletterConsent: s.Settings.Pages.NewsletterConsent,
			},
			ShowPaymentImages:   s.Settings.ShowPaymentImages,
			Logo:                s.Settings.Logo,
			GoogleAnalyticsCode: s.Settings.GoogleAnalyticsCode,
			DisableRegistration: s.Settings.DisableRegistration,
		},
	}
}

func SettingsToResponse(s domain.Settings) SettingsResponse {
	return SettingsResponse{
		Color:   s.Color,
		Domains: s.Domains,
		ContactInfo: ContactInfoResp{
			BusinessName:       s.ContactInfo.BusinessName,
			RegistrationNumber: s.ContactInfo.RegistrationNumber,
			Address:            s.ContactInfo.Address,
			Email:              s.ContactInfo.Email,
			Phone:              s.ContactInfo.Phone,
		},
		Pages: PagesResp{
			Confidential:      s.Pages.Confidential,
			ServiceAgreement:  s.Pages.ServiceAgreement,
			NewsletterConsent: s.Pages.NewsletterConsent,
		},
		ShowPaymentImages:   s.ShowPaymentImages,
		Logo:                s.Logo,
		GoogleAnalyticsCode: s.GoogleAnalyticsCode,
		DisableRegistration: s.DisableRegistration,
	}
}

// ===== Survey DTOs =====

type SurveyResponse struct {
	ID        uint             `json:"id"`
	Questions []QuestionResp   `json:"questions"`
}

type QuestionResp struct {
	ID      uint   `json:"id"`
	Text    string `json:"text"`
	Type    string `json:"type"`
	Options []string `json:"options"`
}

func SurveyToResponse(s domain.Survey) SurveyResponse {
	questions := make([]QuestionResp, len(s.Questions))
	for i, q := range s.Questions {
		questions[i] = QuestionResp{
			ID:      q.ID,
			Text:    q.Text,
			Type:    q.Type,
			Options: q.Options,
		}
	}
	return SurveyResponse{
		ID:        s.ID,
		Questions: questions,
	}
}

type SurveyResultResponse struct {
	ID          uint            `json:"id"`
	Student     StudentInfoResp `json:"student"`
	ModuleID    uint            `json:"moduleId"`
	SubmittedAt int64           `json:"submittedAt"`
	Answers     []AnswerResp    `json:"answers"`
}

type AnswerResp struct {
	QuestionID uint   `json:"questionId"`
	Answer     string `json:"answer"`
}

func SurveyResultToResponse(r domain.SurveyResult) SurveyResultResponse {
	answers := make([]AnswerResp, len(r.Answers))
	for i, a := range r.Answers {
		answers[i] = AnswerResp{
			QuestionID: a.QuestionID,
			Answer:     a.Answer,
		}
	}
	return SurveyResultResponse{
		ID: r.ID,
		Student: StudentInfoResp{
			ID:    r.Student.ID,
			Name:  r.Student.Name,
			Email: r.Student.Email,
		},
		ModuleID:    r.ModuleID,
		SubmittedAt: r.SubmittedAt,
		Answers:     answers,
	}
}

func SurveyResultsToResponse(results []domain.SurveyResult) []SurveyResultResponse {
	res := make([]SurveyResultResponse, len(results))
	for i, r := range results {
		res[i] = SurveyResultToResponse(r)
	}
	return res
}

// ===== File DTOs =====

type FileResponse struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Size             int64  `json:"size"`
	ContentType      string `json:"contentType"`
	URL              string `json:"url"`
	Status           string `json:"status"`
	UploadStartedAt  int64  `json:"uploadStartedAt"`
	UploadFinishedAt int64  `json:"uploadFinishedAt"`
}

func FileToResponse(f domain.File) FileResponse {
	return FileResponse{
		ID:               f.ID,
		Name:             f.Name,
		Type:             string(f.Type),
		Size:             f.Size,
		ContentType:      f.ContentType,
		URL:              f.URL,
		Status:           string(f.Status),
		UploadStartedAt:  f.UploadStartedAt,
		UploadFinishedAt: f.UploadFinishedAt,
	}
}

func FilesToResponse(files []domain.File) []FileResponse {
	res := make([]FileResponse, len(files))
	for i, f := range files {
		res[i] = FileToResponse(f)
	}
	return res
}
