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
	ID                   uuid.UUID           `db:"id" json:"id,omitempty"`
	Name                 string              `db:"name" json:"name,omitempty"`
	BorrowerID           string              `db:"borrower_id" json:"borrower_id,omitempty"`
	BorrowerName         string              `db:"borrower_name" json:"borrower_name,omitempty"`
	BorrowerMail         string              `db:"borrower_mail" json:"borrower_mail,omitempty"`
	BorrowerRate         float64             `db:"borrower_rate" json:"borrower_rate,omitempty"`
	BorrowerAgreementURL null.String         `db:"borrower_agreement_url" json:"borrower_agreement_url"`
	CurrentStatus        types.ProjectStatus `db:"current_status" json:"current_status,omitempty"`
	CurrentPICName       string              `db:"current_pic_name" json:"current_pic_name,omitempty"`
	CurrentPICMail       string              `db:"current_pic_mail" json:"current_pic_mail,omitempty"`
	LoanPrincipalAmount  float64             `db:"loan_principal_amount" json:"loan_principal_amount,omitempty"`
	TotalInvestedAmount  float64             `db:"total_invested_amount" json:"total_invested_amount,omitempty"`
	ROIRate              float64             `db:"roi_rate" json:"roi_rate,omitempty"`
	ApprovedAt           null.Time           `db:"approved_at" json:"approved_at"`
	DisbursedAt          null.Time           `db:"disbursed_at" json:"disbursed_at"`
	CreatedAt            time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time           `db:"updated_at" json:"updated_at"`
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

type ProjectInvestorItem struct {
	ID               uuid.UUID `db:"id"`
	ProjectID        uuid.UUID `db:"project_id"`
	InvestorID       uuid.UUID `db:"investor_id"`
	InvestorName     string    `db:"investor_name"`
	InvestorMail     string    `db:"investor_mail"`
	InvestmentAmount float64   `db:"investment_amount"`
}

type ProjectItem struct {
	ID                  uuid.UUID           `json:"id,omitempty" db:"id"`
	Name                string              `json:"name,omitempty" db:"name"`
	CurrentStatus       types.ProjectStatus `json:"current_status,omitempty" db:"current_status"`
	LoanPrincipalAmount float64             `json:"loan_principal_amount,omitempty" db:"loan_principal_amount"`
	TotalInvestedAmount float64             `json:"total_invested_amount,omitempty" db:"total_invested_amount"`
	ROIRate             float64             `json:"roi_rate,omitempty" db:"roi_rate"`
	CreatedAt           time.Time           `json:"created_at,omitempty" db:"created_at"`
}
