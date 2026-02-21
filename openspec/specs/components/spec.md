# Components Spec

## 개요

냉털마스터 프론트엔드 컴포넌트 스펙.
재사용 가능한 UI 컴포넌트의 인터페이스, 스타일, 동작을 정의한다.

---

## UI 프리미티브 (Base Components)

### BaseButton

**목적**: 프로젝트 전반에 쓰이는 버튼 표준화

#### Props

| prop | type | 기본값 | 설명 |
|------|------|--------|------|
| `variant` | `'primary' \| 'secondary' \| 'ghost' \| 'danger'` | `'primary'` | 버튼 스타일 |
| `size` | `'sm' \| 'md' \| 'lg'` | `'md'` | 크기 |
| `loading` | `boolean` | `false` | 로딩 스피너 표시 |
| `disabled` | `boolean` | `false` | 비활성화 |
| `fullWidth` | `boolean` | `false` | 전체 너비 |

#### 스타일

- Primary: `bg-primary-600 text-white hover:bg-primary-700`
- Secondary: `bg-gray-200 text-gray-800 hover:bg-gray-300`
- Ghost: `bg-transparent text-primary-600 hover:bg-primary-50`
- Danger: `bg-error text-white hover:bg-red-600`
- Size sm: `px-3 py-1.5 text-sm`
- Size md: `px-4 py-2 text-base`
- Size lg: `px-6 py-3 text-lg`

#### 사용 예시

```vue
<BaseButton variant="primary" :loading="isSubmitting">
  로그인
</BaseButton>
```

---

### BaseInput

**목적**: 입력 필드 공통화 (label, error 포함)

#### Props

| prop | type | 기본값 | 설명 |
|------|------|--------|------|
| `modelValue` | `string` | - | v-model 바인딩 |
| `label` | `string` | - | 필드 레이블 |
| `error` | `string` | - | 에러 메시지 |
| `placeholder` | `string` | - | 플레이스홀더 |
| `type` | `string` | `'text'` | input type |
| `disabled` | `boolean` | `false` | 비활성화 |
| `required` | `boolean` | `false` | 필수 필드 |

#### 스타일

- 기본: `border-gray-300 focus:border-primary-600 focus:ring-primary-600`
- 에러: `border-error focus:border-error focus:ring-error`

#### 사용 예시

```vue
<BaseInput
  v-model="email"
  label="이메일"
  type="email"
  placeholder="your@email.com"
  :error="errors.email"
/>
```

---

### BaseCard

**목적**: 카드 레이아웃 래퍼

#### Props

| prop | type | 기본값 | 설명 |
|------|------|--------|------|
| `hover` | `boolean` | `false` | 호버 효과 |
| `padding` | `'sm' \| 'md' \| 'lg'` | `'md'` | 패딩 크기 |

#### Slots

- `default`: 카드 본문
- `header`: 카드 헤더 (선택)
- `footer`: 카드 푸터 (선택)

#### 스타일

- 기본: `bg-white rounded-xl shadow-lg`
- 패딩 md: `p-4`
- 호버: `hover:shadow-xl transition-shadow`

#### 사용 예시

```vue
<BaseCard hover>
  <template #header>
    <h3>레시피 제목</h3>
  </template>
  <p>레시피 설명</p>
</BaseCard>
```

---

## 레이아웃 컴포넌트

### AppLayout

**목적**: 인증 후 공통 레이아웃 (헤더 + 콘텐츠 영역)

#### 구조

```vue
<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader />
    <main class="container mx-auto px-4 py-6 max-w-7xl">
      <slot />
    </main>
  </div>
</template>
```

#### 기능

- 라우터 가드와 연동: 비인증 접근 시 `/login` 리다이렉트
- `AppHeader` 포함
- 콘텐츠 영역: 최대 너비 1200px, 중앙 정렬

---

### AppHeader

**목적**: 상단 네비게이션 바

#### 구조

```
┌─────────────────────────────────────────┐
│ [로고]  [대시보드] [레시피] [요리책]  [👤 닉네임 ▼] │
└─────────────────────────────────────────┘
```

#### Props

없음 (useAuthStore에서 user 정보 가져옴)

#### 기능

- 로고 클릭 → `/dashboard`
- 네비게이션 메뉴
- 사용자 드롭다운: 프로필, 로그아웃

#### 스타일

- 고정 상단: `fixed top-0 left-0 right-0 z-50`
- 배경: `bg-white shadow-sm`
- 높이: `h-16`

---

## 재료 관련 컴포넌트

### IngredientTag

**목적**: 입력된 재료 하나를 태그로 표시, 삭제 가능

#### Props

| prop | type | 설명 |
|------|------|------|
| `name` | `string` | 재료명 |
| `expiringSoon` | `boolean` | 유통기한 임박 강조 |
| `category` | `string` | 카테고리 (선택) |

#### Emits

- `remove`: 삭제 버튼 클릭 시

#### 스타일

- 기본: `bg-primary-100 text-primary-700`
- 임박: `bg-warning-100 text-warning-700 border-2 border-warning-500`

#### 사용 예시

```vue
<IngredientTag
  name="당근"
  :expiringSoon="true"
  @remove="handleRemove"
/>
```

---

### IngredientInput

**목적**: 재료 입력창 + 현재 보유 재료 태그 목록 통합 관리

#### Props

| prop | type | 기본값 | 설명 |
|------|------|--------|------|
| `modelValue` | `string[]` | `[]` | 보유 재료 목록 (v-model) |

