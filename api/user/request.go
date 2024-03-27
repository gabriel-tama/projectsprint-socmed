package user

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
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
	var phoneRegex = regexp.MustCompile(`^\+\d{6,12}$`)
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

type LoginUserPayload struct {
	CredentialType  string `json:"credentialType" binding:"required,oneof=email phone"`
	CredentialValue string `json:"credentialValue" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

func (p LoginUserPayload) Validate() error {
	// Validate CredentialType
	validCredentialTypes := map[string]bool{"phone": true, "email": true}
	if !validCredentialTypes[p.CredentialType] {
		return fmt.Errorf("invalid credential type: %s", p.CredentialType)
	}

	// Validate CredentialValue
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	var phoneRegex = regexp.MustCompile(`^\+\d{6,12}$`)
	if p.CredentialType == "email" {
		if !emailRegex.MatchString(p.CredentialValue) {
			return fmt.Errorf("invalid email format")
		}
	} else {
		if !phoneRegex.MatchString(p.CredentialValue) {
			return fmt.Errorf("invalid phone format")
		}
	}

	// Validate Password
	if len(p.Password) < 5 || len(p.Password) > 15 {
		return fmt.Errorf("password length should be between 5 and 15 characters")
	}

	return nil
}

func (p LoginUserPayload) CredentialByEmail() bool {
	if (p.CredentialType) == "email" {
		return true
	}
	return false
}

type LinkEmailPayload struct {
	Email string `json:"email" binding:"required"`
}

func (p LinkEmailPayload) Validate() error {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(p.Email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

type LinkPhonePayload struct {
	Phone string `json:"phone" binding:"required"`
}

func (p LinkPhonePayload) Validate() error {
	var phoneRegex = regexp.MustCompile(`^\+\d{6,12}$`)
	if !phoneRegex.MatchString(p.Phone) {
		fmt.Println("invalid phone")
		return fmt.Errorf("invalid phone format")
	}
	return nil
}

type UpdateAccountPayload struct {
	ImageURL string `json:"imageUrl" binding:"required"`
	Name     string `json:"name" binding:"required,min=5,max=50"`
}

func (p UpdateAccountPayload) Validate() error {
	url, err := url.ParseRequestURI(p.ImageURL)
	if err != nil {
		return err
	}

	address := net.ParseIP(url.Host)

	if address == nil {

		if !strings.Contains(url.Host, ".") {
			return ErrValidationFailed
		}
	}

	return nil
}
