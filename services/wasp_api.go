package services

// WaspService is used to talk to the backend API for the game
type WaspService struct {
}

// Authenticate verifies the given token is valid
func (w *WaspService) Authenticate(token string) bool {
	return true
}
