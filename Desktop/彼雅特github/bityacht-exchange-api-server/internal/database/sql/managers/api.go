package managers

// Manager for GET managers API
type Manager struct {
	ID              int64  `json:"id" binding:"required"`
	Account         string `json:"account" binding:"required"`
	Name            string `json:"name" binding:"required"`
	ManagersRolesID int64  `json:"managersRolesID" binding:"required"`
	Status          Status `json:"status"`
}

// UpdateRequest for PATCH managers API
type UpdateRequest struct {
	ID              int64  `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`
	Account         string `json:"account" binding:"omitempty,email|eq=admin"`
	Password        string `json:"password" binding:"omitempty"`
	Name            string `json:"name" binding:"omitempty"`
	Status          Status `json:"status" binding:"omitempty"`
	ManagersRolesID int64  `json:"managersRolesID" binding:"omitempty"`
}
