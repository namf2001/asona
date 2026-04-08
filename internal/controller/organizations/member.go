package organizations

import (
	"context"
)

func (i impl) AddMember(ctx context.Context, oid, uid int64, role string) error {
	return nil
}

func (i impl) RemoveMember(ctx context.Context, oid, uid int64) error {
	return nil
}
