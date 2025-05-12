gen_grpc:
	protoc --go_out=. --go-grpc_out=. internal/adapter/grpc/proto/user.proto
	protoc --go_out=. --go-grpc_out=. internal/adapter/grpc/proto/session.proto

gen_mocks:
	mockgen -source=internal/adapter/database/repository.go -destination=mocks/mock_repository.go -package=mocks
	mockgen -source=internal/adapter/database/user.go -destination=mocks/mock_user_repository.go -package=mocks
	mockgen -source=internal/adapter/database/session.go -destination=mocks/mock_session_repository.go -package=mocks