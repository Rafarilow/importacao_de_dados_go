package utils

import (
	"errors"
	"regexp"
)

// Valida se todos os campos obrigatorios estão preenchidos
func ValidateFields(contact map[string]string) error {
	requiredFields := []string{"Nome", "Email", "Telefone"}
	for _, field := range requiredFields {
		if contact[field] == "" {
			return errors.New("campo " + field + "não pode ser vazio")
		}
	}
	return nil
}

// Valida o formato do e-mail
func ValidateEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
