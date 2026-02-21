# Auth API Contract

> BE 작성 기준 계약서. 변경 시 FE와 협의 후 수정.

## 공통 에러 형식

```json
{
  "error": "ERROR_CODE",
  "message": "사람이 읽을 수 있는 메시지"
}
```

유효성 검사 오류는 `field` 추가:

```json
{
  "error": "VALIDATION_ERROR",
  "message": "비밀번호는 최소 8자, 영문+숫자 조합이어야 합니다",
  "field": "password"
}
```

## 인증 헤더

인증 필요 API는 아래 헤더 필수:

```
Authorization: Bearer {access_token}
```

---

## POST /api/v1/auth/signup

회원가입. 성공 시 토큰 즉시 발급.

**Request**

```json
{
  "email": "user@example.com",
  "password": "password123",
  "nickname": "냉털러"
}
```

**Response 201**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2g...",
  "expires_in": 3600,
  "user": {
    "id": "64f1a2b3c4d5e6f7a8b9c0d1",
    "email": "user@example.com",
    "nickname": "냉털러",
    "dietary_prefs": [],
    "allergens": []
  }
}
```

**Errors**

| Status | error | 상황 |
|--------|-------|------|
| 400 | `VALIDATION_ERROR` | 이메일 형식 오류 / 비밀번호 규칙 위반 / 닉네임 길이 초과 |
| 409 | `DUPLICATE_EMAIL` | 이미 가입된 이메일 |

```json
// 409
{
  "error": "DUPLICATE_EMAIL",
  "message": "이미 사용 중인 이메일입니다"
}
```

---

## POST /api/v1/auth/login

로그인. 성공 시 토큰 발급.

**Request**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response 200**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2g...",
  "expires_in": 3600,
  "user": {
    "id": "64f1a2b3c4d5e6f7a8b9c0d1",
    "email": "user@example.com",
    "nickname": "냉털러",
    "dietary_prefs": ["vegetarian"],
    "allergens": ["peanut"]
  }
}
```

**Errors**

| Status | error | 상황 |
|--------|-------|------|
| 400 | `VALIDATION_ERROR` | 필수 필드 누락 |
| 401 | `INVALID_CREDENTIALS` | 이메일 또는 비밀번호 불일치 (구분 없이 동일 메시지) |

```json
// 401
{
  "error": "INVALID_CREDENTIALS",
  "message": "이메일 또는 비밀번호가 올바르지 않습니다"
}
```

---

## POST /api/v1/auth/refresh

Access Token 갱신. Refresh Token rotation 적용 (기존 refresh_token 무효화).

**Request**

```json
{
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2g..."
}
```

**Response 200**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...(new)",
  "refresh_token": "bmV3UmVmcmVzaFRva2Vu...(new)",
  "expires_in": 3600
}
```

**Errors**

| Status | error | 상황 |
|--------|-------|------|
| 401 | `INVALID_REFRESH_TOKEN` | 만료되었거나 유효하지 않은 refresh_token |

```json
// 401
{
  "error": "INVALID_REFRESH_TOKEN",
  "message": "유효하지 않은 Refresh Token입니다. 다시 로그인해주세요."
}
```

---

## POST /api/v1/auth/logout

로그아웃. 서버에서 refresh_token 무효화.

**Headers:** `Authorization: Bearer {access_token}` 필수

**Request**

```json
{
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2g..."
}
```

**Response 204**

본문 없음.

**Errors**

| Status | error | 상황 |
|--------|-------|------|
| 401 | `UNAUTHORIZED` | 유효하지 않은 access_token |

---

## GET /api/v1/users/me

내 프로필 조회.

**Headers:** `Authorization: Bearer {access_token}` 필수

**Response 200**

```json
{
  "id": "64f1a2b3c4d5e6f7a8b9c0d1",
  "email": "user@example.com",
  "nickname": "냉털러",
  "dietary_prefs": ["vegetarian"],
  "allergens": ["peanut", "gluten"],
  "created_at": "2025-03-01T09:00:00Z"
}
```

**Errors**

| Status | error | 상황 |
|--------|-------|------|
| 401 | `UNAUTHORIZED` | 토큰 없음 / 만료 |

---

## PATCH /api/v1/users/me

프로필 수정. 변경할 필드만 포함.

**Headers:** `Authorization: Bearer {access_token}` 필수

**Request**

```json
{
  "nickname": "새닉네임",
  "dietary_prefs": ["vegan"],
  "allergens": []
}
```

**Response 200**

```json
{
  "id": "64f1a2b3c4d5e6f7a8b9c0d1",
  "email": "user@example.com",
  "nickname": "새닉네임",
  "dietary_prefs": ["vegan"],
  "allergens": [],
  "updated_at": "2025-03-10T12:00:00Z"
}
```

**Errors**

| Status | error | 상황 |
|--------|-------|------|
| 400 | `VALIDATION_ERROR` | 닉네임 2자 미만 / 20자 초과, 허용되지 않는 dietary_prefs/allergens 값 |
| 401 | `UNAUTHORIZED` | 토큰 없음 / 만료 |

```json
// 400 허용되지 않는 값
{
  "error": "VALIDATION_ERROR",
  "message": "허용되지 않는 dietary_prefs 값입니다: unknown_diet",
  "field": "dietary_prefs"
}
```