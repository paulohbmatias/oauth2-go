package user

import (
	"database/sql"
	"github.com/paulohbmatias/oauth2-go/models"
	"log"
)

type UserRepository struct {}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (u UserRepository) SignUp(db *sql.DB, user models.User) (models.User, error) {
	stmt := "insert into users (email, password) values($1, $2) RETURNING id;"
	err := db.QueryRow(stmt, user.Email, user.Password).Scan(&user.ID)

	logFatal(err)

	user.Password = ""
	return user, nil
}

func (u UserRepository) Login(db *sql.DB, user models.User) (models.User, error) {
	row := db.QueryRow("select * from users where email=$1", user.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		return user, err
	}

	return user, nil
}