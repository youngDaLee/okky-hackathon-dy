## 1. 사전 준비

- [x] 1.1 `lucide-vue-next` 패키지 설치
  ```bash
  cd frontend && npm install lucide-vue-next
  ```

## 2. 레이아웃 컴포넌트

- [x] 2.1 `frontend/src/layouts/` 디렉토리 생성
- [x] 2.2 `AppLayout.vue` 생성
  - `RouterView` + `AppBottomNav` + `FloatingAddButton` 포함
  - `isAddIngredientPage` computed로 FAB 표시 제어
  - 라우터 가드: `access_token` 없으면 `/login` 리다이렉트

## 3. 공통 컴포넌트

- [x] 3.1 `AppBottomNav.vue` 생성
  - 홈(`/`) · 추가(`/add-ingredient`) · 레시피(`/recipes`) 탭 3개
  - 가운데 추가 버튼: 파란 원형 FAB, 위로 돌출 스타일
  - `useRoute()`로 현재 경로 감지하여 활성 탭 표시
  - `lucide-vue-next`에서 `Home`, `Plus`, `ChefHat` 아이콘 사용

- [x] 3.2 `FloatingAddButton.vue` 생성
  - `fixed bottom-24 right-4` 위치
  - 클릭 시 `/add-ingredient` 이동
  - 파란 배경 + `+` 텍스트 또는 아이콘

## 4. 뷰 스캐폴딩

- [x] 4.1 `DashboardView.vue` 생성
  - 헤더 텍스트만 포함 (`🍳 나의 냉장고`)
  - Phase 2에서 실제 콘텐츠 채울 자리

- [x] 4.2 `RecipesView.vue` 빈 스캐폴딩 생성 (라우터 연결용)
- [x] 4.3 `AddIngredientView.vue` 빈 스캐폴딩 생성 (라우터 연결용)
- [x] 4.4 `RecipeDetailView.vue` 빈 스캐폴딩 생성 (라우터 연결용)

## 5. 라우터 재구성

- [x] 5.1 `router/index.js` 수정
  - 인증 전 라우트: `/home`, `/login`, `/signup`
  - 인증 후 라우트: `AppLayout` 하위 nested routes (`/`, `/recipes`, `/recipes/:id`, `/add-ingredient`)
  - `meta: { requiresAuth: true }` 추가

- [x] 5.2 `router.beforeEach` 네비게이션 가드 추가
  - `access_token` 없을 때 보호된 라우트 → `/login` 리다이렉트
  - 이미 로그인 상태에서 `/login`, `/signup` 접근 → `/` 리다이렉트
  - ⚠️ **개발 중 임시 비활성화** (`router/index.js:34-43` 주석 처리) — 실제 인증 API 연동 시 복원 필요

- [x] 5.3 기존 `HomeView.vue` 라우트 경로 `/` → `/home` 변경 확인

## 6. 검증

- [x] 6.1 `npm run build` 빌드 성공 확인 (1.52s, 오류 없음)
- [ ] 6.2 미로그인 상태에서 `/` 접근 → `/login` 리다이렉트 확인
- [ ] 6.3 로그인 상태(localStorage에 access_token 수동 설정)에서 Dashboard 진입 확인
- [x] 6.4 BottomNav 3탭 클릭 → 각 라우트 이동 확인
- [x] 6.5 `/add-ingredient` 페이지에서 FAB 숨김 확인
