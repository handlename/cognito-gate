package gate

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestParseConfig(t *testing.T) {
	err := parseConfig("config.example.yaml")
	if err != nil {
		t.Errorf("failed to parse config: %s", err)
	}
}

func TestHandlerExactMatch(t *testing.T) {
	config = configRoot{
		Pools: []configPool{
			configPool{
				ID: "foobar",
				Allows: []string{
					"alice@example.com",
					"example.net",
				},
			},
		},
	}

	makeEvent := func(email string) events.CognitoEventUserPoolsPreSignup {
		event := events.CognitoEventUserPoolsPreSignup{}
		event.UserPoolID = "foobar"
		event.Request = events.CognitoEventUserPoolsPreSignupRequest{
			UserAttributes: map[string]string{
				"email": email,
			},
		}

		return event
	}

	for _, c := range []struct {
		email string
		err   error
	}{
		{email: "alice@example.com", err: nil},
		{email: "bob@example.net", err: nil},
		{email: "charlie@example.com", err: ErrNotAllowed},
		{email: "devola@example.org", err: ErrNotAllowed},
		{email: "eve@sub.example.net", err: ErrNotAllowed},
	} {
		_, err := handler(makeEvent(c.email))

		if err != c.err {
			t.Errorf("unexpected error: %s", err)
		}
	}
}
