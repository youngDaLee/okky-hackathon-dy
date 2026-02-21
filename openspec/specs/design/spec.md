# Design System Spec

## 개요

냉털마스터 프론트엔드 디자인 시스템 스펙.
일관된 UI/UX를 위한 색상, 타이포그래피, 간격, 컴포넌트 스타일을 정의한다.

---

## 디자인 토큰

### 색상 (Colors)

#### Primary (주요 색상)
- `primary-50`: #f0f9ff
- `primary-100`: #e0f2fe
- `primary-200`: #bae6fd
- `primary-300`: #7dd3fc
- `primary-400`: #38bdf8
- `primary-500`: #0ea5e9
- `primary-600`: #0284c7 (메인 Primary)
- `primary-700`: #0369a1
- `primary-800`: #075985
- `primary-900`: #0c4a6e

#### Semantic Colors
- `success`: #10b981 (성공 메시지, 완료 상태)
- `warning`: #f59e0b (경고, 유통기한 임박)
- `error`: #ef4444 (에러, 삭제)
- `info`: #3b82f6 (정보)

#### Neutral (중립 색상)
- `gray-50`: #f9fafb (배경)
- `gray-100`: #f3f4f6
- `gray-200`: #e5e7eb
- `gray-300`: #d1d5db
- `gray-400`: #9ca3af
- `gray-500`: #6b7280
- `gray-600`: #4b5563
- `gray-700`: #374151
- `gray-800`: #1f2937 (본문 텍스트)
- `gray-900`: #111827 (제목 텍스트)

#### Background
- `bg-primary`: #ffffff (카드, 모달 배경)
- `bg-secondary`: #f9fafb (페이지 배경)
- `bg-overlay`: rgba(0, 0, 0, 0.5) (오버레이)

---

### 타이포그래피 (Typography)

#### 폰트 패밀리
- `font-sans`: 시스템 폰트 스택 (Pretendard 우선, fallback: -apple-system, sans-serif)

#### 폰트 크기
- `text-xs`: 12px (캡션, 보조 텍스트)
- `text-sm`: 14px (작은 본문)
- `text-base`: 16px (본문, 기본)
- `text-lg`: 18px (강조 본문)
- `text-xl`: 20px (소제목)
- `text-2xl`: 24px (섹션 제목)
- `text-3xl`: 30px (페이지 제목)
- `text-4xl`: 36px (대형 제목)

#### 폰트 굵기
- `font-normal`: 400 (본문)
- `font-medium`: 500 (강조)
- `font-semibold`: 600 (소제목)
- `font-bold`: 700 (제목)

#### 행간 (Line Height)
- `leading-tight`: 1.25
- `leading-normal`: 1.5 (본문)
- `leading-relaxed`: 1.75

---

### 간격 (Spacing)

#### 기본 간격 단위
- `space-1`: 4px
- `space-2`: 8px
- `space-3`: 12px
- `space-4`: 16px (기본)
- `space-5`: 20px
- `space-6`: 24px
- `space-8`: 32px
- `space-10`: 40px
- `space-12`: 48px
- `space-16`: 64px

#### 컴포넌트 간격
- `card-padding`: 16px
- `section-gap`: 24px
- `page-padding`: 16px (모바일), 24px (데스크톱)

---

### 그림자 (Shadows)

- `shadow-sm`: 0 1px 2px 0 rgba(0, 0, 0, 0.05)
- `shadow`: 0 1px 3px 0 rgba(0, 0, 0, 0.1)
- `shadow-md`: 0 4px 6px -1px rgba(0, 0, 0, 0.1)
- `shadow-lg`: 0 10px 15px -3px rgba(0, 0, 0, 0.1) (카드)
- `shadow-xl`: 0 20px 25px -5px rgba(0, 0, 0, 0.1)

---

### 둥근 모서리 (Border Radius)

- `rounded-sm`: 2px
- `rounded`: 4px
- `rounded-md`: 6px
- `rounded-lg`: 8px (버튼, 입력 필드)
- `rounded-xl`: 12px (카드)
- `rounded-2xl`: 16px
- `rounded-full`: 9999px (원형)

---

### 전환 효과 (Transitions)

- `transition-fast`: 150ms
- `transition-base`: 200ms (기본)
- `transition-slow`: 300ms

---

## 컴포넌트 스타일 가이드

### 버튼 (Button)

#### Primary Button
- 배경: `primary-600`
- 텍스트: 흰색
- 패딩: `12px 24px`
- 둥근 모서리: `rounded-lg`
- 호버: `primary-700`
- 비활성화: `gray-300`, `cursor-not-allowed`

#### Secondary Button
- 배경: `gray-200`
- 텍스트: `gray-800`
- 호버: `gray-300`

#### Ghost Button
- 배경: 투명
- 텍스트: `primary-600`
- 호버: `primary-50` 배경

---

### 입력 필드 (Input)

- 배경: 흰색
- 테두리: `gray-300`
- 패딩: `12px 16px`
- 둥근 모서리: `rounded-lg`
- 포커스: `primary-600` 테두리, 그림자
- 에러: `error` 테두리

---

### 카드 (Card)

- 배경: 흰색
- 패딩: `16px`
- 둥근 모서리: `rounded-xl`
- 그림자: `shadow-lg`
- 호버: 그림자 강화 (선택적)

---

### 태그 (Tag)

- 배경: `primary-100`
- 텍스트: `primary-700`
- 패딩: `6px 12px`
- 둥근 모서리: `rounded-full`
- 폰트 크기: `text-sm`

---

## 반응형 브레이크포인트

- `sm`: 640px (모바일 가로)
- `md`: 768px (태블릿)
- `lg`: 1024px (데스크톱)
- `xl`: 1280px (큰 데스크톱)
- `2xl`: 1536px (초대형)

---

## 접근성 (Accessibility)

- 최소 터치 영역: 44x44px
- 색상 대비: WCAG AA 기준 준수 (4.5:1 이상)
- 포커스 표시: 명확한 포커스 링 (outline)

---

## 제약 조건

- 모든 색상은 Tailwind CSS 클래스로 사용
- 커스텀 색상은 `tailwind.config.js`에 정의
- 컴포넌트는 재사용 가능하도록 설계
- 모바일 퍼스트 접근 방식

---

## 의존성

- Tailwind CSS v3.4+
- Vue 3 Composition API
