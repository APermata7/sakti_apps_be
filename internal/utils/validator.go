package utils

func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password minimal 8 karakter"
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	specialChars := "!@#$%^&*()_+-=[]{}|;:',.<>?/`~"

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		} else {
			for _, special := range specialChars {
				if char == special {
					hasSpecial = true
					break
				}
			}
		}
	}

	if !hasUpper {
		return false, "Password harus mengandung minimal 1 huruf besar"
	}
	if !hasLower {
		return false, "Password harus mengandung minimal 1 huruf kecil"
	}
	if !hasDigit {
		return false, "Password harus mengandung minimal 1 angka"
	}
	if !hasSpecial {
		return false, "Password harus mengandung minimal 1 karakter khusus (!@#$%^&*()_+-=[]{}|;:',.<>?/`~)"
	}

	return true, ""
}