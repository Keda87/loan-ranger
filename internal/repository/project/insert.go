package project

import (
	"context"
	"github.com/google/uuid"
	"loan-ranger/internal/model/payload"
)

func (r Repository) Insert(ctx context.Context, data payload.CreateProjectPayload) (insertedID uuid.UUID, err error) {
	return insertedID, nil
}
