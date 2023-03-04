# Simple E-commerce

## Project structure
base on go-fiber with microservices and clean architecture
- `docs` postman collection
- `cmd` contains microservices main.go for service's entrypoint
- `configs` contains service's config files
- `internal` contains service's codes
  - `handlers` contains routes / represent layer
  - `repository` contains datasources layer
  - `usecase` contains business logic layer
- `local-data` just for local data testing
- `pkg` contains shareable functions
  - `config` config loader
  - `database` database connections builder
  - `domain` service's interfaces
  - `logger` log builder / formatter
  - `middleware` shareable middleware
  - `modles` shareable service's data models
  - `present` project's request / response models
  - `utils` shareable helper functions

## Services
### auth
- run: `go run ./cmd/auth`
- listen: `http://localhost:8001`
- base path: `/auth`
  - POST `/register`
  - POST `/login`
### user
- run: `go run ./cmd/user`
- listen: `http://localhost:8002`
- base path: `/api/v1/users`
  - GET `/:email?`
  - PATCH `/:email?`
  - DELETE `/:email?`
### product
- run: `go run ./cmd/product`
- listen: `http://localhost:8003`
- base path: `/api/v1/products`
  - POST `/`
  - GET `/:product_id?`
  - PATCH `/:product_id?`
  - DELETE `/:product_id?`
### order
- run: `go run ./cmd/order`
- listen: `http://localhost:8004`
- base path: `/api/v1/orders`
  - POST `/`
  - GET `/:order_id?`
  - PATCH `/:order_id?`
  - DELETE `/:order_id?`