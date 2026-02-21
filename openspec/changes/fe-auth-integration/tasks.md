# fe-auth-integration Tasks

> 백엔드 auth API 연동 및 AddIngredientView 필드 수정, 로그아웃 버튼 추가

## 1. Auth API 클라이언트

- [x] 1.1 `frontend/src/api/auth.js` 생성
  - signup, login, refresh, logout, getMe 엔드포인트

## 2. LoginView 연동

- [x] 2.1 `LoginView.vue` — `useAuthStore.login()` 호출로 백엔드 연동
- [x] 2.2 에러 처리: `INVALID_CREDENTIALS` → "이메일 또는 비밀번호가 올바르지 않아요"

## 3. SignupView 연동

- [x] 3.1 `SignupView.vue` — `useAuthStore.signup()` 호출로 백엔드 연동
- [x] 3.2 에러 처리: `DUPLICATE_EMAIL`, `VALIDATION_ERROR` (field별 메시지)

## 4. AddIngredientView 필드 수정

- [x] 4.1 카테고리 값 한국어 → 영문 Enum(`VEGETABLE`, `FRUIT` 등)으로 수정
- [x] 4.2 필드명 `expiry_date` → `expiryDate` (camelCase) 수정
- [x] 4.3 날짜 ISO datetime 변환 (`new Date(...).toISOString()`)
- [x] 4.4 수량 min 0 → 1 (`gt=0` 백엔드 검증 맞춤)

## 5. MiniInventory 버그 수정

- [x] 5.1 더보기 링크 `/add-ingredient` → `/fridge` 수정
- [x] 5.2 카테고리 이모지 맵 키 한국어 → 영문 Enum으로 수정

## 6. 로그아웃 버튼

- [x] 6.1 `DashboardView.vue` 헤더 우측에 로그아웃 버튼 추가
- [x] 6.2 클릭 시 `useAuthStore.logout()` → `/login` 리다이렉트
