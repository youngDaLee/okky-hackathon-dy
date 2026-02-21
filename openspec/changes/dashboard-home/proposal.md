## Why

layout-bottom-nav에서 스캐폴딩만 만들어둔 DashboardView에 실제 콘텐츠를 구현한다.
냉장고 재료 현황과 레시피 추천을 한눈에 볼 수 있는 홈 화면이 필요하다.

## What Changes

- `frontend/src/views/DashboardView.vue` 수정 — AlertBar + RecipeCarousel + MiniInventory + Cookbook 조립
- `frontend/src/components/AlertBar.vue` 신규 생성 — 유통기한 임박 경고 바
- `frontend/src/components/MiniInventory.vue` 신규 생성 — 냉장고 현황 미니 뷰 (최대 5개)
- `frontend/src/components/RecipeCarousel.vue` 신규 생성 — 100% 매칭 레시피 캐러셀
- `frontend/src/components/Cookbook.vue` 신규 생성 — 저장된 요리책 (북마크)
- `frontend/src/stores/ingredient.js` 신규 생성 — Pinia 재료 스토어
- `frontend/src/api/fridge.js` 신규 생성 — 냉장고 API 클라이언트

## Capabilities

### New Capabilities

- `fe-dashboard`: Dashboard 홈 화면 (AlertBar, RecipeCarousel, MiniInventory, Cookbook)
- `fe-ingredient-store`: Pinia 재료 스토어 + 냉장고 API 클라이언트

### Modified Capabilities

- `fe-dashboard`: DashboardView 스캐폴딩 → 실제 컴포넌트 구현

## Impact

- `frontend/src/views/DashboardView.vue` 전면 수정
- `frontend/src/components/` 4개 컴포넌트 추가
- `frontend/src/stores/ingredient.js` 신규
- `frontend/src/api/fridge.js` 신규

### Non-goals

- RecipeCarousel 실제 API 연동 (레시피 API 미구현)
- Cookbook 실제 API 연동 (북마크 API 미구현)
- Dashboard 무한 스크롤, 새로고침 등 UX 개선
