package verification_tokens

import "errors"

var ErrTokenNotFoundOrExpired = errors.New("verification token not found or already expired")
