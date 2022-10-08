package sqlmgmt

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
