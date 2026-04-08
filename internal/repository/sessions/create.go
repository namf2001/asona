package sessions

import (
	"asona/internal/model"
	"context"
)

func (i impl) Create(ctx context.Context, session model.Session) (model.Session, error) {
	return session, nil
}
