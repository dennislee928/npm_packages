package managers

type CreateRequest struct {
	Account         string `json:"account" binding:"required,email"`   // Email
	Name            string `json:"name" binding:"required"`            // 管理者名稱
	ManagersRolesID int64  `json:"managersRolesID" binding:"required"` // 管理者角色
}

type DeleteRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0"`
}
