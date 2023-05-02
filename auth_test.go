package twitterclone

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegisterInput_Sanitize(t *testing.T) {
	actual := RegisterInput{
		Username:        "  bob  ",
		Email:           " bob@gmail.com ",
		Password:        "password",
		ConfirmPassword: "password",
	}

	actual.Sanitize()

	expected := RegisterInput{
		Username:        "bob",
		Email:           "bob@gmail.com",
		Password:        "password",
		ConfirmPassword: "password",
	}

	require.Equal(t, expected, actual)
}

func TestRegisterInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input RegisterInput
		err   error
	}{
		{
			name: "valid",
			input: RegisterInput{
				Username:        "bob",
				Email:           "bob@gmail.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: nil,
		},
		{
			name: "invalid email",
			input: RegisterInput{
				Username:        "bob",
				Email:           "bad email",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: ErrValidation,
		},
		{
			name: "invalid username",
			input: RegisterInput{
				Username:        "b",
				Email:           "bob@gmail.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			err: ErrValidation,
		},
		{
			name: "invalid password length",
			input: RegisterInput{
				Username:        "b",
				Email:           "bob@gmail.com",
				Password:        "p",
				ConfirmPassword: "p",
			},
			err: ErrValidation,
		},

		{
			name: "invalid password confirm",
			input: RegisterInput{
				Username:        "b",
				Email:           "bob@gmail.com",
				Password:        "password",
				ConfirmPassword: "p",
			},
			err: ErrValidation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if err != nil {
				require.ErrorIs(t, ErrValidation, tc.err)
			} else {
				require.NoError(t, tc.err)
			}
		})
	}
}

func TestLoginInput_Sanitize(t *testing.T) {
	actual := LoginInput{
		Email:    " bob@gmail.com ",
		Password: "password",
	}

	actual.Sanitize()

	expected := LoginInput{
		Email:    "bob@gmail.com",
		Password: "password",
	}

	require.Equal(t, expected, actual)
}

func TestLoginInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input LoginInput
		err   error
	}{
		{
			name: "valid",
			input: LoginInput{
				Email:    "bob@gmail.com",
				Password: "password",
			},
			err: nil,
		},
		{
			name: "invalid email",
			input: LoginInput{
				Email:    "bad email",
				Password: "password",
			},
			err: ErrValidation,
		},
		{
			name: "invalid password length",
			input: LoginInput{
				Email:    "bob@gmail.com",
				Password: "",
			},
			err: ErrValidation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if err != nil {
				require.ErrorIs(t, ErrValidation, tc.err)
			} else {
				require.NoError(t, tc.err)
			}
		})
	}
}
