# svc-template

Small HTTP service with health endpoints and an echo API, built with Go. Includes unit tests, linting, and production-ready settings.

## Quickstart
```bash
make test
make run
# open http://localhost:8080/healthz
# open "http://localhost:8080/v1/echo?msg=hello"
