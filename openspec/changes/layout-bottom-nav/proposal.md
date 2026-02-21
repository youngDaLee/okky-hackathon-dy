## Why

냉털마스터 Vue 3 프론트엔드의 인증 후 공통 레이아웃을 구축한다.
현재 로그인/회원가입/홈 뷰만 존재하며, 인증 후 진입하는 Dashboard·Recipes·AddIngredient 페이지를 위한
공통 레이아웃(BottomNav + FloatingAddButton)과 라우터 구조가 없다.

## What Changes

- `frontend/src/layouts/AppLayout.vue` 신규 생성 — 인증 후 공통 레이아웃 (BottomNav + RouterView)
- `frontend/src/components/AppBottomNav.vue` 신규 생성 — 하단 고정 네비게이션 (홈/추가/레시피 3탭)
- `frontend/src/components/FloatingAddButton.vue` 신규 생성 — 재료 추가 FAB
- `frontend/src/router/index.js` 수정 — 인증 후 라우트를 AppLayout 하위로 재구성, 라우터 가드 추가
- `frontend/src/views/DashboardView.vue` 신규 생성 — 대시보드 스캐폴딩 (빈 뷰)

## Capabilities

### New Capabilities

- `fe-layout`: 인증 후 공통 레이아웃 컴포넌트 (AppLayout, AppBottomNav, FloatingAddButton)

### Modified Capabilities

- `fe-router`: 인증 후 라우트 구조 재편 (AppLayout nested routes + 라우터 가드)

## Impact

- `frontend/src/layouts/` 디렉토리 신규 생성
- `frontend/src/components/` 컴포넌트 2개 추가
- `frontend/src/router/index.js` 라우트 구조 변경 (기존 로그인/회원가입 라우트 유지)
- `frontend/src/views/DashboardView.vue` 신규 추가

### Non-goals

- Dashboard 실제 콘텐츠 구현 (AlertBar, RecipeCarousel 등은 Phase 2)
- 실제 인증 API 연동 (라우터 가드는 로컬 토큰 존재 여부만 확인)
- Recipes, AddIngredient 뷰 구현 (Phase 3, 4)
- 상단 헤더(AppHeader) 구현 — 와이어프레임 기준 BottomNav 사용
