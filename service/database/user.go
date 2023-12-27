package database

import "database/sql"

func (db *appdbimpl) LoginUser(name string) (User, error, bool) {
	var nu User
	nu.Username = name

	res, err := db.c.Exec("INSERT INTO user(username) VALUES (?)", nu.Username)
	if err != nil {
		var u User
		err := db.c.QueryRow("SELECT UserId, Username FROM user WHERE Username = ?", nu.Username).Scan(u.UserId, u.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				return u, ErrorUserDoesNotExist, false
			}
			return u, err, false
		}
		return u, nil, false
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nu, err, false
	}
	nu.UserId = uint64(id)
	return nu, nil, true
}

func (db *appdbimpl) setUsername(name string)
