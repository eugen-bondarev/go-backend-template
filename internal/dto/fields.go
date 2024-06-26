package dto

type WithToken struct {
	Token string `json:"token"`
}

type WithRefreshToken struct {
	RefreshToken string `json:"refreshToken"`
}

type WithEmail struct {
	Email string `json:"email" validate:"email"`
}

type WithPassword struct {
	Password string `json:"password" validate:"gte=6,lte=24"`
}
