# cognito-gate

`cognito-gate` is a application to simplify make permission trigger for Amazon Cognito, runs on AWS Lambda.

## Usage

Download binary from release, deploy as Lambda function as your favorite way.

`cognito-gate` is for runtime `provided.al2`.
bootstrap file may be like this:

```sh
#!/bin/sh

cd $LAMBDA_TASK_ROOT
./cognito-gate
```

## Configuration

in YAML format.

```yaml
pools:
  - id: <userPoolId>
    allows:
      - <pattern ::= 'domain' | 'email'>
      - ...
```

example:

```yaml
pools:
  - id: ap-northeast-1_XXXX
    allows:
      - "example.com"
      - "alice@example.net"
```

## Lisence

MIT

## Author

@handlename (https://github.com/handlename)
