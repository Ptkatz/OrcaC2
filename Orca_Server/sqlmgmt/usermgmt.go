package sqlmgmt

import (
	"gorm.io/gorm"
)

func InitUser() {
	Db = GetDb()
	Db.Exec("UPDATE users_lists SET online = ?", "Off")
	Db.Exec("UPDATE users_lists SET login_time = ?", "")
}
func Login(username, password string) UsersList {
	Db = GetDb()
	user := UsersList{}
	Db.Where("username = ? and password = ?", username, password).Take(&user)
	return user
}

func GetUsernames() []string {
	Db = GetDb()
	users := []UsersList{}
	Db.Find(&users)
	var usernames []string
	for _, username := range users {
		usernames = append(usernames, username.Username)
	}
	return usernames
}

func AddUser(username, password string) {
	Db = GetDb()
	newUser := UsersList{
		Model:     gorm.Model{},
		Username:  username,
		Password:  password,
		LoginIp:   "",
		LoginTime: "",
		Online:    "Off",
	}
	Db.Create(&newUser)
}

func DelUser(username string) {
	Db = GetDb()
	Db.Exec("DELETE FROM users_lists where username = ?", username)
}

func ModUserPwd(username, password string) {
	Db = GetDb()
	Db.Exec("UPDATE users_lists SET password = ? where username = ?", password, username)
}
