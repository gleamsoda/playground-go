package internal

import (
	"net/mail"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FuzzValidateString(f *testing.F) {
	// Seed corpus
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, value string) {
		err := ValidateString(value, 3, 100)

		if len(value) < 3 || len(value) > 100 {
			assert.Error(t, err, "expected error for value: %s", value)
		} else {
			assert.NoError(t, err, "unexpected error for value: %s", value)
		}
	})
}

func FuzzValidateUsername(f *testing.F) {
	testcases := []string{"test", "", "username", "---"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, value string) {
		err := ValidateUsername(value)

		if len(value) < 3 || len(value) > 100 {
			if assert.Error(t, err, "expected error for value: %s", value) {
				return
			}
		}

		if !isValidUsername(value) {
			assert.Error(t, err, "expected error for invalid username: %s", value)
		} else {
			assert.NoError(t, err, "unexpected error for valid username: %s, err: %s", value, err)
		}
	})
}

func FuzzValidateFullName(f *testing.F) {
	testcases := []string{"John Doe", "Alice", "1234567890", "Yo"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, value string) {
		err := ValidateFullName(value)

		if len(value) < 3 || len(value) > 100 {
			assert.Error(t, err, "expected error for value: %s", value)
		} else {
			assert.NoError(t, err, "unexpected error for value: %s, err: %s", value, err)
		}
	})
}

func FuzzValidatePassword(f *testing.F) {
	testcases := []string{"password", "12345", "strong_password", "pass"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, value string) {
		err := ValidatePassword(value)

		if len(value) < 6 || len(value) > 100 {
			assert.Error(t, err, "expected error for value: %s", value)
		} else {
			assert.NoError(t, err, "unexpected error for value: %s, err: %s", value, err)
		}
	})
}

func FuzzValidateEmail(f *testing.F) {
	testcases := []string{"test@example.com", "invalid_email", "user@domain.com", "yo"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, value string) {
		err := ValidateEmail(value)

		if len(value) < 3 || len(value) > 200 {
			if assert.Error(t, err, "expected error for value: %s", value) {
				return
			}
		}

		_, parseErr := mail.ParseAddress(value)
		if parseErr != nil {
			assert.Error(t, err, "expected error for invalid email address: %s", value)
		} else {
			assert.NoError(t, err, "unexpected error for valid email address: %s, err: %s", value, err)
		}
	})
}
