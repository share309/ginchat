## Getting started

### Prerequisites

- **[Go](https://go.dev/)**: any one of the **three latest major** [releases](https://go.dev/doc/devel/release) (we test it with these).

### Getting ginchat

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```
import "github.com/share309/ginchat"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `gin_chat` package:

```sh
$ go get -u github.com/share309/ginchat
```