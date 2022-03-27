# Introduction

Password-Service is a tool to verify whether if a password is secure.

Given a set of rules, the service verifies if the password meets the configured validations and if so, performs a GET request to the [HaveIBeenPwned](https://haveibeenpwned.com/) service, which will tell if the password is present in any password dictionary.

# Prerequisites

- [Docker & Docker Compose](https://docs.docker.com/compose/install/)

# Service Configuration

The following is an example of the `pwned-config.yml` file, which is where the password validations are configured.

```
password:
  length:
    enabled: true
    min: 10
    max: 15
  symbols:
    enabled: true
    allowSymbols: true        # If false, any other than [a-zA-Z0-9] will be invalid
    allowedSymbols: "~!@#$^&*()_-+={}[]|:;<>,.?/%"
    min: 1
  numbers:
    enabled: true
    allowNumbers: true
    onlyNumbers: false
    min: 1
  case:
    enabled: true
    onlyUpper: false
    onlyLower: false
    minLower: 5
    minUpper: 2
pwned:
    enabled: true
    timeoutSeconds: 2
    url: "https://api.pwnedpasswords.com/range/"
```

These rules are meant to run simultaneously by using `goroutines`, and only if all of them are successfull a request is made to the `Pwned` endpoint.
The `pwned` endpoint allows to be enabled/disabled and also can be configured to timeout if the requests lasts longer than the specified `timeoutSeconds`

The application uses by default the config file located under: `/config/pwned-config.yml`. This is file is then mounted as a volume in `go-pwned` container as it is read by the application at startup time. If another config is used, remember to modify this config path in the `go-pwned` container.

## Running password-service

The service itself is ready to be containerised with `Docker-Compose` and comes configured with `Prometheus` and `Grafana`.

- The password-service (`go-pwned` container) runs in port `2112` by default
- `Prometheus` runs in port `9090`
- `Grafana` is reachable in port `3000` and already exposes a set of dashboards for request latency, successful, bad requests and server errors.

The password is sent in a `POST` request to the `/validate` endpoint using json format and encoded in Base64.

Example; If user wants to verify if whether if password `Passw0rd!` is secure, then it will encode the password in `Base64` and wrap it in the following json body:

```
{
    "password": "UGFzc3cwcmQh"
}
```

If the password is considered secure, the service replies with `200 - Ok`, otherwise will respond with `400 - Bad Request`.

The service exposes the following endpoints:

- `/validate`: Accepts `POST` requests with the json body already specified above.
- `/healthz`: Accepts `GET` requests and will return `200 (Ok)` if service is reachable
- `/metrics`: Accepts `GET` requests and exposes prometheus histogram `http_request_duration` and other `golang` metrics.
