## 1. 컴포넌트

- [x] 1.1 `CategoryButton.vue` 생성
  - props: label, selected
  - emit: select
  - 선택/미선택 스타일 토글

## 2. AddIngredientView 구현

- [x] 2.1 재료 이름 입력 필드 (필수)
- [x] 2.2 카테고리 선택 (CategoryButton 6개)
- [x] 2.3 수량 + 단위 선택 (개/g/kg/ml/L/팩)
- [x] 2.4 유통기한 date 입력 (필수, min=오늘)
- [x] 2.5 폼 제출: POST /fridge 연동 (`useIngredientStore.addIngredient`)
- [x] 2.6 에러 처리: DUPLICATE_INGREDIENT(409), FRIDGE_LIMIT_EXCEEDED(400)
- [x] 2.7 성공 시 `/` 리다이렉트, 취소 시 `router.back()`
- [x] 2.8 사진 추가 버튼 UI (disabled, 기능 미구현)

## 3. 검증

- [x] 3.1 `npm run build` 빌드 성공 확인
- [ ] 3.2 재료 이름 + 유통기한 입력 후 추가 → Dashboard 이동 확인
- [ ] 3.3 중복 재료 등록 시 에러 메시지 표시 확인
- [ ] 3.4 카테고리 선택/해제 토글 동작 확인
