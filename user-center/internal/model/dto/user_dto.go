package dto

type UserResp struct {
	ID     int64    `json:"id"`
	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Avatar string   `json:"avatar"`
	Status string   `json:"status"`
	Roles  []string `json:"roles"`
}

type LoginResp struct {
	Token string   `json:"token"`
	User  UserResp `json:"user"`
}

type UpdateUserStatusReq struct {
	Status string `json:"status"`
}

type AssignRoleReq struct {
	RoleName string `json:"roleName"`
}

type InitSuperAdminReq struct {
	Email string `json:"email"`
}

type UserFilter struct {
	Status   *string  `json:"status,omitempty"`
	Keyword  *string  `json:"keyword,omitempty"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
}

type CheckInitResp struct {
	NeedInit bool `json:"needInit"`
}
