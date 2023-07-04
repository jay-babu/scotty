watch:
	air

run:
	go run main.go

generate:
	oapi-codegen -generate types -o openapi/openapi_types.gen.go -package openapi openapi/model.yml
	oapi-codegen -generate gin -o openapi/openapi_server.gen.go -package openapi openapi/model.yml
	oapi-codegen -generate spec -o openapi/openapi_server.spec.go -package openapi openapi/model.yml
