package user

type User struct {
	ID                    uint64
	Name                  string
	CredentialType        string
	EmailCredential       string
	PhoneNumberCredential string
	Password              string
}

type CredentialType string

const (
	Email CredentialType = "email"
	Phone CredentialType = "phone"
)

var CredentialTypes []interface{} = []interface{}{Email, Phone}
