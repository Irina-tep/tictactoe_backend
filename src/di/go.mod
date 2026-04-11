module main

go 1.25.5

replace (
	datasource => ../datasource
	domain => ../domain
	web => ../web
)

require (
	datasource v0.0.0
	go.uber.org/fx v1.24.0
	web v0.0.0
)

require (
	domain v0.0.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
)
