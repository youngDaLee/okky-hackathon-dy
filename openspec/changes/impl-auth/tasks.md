## 1. 의존성 추가

- [x] 1.1 `go.mod` — `github.com/golang-jwt/jwt/v5` 추가 (`go get`)
- [x] 1.2 `go.mod` — `golang.org/x/crypto` 추가 (bcrypt, 이미 indirect면 direct로 승격)

## 2. Model

- [x] 2.1 `internal/auth/model.go` — `User` 구조체 정의 (id, email, password_hash, nickname, dietary_prefs, allergens, created_at, updated_at)
- [x] 2.2 `internal/auth/model.go` — `RefreshToken` 구조체 정의 (token, user_id, expires_at, created_at)
- [x] 2.3 `internal/auth/model.go` — 요청/응답 타입 정의 (SignupReq, LoginReq, RefreshReq, UpdateProfileReq, TokenResponse, UserResponse)
- [x] 2.4 `internal/auth/model.go` — dietary_prefs / allergens 허용 enum 상수 정의 및 검증 함수

## 3. Repository

- [x] 3.1 `internal/auth/repository.go` — `UserRepository` 인터페이스 및 구조체 정의, `NewUserRepository` 생성자
- [x] 3.2 `internal/auth/repository.go` — `EnsureIndexes` 구현 (email unique index, refresh_token token unique + TTL index)
- [x] 3.3 `internal/auth/repository.go` — `CreateUser` 구현 (InsertOne, 생성된 ID 반환)
- [x] 3.4 `internal/auth/repository.go` — `FindUserByEmail` 구현 (single lookup)
- [x] 3.5 `internal/auth/repository.go` — `FindUserByID` 구현 (ObjectID 파싱 포함)
- [x] 3.6 `internal/auth/repository.go` — `UpdateUser` 구현 ($set partial update)
- [x] 3.7 `internal/auth/repository.go` — `SaveRefreshToken` 구현 (InsertOne)
- [x] 3.8 `internal/auth/repository.go` — `FindRefreshToken` 구현 (token 값으로 조회)
- [x] 3.9 `internal/auth/repository.go` — `DeleteRefreshToken` 구현 (DeleteOne)

## 4. Service

- [x] 4.1 `internal/auth/service.go` — `AuthService` 인터페이스 및 구조체 정의, `NewAuthService` 생성자 (repo + jwtSecret + accessTTL + refreshTTL 주입)
- [x] 4.2 `internal/auth/service.go` — `generateAccessToken` 헬퍼 (JWT HS256, sub=userID, exp=now+1h)
- [x] 4.3 `internal/auth/service.go` — `generateRefreshToken` 헬퍼 (crypto/rand 32바이트 base64)
- [x] 4.4 `internal/auth/service.go` — `Signup` 구현 (이메일 중복 확인, bcrypt 해시, User 저장, 토큰 발급)
- [x] 4.5 `internal/auth/service.go` — `Login` 구현 (이메일 조회, bcrypt 비교, 토큰 발급)
- [x] 4.6 `internal/auth/service.go` — `Refresh` 구현 (DB에서 refresh_token 조회, 만료 확인, rotation — 기존 삭제 후 신규 저장)
- [x] 4.7 `internal/auth/service.go` — `Logout` 구현 (refresh_token DB 삭제)
- [x] 4.8 `internal/auth/service.go` — `GetMe` 구현 (userID로 User 조회)
- [x] 4.9 `internal/auth/service.go` — `UpdateMe` 구현 (소유권 불필요, partial update, enum 검증 포함)
- [x] 4.10 `internal/auth/service.go` — `ValidateAccessToken` 구현 (JWT 파싱, claims.Subject 반환) — 미들웨어에서 사용

## 5. Handler

- [x] 5.1 `internal/auth/handler.go` — `AuthHandler` 구조체 및 `NewAuthHandler` 생성자
- [x] 5.2 `internal/auth/handler.go` — `Signup` 핸들러 (binding → service.Signup → 201 응답)
- [x] 5.3 `internal/auth/handler.go` — `Login` 핸들러 (binding → service.Login → 200 응답)
- [x] 5.4 `internal/auth/handler.go` — `Refresh` 핸들러 (binding → service.Refresh → 200 응답)
- [x] 5.5 `internal/auth/handler.go` — `Logout` 핸들러 (binding → service.Logout → 204 응답)
- [x] 5.6 `internal/auth/handler.go` — `GetMe` 핸들러 (c.GetString("userID") → service.GetMe → 200 응답)
- [x] 5.7 `internal/auth/handler.go` — `UpdateMe` 핸들러 (binding → service.UpdateMe → 200 응답)

## 6. 미들웨어 교체 & 라우터

- [x] 6.1 `internal/server/middleware/auth.go` — 스텁 제거, `AuthService.ValidateAccessToken` 호출로 실제 JWT 검증 구현
- [x] 6.2 `internal/server/router.go` — auth 라우트 4개 등록 (`/auth/signup`, `/auth/login`, `/auth/refresh`, `/auth/logout`)
- [x] 6.3 `internal/server/router.go` — users 라우트 2개 등록 (`/users/me` GET/PATCH, auth 미들웨어 적용)

## 7. DI 와이어링

- [x] 7.1 `cmd/server/main.go` — auth 컬렉션 2개 (users, refresh_tokens) 가져오기
- [x] 7.2 `cmd/server/main.go` — `NewUserRepository`, `NewAuthService`, `NewAuthHandler` 생성 및 RouterDeps에 주입
- [x] 7.3 `cmd/server/main.go` — `authHandler.Service.EnsureIndexes` 호출 (startup)
- [x] 7.4 `internal/server/router.go` — `RouterDeps`에 `AuthHandler *auth.AuthHandler` 필드 추가
