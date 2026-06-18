package service

import (
	"context"
	"fmt"

	"ultrathreads/internal/config"
	"ultrathreads/internal/domain"
	"ultrathreads/pkg/cache"
	emailProvider "ultrathreads/pkg/email"
	"ultrathreads/pkg/email/sendpulse"
)

const (
	verificationLinkTmpl = "https://%s/verification?code=%s" // https://<school host>/verification?code=<verification_code>
)

type EmailService struct {
	sender  emailProvider.Sender
	config  config.EmailConfig
	schools SchoolsService

	cache cache.Cache

	sendpulseClients map[uint]*sendpulse.Client
}

// Structures used for templates.
type verificationEmailInput struct {
	VerificationLink string
}

type purchaseSuccessfulEmailInput struct {
	Name       string
	CourseName string
}

func NewEmailsService(sender emailProvider.Sender, config config.EmailConfig, schools SchoolsService, cache cache.Cache) *EmailService {
	return &EmailService{
		sender:           sender,
		config:           config,
		schools:          schools,
		cache:            cache,
		sendpulseClients: make(map[uint]*sendpulse.Client),
	}
}

func (s *EmailService) SendStudentVerificationEmail(input domain.VerificationEmailInput) error {
	subject := fmt.Sprintf(s.config.Subjects.Verification, input.Name)

	templateInput := verificationEmailInput{s.createVerificationLink(input.Domain, input.VerificationCode)}
	sendInput := emailProvider.SendEmailInput{Subject: subject, To: input.Email}

	if err := sendInput.GenerateBodyFromHTML(s.config.Templates.Verification, templateInput); err != nil {
		return err
	}

	return s.sender.Send(sendInput)
}

func (s *EmailService) SendStudentPurchaseSuccessfulEmail(input domain.StudentPurchaseSuccessfulEmailInput) error {
	templateInput := purchaseSuccessfulEmailInput{Name: input.Name, CourseName: input.CourseName}
	sendInput := emailProvider.SendEmailInput{Subject: s.config.Subjects.PurchaseSuccessful, To: input.Email}

	if err := sendInput.GenerateBodyFromHTML(s.config.Templates.PurchaseSuccessful, templateInput); err != nil {
		return err
	}

	return s.sender.Send(sendInput)
}

func (s *EmailService) SendUserVerificationEmail(input domain.VerificationEmailInput) error {
	// todo implement
	return nil
}

func (s *EmailService) createVerificationLink(domain, code string) string {
	return fmt.Sprintf(verificationLinkTmpl, domain, code)
}

func (s *EmailService) AddStudentToList(ctx context.Context, email, name string, schoolID uint) error {
	// TODO refactor
	school, err := s.schools.GetById(ctx, schoolID)
	if err != nil {
		return err
	}

	if !school.Settings.SendPulse.Connected {
		return domain.ErrSendPulseIsNotConnected
	}

	client, ex := s.sendpulseClients[schoolID]
	if !ex {
		client = sendpulse.NewClient(school.Settings.SendPulse.ID, school.Settings.SendPulse.Secret, s.cache)
		s.sendpulseClients[schoolID] = client
	}

	return client.AddEmailToList(emailProvider.AddEmailInput{
		Email:  email,
		ListID: school.Settings.SendPulse.ListID,
		Variables: map[string]string{
			"Name":   name,
			"source": "registration",
		},
	})
}
