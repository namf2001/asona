package sessions

import (
	"asona/internal/model"
	"context"
)

func (i impl) GetByToken(ctx context.Context, token string) (model.Session, error) {
	return model.Session{}, nil
}
