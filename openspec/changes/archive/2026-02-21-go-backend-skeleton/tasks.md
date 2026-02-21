## 1. Go 모듈 초기화

- [x] 1.1 `backend/` 디렉토리 생성
- [x] 1.2 `go mod init github.com/okky-hackathon/fridge-master-backend` 실행
- [x] 1.3 `go.mod`에 주요 의존성 추가 (gin, mongo-driver, jwt, bcrypt, cron, godotenv, gcs, vision)
- [x] 1.4 `go mod tidy` 실행하여 `go.sum` 생성

## 2. 진입점 및 서버 패키지

- [x] 2.1 `cmd/server/main.go` 생성 — `package main` + 빈 `main()` 함수
- [x] 2.2 `internal/server/router.go` 생성 — `package server` + 빈 `NewRouter()` 함수
- [x] 2.3 `internal/server/middleware/auth.go` 생성 — `package middleware` + 빈 `Auth()` 함수
- [x] 2.4 `internal/server/middleware/cors.go` 생성 — `package middleware` + 빈 `CORS()` 함수
- [x] 2.5 `internal/server/middleware/logger.go` 생성 — `package middleware` + 빈 `Logger()` 함수

## 3. 공유 유틸 패키지 (pkg/)

- [x] 3.1 `pkg/config/config.go` 생성 — `Config` 구조체 + 빈 `Load()` 함수
- [x] 3.2 `pkg/database/mongo.go` 생성 — 빈 `Connect()`, `GetCollection()` 함수
- [x] 3.3 `pkg/jwt/jwt.go` 생성 — 빈 `GenerateAccessToken()`, `GenerateRefreshToken()`, `ValidateAccessToken()`, `HashToken()` 함수
- [x] 3.4 `pkg/gcs/client.go` 생성 — `Client` 구조체 + 빈 `Upload()`, `GenerateSignedURL()`, `Delete()` 함수
- [x] 3.5 `pkg/vertexai/client.go` 생성 — `Client` 구조체 + 빈 `AnalyzeReceipt()`, `AnalyzeFridge()` 함수

## 4. auth 도메인 스켈레톤

- [x] 4.1 `internal/auth/model.go` — `User`, `RefreshToken`, 요청/응답 DTO 구조체 선언
- [x] 4.2 `internal/auth/repository.go` — `UserRepository` 인터페이스 + 빈 구현체
- [x] 4.3 `internal/auth/service.go` — `AuthService` 인터페이스 + 빈 구현체
- [x] 4.4 `internal/auth/handler.go` — 빈 `AuthHandler` 구조체 + 빈 핸들러 메서드

## 5. fridge 도메인 스켈레톤

- [x] 5.1 `internal/fridge/model.go` — `Ingredient`, DTO, enum 상수 선언
- [x] 5.2 `internal/fridge/repository.go` — `IngredientRepository` 인터페이스 + 빈 구현체
- [x] 5.3 `internal/fridge/service.go` — `FridgeService` 인터페이스 + 빈 구현체
- [x] 5.4 `internal/fridge/handler.go` — 빈 `FridgeHandler` 구조체 + 빈 핸들러 메서드

## 6. vision 도메인 스켈레톤

- [x] 6.1 `internal/vision/model.go` — `VisionJob`, `ExtractedItem`, DTO, Status enum 선언
- [x] 6.2 `internal/vision/repository.go` — `VisionRepository` 인터페이스 + 빈 구현체
- [x] 6.3 `internal/vision/service.go` — `VisionService` 인터페이스 + 빈 구현체
- [x] 6.4 `internal/vision/handler.go` — 빈 `VisionHandler` 구조체 + 빈 핸들러 메서드
- [x] 6.5 `internal/vision/worker.go` — 빈 `StartWorker()` 함수 + `jobCh` 채널 선언

## 7. recommendation 도메인 스켈레톤

- [x] 7.1 `internal/recommendation/model.go` — `Recipe`, `RecommendationResult`, DTO 선언
- [x] 7.2 `internal/recommendation/repository.go` — `RecipeRepository` 인터페이스 + 빈 구현체
- [x] 7.3 `internal/recommendation/service.go` — `RecommendationService` 인터페이스 + 빈 구현체
- [x] 7.4 `internal/recommendation/handler.go` — 빈 `RecommendationHandler` 구조체 + 빈 핸들러 메서드
- [x] 7.5 `internal/recommendation/external.go` — `ExternalSearcher` 인터페이스 + 빈 `ExternalClient` 구조체

## 8. cookbook 도메인 스켈레톤

- [x] 8.1 `internal/cookbook/model.go` — `SavedRecipe`, `RecipeSnap`, `LabelSummary`, DTO 선언
- [x] 8.2 `internal/cookbook/repository.go` — `CookbookRepository` 인터페이스 + 빈 구현체
- [x] 8.3 `internal/cookbook/service.go` — `CookbookService` 인터페이스 + 빈 구현체
- [x] 8.4 `internal/cookbook/handler.go` — 빈 `CookbookHandler` 구조체 + 빈 핸들러 메서드

## 9. notification 도메인 스켈레톤

- [x] 9.1 `internal/notification/model.go` — `ExpiryAlert`, DTO, AlertType enum 선언
- [x] 9.2 `internal/notification/repository.go` — `NotificationRepository` 인터페이스 + 빈 구현체
- [x] 9.3 `internal/notification/service.go` — `NotificationService` 인터페이스 + 빈 구현체
- [x] 9.4 `internal/notification/handler.go` — 빈 `NotificationHandler` 구조체 + 빈 핸들러 메서드
- [x] 9.5 `internal/notification/scheduler.go` — 빈 `StartScheduler()` 함수

## 10. 최종 검증

- [x] 10.1 `go build ./...` 실행하여 컴파일 오류 없음 확인
- [x] 10.2 `go vet ./...` 실행하여 정적 분석 오류 없음 확인
