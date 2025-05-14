gen_grpc:
	protoc --go_out=. --go-grpc_out=. internal/adapter/grpc/proto/user.proto
	protoc --go_out=. --go-grpc_out=. internal/adapter/grpc/proto/session.proto

grpc_docs:
	mkdir -p grpc-docs
	protoc \
		--proto_path=. \
		--doc_out=grpc-docs \
		--doc_opt=markdown,user-api.md \
		internal/adapter/grpc/proto/user.proto
	protoc \
		--proto_path=. \
		--doc_out=grpc-docs \
		--doc_opt=markdown,session-api.md \
		internal/adapter/grpc/proto/session.proto

gen_mocks:
	mockgen -source=internal/adapter/database/repository.go -destination=mocks/mock_repository.go -package=mocks
	mockgen -source=internal/adapter/database/user.go -destination=mocks/mock_user_repository.go -package=mocks
	mockgen -source=internal/adapter/database/session.go -destination=mocks/mock_session_repository.go -package=mocks
	mockgen -source=internal/core/services/user_services/user.go -destination=mocks/mock_user_services.go -package=mocks
	mockgen -source=internal/core/services/session_services/session.go -destination=mocks/mock_session_services.go -package=mocks
	mockgen -source=internal/core/services/auth_services/auth.go -destination=mocks/mock_auth_services.go -package=mocks

UNIT_TEST_DIRS = \
	./internal/core/services/session_services \
	./internal/core/services/user_services \
	./internal/core/services/auth_services \
	./internal/adapter/http/handlers/users

unit_test:
	go test -v -coverprofile=coverage.out $(UNIT_TEST_DIRS)
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out