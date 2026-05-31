package repository

type UmiMemberRole string

const (
	UmiMemberRoleOwner  UmiMemberRole = "owner"
	UmiMemberRoleAdmin  UmiMemberRole = "admin"
	UmiMemberRoleMember UmiMemberRole = "member"
)

type UmiMemberPermission int

const (
	UmiMemberPermissionRead UmiMemberPermission = 1 << iota
	UmiMemberPermissionWrite
	UmiMemberPermissionAdmin
)

var RolePermissions = map[UmiMemberRole]UmiMemberPermission{
	UmiMemberRoleOwner:  UmiMemberPermissionRead | UmiMemberPermissionWrite | UmiMemberPermissionAdmin,
	UmiMemberRoleAdmin:  UmiMemberPermissionRead | UmiMemberPermissionWrite,
	UmiMemberRoleMember: UmiMemberPermissionRead,
}

func (r UmiMemberRole) Permissions() UmiMemberPermission {
	return RolePermissions[r]
}

func (r UmiMemberRole) HasPermission(permission UmiMemberPermission) bool {
	return r.Permissions()&permission != 0
}

func HasMinRole(role UmiMemberRole, minRole UmiMemberRole) bool {
	return role.Permissions()&minRole.Permissions() == minRole.Permissions()
}
