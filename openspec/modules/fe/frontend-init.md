# 프론트엔드 초기 설정

## 개요

냉털마스터 프론트엔드는 **Vite + Vue 3 + Tailwind CSS**로 구성된 모던 웹 애플리케이션입니다.

## 기술 스택

| 기술 | 버전 | 용도 |
|------|------|------|
| Vite | ^5.1.6 | 빌드 도구 및 개발 서버 |
| Vue | ^3.4.21 | 프론트엔드 프레임워크 |
| Vue Router | ^4.3.0 | 클라이언트 사이드 라우팅 |
| Pinia | ^2.1.7 | 상태 관리 |
| Tailwind CSS | ^3.4.1 | 유틸리티 CSS 프레임워크 |
| Axios | ^1.6.7 | HTTP 클라이언트 |
| ESLint | ^8.57.0 | 코드 린팅 |

---

## 프로젝트 구조

```
frontend/
├── src/
│   ├── api/
│   │   └── index.js          # Axios API 클라이언트 (인터셉터 포함)
│   ├── components/           # 재사용 가능한 Vue 컴포넌트
│   ├── router/
│   │   └── index.js          # Vue Router 설정
│   ├── stores/
│   │   └── auth.js           # Pinia 인증 스토어
│   ├── utils/
│   │   └── index.js          # 유틸리티 함수
│   ├── views/
│   │   ├── HomeView.vue      # 홈 페이지
│   │   ├── LoginView.vue     # 로그인 페이지
│   │   └── SignupView.vue    # 회원가입 페이지
│   ├── App.vue               # 루트 컴포넌트
│   ├── main.js               # 애플리케이션 진입점
│   └── style.css             # 전역 스타일 (Tailwind CSS)
├── public/                   # 정적 파일
├── index.html                # HTML 템플릿
├── package.json              # 의존성 및 스크립트
├── vite.config.js            # Vite 설정
├── tailwind.config.js        # Tailwind CSS 설정
├── postcss.config.js         # PostCSS 설정
├── .eslintrc.cjs             # ESLint 설정
└── .gitignore
```

---

## 주요 설정 파일

### 1. `package.json`

**의존성**:
- `vue`, `vue-router`, `pinia`: Vue 생태계 핵심 라이브러리
- `axios`: HTTP 클라이언트
- `tailwindcss`, `autoprefixer`, `postcss`: CSS 프레임워크 및 처리

**스크립트**:
```json
{
  "dev": "vite --port 3000 --host 0.0.0.0",
  "build": "vite build",
  "preview": "vite preview",
  "lint": "eslint . --ext .vue,.js,.jsx,.cjs,.mjs --fix --ignore-path .gitignore"
}
```

### 2. `vite.config.js`

**주요 설정**:
- **경로 별칭**: `@` → `./src` (절대 경로 import 지원)
- **프록시**: `/api` → `http://localhost:8080` (백엔드 API 프록시)
- **개발 서버**: 포트 3000

```javascript
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

### 3. `tailwind.config.js`

**커스텀 설정**:
- Primary 색상 팔레트 정의 (50-900)
- Vue 파일 및 HTML 파일에서 Tailwind 클래스 검색

```javascript
theme: {
  extend: {
    colors: {
      primary: {
        50: '#f0f9ff',
        // ... 900까지
      }
    }
  }
}
```

### 4. `postcss.config.js`

Tailwind CSS와 Autoprefixer를 PostCSS 플러그인으로 설정.

---

## 핵심 모듈

### 1. API 클라이언트 (`src/api/index.js`)

**기능**:
- Axios 인스턴스 생성 (`baseURL: '/api/v1'`)
- **Request 인터셉터**: 모든 요청에 `Authorization: Bearer {token}` 헤더 자동 추가
- **Response 인터셉터**: 
  - 401 에러 시 자동 토큰 갱신 시도
  - 토큰 갱신 실패 시 로그인 페이지로 리다이렉트

**사용 예시**:
```javascript
import api from '@/api'

