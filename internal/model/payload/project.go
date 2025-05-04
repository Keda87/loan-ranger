package payload

import "github.com/google/uuid"

type CreateProjectPayload struct {
	Name                string    `json:"name" validate:"required,max=255"`
	BorrowerID          uuid.UUID `json:"borrower_id" validate:"required"`
	BorrowerName        string    `json:"borrower_name" validate:"required,max=255"`
	BorrowerMail        string    `json:"borrower_mail" validate:"required,email"`
	LoanPrincipalAmount float64   `json:"loan_principal_amount" validate:"required"`
	BorrowerRate        float64   `json:"borrower_rate" validate:"required,max=100"`
	ROIRate             float64   `json:"roi_rate" validate:"required"`
	ActorName           string    `json:"actor_name"`
	ActorMail           string    `json:"actor_mail"`
}
