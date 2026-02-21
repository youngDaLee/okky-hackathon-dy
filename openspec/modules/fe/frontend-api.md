# Frontend API Integration Spec

## 개요

백엔드 API와의 통합 스펙.
각 도메인별 API 엔드포인트, 요청/응답 형식, 에러 처리, 상태 관리 방법을 정의한다.

---

## API 클라이언트 설정

### Base Configuration

```javascript
// src/api/index.js
import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json'
  }
})
```

### 인증 인터셉터

- **Request**: 모든 요청에 `Authorization: Bearer {access_token}` 자동 추가
- **Response**: 401 에러 시 자동 토큰 갱신 시도, 실패 시 로그인 페이지 리다이렉트

---

## Auth API

### POST /api/v1/auth/signup

**요청**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "nickname": "사용자"
}
```

**응답**:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "opaque_token_string",
  "expires_in": 3600
}
```

**에러**:
- `400`: 요청 형식 오류
- `409`: 이메일 중복

### POST /api/v1/auth/login

**요청**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**응답**: signup과 동일

**에러**:
- `401`: 잘못된 자격증명 (이메일/패스워드 구분 없이 동일 메시지)

### POST /api/v1/auth/refresh

**요청**:
```json
{
  "refresh_token": "opaque_token_string"
}
```

**응답**: signup과 동일 (새로운 토큰 쌍)

**에러**:
- `401`: 유효하지 않은 refresh_token

### POST /api/v1/auth/logout

**요청**:
```json
{
  "refresh_token": "opaque_token_string"
}
```

**응답**: `204 No Content`

### GET /api/v1/users/me

**요청**: 헤더에 `Authorization: Bearer {access_token}`

