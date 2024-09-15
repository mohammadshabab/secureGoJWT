package models

import (
	"fmt"
	"securegojwt/utils"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	UserId uint   `json:"user_id"`
}

// validate the required parameters sent through the http request body
// returns message and true if the requirement is met
func (contact *Contact) Validate() (map[string]interface{}, bool) {
	if contact.Name == "" {
		return utils.Message(false, "Contact name should be on the payload"), false
	}
	if contact.Phone == "" {
		return utils.Message(false, "Phone number should be on the payload"), false
	}
	if contact.UserId <= 0 {
		return utils.Message(false, "User is not recognised"), false
	}

	return utils.Message(true, "success"), true
}

func (c *Contact) Create() map[string]interface{} {
	if resp, ok := c.Validate(); !ok {
		return resp
	}
	GetDB().Create(c)

	resp := utils.Message(true, "success")
	resp["contact"] = c
	return resp
}

func GetContact(id uint) *Contact {
	contact := &Contact{}
	err := GetDB().Table("contacts").Where("id = ?", id).First(contact).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return contact
}

func GetContacts(user uint) []*Contact {
	contacts := make([]*Contact, 0)
	err := GetDB().Table("contacts").Where("user_id = ?", user).Find(contacts).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return contacts
}
