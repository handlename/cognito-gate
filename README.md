# cognito-gate

`cognito-gate` is a go package to simplify make permission trigger for Amazon Cognito.

## Usage

```go
package main

import (
    "log"

    "github.com/handlename/cognito-gate"
)

func main() {
    if err := gate.Run(os.Getenv("GATE_CONFIG_PATH")); err != nil {
        log.Println(err)
        os.Exit(1)
    }
}
```

## Configuration

in YAML format.

```yaml
pools:
  - id: <userPoolId>
    allows:
      - key:   "<target key of user attribute>"
        value: "<expected value>"
        rule:  "<matching rule> ::= 'exact_match'(default) | 'forward_match' | 'backward_match'"
```

example:

```yaml
pools:
  - id: ap-northeast-1_XXXX
    allows:
      - key:   "email"
        value: "@example.com"
        rule:  "backward_match"
```

## Lisence

MIT

## Author

@handlename (https://github.com/handlename)
