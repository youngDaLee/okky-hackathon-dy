# 냉털마스터 프론트엔드

Vite + Vue 3 + Tailwind CSS로 구성된 프론트엔드 프로젝트입니다.

## 기술 스택

- **Vite**: 빌드 도구
- **Vue 3**: 프론트엔드 프레임워크
- **Vue Router**: 라우팅
- **Pinia**: 상태 관리
- **Tailwind CSS**: 유틸리티 CSS 프레임워크
- **Axios**: HTTP 클라이언트

## 설치

```bash
npm install
```

## 개발 서버 실행

```bash
npm run dev
```

개발 서버는 `http://localhost:3000`에서 실행됩니다.

## 빌드

```bash
npm run build
```

## 프로젝트 구조

```
frontend/
├── src/
│   ├── api/          # API 클라이언트
│   ├── components/   # Vue 컴포넌트
│   ├── router/       # 라우터 설정
│   ├── stores/       # Pinia 스토어
│   ├── utils/        # 유틸리티 함수
│   ├── views/        # 페이지 컴포넌트
│   ├── App.vue       # 루트 컴포넌트
│   ├── main.js       # 진입점
│   └── style.css     # 전역 스타일
├── public/           # 정적 파일
├── index.html        # HTML 템플릿
├── vite.config.js    # Vite 설정
├── tailwind.config.js # Tailwind 설정
└── package.json
```
