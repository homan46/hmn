package dto

type userUserNameField struct {
	UserName string `db:"user_name" json:"userName"`
}

func (u *userUserNameField) GetUserName() string {
	return u.UserName
}
func (u *userUserNameField) SetUserName(userName string) {
	u.UserName = userName
}

type userPasswordField struct {
	Password string `db:"password" json:"password"`
}

func (u *userPasswordField) GetPassword() string {
	return u.Password
}
func (u *userPasswordField) SetPassword(password string) {
	u.Password = password
}
