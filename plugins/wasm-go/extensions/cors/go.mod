module cors

go 1.19

replace github.com/alibaba/higress/plugins/wasm-go => ../..

require (
	github.com/alibaba/higress/plugins/wasm-go v0.0.0-20230519024024-625c06e58f91
	github.com/higress-group/proxy-wasm-go-sdk v0.0.0-20240318034951-d5306e367c43
	github.com/stretchr/testify v1.8.4
	github.com/tidwall/gjson v1.14.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/higress-group/nottinygc v0.0.0-20231101025119-e93c4c2f8520 // indirect
	github.com/magefile/mage v1.14.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/resp v0.1.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
