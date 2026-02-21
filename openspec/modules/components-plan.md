# 컴포넌트 구현 계획 (와이어프레임 기반)

> **와이어프레임 출처**: `/Users/munggae/workspace/Okkyhackathonwireframe`
> Figma의 **Make** 기능(Figma-to-Code)으로 자동 생성된 React + TypeScript 결과물입니다.
> 해당 코드를 그대로 사용하지 않고, 구조와 UI를 참고하여 Vue 3으로 재구현하는 계획입니다.

## 1. 페이지 구성 (라우터)

### 와이어프레임 라우트 → Vue 라우터 대응

| 경로 | 와이어프레임 | Vue 뷰 | 인증 필요 |
|------|------------|--------|----------|
| `/` | Dashboard | DashboardView | ✅ |
| `/recipes` | Recipes | RecipesView | ✅ |
| `/recipes/:id` | RecipeDetail | RecipeDetailView | ✅ |
| `/add-ingredient` | AddIngredient | AddIngredientView | ✅ |
| `/login` | (없음) | LoginView | ❌ |
| `/signup` | (없음) | SignupView | ❌ |

### 라우터 구조

```js
// 인증 전: 기존 라우트 유지
{ path: '/login', component: LoginView }
{ path: '/signup', component: SignupView }

// 인증 후: AppLayout 아래
{
  path: '/',
  component: AppLayout,       // BottomNav + RouterView
  meta: { requiresAuth: true },
  children: [
    { path: '',        name: 'dashboard',      component: DashboardView },
    { path: 'recipes', name: 'recipes',        component: RecipesView },
    { path: 'recipes/:id', name: 'recipe-detail', component: RecipeDetailView },
    { path: 'add-ingredient', name: 'add-ingredient', component: AddIngredientView },
  ]
}
```

---

## 2. 전체 컴포넌트 트리

```
App.vue
├── [인증 전]
│   ├── HomeView (랜딩 페이지 - 기존)
│   ├── LoginView (기존, API 연동 필요)
│   └── SignupView (기존, API 연동 필요)
│
└── AppLayout.vue              ← Root.tsx 역할
    ├── RouterView
    │   ├── DashboardView      ← Dashboard.tsx
    │   │   ├── AlertBar
    │   │   ├── RecipeCarousel
    │   │   ├── MiniInventory
    │   │   └── Cookbook
    │   │
    │   ├── RecipesView        ← Recipes.tsx
    │   │   └── RecipeListItem (n개)
    │   │
    │   ├── RecipeDetailView   ← RecipeDetail.tsx
    │   │   ├── IngredientCheckItem (n개)
    │   │   └── CookingStep (n개)
    │   │
    │   └── AddIngredientView  ← AddIngredient.tsx
    │       └── CategoryButton (n개)
    │
    ├── AppBottomNav           ← BottomNav.tsx
    └── FloatingAddButton      ← FloatingAddButton.tsx
```

---

## 3. 데이터 모델

### Ingredient (재료)

```js
{
  id: string,
  name: string,
  quantity: number,
  unit: '개' | 'g' | 'kg' | 'ml' | 'L' | '팩',
  expiryDate: string,       // ISO date string
  category: '채소' | '과일' | '육류' | '해산물' | '유제품' | '기타',
  daysLeft: number          // computed: 만료까지 남은 일수
}
```

### Recipe (레시피 목록용)

```js
{
  id: string,
  title: string,
  description: string,
  ingredients: string[],
  cookTime: string,         // '30분'
  servings: string,         // '2인분'
  matchPercentage: number,  // 0~100
  missingIngredients: string[]
}
```

### RecipeDetail (레시피 상세용)

```js
{
  id: string,
  title: string,
  description: string,
  ingredients: Array<{
    name: string,
    amount: string,
    hasIt: boolean          // 보유 여부
  }>,
  steps: string[],
  cookTime: string,
  servings: string,
  matchPercentage: number,
  difficulty: string,
  youtubeUrl?: string,
  blogUrl?: string
}
```

### SavedRecipe (요리책 북마크)

```js
{
  id: string,
  title: string,
  source: string,
  type: 'youtube' | 'blog',
  url: string,
  savedDate: string
}
```

---

## 4. 컴포넌트 상세 명세

### 레이아웃

