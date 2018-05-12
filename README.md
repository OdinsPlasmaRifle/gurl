# gurl

Go Lang cURL-like tool for making http requests.

```bash
./gurl -U="http://test.example.com" -X="POST" -d='{"hello": "hello"}' -H="Content-Type: application/json" -H="Authorization: Token {token}" -interval=2 -repeat=2 -batch=2 -file="log.txt"
```

## Roadmap

- Add more cURL features.
- Improve logging format (ensure all http details are logged, add `-verbose` for extra details).
- Add tests.
- Find out why gurl sometimes gets a fatal error.
