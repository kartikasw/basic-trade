test:
	go test ./... -v -cover

server:
	go run main.go

sqlc:
	sqlc generate

mockRepo:
	mockgen -package mockRepo -destination internal/repository/mock/admin_repository.go basic-trade/internal/repository AdminRepository

mockService:
	mockgen -package mockService -destination internal/service/mock/auth_service.go basic-trade/internal/service AuthService && \
	mockgen -package mockService -destination internal/service/mock/product_service.go basic-trade/internal/service ProductService && \
	mockgen -package mockService -destination internal/service/mock/variant_service.go basic-trade/internal/service VariantService

.PHONY: test server sqlc mockRepo mockService
