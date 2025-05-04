package db

import (
	"github.com/google/uuid"
	"loan-ranger/internal/pkg/types"
)

type CreateProjectHistory struct {
	ProjectID uuid.UUID           `json:"project_id"`
	Status    types.ProjectStatus `json:"status"`
	PICName   string              `json:"pic_name"`
	PICMail   string              `json:"pic_mail"`
	Extra     map[string]string   `json:"extra,omitempty"`
}
