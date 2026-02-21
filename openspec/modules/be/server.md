# Module: internal/server + pkg/

## 패키지 경로
- `backend/internal/server`
- `backend/pkg/config`
- `backend/pkg/database`
- `backend/pkg/jwt`
- `backend/pkg/gcs`
- `backend/pkg/vertexai`
- `backend/cmd/server`

## 책임
HTTP 서버 초기화, 라우터 구성, 미들웨어 등록, 공유 인프라 유틸리티.

---

## 파일 구조

```
backend/
├── cmd/server/
│   └── main.go
│
├── internal/server/
│   ├── router.go
│   └── middleware/
│       ├── auth.go
│       ├── cors.go
│       └── logger.go
│
└── pkg/
    ├── config/config.go
    ├── database/mongo.go
    ├── jwt/jwt.go
    ├── gcs/client.go
    └── vertexai/client.go
```

---

## cmd/server/main.go

### 역할
- 의존성 주입(DI) 조립점
- 실행 순서:
  1. config 로드
  2. MongoDB 연결
  3. pkg 클라이언트 초기화 (jwt, gcs, vertexai)
  4. 각 도메인 repository → service → handler 초기화
  5. router 구성
  6. notification scheduler 시작 (goroutine)
  7. vision worker 시작 (goroutine)
  8. HTTP 서버 시작 (Graceful shutdown 포함)

---

## internal/server/router.go

### 역할
- Gin 엔진 생성
- 전역 미들웨어 등록 (CORS, Logger, Recovery)
- 도메인 핸들러 라우트 등록
- 인증 필요 그룹과 공개 그룹 분리

### 라우트 그룹 구조
```
/api/v1/
├── (public)
│   ├── POST /auth/signup
│   ├── POST /auth/login
│   └── POST /auth/refresh
│
└── (JWT 미들웨어 적용)
    ├── POST /auth/logout
    ├── GET  /users/me
    ├── PATCH /users/me
    ├── /fridge/**
    ├── /vision/**
    ├── /recommendations/**
    ├── /recipes/**
    ├── /cookbook/**
    └── /notifications/**
```

---

## internal/server/middleware/

### auth.go
- Authorization 헤더에서 Bearer 토큰 추출
- `pkg/jwt`로 검증 → userID를 Gin Context에 주입
- 검증 실패 시 401 반환

### cors.go
- 개발 환경: `*` (전체 허용)
- 프로덕션: 프론트엔드 도메인(Vue.js) 한정

### logger.go
- 요청 메서드, 경로, 상태코드, 처리 시간 구조화 로그

---

## pkg/config/config.go

### Config 구조체 필드
```
Port              string
MongoURI          string
MongoDB           string
JWTSecret         string
JWTAccessTTL      int    // 초
JWTRefreshTTL     int    // 초
GCSBucket         string
GCPProjectID      string
VisionLocation    string
YouTubeAPIKey     string
GoogleSearchKey   string
GoogleSearchCX    string
```

### 로드 방식
- `os.Getenv` 기반 (도커/쿠버네티스 친화적)
- `.env` 파일 지원 (로컬 개발, godotenv)

---

## pkg/database/mongo.go

### 제공 함수
```
Connect(ctx, uri string) (*mongo.Client, error)
GetCollection(client *mongo.Client, db, collection string) *mongo.Collection
```

---

## pkg/jwt/jwt.go

### 제공 함수
```
GenerateAccessToken(userID string, secret string, ttl int) (string, error)
GenerateRefreshToken() (string, error)       // opaque random token
ValidateAccessToken(token, secret string) (userID string, error)
HashToken(token string) string               // refresh token 저장용 hash
```

---

## pkg/gcs/client.go

### 제공 함수
```
Upload(ctx, bucket, objectName string, file io.Reader, contentType string) (url string, error)
GenerateSignedURL(ctx, bucket, objectName string, expiry time.Duration) (string, error)
Delete(ctx, bucket, objectName string) error
```

---

## pkg/vertexai/client.go

### 제공 함수
```
AnalyzeReceipt(ctx, imageURL string) (*VisionResult, error)    // OCR
AnalyzeFridge(ctx, imageURL string) (*VisionResult, error)     // Object Detection
```

### VisionResult 구조체
```
Labels    []DetectedLabel
FullText  string           // OCR 전문
```

### DetectedLabel 구조체
```
Name        string
Confidence  float64
```

---

## 외부 라이브러리 (go.mod 예상 의존성)

| 패키지                              | 용도                    |
|-------------------------------------|-------------------------|
| `github.com/gin-gonic/gin`          | HTTP 프레임워크         |
| `go.mongodb.org/mongo-driver`       | MongoDB 드라이버        |
| `golang.org/x/crypto`               | bcrypt                  |
| `github.com/golang-jwt/jwt/v5`      | JWT                     |
| `cloud.google.com/go/storage`       | GCS                     |
| `cloud.google.com/go/vision/apiv1`  | Cloud Vision API        |
| `github.com/robfig/cron/v3`         | cron 스케줄러           |
| `github.com/joho/godotenv`          | .env 로드               |
| `github.com/google/uuid`            | refresh token 생성      |
