## 컴포넌트 설계

### AppLayout.vue

**경로**: `frontend/src/layouts/AppLayout.vue`

**역할**: 인증 후 모든 페이지의 공통 래퍼. RouterView + AppBottomNav + FloatingAddButton 포함.

**라우터 가드 로직**:
- `onMounted` 또는 `beforeEach` 네비게이션 가드에서 `localStorage.access_token` 존재 여부 확인
- 없으면 `/login`으로 리다이렉트

**템플릿 구조**:
```
<div id="app-layout" class="min-h-screen bg-gray-50 pb-20">
  <RouterView />
  <AppBottomNav />
  <FloatingAddButton v-if="!isAddIngredientPage" />
</div>
```

**계산 속성**:
- `isAddIngredientPage`: 현재 라우트가 `/add-ingredient`이면 FAB 숨김

---

### AppBottomNav.vue

**경로**: `frontend/src/components/AppBottomNav.vue`

**역할**: 하단 고정 네비게이션 바. 3개 탭 (홈, 추가, 레시피).

**탭 구성**:
```
[홈 🏠]   [➕]   [레시피 🍳]
   /     /add-ingredient  /recipes
```

**가운데 추가 버튼 스타일**:
- 파란 원형 (`bg-blue-600 rounded-full p-4`)
- 위로 돌출 (`-mt-8`)
- 현재 라우트와 무관하게 항상 강조

**활성 탭 표시**:
- `useRoute().path === item.path` 비교
- 활성: `text-blue-600`, 비활성: `text-gray-500`

**아이콘**: `lucide-vue-next` 패키지 사용 (`Home`, `Plus`, `ChefHat`)

---

### FloatingAddButton.vue

**경로**: `frontend/src/components/FloatingAddButton.vue`

**역할**: 재료 추가 진입 플로팅 버튼.

**위치**: `fixed bottom-24 right-4 z-30`

**클릭 동작**: `router.push('/add-ingredient')`

**표시 조건**: AppLayout에서 `/add-ingredient` 페이지가 아닐 때만 렌더

---

### DashboardView.vue

**경로**: `frontend/src/views/DashboardView.vue`

**역할**: Phase 2 구현 전 스캐폴딩. 헤더 텍스트만 포함.

```html
<template>
  <div class="max-w-4xl mx-auto px-4 py-4">
    <header>
      <h1>🍳 나의 냉장고</h1>
      <p>똑똑한 재료 관리와 맞춤 레시피 추천</p>
    </header>
  </div>
</template>
```

---

## 라우터 재구성

**변경 전**:
```js
routes: [
  { path: '/',        component: HomeView },
  { path: '/login',   component: LoginView },
  { path: '/signup',  component: SignupView },
]
```

**변경 후**:
```js
routes: [
  // 인증 전 (공개)
  { path: '/home',    name: 'home',    component: HomeView },
  { path: '/login',   name: 'login',   component: LoginView },
  { path: '/signup',  name: 'signup',  component: SignupView },

  // 인증 후 (AppLayout 하위)
  {
    path: '/',
    component: AppLayout,
    meta: { requiresAuth: true },
    children: [
      { path: '',                name: 'dashboard',      component: DashboardView },
      { path: 'recipes',         name: 'recipes',        component: () => import('@/views/RecipesView.vue') },
      { path: 'recipes/:id',     name: 'recipe-detail',  component: () => import('@/views/RecipeDetailView.vue') },
      { path: 'add-ingredient',  name: 'add-ingredient', component: () => import('@/views/AddIngredientView.vue') },
    ]
  },
  // 기본 리다이렉트
  { path: '/:pathMatch(.*)*', redirect: '/' },
]
```

**네비게이션 가드** (`router.beforeEach`):
```js
router.beforeEach((to) => {
  const token = localStorage.getItem('access_token')
  if (to.meta.requiresAuth && !token) {
    return { name: 'login' }
  }
  if (!to.meta.requiresAuth && token && (to.name === 'login' || to.name === 'signup')) {
    return { name: 'dashboard' }
  }
})
```

---

## 의존성

- `lucide-vue-next`: BottomNav 아이콘 (npm install 필요)
- `vue-router`: 이미 설치됨
- `pinia`: auth 스토어 연동 (이미 설치됨)

---

## 파일 구조 변경

```
frontend/src/
├── layouts/             ← 신규
│   └── AppLayout.vue
├── components/          ← 신규 파일 추가
│   ├── AppBottomNav.vue
│   └── FloatingAddButton.vue
├── views/
│   ├── HomeView.vue     ← 유지
│   ├── LoginView.vue    ← 유지
│   ├── SignupView.vue   ← 유지
│   └── DashboardView.vue ← 신규
└── router/
    └── index.js         ← 수정
```
