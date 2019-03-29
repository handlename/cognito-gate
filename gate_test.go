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
				Allows: []configAllow{
					configAllow{
						Key:   "Email",
						Value: "alice@example.com",
						Rule:  "exact_match",
					},
					configAllow{
						Key:   "Email",
						Value: "bob@",
						Rule:  "forward_match",
					},
					configAllow{
						Key:   "Email",
						Value: "@example.org",
						Rule:  "backward_match",
					},
				},
			},
		},
	}

	makeEvent := func(email string) events.CognitoEventUserPoolsPreSignup {
		event := events.CognitoEventUserPoolsPreSignup{}
		event.UserPoolID = "foobar"
		event.Request = events.CognitoEventUserPoolsPreSignupRequest{
			UserAttributes: map[string]string{
				"Email": email,
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
		{email: "charlie@example.org", err: nil},
		{email: "devola@example.com", err: ErrNotAllowed},
	} {
		_, err := handler(makeEvent(c.email))

		if err != c.err {
			t.Errorf("unexpected error: %s", err)
		}
	}
}
