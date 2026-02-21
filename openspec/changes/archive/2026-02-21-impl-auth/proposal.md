## Why

스켈레톤으로만 존재하는 auth 도메인에 실제 인증 로직이 없어 다른 모든 도메인의 `userID` 기반 동작이 불가하다. Issue #2: 회원가입/로그인/토큰 관리 API가 MVP의 선행 조건.

## What Changes

- `POST /api/v1/auth/signup` — 이메일 중복 확인, bcrypt 해시, User 저장, 토큰 즉시 발급
- `POST /api/v1/auth/login` — 자격증명 검증 후 JWT access_token(1h) + refresh_token(7d) 발급
- `POST /api/v1/auth/refresh` — refresh_token 검증 및 rotation (기존 토큰 무효화)
- `POST /api/v1/auth/logout` — refresh_token DB 삭제로 무효화
- `GET /api/v1/users/me` — 내 프로필 조회
- `PATCH /api/v1/users/me` — nickname, dietary_prefs, allergens 수정
- JWT Bearer 인증 미들웨어 실제 구현 (현재 스텁 → 실제 검증으로 교체)

### Non-goals

- 소셜 로그인 (카카오, 구글)
- 이메일 인증 / 비밀번호 재설정
- 관리자 권한 분리
- Rate limiting

## Capabilities

### New Capabilities

- `auth`: 이메일/패스워드 기반 회원가입·로그인, JWT 토큰 발급·갱신·무효화, 프로필 관리

### Modified Capabilities

- (없음 — auth 스펙은 이미 정의되어 있고 요구사항 변경 없음)

## Impact

- `backend/internal/auth/` — 신규 도메인 패키지 생성 (model, repository, service, handler)
- `backend/internal/server/middleware/auth.go` — 스텁 미들웨어를 실제 JWT 검증으로 교체
- `backend/internal/server/router.go` — auth, users 라우트 추가
- `backend/cmd/server/main.go` — auth DI 와이어링 추가
- `backend/go.mod` — `golang.org/x/crypto` (bcrypt), JWT 라이브러리 추가
- MongoDB `users` 컬렉션, `refresh_tokens` 컬렉션 신규 생성