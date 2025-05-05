package payload

import (
	"github.com/google/uuid"
	"mime/multipart"
)

type CreateProjectPayload struct {
	Name                string  `json:"name" validate:"required,max=255"`
	BorrowerID          string  `json:"borrower_id" validate:"required,max=10"`
	BorrowerName        string  `json:"borrower_name" validate:"required,max=255"`
	BorrowerMail        string  `json:"borrower_mail" validate:"required,email"`
	LoanPrincipalAmount float64 `json:"loan_principal_amount" validate:"required"`
	BorrowerRate        float64 `json:"borrower_rate" validate:"required,max=100"`
	ROIRate             float64 `json:"roi_rate" validate:"required"`
	ActorName           string  `json:"actor_name" validate:"required"`
	ActorMail           string  `json:"actor_mail" validate:"required"`
}

type ApproveProject struct {
	ProjectID          uuid.UUID `json:"project_id"`
	FieldVisitPICID    uuid.UUID `json:"field_visit_pic_id" validate:"required"`
	FieldVisitPICName  string    `json:"field_visit_pic_name" validate:"required,max=255"`
	FieldVisitPICMail  string    `json:"field_visit_pic_mail" validate:"required,email"`
	FieldVisitProofURL string    `json:"field_visit_proof_url" validate:"required"`
	ActorName          string    `json:"actor_name" validate:"required"`
	ActorMail          string    `json:"actor_mail" validate:"required"`
}

type InvestProject struct {
	ProjectID        uuid.UUID `json:"project_id"`
	InvestorID       uuid.UUID `json:"investor_id" validate:"required"`
	InvestorName     string    `json:"investor_name" validate:"required"`
	InvestorMail     string    `json:"investor_mail" validate:"required,email"`
	InvestmentAmount float64   `json:"investment_amount" validate:"required,gt=0"`
}

type DisburseProject struct {
	ProjectID               uuid.UUID
	FieldVisitPICID         uuid.UUID       `form:"field_visit_pic_id" validate:"required"`
	FieldVisitPICName       string          `form:"field_visit_pic_name" validate:"required,max=255"`
	FieldVisitPICMail       string          `form:"field_visit_pic_mail" validate:"required,email"`
	ActorName               string          `form:"actor_name" validate:"required"`
	ActorMail               string          `form:"actor_mail" validate:"required"`
	SignedAgreementDocument *multipart.File `form:"-"`
	DocumentExtension       string          `form:"-"`
}
