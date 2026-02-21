# fe-recipes-view Tasks

> 레시피 목록/상세 화면 구현 (더미 데이터, 백엔드 API 미연동)

## 1. 공유 데이터 모듈

- [x] 1.1 `frontend/src/data/recipes.js` 생성 — 더미 레시피 8종
- [x] 1.2 `calcMatchRate(recipe, fridgeNames)` 헬퍼 — 냉장고 재료명 fuzzy 매칭율 계산

## 2. RecipesView 구현 (`/recipes`)

- [x] 2.1 헤더 — "🍲 레시피 찾기"
- [x] 2.2 검색 입력 (레시피명 / 재료명)
- [x] 2.3 필터 탭 (전체 / 100% 매칭 / 75% 이상)
- [x] 2.4 냉장고 재료 수 요약 배너
- [x] 2.5 레시피 카드: 제목, 매칭율 뱃지, 설명, 재료 태그 (보유 사항 초록/미보유 회색), 요리정보
- [x] 2.6 하단 매칭율 프로그레스 바
- [x] 2.7 카드 클릭 → `/recipes/:id` 이동
- [x] 2.8 `onMounted` — 냉장고 재료 없으면 `fetchIngredients()` 호출

## 3. RecipeDetailView 구현 (`/recipes/:id`)

- [x] 3.1 뒤로 가기 버튼
- [x] 3.2 헤더 카드 (제목, 매칭율 뱃지, 설명, 요리정보)
- [x] 3.3 재료 섹션 — 냉장고 보유 여부 체크 아이콘 + "냉장고 있음" 뱃지
- [x] 3.4 매칭율 프로그레스 바 (N/전체개 표시)
- [x] 3.5 조리 순서 — 번호 원형 + 단계 텍스트
- [x] 3.6 미구현 기능 안내 문구 (북마크/평점)
- [x] 3.7 레시피 없을 때 빈 상태 처리

## 4. DashboardView 연결

- [x] 4.1 `RecipeCarousel` — 빈 배열 → 100% 매칭 레시피 자동 전달 (최대 5개)

## 5. 비고

- [ ] 5.1 백엔드 레시피 API 구현 시 더미 데이터 → 실제 API로 교체 필요
  - recommendation 패키지 handler/service 미구현 상태 (라우터 미등록)
