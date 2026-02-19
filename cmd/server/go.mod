module cmd/server

go 1.24.5

require internal/handler v1.0.0
replace internal/handler => ../../internal/handler
require internal/controller v1.0.0
replace internal/controller => ../../internal/controller
require internal/storage v1.0.0
replace internal/storage => ../../internal/storage
require internal/logger v1.0.0
replace internal/logger => ../../internal/logger
