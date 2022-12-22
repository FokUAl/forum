package internal

import (
	"database/sql"
	"forumAA/database"
	"net/mail"
	"unicode"
)

func ValidMailAddress(address string) (string, bool) {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return "", false
	}
	return addr.Address, true
}

func CheckPassword(pass string) bool {
	var (
		upp, low, num, sym bool
		tot                uint8
	)

	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			return false
		}
	}

	if !upp || !low || !num || !sym || tot < 8 {
		return false
	}

	return true
}

func CheckInput(db *sql.DB, user database.User) (string, bool) {
	if _, ok := ValidMailAddress(user.Email); !ok {
		return "Email is invalid", false
	}

	if !CheckPassword(user.Password) {
		return "Invalid password: At least 1 upper case letter," +
			"1 lowercase letter, 1 digit, 1 digit and 8 characters long", false
	}

	_, err := database.GetUser(db, user.Nickname)
	if err == nil {
		return "A user with this nickname exists", false
	}

	err = database.CheckEmail(db, user.Email)
	if err == nil {
		return "A user with this nickname exists", false
	}

	return "", true
}