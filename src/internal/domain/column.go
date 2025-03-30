package domain

type Column struct {
	Id        uint
	ProjectId uint
	Name      string
}

type CreateColumnInput struct {
	Name string `json:"username" binding:"required,min=3,max=50"`
}

type UpdateColumnInput struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}
