package auth

type UserCreation struct {
	Email    string
	Username string
	Password string
}

type LogInCredentials struct {
	Email    string
	Password string
}
