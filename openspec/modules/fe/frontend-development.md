# Frontend Development Guide

## 개요

냉털마스터 프론트엔드 개발 가이드.
기술 스택, 프로젝트 구조, 개발 워크플로우, 구현 가이드를 정리한다.

---

## 기술 스택

| 기술 | 버전 | 용도 |
|------|------|------|
| Vue | ^3.4.21 | 프론트엔드 프레임워크 |
| Vue Router | ^4.3.0 | 라우팅 |
| Pinia | ^2.1.7 | 상태 관리 |
| Vite | ^5.1.6 | 빌드 도구 |
| Tailwind CSS | ^3.4.1 | 스타일링 |
| Axios | ^1.6.7 | HTTP 클라이언트 |

---

## 프로젝트 구조

```
frontend/
├── src/
│   ├── api/
│   │   └── index.js              # Axios 인스턴스 및 인터셉터
│   ├── components/
│   │   ├── base/                # Base 컴포넌트 (Button, Input, Card)
│   │   ├── layout/              # 레이아웃 컴포넌트 (Header, BottomNav)
│   │   ├── fridge/              # 재료 관련 컴포넌트
│   │   ├── recipe/              # 레시피 관련 컴포넌트
│   │   └── vision/              # Vision 관련 컴포넌트
│   ├── stores/
│   │   ├── auth.js             # 인증 스토어
│   │   ├── fridge.js           # 재료 스토어
│   │   ├── recipe.js            # 레시피 스토어
│   │   ├── vision.js            # Vision 스토어
│   │   ├── cookbook.js          # 요리책 스토어
│   │   └── notification.js      # 알림 스토어
│   ├── views/
│   │   ├── HomeView.vue         # 홈 (랜딩)
│   │   ├── LoginView.vue        # 로그인
│   │   ├── SignupView.vue       # 회원가입
│   │   ├── DashboardView.vue    # 대시보드
│   │   ├── RecipesView.vue      # 레시피 목록
│   │   ├── RecipeDetailView.vue  # 레시피 상세
│   │   └── AddIngredientView.vue # 재료 추가
│   ├── router/
│   │   └── index.js             # 라우터 설정
│   ├── utils/
│   │   └── index.js             # 유틸리티 함수
│   ├── App.vue                  # 루트 컴포넌트
│   ├── main.js                  # 진입점
│   └── style.css                # 전역 스타일
├── public/                      # 정적 파일
├── index.html
├── package.json
├── vite.config.js
└── tailwind.config.js
```

---

## 개발 워크플로우

### 1. 환경 설정

```bash
cd frontend
npm install
npm run dev
```

개발 서버: `http://localhost:3000`

### 2. 컴포넌트 개발 순서

1. **Base 컴포넌트** (Button, Input, Card)
2. **레이아웃 컴포넌트** (Header, BottomNav)
3. **도메인 컴포넌트** (IngredientTag, RecipeCard)
4. **페이지 컴포넌트** (Dashboard, Recipes)

### 3. 상태 관리 패턴

```javascript
// stores/fridge.js
import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'

export const useFridgeStore = defineStore('fridge', () => {
  const ingredients = ref([])
  const loading = ref(false)
  const error = ref(null)

  const fetchIngredients = async (filter = {}) => {
    loading.value = true
    error.value = null
    try {
      const { data } = await api.get('/fridge', { params: filter })
      ingredients.value = data
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    ingredients,
    loading,
    error,
    fetchIngredients
  }
})
```

### 4. 컴포넌트 사용 패턴

```vue
<template>
  <div>
    <BaseButton 
      :loading="store.loading" 
      @click="handleClick"
    >
      추가하기
    </BaseButton>
  </div>
</template>

<script setup>
import { useFridgeStore } from '@/stores/fridge'

const store = useFridgeStore()

const handleClick = async () => {
  await store.fetchIngredients()
}
</script>
```

---

## 주요 구현 가이드

### 인증 플로우

1. **로그인/회원가입**
   - `useAuthStore.login()` 또는 `signup()` 호출
   - 토큰을 localStorage에 저장
   - 사용자 정보 조회 (`fetchUser()`)
   - 대시보드로 리다이렉트

2. **토큰 갱신**
   - API 인터셉터에서 자동 처리
   - 401 에러 시 refresh_token으로 갱신 시도
   - 실패 시 로그인 페이지로 리다이렉트

3. **로그아웃**
   - `useAuthStore.logout()` 호출
   - localStorage 토큰 삭제
   - 홈으로 리다이렉트

### 재료 관리 플로우

