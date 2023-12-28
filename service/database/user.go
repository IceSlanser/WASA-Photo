package database

import "database/sql"

func (db *appdbimpl) LoginUser(name string) (User, bool, error) {
	var nu User
	nu.Username = name

	res, err := db.c.Exec("INSERT INTO users(username) VALUES (?)", nu.Username)
	if err != nil {
		var u User
		err := db.c.QueryRow("SELECT UserId, Username FROM users WHERE Username = ?", nu.Username).Scan(u.UserId, u.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				return u, false, ErrorUserDoesNotExist
			}
		}
		return u, false, nil
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nu, false, err
	}
	nu.UserId = uint64(id)
	return nu, true, nil
}

func (db *appdbimpl) SetUsername(uId uint64, nname string) error {
	_, err := db.c.Exec("UPDATE users SET Username = ? WHERE UserId = ?", nname, uId)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) IsAvailable(nname string) (bool, error) {
	var username string
	err := db.c.QueryRow("SELECT Username FROM users WHERE Username = ?", nname).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, err
		}
	}
	return false, err
}
