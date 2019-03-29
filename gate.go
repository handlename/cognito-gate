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
	ID     string        `yaml:"id"`
	Allows []configAllow `yaml:"allows"`
}

type configAllow struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
	Rule  string `yaml:"rule"`
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
	for _, pool := range config.Pools {
		if event.UserPoolID != pool.ID {
			continue
		}

		for _, allow := range pool.Allows {
			value, ok := event.Request.UserAttributes[allow.Key]
			if !ok {
				continue
			}

			switch allow.Rule {
			case "exact_match":
				if value == allow.Value {
					return event, nil
				}
			case "forward_match":
				if strings.HasPrefix(value, allow.Value) {
					return event, nil
				}
			case "backward_match":
				if strings.HasSuffix(value, allow.Value) {
					return event, nil
				}
			default:
				log.Println("unknown rule:", allow.Rule)
			}
		}
	}

	return event, ErrNotAllowed
}
