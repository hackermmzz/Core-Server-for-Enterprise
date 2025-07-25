package Controller

type User struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RegistDate string `json:"registDate"`
	Avatar     string `json:"avatar"`
	Vip        int    `json:"vip"`
}

func (user *User) AddUser(email string, password string) {

}
