module web

go 1.25.5

replace domain => ../domain

replace datasource => ../datasource

require (
	datasource v0.0.0
	domain v0.0.0
	github.com/google/uuid v1.6.0
)
