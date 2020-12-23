package enums

type UserStatus string

const (
	UserStatusPending UserStatus = "Pending"
	UserStatusCreated UserStatus = "Created"
	UserStatusDeleted UserStatus = "Deleted"
	UserStatusBanned  UserStatus = "Banned"
)
