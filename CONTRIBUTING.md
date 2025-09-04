# ZobTube contributing guide

Thank you for investing your time in contributing to our project!

Any contribution is welcomed :sparkles: !

## :hammer_and_wrench: Development environement setup

This project is written in Go and will require the `go` binary to work.

If you want to implement a new feature or just try to run the raw source code, you can use the following commnand.

__Linux__
```
tools/dev-linux.sh
```

__Windows__
```
tools/dev-win.sh
```

This will start the build and live reload using `air` ([a golang tool to help build binaries in developement](https://github.com/air-verse/air)) and it will be accessible on [http://127.0.0.1:8069](http://127.0.0.1:8069).
