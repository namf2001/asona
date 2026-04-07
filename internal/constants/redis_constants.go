package constants

// Redis database indices for different purposes.
const (
	// RedisDBSession is used for session/authentication data.
	RedisDBSession = 0

	// RedisDBWebSocket is used for WebSocket pub/sub.
	RedisDBWebSocket = 1

	// RedisDBCache is used for general caching.
	RedisDBCache = 2
)

// Redis key prefixes.
const (
	// SessionKeyPrefix is the prefix for session keys.
	SessionKeyPrefix = "session:"

	// UserOnlineKeyPrefix is the prefix for user online status.
	UserOnlineKeyPrefix = "user:online:"

	// RoomKeyPrefix is the prefix for room-related keys.
	RoomKeyPrefix = "room:"

	// TypingKeyPrefix is the prefix for typing indicator keys.
	TypingKeyPrefix = "typing:"
)
