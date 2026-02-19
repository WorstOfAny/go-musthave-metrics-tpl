module internal/handler

go 1.24.5

require (
	github.com/stretchr/testify v1.11.1
	internal/logger v1.0.0
	internal/storage v1.0.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	internal/logger => ../logger
	internal/storage => ../storage
)