#### `AppLayout.vue`
**역할**: 인증 후 공통 레이아웃. `BottomNav` + `RouterView` + `FloatingAddButton`을 포함.
**라우터 가드**: `meta.requiresAuth` 감지 → 미인증 시 `/login` 리다이렉트

```
┌─────────────────────────┐
│      RouterView          │ ← 페이지 콘텐츠
│                          │
│                          │
│                          │
├─────────────────────────┤
│   🏠   ➕   🍳          │ ← AppBottomNav (fixed bottom)
└─────────────────────────┘
         [➕]              ← FloatingAddButton (fixed, add-ingredient 페이지에서 숨김)
```

---

### 공통 컴포넌트

#### `AppBottomNav.vue`
**역할**: 하단 네비게이션 바 (fixed)
**참고**: `BottomNav.tsx`

**탭 구성:**
```
[홈 🏠]  [추가 ➕]  [레시피 🍳]
  /        /add-ingredient   /recipes
```

- 가운데 추가 버튼: 파란 원형 FAB 스타일, 위로 돌출
- 현재 경로 활성 표시

**props**: 없음 (내부에서 `useRoute()` 사용)

---

#### `FloatingAddButton.vue`
**역할**: 재료 추가 플로팅 버튼 (Dashboard에서만 표시)
**참고**: `FloatingAddButton.tsx`

- 오른쪽 하단 고정 (`fixed bottom-24 right-4`)
- `/add-ingredient`로 라우팅
- `+` 아이콘 + "재료 추가" 텍스트

---

#### `AlertBar.vue`
**역할**: 유통기한 임박 재료 경고 알림
**참고**: `AlertBar.tsx`

```
⚠️ 유통기한 임박 재료가 [N]개 있어요! 빨리 요리하세요.
```

**props:**
| prop | type | 설명 |
|------|------|------|
| `count` | `number` | 임박 재료 수 |

- `count === 0`이면 숨김
- 붉은 좌측 보더 스타일

---

### Dashboard 컴포넌트

#### `RecipeCarousel.vue`
**역할**: 100% 매칭 레시피 카드 캐러셀 슬라이더
**참고**: `RecipeCarousel.tsx`

```
┌──────────────────────────────┐
│  100% 매칭 레시피             │
│  추가 구매 없이 바로 만들 수 있어요│
│ ┌───────────────────────────┐ │
│ │ [100% 매칭 배지]           │ │
│ │ 감자짜글이                 │ │
│ │ 지금 가진 감자와 양파로...  │ │
│ │ [감자] [양파]              │ │
│ │ ⏱ 30분  👥 2인분          │ │
│ └───────────────────────────┘ │
│  ← [●○○] →                   │
└──────────────────────────────┘
```

**내부 상태:** `currentIndex: number`

**props:** 없음 (스토어에서 100% 매칭 레시피 가져옴)

**emit:** 없음 (내부에서 라우터 push)

---

#### `MiniInventory.vue`
**역할**: 냉장고 현황 미니 뷰 (재료 5개 + 유통기한)
**참고**: `MiniInventory.tsx`

```
┌──────────────────────────────┐
│  나의 냉장고 현황   [더보기 >]│
│  🥕 감자           [7일]     │
│  🍎 양파           [5일]     │
│  🥛 우유           [2일] 🔴  │
│  🥚 계란           [10일]    │
│  🐟 고등어         [1일] 🔴  │
└──────────────────────────────┘
```

**유통기한 색상 규칙:**
- `daysLeft <= 2`: 빨간 배지
- `daysLeft <= 5`: 노란 배지
- `daysLeft > 5`: 초록 배지

**props:** 없음 (스토어에서 재료 가져옴)

---

#### `Cookbook.vue`
**역할**: 저장된 요리책 (북마크된 외부 레시피)
**참고**: `Cookbook.tsx`

```
┌──────────────────────────────┐
│  📖 저장된 요리책             │
│  ★ 10분 완성 간단 감자전     │
│    [YouTube] 백종원의 요리비책↗│
│  ★ 집에서 만드는 완벽한 오므라이스│
│    [YouTube] 요리하는 유튜버  ↗│
└──────────────────────────────┘
```

**props:** 없음 (스토어에서 북마크 가져옴)

---

### RecipesView 컴포넌트

#### `RecipeListItem.vue`
**역할**: 레시피 목록의 각 카드 아이템
**참고**: `Recipes.tsx` 내부 렌더링

