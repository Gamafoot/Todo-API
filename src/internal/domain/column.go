package domain

type Column struct {
	Id        uint
	ProjectId uint
	Name      string
}

type CreateColumnInput struct {
	ProjectId uint   `json:"project_id" binding:"required"`
	Name      string `json:"name" binding:"required,min=3,max=50"`
}

type UpdateColumnInput struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}
