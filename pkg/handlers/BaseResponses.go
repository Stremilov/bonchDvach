package handlers

type InternalServerErrorResponse struct {
	Status  string `json:"error" example:"Непредвиденная ошибка"`
	Details string `json:"details" example:"Какая-то ошибка"`
}

type BadRequestResponse struct {
	Status  string `json:"error" example:"Ошибка при получении данных"`
	Details string `json:"details" example:"Какая-то ошибка"`
}

type SuccessCreatingResponse struct {
	Status string `json:"status" example:"success"`
}