// 자동으로 Authorization 헤더가 추가됨
const response = await api.get('/users/me')
```

### 2. 인증 스토어 (`src/stores/auth.js`)

**Pinia Composition API 스타일**로 작성된 인증 스토어.

**상태**:
- `user`: 현재 로그인한 사용자 정보
- `accessToken`: JWT 액세스 토큰
- `refreshToken`: 리프레시 토큰
- `isAuthenticated`: 로그인 여부 (computed)

**액션**:
- `login(email, password)`: 로그인 및 토큰 저장
- `signup(email, password, nickname)`: 회원가입 및 자동 로그인
- `fetchUser()`: 현재 사용자 정보 조회
- `logout()`: 로그아웃 및 토큰 삭제

**토큰 저장**: `localStorage`에 `access_token`, `refresh_token` 저장

**사용 예시**:
```javascript
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
await authStore.login('user@example.com', 'password')
```

### 3. 라우터 (`src/router/index.js`)

**기본 라우트**:
- `/`: 홈 페이지 (`HomeView`)
- `/login`: 로그인 페이지 (`LoginView`)
- `/signup`: 회원가입 페이지 (`SignupView`)

**History 모드**: `createWebHistory` 사용 (깔끔한 URL)

### 4. 유틸리티 (`src/utils/index.js`)

**함수**:
- `formatDate(date)`: 날짜를 한국어 형식으로 포맷
- `formatExpiryStatus(status)`: 유통기한 상태 한글 변환
- `getExpiryStatusColor(status)`: 상태별 Tailwind CSS 클래스 반환

---

## 기본 컴포넌트

### 1. `App.vue`

루트 컴포넌트. 최소한의 레이아웃만 제공:
```vue
<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <RouterView />
  </div>
</template>
```

### 2. `HomeView.vue`

랜딩 페이지. 로그인/회원가입 링크 제공.

### 3. `LoginView.vue`

로그인 폼:
- 이메일 입력
- 비밀번호 입력
- 제출 시 `authStore.login()` 호출 (현재는 TODO)

### 4. `SignupView.vue`

회원가입 폼:
- 이메일 입력
- 비밀번호 입력 (최소 8자)
- 닉네임 입력 (2-20자)
- 제출 시 `authStore.signup()` 호출 (현재는 TODO)

---

## 스타일링

### Tailwind CSS 사용

**전역 스타일** (`src/style.css`):
```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

**커스텀 색상**:
- `primary-*`: Primary 색상 팔레트 (50-900)
- 예: `bg-primary-600`, `text-primary-700`

**반응형 디자인**: Tailwind의 기본 브레이크포인트 사용
- `sm`: 640px
- `md`: 768px
- `lg`: 1024px
- `xl`: 1280px

---

## 개발 워크플로우

### 1. 개발 서버 실행

```bash
cd frontend
npm install
npm run dev
```

개발 서버는 `http://localhost:3000`에서 실행됩니다.

### 2. 빌드

```bash
npm run build
```

빌드 결과물은 `dist/` 디렉토리에 생성됩니다.

### 3. 린팅

```bash
npm run lint
```

---

## 백엔드 연동

### API 엔드포인트

프론트엔드는 다음 API 엔드포인트를 사용합니다:

| 메서드 | 경로 | 설명 |
|--------|------|------|
| POST | `/api/v1/auth/login` | 로그인 |
| POST | `/api/v1/auth/signup` | 회원가입 |
| POST | `/api/v1/auth/refresh` | 토큰 갱신 |
| POST | `/api/v1/auth/logout` | 로그아웃 |
| GET | `/api/v1/users/me` | 현재 사용자 정보 |

### 프록시 설정

개발 환경에서 `/api` 요청은 자동으로 `http://localhost:8080`으로 프록시됩니다.

---

## 다음 단계

1. **인증 플로우 완성**: LoginView, SignupView에서 실제 API 호출 구현
2. **보호된 라우트**: 인증이 필요한 페이지에 라우트 가드 추가
3. **레이아웃 컴포넌트**: 공통 헤더/네비게이션 바 추가
4. **에러 처리**: API 에러 메시지 표시 컴포넌트
5. **로딩 상태**: 비동기 작업 중 로딩 인디케이터

---

## 참고 사항

- **경로 별칭**: `@/`를 사용하여 `src/` 디렉토리에서 절대 경로로 import 가능
- **토큰 관리**: `localStorage`에 저장되며, API 클라이언트가 자동으로 헤더에 추가
- **자동 토큰 갱신**: 401 에러 시 자동으로 refresh token으로 갱신 시도
- **CORS**: 개발 환경에서는 Vite 프록시를 통해 CORS 문제 해결
