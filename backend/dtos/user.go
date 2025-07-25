package dtos

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleID   uint   `json:"role_id"`
	Role     string `json:"role"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
}

type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	RoleID   uint   `json:"role_id,omitempty"`
}
