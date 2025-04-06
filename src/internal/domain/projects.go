package domain

type Project struct {
	Id     uint
	UserId uint
	Name   string
}

type CreateProjectInput struct {
	Name string `json:"username" binding:"required,min=3,max=50"`
}

type UpdateProjectInput struct {
	Name string `json:"username" binding:"required,min=3,max=50"`
}
