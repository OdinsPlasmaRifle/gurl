# gurl

Go Lang cURL-like tool for making http requests.

````bash
./gurl -U="http://test.example.com" -X="GET" -d="{'hello':'hello'}" -H="Test: 123" -interval=2 -repeat=2
```

## Roadmap

- Add more cURL features.
- Improve repetitive url requests.
- Improve logging format (ensure all http details are logged, add `-verbose` for extra details).
- Allow log destination/filename to be specified.
- Add tests.
- Add concurrency model, multiple gurl tasks at once.
