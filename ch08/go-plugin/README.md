

## Build shared object (`.so`) files.

```bash
$ go build -buildmode=plugin -o duck/duck.so duck/duck.go
$ go build -buildmode=plugin -o frog/frog.so frog/frog.go
```

## Check file types

Just for fun, check the file types.

```bash
$ file duck/duck.so
$ file frog/frog.so
```

Depending on your OS, you'll see something like:

```bash
$ file duck/duck.so
duck/duck.so: Mach-O 64-bit dynamically linked shared library x86_64
```
