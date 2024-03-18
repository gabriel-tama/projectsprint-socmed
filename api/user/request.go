package user

import (
	"fmt"
	"regexp"
)

type CreateUserPayload struct {
	Name            string `json:"name" binding:"required,min=5,max=50"`
	CredentialType  string `json:"credentialType" binding:"required,oneof=email phone"`
	CredentialValue string `json:"credentialValue" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

func (p CreateUserPayload) Validate() error {
	// Validate CredentialType
	validCredentialTypes := map[string]bool{"phone": true, "email": true}
	if !validCredentialTypes[p.CredentialType] {
		return fmt.Errorf("invalid credential type: %s", p.CredentialType)
	}

	// Validate CredentialValue
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	var phoneRegex = regexp.MustCompile(`^\+\d{7,13}$`)
	if p.CredentialType == "email" {
		if !emailRegex.MatchString(p.CredentialValue) {
			return fmt.Errorf("invalid email format")
		}
	} else {
		if !phoneRegex.MatchString(p.CredentialValue) {
			return fmt.Errorf("invalid phone format")
		}
	}

	// Validate Name
	if len(p.Name) < 5 || len(p.Name) > 50 {
		return fmt.Errorf("name length should be between 5 and 50 characters")
	}

	// Validate Password
	if len(p.Password) < 5 || len(p.Password) > 15 {
		return fmt.Errorf("password length should be between 5 and 15 characters")
	}

	return nil
}

func (p CreateUserPayload) CredentialByEmail() bool {
	if (p.CredentialType) == "email" {
		return true
	}
	return false
}