**응답**:
```json
{
  "id": "507f1f77bcf86cd799439011",
  "email": "user@example.com",
  "nickname": "사용자",
  "dietary_prefs": ["vegetarian"],
  "allergens": ["peanut"],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### PATCH /api/v1/users/me

**요청**:
```json
{
  "nickname": "새 닉네임",
  "dietary_prefs": ["vegan"],
  "allergens": ["gluten", "peanut"]
}
```

**응답**: GET /users/me와 동일

---

## Fridge API

### GET /api/v1/fridge

**쿼리 파라미터**:
- `category`: VEGETABLE | FRUIT | MEAT | SEAFOOD | DAIRY | GRAIN | CONDIMENT | FROZEN | OTHER
- `expiry_status`: URGENT | SOON | NORMAL | NO_EXPIRY
- `search`: 재료명 검색 (prefix)

**응답**:
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "name": "당근",
    "category": "VEGETABLE",
    "quantity": 3,
    "unit": "개",
    "expiry_date": "2024-01-10T00:00:00Z",
    "expiry_status": "URGENT",
    "source": "manual",
    "added_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### POST /api/v1/fridge

**요청**:
```json
{
  "name": "당근",
  "category": "VEGETABLE",
  "quantity": 3,
  "unit": "개",
  "expiry_date": "2024-01-10T00:00:00Z"
}
```

**응답**: GET /fridge 단건 응답과 동일

**에러**:
- `409`: 동일 이름 재료 이미 존재 (병합 여부는 클라이언트에서 결정)

### GET /api/v1/fridge/:id

**응답**: GET /fridge 단건 응답과 동일

### PATCH /api/v1/fridge/:id

**요청**: POST /fridge와 동일 (모든 필드 선택적)

**응답**: GET /fridge 단건 응답과 동일

### DELETE /api/v1/fridge/:id

**응답**: `204 No Content`

### DELETE /api/v1/fridge

**요청**:
```json
{
  "ids": ["507f1f77bcf86cd799439011", "507f1f77bcf86cd799439012"]
}
```

**응답**: `204 No Content`

### GET /api/v1/fridge/summary

**응답**:
```json
{
  "total_count": 15,
  "urgent_count": 3,
  "soon_count": 5,
  "normal_count": 7
}
```

---

## Recommendation API

### GET /api/v1/recommendations

**쿼리 파라미터**:
- `tier`: 1 | 2 | 3 | all (기본값: all)
- `category`: 한식 | 중식 | 일식 | 양식 | 간식 | 기타
- `max_missing`: 부족 재료 최대 허용 개수
- `limit`: 결과 최대 개수 (기본값: 20)

**응답**:
```json
[
  {
    "recipe": {
      "id": "507f1f77bcf86cd799439011",
      "title": "김치볶음밥",
      "description": "냉장고 속 김치로 만드는 간단한 한끼",
      "required_ingredients": ["밥", "김치", "계란"],
      "optional_ingredients": ["참기름"],
      "main_ingredient": "김치",
      "source_type": "INTERNAL",
      "source_url": "https://youtube.com/watch?v=...",
      "thumbnail_url": "https://example.com/image.jpg",
      "category": "한식",
      "tags": ["초간단", "10분완성"],
      "cooking_time_min": 10,
      "difficulty": "EASY"
    },
    "tier": 1,
    "match_score": 1.0,
    "matched_ingredients": ["밥", "김치", "계란"],
    "missing_ingredients": [],
    "urgency_bonus": true
  }
]
```

### GET /api/v1/recommendations/today

**응답**: GET /recommendations와 동일 (URGENT 재료 기반 Tier 1만)

### GET /api/v1/recipes

**쿼리 파라미터**:
- `keyword`: 검색 키워드
- `category`: 한식 | 중식 | 일식 | 양식 | 간식 | 기타
- `limit`: 결과 최대 개수

**응답**: Recipe 배열 (RecommendationResult가 아닌 Recipe만)

### GET /api/v1/recipes/:id

**응답**: Recipe 단건

---

## Vision API

### POST /api/v1/vision/jobs

**요청**: `multipart/form-data`
- `type`: RECEIPT | FRIDGE
- `image`: 이미지 파일 (최대 10MB, JPEG/PNG/WEBP)

**응답**:
```json
{
  "id": "507f1f77bcf86cd799439011",
  "type": "RECEIPT",
  "status": "PENDING",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### GET /api/v1/vision/jobs/:id

**응답**:
```json
{
  "id": "507f1f77bcf86cd799439011",
  "type": "RECEIPT",
  "status": "DONE",
  "extracted": [
    {
      "name": "계란10구",
      "normalized_name": "계란",
      "quantity": 10,
      "unit": "개",
      "confidence": 0.95,
      "selected": false
    }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "completed_at": "2024-01-01T00:00:05Z"
}
```

**Status 값**:
- `PENDING`: 대기 중
- `PROCESSING`: 처리 중
- `DONE`: 완료
- `FAILED`: 실패

### POST /api/v1/vision/jobs/:id/confirm

**요청**:
```json
{
  "selected_items": [
    {
      "name": "계란",
      "normalized_name": "계란",
      "quantity": 10,
      "unit": "개"
    }
  ]
}
```

**응답**: `204 No Content` (재료는 자동으로 fridge에 등록됨)

---

## Cookbook API

### GET /api/v1/cookbook

**쿼리 파라미터**:
- `label`: 라벨 필터
- `search`: 레시피 제목 검색

**응답**:
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "recipe_id": "507f1f77bcf86cd799439012",
    "recipe_snapshot": {
      "title": "김치볶음밥",
      "source_url": "https://youtube.com/watch?v=...",
      "thumbnail_url": "https://example.com/image.jpg",
      "source_type": "INTERNAL",
      "main_ingredient": "김치",
      "category": "한식"
    },
    "label": "주말요리",
    "note": "맛있었어요",
    "rating": 5,
    "saved_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### POST /api/v1/cookbook

**요청**:
```json
{
  "recipe_id": "507f1f77bcf86cd799439012",
  "recipe_snapshot": {
    "title": "김치볶음밥",
    "source_url": "https://youtube.com/watch?v=...",
    "thumbnail_url": "https://example.com/image.jpg",
    "source_type": "INTERNAL",
    "main_ingredient": "김치",
    "category": "한식"
  },
  "label": "주말요리",
  "note": "맛있었어요"
}
```

**응답**: GET /cookbook 단건 응답과 동일

**에러**:
- `409`: 이미 저장된 레시피 (중복 저장 방지)

### GET /api/v1/cookbook/labels

**응답**:
```json
[
  {
    "label": "주말요리",
    "count": 5
  },
  {
    "label": "다이어트",
    "count": 3
  }
]
```

### GET /api/v1/cookbook/:id

**응답**: GET /cookbook 단건 응답과 동일

### PATCH /api/v1/cookbook/:id

**요청**:
```json
{
  "label": "새 라벨",
  "note": "수정된 메모",
  "rating": 4
}
```

**응답**: GET /cookbook 단건 응답과 동일

### DELETE /api/v1/cookbook/:id

**응답**: `204 No Content`

---

## Notification API

### GET /api/v1/notifications

**쿼리 파라미터**:
- `limit`: 최대 개수 (기본값: 50)

**응답**:
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "ingredient_id": "507f1f77bcf86cd799439012",
    "ingredient_name": "당근",
    "type": "D_DAY",
    "scheduled_date": "2024-01-01T00:00:00Z",
    "sent_at": "2024-01-01T09:00:00Z",
    "is_read": false,
    "created_at": "2023-12-29T00:00:00Z"
  }
]
```

### GET /api/v1/notifications/count

**응답**:
```json
{
  "unread_count": 5
}
```

### PATCH /api/v1/notifications/:id/read

**응답**: `204 No Content`

### POST /api/v1/notifications/read-all

**응답**: `204 No Content`

---

## 에러 처리

### 공통 에러 응답 형식

```json
{
  "error": "에러 메시지",
  "code": "ERROR_CODE"
}
```

### HTTP 상태 코드

- `200`: 성공
- `201`: 생성 성공
- `204`: 성공 (응답 본문 없음)
- `400`: 잘못된 요청
- `401`: 인증 필요
- `403`: 권한 없음
- `404`: 리소스 없음
- `409`: 충돌 (중복 등)
- `500`: 서버 오류

---

## 상태 관리 (Pinia Stores)

### useAuthStore

```javascript
{
  user: User | null,
  accessToken: string | null,
  refreshToken: string | null,
  isAuthenticated: computed,
  login(email, password): Promise<TokenResponse>,
  signup(email, password, nickname): Promise<TokenResponse>,
  logout(): Promise<void>,
  fetchUser(): Promise<User>,
  updateProfile(data): Promise<User>
}
```

### useFridgeStore

```javascript
{
  ingredients: Ingredient[],
  summary: FridgeSummary | null,
  fetchIngredients(filter): Promise<Ingredient[]>,
  addIngredient(data): Promise<Ingredient>,
  updateIngredient(id, data): Promise<Ingredient>,
  removeIngredient(id): Promise<void>,
  bulkRemove(ids): Promise<void>,
  fetchSummary(): Promise<FridgeSummary>
}
```

### useRecipeStore

```javascript
{
  recommendations: RecommendationResult[],
  todayRecommendations: RecommendationResult[],
  fetchRecommendations(params): Promise<RecommendationResult[]>,
  fetchTodayRecommendations(): Promise<RecommendationResult[]>,
  searchRecipes(keyword, category): Promise<Recipe[]>,
  getRecipeById(id): Promise<Recipe>
}
```

### useVisionStore

```javascript
{
  jobs: VisionJob[],
  createJob(type, file): Promise<VisionJob>,
  getJob(id): Promise<VisionJob>,
  confirmJob(id, selectedItems): Promise<void>
}
```

### useCookbookStore

```javascript
{
  savedRecipes: SavedRecipe[],
  labels: LabelSummary[],
  fetchSavedRecipes(filter): Promise<SavedRecipe[]>,
  saveRecipe(data): Promise<SavedRecipe>,
  updateSavedRecipe(id, data): Promise<SavedRecipe>,
  removeSavedRecipe(id): Promise<void>,
  fetchLabels(): Promise<LabelSummary[]>
}
```

### useNotificationStore

```javascript
{
  notifications: Alert[],
  unreadCount: number,
  fetchNotifications(limit): Promise<Alert[]>,
  fetchUnreadCount(): Promise<number>,
  markAsRead(id): Promise<void>,
  markAllAsRead(): Promise<void>
}
```

---

## 로딩 및 에러 상태

### 로딩 상태 관리

각 Store에 `loading` 상태 추가:
```javascript
{
  loading: boolean,
  error: string | null
}
```

### 에러 처리 패턴

```javascript
try {
  await store.fetchData()
} catch (error) {
  if (error.response?.status === 401) {
    // 자동으로 인터셉터가 처리
  } else if (error.response?.status === 409) {
    // 중복 에러 - 사용자에게 확인
    showConfirmDialog('이미 존재하는 재료입니다. 병합하시겠습니까?')
  } else {
    // 일반 에러
    showErrorMessage(error.message)
  }
}
```

---

## 제약 조건

- 모든 인증 필요 API는 JWT 토큰 필수
- 토큰 만료 시 자동 갱신 시도
- 파일 업로드는 multipart/form-data 사용
- 날짜 형식: ISO 8601 (YYYY-MM-DDTHH:mm:ssZ)
- 페이지네이션: limit 파라미터 사용 (offset 기반 아님)

---

## 의존성

- `axios`: HTTP 클라이언트
- `pinia`: 상태 관리
- 백엔드 API: `/api/v1/*`
