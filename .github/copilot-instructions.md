<!-- .github/copilot-instructions.md - targeted guidance for AI coding agents -->
# ProiectColectivBackEnd — Copilot instructions (concise)

Goal: be productive quickly. Below are the concrete architecture, patterns, commands and examples discovered in this repo.

- Project type: Go web service using Gin + Firebase Realtime Database. Entry: `main.go`.
- Key packages/directories:
  - `config/` — Firebase initialization (`InitFirebase` in `db_config.go`) using `godotenv` and env vars:
    - Required env vars: `FIREBASE_DATABASE_URL`, `FIREBASE_CREDENTIALS_PATH` (path to service account JSON).
  - `routes/` — HTTP routing. `routes.SetupRoutes()` wires `SetupUserRoutes` and the root route.
  - `controller/` — HTTP controllers (example: `controller/user_controller.go`). Controllers declare small interfaces (e.g. `UserServiceInterface`) to enable DI in tests.
  - `service/` — Business logic (example: `service/user_service.go`). Services use repository interfaces and provide test constructors like `NewUserServiceWithRepo`.
  - `persistence/` — Data access backed by Firebase Realtime DB (example: `persistence/user_repository.go`). Uses `config.FirebaseDB.NewRef("users/...")` and firebase queries (`OrderByChild`, `GetOrdered`).
  - `model/entity`, `model/dto` — data shapes and request/response DTOs.
  - `tests/` — unit-test helpers and mocks (`tests/mocks.go`, `tests/test_data.go`). Tests use `testify/mock` and mock constructors living in `tests` package.

Notable patterns and conventions (explicit, reproducible):
- Dependency injection is accomplished by providing alternate constructors. Examples:
  - `service.NewUserServiceWithRepo(repo interface{})` — used by tests to inject `tests.MockUserRepository`.
  - `controller.NewUserControllerWithService(userService UserServiceInterface)` — used by controller tests to inject mock services.
- Repositories return concrete `*entity.User` and errors; service methods wrap repository calls and translate validation/uniqueness checks into errors (e.g. "username already exists").
- ID generation in `service.generateID()` uses crypto-rand + hex; passwords are hashed with `bcrypt.GenerateFromPassword`.
- Controllers map errors to HTTP status codes in a simple way (e.g. not-found -> 404, validation/service error -> 500 for SignUp). Keep this behavior when making changes.

Developer workflows and exact commands
- Install dependencies / setup: `go mod tidy` (README recommends this).
- Run server (local dev):
  1. Create a `.env` with `FIREBASE_DATABASE_URL` and `FIREBASE_CREDENTIALS_PATH` and place the service account JSON at the path.
  2. Start: `go run main.go` — server listens on `:8080`.
- Unit tests (fast, do NOT require Firebase because tests use mocks):
  - Run all unit tests: `go test ./...`
  - Run a package: `go test ./service -v` or `go test ./controller -v`
  - Test files use `github.com/stretchr/testify` mocks — look at `tests/mocks.go` to see available mocks.
- Linting / quick checks: `gofmt -w .`, `go vet ./...` (not configured in repo, but safe defaults).

Integration and external dependencies
- Firebase Realtime DB is the single external datastore. Initialization occurs in `config.InitFirebase()` and uses `godotenv` — the app expects `.env` (or environment) at runtime.
- For local integration testing against Firebase: set env vars and ensure the credentials file exists. Most unit tests avoid hitting Firebase by using mocks.

How to modify code safely (common edits)
- When adding a new controller or service:
  - Add an interface type near the controller/service so tests can inject mocks (this repo consistently declares small interfaces adjacent to the consumer).
  - Provide a `NewXWithY` constructor for tests (pattern used in `service` and `controller`).
- When adding persistence methods:
  - Use `config.FirebaseDB.NewRef("resource/...")` and follow the existing query/unmarshal patterns in `persistence/user_repository.go`.

Testing notes / examples from repo
- Controller tests use Gin test context: create recorder with `httptest.NewRecorder()` and `gin.CreateTestContext(w)` — see `tests/controller/user_controller_test.go`.
- Service tests inject a mock repo (see `tests.MockUserRepository`) with expectations on `GetByUsername`, `GetByEmail`, `Create`.

Files to inspect for details when editing behavior
- `main.go` — app entry (calls `config.InitFirebase()` and `routes.SetupRoutes()`).
- `config/db_config.go` — env var names and Firebase initialization.
- `persistence/user_repository.go` — Firebase read/write/query patterns.
- `service/user_service.go` — validation, password hashing, ID creation, and returned DTOs.
- `controller/user_controller.go` — HTTP bindings, JSON binding patterns, and status code mapping.
- `tests/` — mocks and sample test data used by unit tests.

If you need more context or CI commands
- I merged repository-discoverable facts only. If there are CI workflows, secret management, or additional integration tests you want included, point me to `.github/workflows` or share the CI pipeline and I'll incorporate exact commands.

Short checklist for an AI making edits
1. Run `go test ./...` locally after changes (most unit tests use mocks).
2. Preserve constructor-with-mock patterns (`NewXWith...`) to keep tests working.
3. Keep Firebase access isolated to `persistence/` and `config/` so unit tests can use mocks.
4. Update `README.md` only if runtime steps or env var names change.

Please review: is there any CI, secret storage, or deployment commands I missed that should be added here?
