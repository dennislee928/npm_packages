package receipts

type ViewRequest struct {
	ID string `uri:"ID" binding:"required,gt=0"`
}

type IssueRequest struct {
	IDs []string `json:"ids" binding:"required,unique,min=1,dive,required"`
}
