# fe-fridge-view Tasks

> 냉장고 상세 현황 페이지 구현 (`/fridge`)

## 1. 라우터

- [x] 1.1 `router/index.js` — `/fridge` 라우트 추가 (`FridgeView.vue` lazy import)

## 2. FridgeView 구현

- [x] 2.1 요약 카드 (전체 / 🔴임박 / 🟡주의)
- [x] 2.2 검색 입력 (재료명 필터)
- [x] 2.3 카테고리 필터 버튼 (전체/채소/과일/육류/해산물/유제품/기타)
- [x] 2.4 유통기한 필터 버튼 (전체/🔴임박/🟡주의/🟢여유)
- [x] 2.5 재료 목록: 이모지, 이름, 수량+단위, D-Day 배지
- [x] 2.6 삭제 버튼 — confirm 후 `store.removeIngredient(id)` 호출
- [x] 2.7 빈 상태 메시지 + 재료 추가 링크
- [x] 2.8 `onMounted` — `fetchIngredients()` + `fetchSummary()` 병렬 호출

## 3. 카테고리/상태 헬퍼

- [x] 3.1 영문 Enum 기반 이모지 맵 (`VEGETABLE` → 🥕 등)
- [x] 3.2 영문 Enum 기반 라벨 맵 (`VEGETABLE` → '채소' 등)
- [x] 3.3 `expiryBadgeClass` — URGENT/SOON/NORMAL 색상
- [x] 3.4 `expiryLabel` — D-Day 계산 (만료됨 / D-Day / D-N)
