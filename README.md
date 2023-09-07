Run from terminal:
```bash
$ go run cmd/parserd/*.go -api.port 3000 -rpc.endpoint https://cloudflare-eth.com -workers 2 -interval 10
```

Run with docker:
```bash
$ docker build -f Dockerfile . -t parserd:latest

$ docker run --rm -p 3000:3000 parserd:latest
```


Test curl:
```bash
# Subscribe an address
$ curl -X POST 'http://localhost:3000/subscribe' -d '{"address":"0xf584F8728B874a6a5c7A8d4d387C9aae9172D621"}'

# Get parsed transactions:
$ curl 'http://localhost:3000/transactions?address=0xf584F8728B874a6a5c7A8d4d387C9aae9172D621'
```