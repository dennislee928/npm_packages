package modelpkg

type GetResponse struct {
	Data interface{} `json:"data" binding:"required"`
	Paginator
}
