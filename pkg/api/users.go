package api

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
	"taskRumbler/pkg/models"
	u "taskRumbler/pkg/utils"
)

func (api *api) validate(user *models.User) (map[string]interface{}, bool) {
	if !strings.Contains(user.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	//Все поля доджны быть заполнены
	if user.Email == "" || user.Password == "" || user.Name == "" || user.Surname == "" || user.Patronymic == "" {
		return u.Message(false, "Not all user parameters are filled in"), false
	}

	//Электронная почта должна быть уникальной
	temp, _ := api.db.GetUserByEmail(user.Email)
	//if err != nil {
	//	return u.Message(false, "Connection error. Please retry"), false
	//}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user"), false
	}
	return u.Message(false, "Requirement passed"), true
}

func (api *api) createUser(w http.ResponseWriter, r *http.Request) {

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request "+err.Error()))
		return
	}
	if resp, ok := api.validate(&user); !ok {
		u.Respond(w, resp)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Respond(w, u.Message(false, "Failed to create users, not work hashed password"))
		return
	}

	user.Password = string(hashedPassword)
	user.Id, err = api.db.CreateUser(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Failed to create users,"+err.Error()))
		return
	}

	//Создаем новый токен JWT для вновь зарегистрированной учетной записи
	tk := &models.Token{UserId: user.Id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = "" //удаляем пароль
	response := u.Message(true, "Account has been created")
	response["user"] = user

	u.Respond(w, response)
	return
}

func (api *api) getUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := api.db.GetUsers()
	if err != nil {
		u.Respond(w, u.Message(false, "Connection error"))
		return
	}
	response := u.Message(true, "Account has been created")
	response["users"] = resp
	u.Respond(w, response)
	return
}

func (api *api) authenticate(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) // декодируем тело запроса в структуру и завершаемся неудачно, если возникает какая-либо ошибка
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := api.login(user.Email, user.Password)
	u.Respond(w, resp)
	return
}

func (api *api) login(email string, password string) map[string]interface{} {
	user, err := api.db.GetUserByEmail(email)

	if err != nil {
		return u.Message(false, "Connection error. Please retry")
	}

	if user.Email == "" {
		return u.Message(false, "Email address not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Пароль не подходит!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Работал! Вошёл в систему
	user.Password = ""

	//Создаем токен JWT
	tk := &models.Token{UserId: user.Id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString //Сохраняем токен в ответе

	resp := u.Message(true, "Logged In")
	resp["account"] = user
	return resp
}
