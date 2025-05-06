package domain

type Column struct {
	Id        uint   `json:"id"`
	ProjectId uint   `json:"project_id"`
	Name      string `json:"name"`
	Position  int    `json:"position"`
}

type CreateColumnInput struct {
	ProjectId uint   `json:"project_id" validate:"required"`
	Name      string `json:"name" validate:"required,gte=3,lte=50"`
}

type UpdateColumnInput struct {
	Name     string `json:"name" validate:"omitempty,gte=3,lte=50"`
	Position int    `json:"position" validate:"omitempty,gte=1"`
}
