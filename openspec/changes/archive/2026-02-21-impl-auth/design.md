## Context

auth 스켈레톤(`internal/auth/`)이 존재하나 모든 메서드가 스텁 상태.
현재 `middleware/auth.go`는 Bearer 토큰을 파싱만 하고 JWT 검증 없이 token 값 자체를 userID로 주입하는 개발용 스텁이다.
다른 모든 도메인 API가 `c.Get("userID")`에 의존하므로 auth 미들웨어가 먼저 완성되어야 한다.

MongoDB `users`, `refresh_tokens` 두 컬렉션 신규 생성.

## Goals / Non-Goals

**Goals:**
- email/password 기반 회원가입·로그인 구현
- JWT access_token(1h) + opaque refresh_token(7d) 발급 및 rotation
- JWT 미들웨어 실제 검증 구현 (스텁 교체)
- 프로필 조회·수정 (nickname, dietary_prefs, allergens)

**Non-Goals:**
- 소셜 로그인, 이메일 인증, 비밀번호 재설정
- Rate limiting, 관리자 권한

## Decisions

### 1. refresh_token 저장 방식: DB 저장 (opaque token)

JWT로 만들어도 되지만 rotation + 즉시 무효화를 위해 opaque token을 MongoDB `refresh_tokens` 컬렉션에 저장한다.
- `token` (string, unique index), `user_id`, `expires_at`, `created_at`
- logout/rotation 시 해당 token document 삭제 → 재사용 방지

**Alternative**: JWT refresh token — 무효화를 위한 blacklist DB가 결국 필요하므로 동일한 복잡도.

### 2. JWT 라이브러리: `github.com/golang-jwt/jwt/v5`

Go 생태계 표준. HS256 알고리즘, `JWTSecret`은 config에서 주입.
Payload: `sub` (user_id), `exp`, `iat`.

### 3. bcrypt cost: 10

개발/해커톤 환경에서 로그인 응답 속도 균형. 프로덕션 배포 시 12로 상향 가능.

### 4. dietary_prefs / allergens 허용 값: 서버 enum 검증

허용 목록을 상수로 정의하여 검증. 빈 배열(`[]`) 허용.

### 5. 미들웨어 교체 전략

현재 스텁은 token 값 = userID로 가정. 실제 구현 후 `c.Set("userID", claims.Subject)` 형태로 교체.
fridge 도메인은 이미 `c.GetString("userID")`를 string으로 사용 중 — ObjectID 파싱은 각 handler에서 처리.

## Risks / Trade-offs

- **[Risk] bcrypt 10이 해커톤 배포 환경에서 느릴 수 있다** → 서버 사양에 따라 cost 조정 가능하도록 config에 추가 고려
- **[Risk] refresh_token MongoDB TTL 인덱스 미설정 시 만료된 토큰 누적** → `expires_at` 필드에 TTL 인덱스 설정 (EnsureIndexes)
- **[Trade-off] opaque refresh_token은 DB 조회 필요** → 해커톤 규모에서는 무시 가능한 성능 비용

## Migration Plan

1. `internal/auth/` 도메인 패키지 구현
2. `middleware/auth.go` 실제 JWT 검증으로 교체
3. `router.go`에 auth 라우트 추가, main.go에 DI 추가
4. `go.mod` 의존성 추가 (`golang-jwt/jwt`, `golang.org/x/crypto`)

롤백: middleware를 스텁으로 되돌리면 이전 상태로 복원 가능.

## Open Questions

- dietary_prefs / allergens enum 목록을 config으로 외부화할지, 코드 상수로 둘지 → 해커톤 기간 동안 코드 상수로 고정