1. **재료 목록 조회**
   - `useFridgeStore.fetchIngredients(filter)` 호출
   - 필터: category, expiry_status, search

2. **재료 추가**
   - 수동 입력: `useFridgeStore.addIngredient(data)`
   - Vision: `useVisionStore.createJob()` → `confirmJob()`

3. **재료 수정/삭제**
   - `useFridgeStore.updateIngredient(id, data)`
   - `useFridgeStore.removeIngredient(id)`

### 레시피 추천 플로우

1. **추천 조회**
   - `useRecipeStore.fetchRecommendations(params)` 호출
   - 파라미터: tier, category, max_missing, limit

2. **오늘의 추천**
   - `useRecipeStore.fetchTodayRecommendations()` 호출
   - URGENT 재료 기반 Tier 1 레시피

3. **레시피 상세**
   - `useRecipeStore.getRecipeById(id)` 호출
   - 보유/부족 재료 표시

### Vision 플로우

1. **이미지 업로드**
   ```javascript
   const formData = new FormData()
   formData.append('type', 'RECEIPT')
   formData.append('image', file)
   
   const job = await useVisionStore.createJob('RECEIPT', file)
   ```

2. **결과 폴링**
   ```javascript
   const pollJob = async (jobId) => {
     const job = await useVisionStore.getJob(jobId)
     if (job.status === 'DONE') {
       // 결과 표시
     } else if (job.status === 'PROCESSING') {
       setTimeout(() => pollJob(jobId), 2000)
     }
   }
   ```

3. **확인 및 등록**
   ```javascript
   await useVisionStore.confirmJob(jobId, selectedItems)
   // 자동으로 fridge에 등록됨
   ```

### 알림 플로우

1. **알림 목록**
   - `useNotificationStore.fetchNotifications(limit)` 호출
   - 읽지 않은 알림 강조

2. **읽음 처리**
   - 단건: `useNotificationStore.markAsRead(id)`
   - 전체: `useNotificationStore.markAllAsRead()`

---

## 스타일링 가이드

### Tailwind CSS 사용

- 디자인 시스템 스펙 참조: `openspec/specs/design/spec.md`
- Primary 색상: `bg-primary-600`, `text-primary-600`
- 반응형: `md:`, `lg:` prefix 사용

### 컴포넌트 스타일

```vue
<template>
  <div class="bg-white rounded-xl shadow-lg p-4">
    <h3 class="text-xl font-bold text-gray-900 mb-2">
      {{ title }}
    </h3>
    <p class="text-sm text-gray-600">
      {{ description }}
    </p>
  </div>
</template>
```

---

## 에러 처리

### 공통 에러 컴포넌트

```vue
<template>
  <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
    <p class="text-red-700 text-sm">{{ error }}</p>
  </div>
</template>
```

### 에러 처리 패턴

```javascript
try {
  await store.fetchData()
} catch (error) {
  if (error.response?.status === 409) {
    // 중복 에러 - 사용자 확인
    const confirmed = await showConfirmDialog('이미 존재합니다. 계속하시겠습니까?')
    if (confirmed) {
      // 병합 로직
    }
  } else {
    showErrorMessage(error.message)
  }
}
```

---

## 테스트 전략

### 단위 테스트

- 컴포넌트: Vue Test Utils
- Store: Pinia 테스트 유틸리티
- 유틸리티: Jest

### 통합 테스트

- API 통합: MSW (Mock Service Worker)
- E2E: Playwright (선택적)

---

## 성능 최적화

### 코드 스플리팅

- 라우트별 lazy loading
- 컴포넌트 lazy loading

### 이미지 최적화

- WebP 포맷 사용
- Lazy loading

### 상태 관리 최적화

- 필요한 데이터만 fetch
- 캐싱 전략 적용

---

## 접근성

- ARIA 속성 사용
- 키보드 네비게이션 지원
- 색상 대비 WCAG AA 준수
- 최소 터치 영역 44x44px

---

## 참고 문서

- [디자인 시스템 스펙](../../specs/design/spec.md)
- [화면 스펙](../../specs/screens/spec.md)
- [컴포넌트 스펙](../../specs/components/spec.md)
- [API 통합 스펙](./frontend-api.md)
- [프론트엔드 초기 설정](./frontend-init.md)
- [백엔드 API 스펙](../../specs/auth/spec.md)

---

## 다음 단계

1. Base 컴포넌트 구현
2. Pinia Store 구현
3. 페이지 컴포넌트 구현
4. API 통합
5. 테스트 작성