```
┌──────────────────────────────┐
│  감자짜글이  [100%]          │
│  지금 가진 감자와 양파로...   │
│  ⏱ 30분  👥 2인분           │
│  [감자✅] [양파✅] [고추가루❌]│
│                              >│
└──────────────────────────────┘
```

**props:**
| prop | type | 설명 |
|------|------|------|
| `recipe` | `Recipe` | 레시피 데이터 |

**emit:** 없음 (클릭 시 라우터 push)

**재료 태그 색상:**
- 보유 재료: 초록 (`bg-green-50 text-green-700`)
- 부족 재료: 빨간 (`bg-red-50 text-red-700`)

---

### RecipeDetailView 컴포넌트

#### `IngredientCheckItem.vue`
**역할**: 레시피 상세의 재료 한 줄 (보유/미보유 표시)
**참고**: `RecipeDetail.tsx` 재료 섹션

```
✅ 감자    2개       (보유 - 초록 배경)
❌ 고추가루  2큰술    (미보유 - 빨간 배경)
```

**props:**
| prop | type | 설명 |
|------|------|------|
| `name` | `string` | 재료명 |
| `amount` | `string` | 수량 |
| `hasIt` | `boolean` | 보유 여부 |

---

#### `CookingStep.vue`
**역할**: 조리 순서 한 단계
**참고**: `RecipeDetail.tsx` 조리방법 섹션

```
[1] 감자는 껍질을 벗기고 한입 크기로 자릅니다.
[2] 양파는 채썰어 준비합니다.
```

**props:**
| prop | type | 설명 |
|------|------|------|
| `step` | `number` | 순서 번호 |
| `text` | `string` | 조리 내용 |

---

### AddIngredientView 컴포넌트

#### `CategoryButton.vue`
**역할**: 카테고리 선택 버튼 (채소/과일/육류 등)
**참고**: `AddIngredient.tsx` 카테고리 섹션

**props:**
| prop | type | 설명 |
|------|------|------|
| `label` | `string` | 카테고리명 |
| `selected` | `boolean` | 선택 여부 |

**emit:** `select`

---

## 5. Pinia 스토어 설계

### `useIngredientStore.js`

```js
state: {
  ingredients: Ingredient[]
}

getters: {
  expiringSoon: // daysLeft <= 2 인 재료
  sortedByExpiry: // daysLeft 오름차순 정렬
}

actions: {
  fetchIngredients(),   // GET /api/ingredients
  addIngredient(form),  // POST /api/ingredients
  removeIngredient(id)  // DELETE /api/ingredients/:id
}
```

---

### `useRecipeStore.js`

```js
state: {
  recommendations: Recipe[],    // 매칭 레시피 목록
  currentDetail: RecipeDetail | null,
  bookmarks: SavedRecipe[]
}

getters: {
  perfectMatches: // matchPercentage === 100
}

actions: {
  fetchRecommendations(),       // GET /api/recipes/recommendations
  fetchRecipeDetail(id),        // GET /api/recipes/:id
  fetchBookmarks(),             // GET /api/bookmarks
  addBookmark(recipeId),        // POST /api/bookmarks
  removeBookmark(id)            // DELETE /api/bookmarks/:id
}
```

---

## 6. 페이지별 UI 스케치

### Dashboard (`/`)

```
┌─────────────────────────────────┐
│ 🍳 나의 냉장고                   │
│ 똑똑한 재료 관리와 맞춤 레시피 추천│
├─────────────────────────────────┤
│ ⚠️ 유통기한 임박 재료 3개        │  ← AlertBar
├─────────────────────────────────┤
│ 100% 매칭 레시피 [캐러셀]        │  ← RecipeCarousel
│ [← 감자짜글이 →]                │
├─────────────────────────────────┤
│ 나의 냉장고 현황                 │  ← MiniInventory
│ 🥕 감자 [7일]                   │
│ 🥛 우유 [2일]🔴                 │
├─────────────────────────────────┤
│ 📖 저장된 요리책                 │  ← Cookbook
│ ★ 간단 감자전 [YouTube]↗        │
└─────────────────────────────────┘
│   🏠        ➕        🍳        │  ← AppBottomNav
└─────────────────────────────────┘
                    [➕재료 추가]  ← FloatingAddButton
```

---

### Recipes (`/recipes`)

