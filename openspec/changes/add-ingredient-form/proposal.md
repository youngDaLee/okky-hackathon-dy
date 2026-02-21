## Why

layout-bottom-nav에서 빈 스캐폴딩으로만 만들어둔 AddIngredientView에
실제 재료 추가 폼을 구현한다.
사용자가 냉장고 재료를 직접 입력하여 등록하는 핵심 기능이다.

## What Changes

- `frontend/src/views/AddIngredientView.vue` 수정 — 재료 추가 폼 완전 구현
- `frontend/src/components/CategoryButton.vue` 신규 생성 — 카테고리 선택 토글 버튼

## Capabilities

### New Capabilities

- `fe-add-ingredient`: 재료 추가 폼 (이름/카테고리/수량/유통기한)

### Modified Capabilities

- `fe-add-ingredient`: AddIngredientView 스캐폴딩 → 실제 폼 구현

## Impact

- `frontend/src/views/AddIngredientView.vue` 전면 수정
- `frontend/src/components/CategoryButton.vue` 신규

### Non-goals

- 사진/영수증 촬영 기능 (UI 버튼만 노출, disabled)
- 재료 수정(PATCH) 기능
- 재료 목록 조회 화면 (MiniInventory 담당)
