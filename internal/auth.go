package internal

import (
	"database/sql"
	"forumAA/database"
	"net/mail"
	"strings"
	"time"
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

	existed_user, _ := database.CheckEmail(db, user.Email)
	if existed_user.Id != 0 {
		return "A user with this email exists", false
	}

	name := strings.Trim(user.Nickname, " ")
	firstname := strings.Trim(user.Firstname, " ")
	lastname := strings.Trim(user.Lastname, " ")
	if name == "" {
		return "Username is empty.", false
	} else if firstname == "" {
		return "Firstname is empty", false
	} else if lastname == "" {
		return "Lastname is empty", false
	}

	_, err := database.GetUser(db, user.Nickname)
	if err == nil {
		return "A user with this nickname exists", false
	}

	if !CheckPassword(user.Password) {
		return "Invalid password: At least 1 upper case letter," +
			"1 lowercase letter, 1 digit, 1 digit and 8 characters long", false
	}

	return "", true
}

func CreateSession(db *sql.DB, nick string, token string, expiry time.Time) (err error) {
	sessionExistence, _ := database.IsSessionExist(db, nick)
	if !sessionExistence {
		err = database.CreateSession(db, nick, token, expiry)
	} else {
		err = database.UpdateSession(db, nick, token, expiry)
	}

	return
}
