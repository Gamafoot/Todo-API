package domain

type Column struct {
	Id        uint   `json:"id"`
	ProjectId uint   `json:"project_id"`
	Name      string `json:"name"`
}

type CreateColumnInput struct {
	ProjectId uint   `json:"project_id" binding:"required"`
	Name      string `json:"name" binding:"required,min=3,max=50"`
}

type UpdateColumnInput struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}
