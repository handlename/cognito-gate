package gate

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type configRoot struct {
	Pools []configPool `yaml:"pools"`
}

type configPool struct {
	ID     string   `yaml:"id"`
	Allows []string `yaml:"allows"`
}

var config configRoot

var ErrNotAllowed = fmt.Errorf("not allowed")

func Run(configPath string) error {
	err := parseConfig(configPath)
	if err != nil {
		return errors.Wrapf(err, "failed to parse config: %s", configPath)
	}

	lambda.Start(handler)
	return nil
}

func parseConfig(configPath string) error {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errors.Wrapf(err, "failed to read config: %s", configPath)
	}

	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal config as YAML: %s", configPath)
	}

	return nil
}

func handler(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	email, ok := event.Request.UserAttributes["email"]
	if !ok {
		return event, ErrNotAllowed
	}

	for _, pool := range config.Pools {
		if event.UserPoolID != pool.ID {
			continue
		}

		for _, allow := range pool.Allows {
			if strings.Contains(allow, "@") && email == allow {
				// exact email address
				return event, nil
			} else if strings.HasSuffix(email, allow) {
				// domain
				return event, nil
			}
		}
	}

	log.Printf("denied %+v\n", event.Request.UserAttributes)

	return event, ErrNotAllowed
}
