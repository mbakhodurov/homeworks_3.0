package input

// RegisterInput — вход use case'а регистрации пользователя.
type RegisterInput struct {
	Login    string
	Password string
}

// LoginInput — вход use case'а входа пользователя.
type LoginInput struct {
	Login    string
	Password string
}
