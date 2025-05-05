package db

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
	"loan-ranger/internal/pkg/types"
	"time"
)

type CreateProjectHistory struct {
	ProjectID uuid.UUID           `json:"project_id"`
	Status    types.ProjectStatus `json:"status"`
	PICName   string              `json:"pic_name"`
	PICMail   string              `json:"pic_mail"`
	Extra     map[string]string   `json:"extra,omitempty"`
}

type ProjectDetail struct {
	ID                   uuid.UUID           `db:"id"`
	Name                 string              `db:"name"`
	BorrowerID           string              `db:"borrower_id"`
	BorrowerName         string              `db:"borrower_name"`
	BorrowerMail         string              `db:"borrower_mail"`
	BorrowerRate         float64             `db:"borrower_rate"`
	BorrowerAgreementURL null.String         `db:"borrower_agreement_url"`
	CurrentStatus        types.ProjectStatus `db:"current_status"`
	CurrentPICName       string              `db:"current_pic_name"`
	CurrentPICMail       string              `db:"current_pic_mail"`
	LoanPrincipalAmount  float64             `db:"loan_principal_amount"`
	TotalInvestedAmount  float64             `db:"total_invested_amount"`
	ROIRate              float64             `db:"roi_rate"`
	ApprovedAt           null.Time           `db:"approved_at"`
	DisbursedAt          null.Time           `db:"disbursed_at"`
	CreatedAt            time.Time           `db:"created_at"`
	UpdatedAt            time.Time           `db:"updated_at"`
}

type UpdateProject struct {
	CurrentStatus        types.ProjectStatus `db:"current_status"`
	LastUpdatedAt        time.Time           `db:"last_updated_at"`
	CurrentPICName       null.String         `db:"current_pic_name"`
	CurrentPICMail       null.String         `db:"current_pic_mail"`
	BorrowerAgreementURL null.String         `db:"borrower_agreement_url"`
	TotalInvestedAmount  null.Float          `db:"total_invested_amount"`
	ApprovedAt           null.Time           `db:"approved_at"`
	DisbursedAt          null.Time           `db:"disbursed_at"`
}

type CreateProjectInvestment struct {
	ProjectID              uuid.UUID   `db:"project_id"`
	InvestorID             uuid.UUID   `db:"investor_id"`
	InvestorName           string      `db:"investor_name"`
	InvestorMail           string      `db:"investor_main"`
	InvestmentAmount       float64     `db:"investment_amount"`
	InvestmentAgreementURL null.String `db:"investment_agreement_url"`
}
