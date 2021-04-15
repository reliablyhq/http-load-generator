# http-load-generator
A simple HTTP load generator that works with the reliably HTTP demo applications

## Running the load tester

```
go run main.go -host {SCHEME}://{HOSTNAME} -rps 20 -for 60 -workers 10 -min-latency 50 -max-latency 1000
```