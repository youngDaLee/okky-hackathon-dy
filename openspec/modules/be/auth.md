# Module: internal/auth

## 패키지 경로
`backend/internal/auth`

## 책임
사용자 인증·인가 및 프로필 관리.
JWT 발급/갱신/무효화, 비밀번호 해시 처리.

---

## 파일 구조

```
internal/auth/
├── model.go
├── repository.go
├── service.go
└── handler.go
```

---

## model.go

### 구조체

| 구조체             | 역할                                        |
|--------------------|---------------------------------------------|
| `User`             | MongoDB `users` 컬렉션 문서                 |
| `RefreshToken`     | MongoDB `refresh_tokens` 컬렉션 문서        |
| `SignupRequest`    | POST /auth/signup 요청 DTO                  |
| `LoginRequest`     | POST /auth/login 요청 DTO                   |
| `TokenResponse`    | 토큰 발급 응답 DTO                          |
| `UpdateMeRequest`  | PATCH /users/me 요청 DTO                    |
| `UserResponse`     | GET /users/me 응답 DTO (password_hash 제외) |

### User 필드
```
ID            primitive.ObjectID  `bson:"_id"`
Email         string              `bson:"email"`
PasswordHash  string              `bson:"password_hash"`
Nickname      string              `bson:"nickname"`
DietaryPrefs  []string            `bson:"dietary_prefs"`
Allergens     []string            `bson:"allergens"`
CreatedAt     time.Time           `bson:"created_at"`
UpdatedAt     time.Time           `bson:"updated_at"`
```

### RefreshToken 필드
```
ID        primitive.ObjectID  `bson:"_id"`
UserID    primitive.ObjectID  `bson:"user_id"`
Token     string              `bson:"token"`      // hashed
ExpiresAt time.Time           `bson:"expires_at"`
CreatedAt time.Time           `bson:"created_at"`
```

---

## repository.go

### 인터페이스 (UserRepository)
```
FindByEmail(ctx, email) (*User, error)
FindByID(ctx, id) (*User, error)
Create(ctx, user) error
Update(ctx, id, update) error

SaveRefreshToken(ctx, token) error
FindRefreshToken(ctx, tokenHash) (*RefreshToken, error)
DeleteRefreshToken(ctx, tokenHash) error
DeleteAllRefreshTokens(ctx, userID) error
```

### MongoDB 컬렉션
- `users` — unique index on `email`
- `refresh_tokens` — TTL index on `expires_at`, index on `token`

---

## service.go

### 인터페이스 (AuthService)
```
Signup(ctx, req SignupRequest) (*TokenResponse, error)
Login(ctx, req LoginRequest) (*TokenResponse, error)
Refresh(ctx, refreshToken string) (*TokenResponse, error)
Logout(ctx, refreshToken string) error
GetMe(ctx, userID) (*UserResponse, error)
UpdateMe(ctx, userID, req UpdateMeRequest) (*UserResponse, error)
```

### 핵심 로직 포인트
- Signup: Email 중복 체크 → bcrypt hash → 저장 → 토큰 발급
- Login: 비밀번호 비교 → 토큰 쌍 발급
- Refresh: RefreshToken 검증 → rotation (기존 삭제 → 신규 발급)
- 토큰 발급은 `pkg/jwt` 패키지에 위임

---

## handler.go

### 엔드포인트 등록
```
POST   /api/v1/auth/signup
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
GET    /api/v1/users/me      (JWT 미들웨어 필요)
PATCH  /api/v1/users/me     (JWT 미들웨어 필요)
```

---

## 외부 의존성

| 의존 대상          | 용도                |
|--------------------|---------------------|
| `pkg/jwt`          | 토큰 생성·검증      |
| `pkg/database`     | MongoDB 컬렉션 접근 |
| `golang.org/x/crypto/bcrypt` | 비밀번호 해시 |
