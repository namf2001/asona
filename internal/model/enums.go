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
	WorkplaceSize2_5     WorkplaceSize = "2-5"
	WorkplaceSize6_10    WorkplaceSize = "6-10"
	WorkplaceSize11_20   WorkplaceSize = "11-20"
	WorkplaceSize21_50   WorkplaceSize = "21-50"
	WorkplaceSize51_100  WorkplaceSize = "51-100"
	WorkplaceSize101_250 WorkplaceSize = "101-250"
	WorkplaceSize250Plus WorkplaceSize = "250-more"
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

// TaskStatus represents task status levels
type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusBacklog    TaskStatus = "backlog"
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

// OnboardingStatus represents the onboarding flow state for a user.
type OnboardingStatus string

const (
	OnboardingStatusPending    OnboardingStatus = "pending"
	OnboardingStatusInProgress OnboardingStatus = "in_progress"
	OnboardingStatusCompleted  OnboardingStatus = "completed"
)
