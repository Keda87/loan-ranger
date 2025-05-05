package types

type ProjectStatus string

func (m ProjectStatus) String() string {
	return string(m)
}

func (m ProjectStatus) Next() ProjectStatus {
	stateMap := map[string]ProjectStatus{
		StatusProposed.String(): StatusApproved,
		StatusApproved.String(): StatusInvested,
		StatusInvested.String(): StatusDisbursed,
	}
	return stateMap[string(m)]
}

const (
	StatusProposed  ProjectStatus = "proposed"
	StatusApproved  ProjectStatus = "approved"
	StatusInvested  ProjectStatus = "invested"
	StatusDisbursed ProjectStatus = "disbursed"
)
