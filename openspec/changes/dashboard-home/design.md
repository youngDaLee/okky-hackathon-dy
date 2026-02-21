## 컴포넌트 설계

### DashboardView.vue

**경로**: `frontend/src/views/DashboardView.vue`

**구조**:
```
<AlertBar :count="store.urgentCount" />
<RecipeCarousel :recipes="[]" />
<MiniInventory :items="store.ingredients.slice(0, 5)" :loading="store.loading" />
<Cookbook :bookmarks="[]" />
```

**데이터**: `onMounted`에서 `fetchIngredients()` + `fetchSummary()` 병렬 호출

---

### AlertBar.vue

**경로**: `frontend/src/components/AlertBar.vue`

**props**: `count: number`
- `count === 0`이면 숨김 (`v-if`)
- 붉은 좌측 보더 스타일 (`border-l-4 border-red-500 bg-red-50`)

---

### MiniInventory.vue

**경로**: `frontend/src/components/MiniInventory.vue`

**props**: `items: Ingredient[]`, `loading: boolean`

**유통기한 배지**:
- `URGENT` → 빨간 (`bg-red-100 text-red-700`)
- `SOON` → 노란 (`bg-yellow-100 text-yellow-700`)
- `NORMAL` → 초록 (`bg-green-100 text-green-700`)
- `NO_EXPIRY` → 회색 텍스트

**D-Day 표시**: expiry_date와 오늘 날짜 차이로 계산

---

### RecipeCarousel.vue

**경로**: `frontend/src/components/RecipeCarousel.vue`

**props**: `recipes: Recipe[]`

**내부 상태**: `currentIndex: number`
- `prev()` / `next()` 메서드로 인덱스 이동
- 빈 상태 메시지: "재료를 추가하면 맞춤 레시피를 추천해드려요 🍳"

---

### Cookbook.vue

**경로**: `frontend/src/components/Cookbook.vue`

**props**: `bookmarks: SavedRecipe[]`
- 빈 상태 메시지: "저장된 레시피가 없어요."

---

### useIngredientStore (Pinia)

**경로**: `frontend/src/stores/ingredient.js`

**state**: `ingredients`, `summary`, `loading`, `error`

**getters**: `expiringSoon` (URGENT/SOON), `urgentCount`

**actions**: `fetchIngredients`, `fetchSummary`, `addIngredient`, `removeIngredient`

---

### fridgeApi

**경로**: `frontend/src/api/fridge.js`

**엔드포인트**: `getList`, `getSummary`, `getById`, `create`, `update`, `deleteOne`, `deleteBulk`
