# gurl

Go Lang cURL-like tool for making http requests.

````bash
./gurl -U="http://test.example.com" -X="GET" -d="{'hello':'hello'}" -H="Test: 123" -interval=2 -repeat=2 -batch=2
```

## Roadmap

- Add more cURL features.
- Improve logging format (ensure all http details are logged, add `-verbose` for extra details).
- Add tests.
