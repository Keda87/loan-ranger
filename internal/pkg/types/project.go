package types

type ProjectStatus string

const (
	StatusProposed  ProjectStatus = "proposed"
	StatusApproved  ProjectStatus = "approved"
	StatusInvested  ProjectStatus = "invested"
	StatusDisbursed ProjectStatus = "disbursed"
)
