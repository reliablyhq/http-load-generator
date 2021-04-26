# http-load-generator
A simple HTTP load generator that works with the reliably HTTP demo applications

## Running the load tester
```
go run main.go -host {SCHEME}://{HOSTNAME}
```

## Controlling the number of workers
```
go run main.go -host {SCHEME}://{HOSTNAME} -workers 10
```

## Controlling the number of requests per second
*this is best effort - its not guaranteed*
```
go run main.go -host {SCHEME}://{HOSTNAME} -rps 20 -for 60
```

## Returning with a certain amount of latency
```
go run main.go -host {SCHEME}://{HOSTNAME} -min-latency 50 -max-latency 1000
```

## Returning a specific status code
```
go run main.go -host {SCHEME}://{HOSTNAME} -status-code 200
```

## Returning random status codes
```
go run main.go -host {SCHEME}://{HOSTNAME} -status-codes 200,400,401,404,500
```

## bringing it together
```
go run main.go -host {SCHEME}://{HOSTNAME} -rps 20 -for 60 -workers 10 -min-latency 50 -max-latency 1000 -status-code 404
```

```
go run main.go -host {SCHEME}://{HOSTNAME} -rps 20 -for 60 -workers 10 -min-latency 50 -max-latency 1000 -status-codes 200,400,500
```