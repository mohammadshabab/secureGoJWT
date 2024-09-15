package models

import (
	"os"
	"securegojwt/utils"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Token struct {
	UserId uint
	jwt.RegisteredClaims
}

type Account struct {
	gorm.Model `json:"-"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Token      string `json:"token" sql:"-"` // ignored by sql but include in json
}

func (account *Account) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return utils.Message(false, "Email address is required"), false
	}
	if len(account.Password) < 6 {
		return utils.Message(false, "Password is required"), false
	}
	//Verify email is unique
	temp := &Account{}
	//check for duplicate email
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.Message(false, "Connection error. Please retry"), false
	}

	if temp.Email != "" {
		return utils.Message(false, "Email address already in use by another user."), false
	}
	return utils.Message(false, "Requirement passed"), true
}

func (account *Account) Create() map[string]interface{} {
	if resp, ok := account.Validate(); !ok {
		return resp
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)
	if account.ID <= 0 {
		return utils.Message(false, "Failed to create account, Connection error")
	}
	//Create JWT token for newly created account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" //deleting the password

	response := utils.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func Login(email, password string) map[string]interface{} {
	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return utils.Message(false, "Email address not found")

		}
		return utils.Message(false, "Connection error. Please retry")
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return utils.Message(false, "Invalid login credentials. Please try again")
	}
	//logged in
	account.Password = ""

	//create jwt token
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString // store the token in response

	resp := utils.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

func GetUser(u uint) *Account {
	account := &Account{}
	GetDB().Table("accounts").Where("id = ?", u)
	if account.Email == "" {
		return nil
	}
	account.Password = ""
	return account
}
