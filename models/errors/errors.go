package errors

import "errors"

var (
	ConnectionError  = errors.New("Грешка при връзката с БД")
	BannedUser       = errors.New("Вие сте баннат!")
	WrongPassword    = errors.New("Грешна парола")
	NoUser           = errors.New("Няма такъв потребител")
	ResultError      = errors.New("Error with result")
	UpdateError      = errors.New("Error with updating db")
	InsertError      = errors.New("Error with inserting into db")
	DuplicationError = errors.New("Duplication error")
	ValidError       = errors.New("Not valid data")
	NotFoundError    = errors.New("Няма резултати")
	DataError        = errors.New("Error with data")
)
