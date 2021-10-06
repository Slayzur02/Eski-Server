package auth

type UserCreation struct {
	Email    string
	Username string
	Password string
}

type LogInCredentials struct {
	Username string
	Password string
}
