## 1. API 클라이언트

- [x] 1.1 `frontend/src/api/fridge.js` 생성 — 7개 엔드포인트 (getList, getSummary, getById, create, update, deleteOne, deleteBulk)

## 2. Pinia 스토어

- [x] 2.1 `frontend/src/stores/ingredient.js` 생성
  - state: ingredients, summary, loading, error
  - getters: expiringSoon, urgentCount
  - actions: fetchIngredients, fetchSummary, addIngredient, removeIngredient

## 3. 컴포넌트

- [x] 3.1 `AlertBar.vue` 생성 — count > 0 시 유통기한 임박 경고 표시
- [x] 3.2 `MiniInventory.vue` 생성 — 재료 5개, D-Day 배지, 카테고리 이모지
- [x] 3.3 `RecipeCarousel.vue` 생성 — 캐러셀 UI (빈 상태 메시지 포함)
- [x] 3.4 `Cookbook.vue` 생성 — 북마크 목록 (빈 상태 메시지 포함)

## 4. DashboardView 조립

- [x] 4.1 `DashboardView.vue` 수정 — 4개 컴포넌트 조립
- [x] 4.2 onMounted에서 fetchIngredients + fetchSummary 병렬 호출

## 5. 검증

- [x] 5.1 `npm run build` 빌드 성공 확인
- [ ] 5.2 백엔드 연동 후 AlertBar urgentCount 표시 확인
- [ ] 5.3 백엔드 연동 후 MiniInventory 재료 목록 표시 확인
- [ ] 5.4 RecipeCarousel / Cookbook 빈 상태 메시지 표시 확인 (API 연동 전)
