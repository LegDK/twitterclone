package twitterclone

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestCreateTweetInput_Sanitize(t *testing.T) {
	input := CreateTweetInput{
		Body: " hello  ",
	}

	want := CreateTweetInput{
		Body: "hello",
	}

	input.Sanitize()

	require.Equal(t, want.Body, input.Body)

}

func TestCreateTweetInput_Validate(t *testing.T) {

	testCases := []struct {
		name  string
		input CreateTweetInput
		err   error
	}{
		{
			name: "valid",
			input: CreateTweetInput{
				Body: "hello",
			},
			err: nil,
		},
		{
			name: "tweet not long enough",
			input: CreateTweetInput{
				Body: "h",
			},
			err: ErrValidation,
		},

		{
			name: "tweet to long enough",
			input: CreateTweetInput{
				Body: randStringRunes(300),
			},
			err: ErrValidation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}

}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
