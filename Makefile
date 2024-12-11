build-mocks:
	@go get github.com/golang/mock/gomock
	@go get github.com/golang/mock/mockgen

	@go run -mod=mod github.com/golang/mock/mockgen -package mocks -destination=repositories/log/mock/mock.go -source=repositories/log/mongodb.go -build_flags=-mod=mod -imports=reflect=reflect,models=service-secret-santa/models,gomock=github.com/golang/mock/gomock
	@go run -mod=mod github.com/golang/mock/mockgen -package mocks -destination=services/log/mock/mock.go -source=services/log/service.go  -build_flags=-mod=mod