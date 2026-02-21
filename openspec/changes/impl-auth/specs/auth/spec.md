## ADDED Requirements

### Requirement: 회원가입
시스템은 이메일/패스워드로 신규 사용자를 등록하고 즉시 토큰을 발급해야 한다.
- 이메일 형식 검증 필수
- 비밀번호는 최소 8자, 영문+숫자 조합 필수
- 닉네임 2~20자 필수
- 이메일 중복 시 409 DUPLICATE_EMAIL 반환
- 성공 시 201 + access_token, refresh_token, expires_in, user 반환

#### Scenario: 정상 회원가입
- **WHEN** 유효한 email, password, nickname으로 POST /api/v1/auth/signup 요청
- **THEN** 201 응답, access_token(JWT 1h) + refresh_token(7d) + user 객체 반환

#### Scenario: 이메일 중복 회원가입
- **WHEN** 이미 가입된 email로 signup 요청
- **THEN** 409 DUPLICATE_EMAIL 에러 반환

#### Scenario: 비밀번호 규칙 위반
- **WHEN** 숫자 없는 8자 이상 문자열로 password 설정
- **THEN** 400 VALIDATION_ERROR, field: "password" 반환

#### Scenario: 닉네임 길이 초과
- **WHEN** 21자 이상 nickname으로 signup 요청
- **THEN** 400 VALIDATION_ERROR, field: "nickname" 반환

### Requirement: 로그인
시스템은 등록된 이메일/패스워드로 인증 후 토큰을 발급해야 한다.
- 이메일 또는 패스워드 불일치 시 동일한 401 에러 반환 (enumeration 방지)
- 성공 시 200 + access_token, refresh_token, expires_in, user 반환

#### Scenario: 정상 로그인
- **WHEN** 올바른 email/password로 POST /api/v1/auth/login 요청
- **THEN** 200 응답, access_token + refresh_token + user 반환

#### Scenario: 잘못된 비밀번호
- **WHEN** 존재하는 email에 틀린 password로 로그인 시도
- **THEN** 401 INVALID_CREDENTIALS 반환 (이메일/비밀번호 오류 구분 없음)

#### Scenario: 존재하지 않는 이메일
- **WHEN** 미등록 email로 로그인 시도
- **THEN** 401 INVALID_CREDENTIALS 반환 (이메일/비밀번호 오류 구분 없음)

### Requirement: 토큰 갱신
시스템은 유효한 refresh_token으로 새 access_token을 발급하고 refresh_token을 rotation해야 한다.
- 기존 refresh_token은 즉시 무효화 (재사용 불가)
- 만료되거나 존재하지 않는 refresh_token은 401 반환

#### Scenario: 정상 토큰 갱신
- **WHEN** 유효한 refresh_token으로 POST /api/v1/auth/refresh 요청
- **THEN** 200 응답, 새 access_token + 새 refresh_token 반환, 기존 refresh_token 무효화

#### Scenario: 만료된 refresh_token 갱신 시도
- **WHEN** 만료된 refresh_token으로 refresh 요청
- **THEN** 401 INVALID_REFRESH_TOKEN 반환

#### Scenario: 이미 사용된 refresh_token 재사용 시도
- **WHEN** rotation으로 무효화된 refresh_token으로 refresh 요청
- **THEN** 401 INVALID_REFRESH_TOKEN 반환

### Requirement: 로그아웃
시스템은 refresh_token을 서버에서 삭제하여 무효화해야 한다.

#### Scenario: 정상 로그아웃
- **WHEN** 유효한 access_token + refresh_token으로 POST /api/v1/auth/logout 요청
- **THEN** 204 응답, refresh_token DB에서 삭제

#### Scenario: 인증 없이 로그아웃 시도
- **WHEN** Authorization 헤더 없이 logout 요청
- **THEN** 401 UNAUTHORIZED 반환

### Requirement: 내 프로필 조회
인증된 사용자는 자신의 프로필 정보를 조회할 수 있어야 한다.

#### Scenario: 정상 프로필 조회
- **WHEN** 유효한 access_token으로 GET /api/v1/users/me 요청
- **THEN** 200 응답, id/email/nickname/dietary_prefs/allergens/created_at 반환

#### Scenario: 미인증 프로필 조회
- **WHEN** Authorization 헤더 없이 GET /api/v1/users/me 요청
- **THEN** 401 UNAUTHORIZED 반환

### Requirement: 프로필 수정
인증된 사용자는 nickname, dietary_prefs, allergens를 수정할 수 있어야 한다.
- 변경할 필드만 포함 가능 (partial update)
- dietary_prefs/allergens는 허용 enum 목록 검증
- 닉네임 2~20자 검증

#### Scenario: 닉네임 수정
- **WHEN** 유효한 새 nickname으로 PATCH /api/v1/users/me 요청
- **THEN** 200 응답, 수정된 프로필 반환

#### Scenario: 허용되지 않는 dietary_prefs 값
- **WHEN** 허용 목록에 없는 dietary_prefs 값으로 PATCH 요청
- **THEN** 400 VALIDATION_ERROR, field: "dietary_prefs" 반환

### Requirement: JWT Bearer 미들웨어
모든 인증 필요 API는 Authorization: Bearer {access_token} 헤더로 보호되어야 한다.
- 유효한 JWT → claims에서 sub(user_id) 추출 → gin context에 "userID" 주입
- 만료·변조된 JWT → 401 UNAUTHORIZED

#### Scenario: 유효한 JWT로 인증 필요 API 접근
- **WHEN** 만료되지 않은 유효한 JWT access_token으로 요청
- **THEN** 미들웨어 통과, c.GetString("userID")에 사용자 ID 설정

#### Scenario: 만료된 JWT로 요청
- **WHEN** 만료된 access_token으로 인증 필요 API 요청
- **THEN** 401 UNAUTHORIZED 반환, 요청 중단
