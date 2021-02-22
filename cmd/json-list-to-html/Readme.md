### Build + install
```
go build -v ./cmd/json-list-to-html && mv json-list-to-html /usr/local/bin
```

### Usage
```
json-list-to-html --dataUrl https://foo.bar/baz.json \
    --template my-custom-html.gohtml \
    --out target-report.html
```
