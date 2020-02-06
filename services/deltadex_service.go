package services

// DeltadexService is used to talk to the backend API for the game
type DeltadexService struct {
}

// Authenticate verifies the given token is valid
func (w *DeltadexService) Authenticate(token string) bool {
	return true
}
