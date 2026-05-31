package repository

type UmiMemberRole string

const (
	MemberRoleOwner  UmiMemberRole = "owner"
	MemberRoleAdmin  UmiMemberRole = "admin"
	MemberRoleMember UmiMemberRole = "member"
)

type UmiMemberPermission int

const (
	UmiMemberPermissionRead UmiMemberPermission = 1 << iota
	UmiMemberPermissionWrite
	UmiMemberPermissionAdmin
)

var RolePermissions = map[UmiMemberRole]UmiMemberPermission{
	MemberRoleOwner:  UmiMemberPermissionRead | UmiMemberPermissionWrite | UmiMemberPermissionAdmin,
	MemberRoleAdmin:  UmiMemberPermissionRead | UmiMemberPermissionWrite,
	MemberRoleMember: UmiMemberPermissionRead,
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
