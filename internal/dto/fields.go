package dto

type WithEmail struct {
	Email string `json:"email" validate:"email"`
}

type WithPassword struct {
	Password string `json:"password" validate:"gte=6,lte=24"`
}
