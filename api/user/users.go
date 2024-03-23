package user

type User struct {
	ID             uint64
	Name           string
	CredentialType string
	Credential     string
	Password       string
	Email          string
	Phone          string
}

type CredentialType string

const (
	Email CredentialType = "email"
	Phone CredentialType = "phone"
)

var CredentialTypes []interface{} = []interface{}{Email, Phone}
