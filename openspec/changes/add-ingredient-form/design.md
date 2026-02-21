## 컴포넌트 설계

### AddIngredientView.vue

**경로**: `frontend/src/views/AddIngredientView.vue`

**폼 필드**:
| 필드 | 타입 | 필수 | 기본값 |
|------|------|------|--------|
| name | text | ✅ | - |
| category | CategoryButton 선택 | ❌ | '기타' |
| quantity | number | ❌ | 0 |
| unit | select | ❌ | '개' |
| expiry_date | date | ✅ | - |

**카테고리**: 채소 / 과일 / 육류 / 해산물 / 유제품 / 기타

**단위**: 개 / g / kg / ml / L / 팩

**제출 동작**:
1. `useIngredientStore.addIngredient(payload)` 호출
2. 성공 → `router.push('/')`
3. 실패 에러 처리:
   - `DUPLICATE_INGREDIENT` → "이미 등록된 재료예요."
   - `FRIDGE_LIMIT_EXCEEDED` → "재료는 최대 200개까지 등록할 수 있어요."
   - 기타 → "재료 추가에 실패했어요. 다시 시도해주세요."

**취소**: `router.back()`

---

### CategoryButton.vue

**경로**: `frontend/src/components/CategoryButton.vue`

**props**: `label: string`, `selected: boolean`

**emit**: `select(label)`

**스타일**:
- 선택: `bg-blue-600 border-blue-600 text-white`
- 미선택: `bg-white border-gray-300 text-gray-600`

---

## API 연동

`useIngredientStore.addIngredient(payload)` 사용:
```js
{
  name: string,       // trim 처리
  category: string,   // 미선택 시 '기타'
  quantity: number,   // 미입력 시 0
  unit: string,
  expiry_date: string | null,
}
```
