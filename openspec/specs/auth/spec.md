# Auth Spec

## 개요

냉털마스터 사용자 인증 및 프로필 관리 도메인.
이메일/패스워드 기반 회원가입·로그인과 JWT 토큰 발급, 개인 식이 설정 관리를 담당한다.

---

## 핵심 엔티티

### User

| 필드              | 타입       | 설명                                   |
|-------------------|------------|----------------------------------------|
| id                | ObjectID   | 사용자 고유 ID                         |
| email             | string     | 로그인 이메일 (unique)                 |
| password_hash     | string     | bcrypt 해시                            |
| nickname          | string     | 표시 이름                              |
| dietary_prefs     | []string   | 식이 제한 (예: vegetarian, vegan)      |
| allergens         | []string   | 알레르기 항목 (예: gluten, peanut)     |
| created_at        | time.Time  | 가입 일시                              |
| updated_at        | time.Time  | 마지막 수정 일시                       |

### Token

| 필드          | 타입      | 설명                       |
|---------------|-----------|----------------------------|
| access_token  | string    | JWT, 만료 1시간            |
| refresh_token | string    | opaque token, 만료 7일     |
| expires_in    | int       | access_token 만료 초       |

---

## 핵심 동작 (Behaviors)

### 회원가입
- 이메일 중복 확인 → 비밀번호 bcrypt 해시 → User 저장
- 가입 직후 access_token + refresh_token 발급

### 로그인
- 이메일/패스워드 검증 → 토큰 쌍 발급
- 잘못된 자격증명: 인증 오류 반환 (이메일/패스워드 구분 없이 동일 메시지)

### 토큰 갱신
- refresh_token 유효성 검증 → 새 access_token 발급
- refresh_token 재사용 방지 (rotation)

### 프로필 수정
- nickname, dietary_prefs, allergens 수정 가능
- email/password 변경은 별도 엔드포인트

### 로그아웃
- refresh_token 무효화 (서버 측 블랙리스트 또는 DB 삭제)

---

## API 엔드포인트

```
POST   /api/v1/auth/signup       회원가입
POST   /api/v1/auth/login        로그인
POST   /api/v1/auth/refresh      토큰 갱신
POST   /api/v1/auth/logout       로그아웃
GET    /api/v1/users/me          내 프로필 조회
PATCH  /api/v1/users/me         프로필 수정
```

---

## 제약 조건

- 비밀번호: 최소 8자, 영문+숫자 조합
- 닉네임: 최소 2자, 최대 20자
- dietary_prefs / allergens 값은 서버에서 허용 목록(enum) 검증
- 모든 인증 필요 API는 Authorization: Bearer {access_token} 헤더 필수

---

## 의존성

- 없음 (다른 도메인의 user_id 참조 기반)