#### Emits

- `update:modelValue`: 재료 목록 변경 시

#### 기능

- 텍스트 입력 후 Enter 또는 추가 버튼으로 재료 등록
- `IngredientTag` 목록 렌더링
- 재료 삭제
- 중복 재료 방지

#### 사용 예시

```vue
<IngredientInput v-model="ingredients" />
```

---

## 레시피 컴포넌트

### RecipeCard

**목적**: 추천 레시피 한 건 표시

#### Props

| prop | type | 설명 |
|------|------|------|
| `recipe` | `Recipe` | 레시피 데이터 |
| `bookmarked` | `boolean` | 북마크 여부 |

#### Recipe 타입

```typescript
interface Recipe {
  id: string
  title: string
  thumbnail: string       // 썸네일 이미지 URL
  tier: 1 | 2 | 3        // Tier 레벨
  matchScore: number      // 보유 재료 일치율 (0~100)
  matchedIngredients: string[]  // 일치한 재료
  missingIngredients: string[]  // 부족한 재료
  sourceType: 'INTERNAL' | 'EXTERNAL'
  sourceUrl: string       // 외부 링크
  category?: string       // 카테고리
  cookingTime?: number    // 조리 시간 (분)
}
```

#### Emits

- `bookmark`: 북마크 버튼 클릭 시
- `unbookmark`: 북마크 해제 시
- `click`: 카드 클릭 시 (상세 페이지 이동)

#### 레이아웃

```
┌─────────────────────────────┐
│  [썸네일 이미지]             │
│  [Tier 배지]                │
├─────────────────────────────┤
│  김치볶음밥                  │
│  재료 일치율: ████████░ 85% │
│  ✅ 보유: 밥, 계란          │
│  ❌ 부족: 김치               │
│  [유튜브 ↗]      [★ 저장]   │
└─────────────────────────────┘
```

#### 스타일

- 카드: `rounded-xl shadow-lg hover:shadow-xl`
- Tier 배지:
  - Tier 1: `bg-green-100 text-green-700`
  - Tier 2: `bg-yellow-100 text-yellow-700`
  - Tier 3: `bg-blue-100 text-blue-700`

---

### RecipeDetailHeader

**목적**: 레시피 상세 페이지 헤더 (썸네일, 제목)

#### Props

| prop | type | 설명 |
|------|------|------|
| `recipe` | `Recipe` | 레시피 데이터 |

#### 구조

- 썸네일 이미지 (전체 너비)
- 제목 (하단 오버레이)
- Tier 배지

---

### IngredientList

**목적**: 보유/부족 재료 목록 표시

#### Props

| prop | type | 설명 |
|------|------|------|
| `matched` | `string[]` | 보유 재료 목록 |
| `missing` | `string[]` | 부족한 재료 목록 |

#### 스타일

- 보유: `text-green-700` + 체크 아이콘
- 부족: `text-gray-600` + X 아이콘

---

## 필터 컴포넌트

### FilterBar

**목적**: 레시피 목록 필터링

#### Props

| prop | type | 기본값 | 설명 |
|------|------|--------|------|
| `selectedTier` | `1 \| 2 \| 3 \| 'all'` | `'all'` | 선택된 Tier |
| `selectedCategory` | `string \| null` | `null` | 선택된 카테고리 |

#### Emits

- `update:tier`: Tier 필터 변경
- `update:category`: 카테고리 필터 변경

#### 필터 옵션

- Tier: 전체, Tier 1, Tier 2, Tier 3
- 카테고리: 전체, 한식, 중식, 일식, 양식, 기타

---

## 폼 컴포넌트

### ManualInputForm

**목적**: 재료 수동 입력 폼

#### Props

없음 (내부 상태 관리)

#### Emits

- `submit`: 재료 추가 완료 시 (Ingredient 데이터)

#### 필드

- 재료명 (필수)
- 카테고리 (선택)
- 수량 (선택)
- 단위 (선택)
- 유통기한 (선택)

---

### SourceSelector

**목적**: 재료 입력 방법 선택

#### Props

| prop | type | 기본값 | 설명 |
|------|------|--------|------|
| `modelValue` | `'manual' \| 'receipt' \| 'vision'` | `'manual'` | 선택된 방법 |

#### Emits

- `update:modelValue`: 방법 변경 시

#### 옵션

- 수동 입력
- 영수증 촬영
- 냉장고 사진

---

## 로딩/에러 컴포넌트

### LoadingSpinner

**목적**: 로딩 상태 표시

#### Props

| prop | type | 기본값 | 설명 |
|------|------|--------|------|
| `size` | `'sm' \| 'md' \| 'lg'` | `'md'` | 크기 |

---

### ErrorMessage

**목적**: 에러 메시지 표시

#### Props

| prop | type | 설명 |
|------|------|------|
| `message` | `string` | 에러 메시지 |
| `retry` | `() => void` | 재시도 함수 (선택) |

---

## 제약 조건

- 모든 컴포넌트는 Vue 3 Composition API 사용
- Props는 TypeScript 타입 정의 (또는 JSDoc)
- 스타일은 Tailwind CSS 클래스 사용
- 접근성: ARIA 속성, 키보드 네비게이션 지원

---

## 의존성

- `design/spec.md`: 디자인 시스템
- `screens/spec.md`: 화면 스펙
- Pinia 스토어: `useAuthStore`, `useIngredientStore`, `useRecipeStore`
