# trickster-ransomware-poc

A proof of concept ransomware example, written in golang language. Uses built in crypto module in golang.

### Compiling
```sh
go build
```

### Cross-compiling
With [xgo](https://github.com/karalabe/xgo)

```sh
# For linux amd64
xgo --targets=linux/amd64 .
```

### Running
```sh
./trickster-ransomware-poc

# or, if cross compiling

./trickster-ransomware-poc-linux-amd64
```