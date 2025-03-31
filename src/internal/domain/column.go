package domain

type Column struct {
	Id        uint
	ProjectId uint
	Name      string
}

type CreateColumnInput struct {
	ProjectId uint   `json:"project_id"`
	Name      string `json:"name" binding:"min=3,max=50"`
}

type UpdateColumnInput struct {
	Name string `json:"name" binding:"min=3,max=50"`
}
