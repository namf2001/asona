package model

// VerificationCodeType represents the purpose of a verification code
type VerificationCodeType string

const (
	VerificationCodeTypeEmailVerification VerificationCodeType = "email_verification"
	VerificationCodeTypeOrganizationJoin  VerificationCodeType = "organization_join"
)

// FriendshipStatus represents the status of a friendship between two users
type FriendshipStatus string

const (
	FriendshipStatusPending  FriendshipStatus = "pending"
	FriendshipStatusAccepted FriendshipStatus = "accepted"
	FriendshipStatusBlocked  FriendshipStatus = "blocked"
)

// WorkplaceSize represents the size category of a workplace
type WorkplaceSize string

const (
	WorkplaceSizeXS WorkplaceSize = "1-10"
	WorkplaceSizeS  WorkplaceSize = "11-50"
	WorkplaceSizeM  WorkplaceSize = "51-200"
	WorkplaceSizeL  WorkplaceSize = "201-500"
	WorkplaceSizeXL WorkplaceSize = "500+"
)

// OrgRole represents roles within an organization
type OrgRole string

const (
	OrgRoleAdmin    OrgRole = "admin"
	OrgRoleSubAdmin OrgRole = "sub_admin"
	OrgRoleMember   OrgRole = "member"
)

// WorkplaceRole represents roles within a workplace
type WorkplaceRole string

const (
	WorkplaceRoleAdmin  WorkplaceRole = "admin"
	WorkplaceRoleMember WorkplaceRole = "member"
)

// ProjectAccess represents project access levels
type ProjectAccess string

const (
	ProjectAccessPublic  ProjectAccess = "public"
	ProjectAccessPrivate ProjectAccess = "private"
)

// ProjectRole represents roles within a project
type ProjectRole string

const (
	ProjectRoleOwner  ProjectRole = "owner"
	ProjectRoleMember ProjectRole = "member"
)

// TaskPriority represents task priority levels
type TaskPriority string

const (
	TaskPriorityHighest TaskPriority = "highest"
	TaskPriorityHigh    TaskPriority = "high"
	TaskPriorityMedium  TaskPriority = "medium"
	TaskPriorityLow     TaskPriority = "low"
	TaskPriorityLowest  TaskPriority = "lowest"
)

// TaskLinkType represents the type of relation between tasks
type TaskLinkType string

const (
	TaskLinkTypeBlocks      TaskLinkType = "blocks"
	TaskLinkTypeIsBlockedBy TaskLinkType = "is_blocked_by"
	TaskLinkTypeDuplicates  TaskLinkType = "duplicates"
	TaskLinkTypeRelatesTo   TaskLinkType = "relates_to"
)

// ChannelType represents the type of communication channel
type ChannelType string

const (
	ChannelTypeGlobal  ChannelType = "global"
	ChannelTypeDM      ChannelType = "dm"
	ChannelTypeGroup   ChannelType = "group"
	ChannelTypeProject ChannelType = "project"
)
