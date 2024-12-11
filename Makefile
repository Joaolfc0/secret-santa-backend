build-mocks:
	@go get github.com/golang/mock/gomock
	@go get github.com/golang/mock/mockgen

	@go run -mod=mod github.com/golang/mock/mockgen -package mocks -destination=repositories/group/mock/mock.go -source=repositories/group/mongodb.go -build_flags=-mod=mod 
	@go run -mod=mod github.com/golang/mock/mockgen -package mocks -destination=services/group/mock/mock.go -source=services/group/service.go  -build_flags=-mod=mod