```
┌─────────────────────────────────┐
│ 🍲 레시피 찾기                   │
│ 가진 재료로 만들 수 있는 요리     │
│ [🔍 레시피 검색...]              │
│ [전체] [100% 매칭] [75% 이상]    │
├─────────────────────────────────┤
│ 감자짜글이 [100%]    →          │
│ 지금 가진 감자와 양파로...        │
│ ⏱30분 👥2인분                   │
│ [감자✅][양파✅]                  │
├─────────────────────────────────┤
│ 김치볶음밥 [75%]     →          │
│ 필요: 밥, 김치                   │
└─────────────────────────────────┘
```

---

### RecipeDetail (`/recipes/:id`)

```
┌─────────────────────────────────┐
│ [←] 감자짜글이          [🔖]   │  ← sticky header
├─────────────────────────────────┤
│ (orange gradient)               │
│ 감자짜글이              [100%] │
│ 지금 가진 감자와 양파로...        │
│ ⏱30분  👥2인분  👨‍🍳쉬움        │
├─────────────────────────────────┤
│ 재료                  3/7 보유  │
│ ✅ 감자   2개                   │
│ ❌ 고추가루  2큰술               │
│ [부족한 재료 장바구니에 추가]     │
├─────────────────────────────────┤
│ 조리 방법                        │
│ [1] 감자는 껍질을 벗기고...      │
│ [2] 양파는 채썰어...             │
├─────────────────────────────────┤
│ 참고 링크                        │
│ YouTube에서 보기         ↗      │
└─────────────────────────────────┘
```

---

### AddIngredient (`/add-ingredient`)

```
┌─────────────────────────────────┐
│ 🥕 재료 추가               [X] │
│ 냉장고에 있는 재료를 등록하세요  │
│ [📷 영수증/재료 사진으로 빠른 추가]│
├─────────────────────────────────┤
│ 재료 이름 *                      │
│ [예: 감자, 양파, 우유           ]│
├─────────────────────────────────┤
│ 카테고리                         │
│ [채소] [과일] [육류]            │
│ [해산물] [유제품] [기타]        │
├─────────────────────────────────┤
│ 수량                             │
│ [0        ][단위 ▼]            │
├─────────────────────────────────┤
│ 유통기한 *                       │
│ [📅 날짜 선택               ]  │
├─────────────────────────────────┤
│ [취소]              [➕ 추가하기]│
└─────────────────────────────────┘
```

---

## 7. 구현 우선순위

```
Phase 1: 레이아웃 & 네비게이션
├── AppLayout.vue (인증 후 레이아웃 + 라우터 가드)
├── AppBottomNav.vue
└── 라우터 재구성

Phase 2: Dashboard 핵심
├── AlertBar.vue
├── RecipeCarousel.vue (더미 데이터 우선)
├── MiniInventory.vue (더미 데이터 우선)
├── Cookbook.vue (더미 데이터 우선)
├── FloatingAddButton.vue
└── DashboardView.vue (조립)

Phase 3: 재료 추가
├── CategoryButton.vue
├── AddIngredientView.vue
└── useIngredientStore.js

Phase 4: 레시피
├── RecipeListItem.vue
├── RecipesView.vue
├── IngredientCheckItem.vue
├── CookingStep.vue
├── RecipeDetailView.vue
└── useRecipeStore.js

Phase 5: API 연동
├── 스토어 더미 데이터 → 실제 API 교체
├── LoginView / SignupView API 연동
└── 북마크 저장/삭제 연동
```

---

## 8. 주요 결정 사항

### 와이어프레임과의 차이점
| 항목 | 와이어프레임 (React) | Vue 구현 |
|------|-------------------|---------|
| 프레임워크 | React + TypeScript | Vue 3 + JavaScript |
| 상태관리 | useState (로컬) | Pinia (전역) |
| 라우팅 | React Router | Vue Router |
| UI 라이브러리 | shadcn/ui | Tailwind CSS (직접 구현) |
| 아이콘 | lucide-react | lucide-vue-next |

### 보류 기능 (MVP 외)
- 영수증/이미지 촬영 (카메라 버튼 UI만 노출, 기능 비활성)
- 부족한 재료 장바구니 추가
- 카테고리별 요리책 필터링

---

## 관련 모듈
- `frontend-init.md` — 기술 스택 및 프로젝트 초기 설정
