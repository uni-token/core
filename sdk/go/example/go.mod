module uni-token-example

go 1.24.0

toolchain go1.24.6

require (
	github.com/openai/openai-go/v2 v2.1.0
	github.com/uni-token/core/sdk/go v0.0.0-00010101000000-000000000000
)

require (
	github.com/tidwall/gjson v1.14.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
)

replace github.com/uni-token/core/sdk/go => ..
