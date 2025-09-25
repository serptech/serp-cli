module github.com/serptech/serp-cli

go 1.24.0

require (
	github.com/rs/zerolog v1.34.0
	github.com/serptech/serp-go v0.3.0
	github.com/spf13/cobra v1.10.1
	github.com/tidwall/pretty v1.2.1
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	golang.org/x/sys v0.36.0 // indirect
)

replace github.com/serptech/serp-go => ../serp-go
