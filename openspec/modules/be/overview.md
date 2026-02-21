# Backend 프로젝트 구조 (Go)

## 설계 원칙

- **Domain-flat 패키지 구조**: 도메인별로 handler/service/repository/model을 하나의 패키지에 모음
- **내부(internal) / 공유(pkg) 분리**: 도메인 로직은 `internal/`, 도메인 무관 유틸은 `pkg/`
- **단방향 의존성**: pkg ← internal ← cmd (역방향 import 금지)
- **HTTP 프레임워크**: Gin

---

## 디렉토리 트리

```
backend/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── server/
│   │   ├── router.go
│   │   └── middleware/
│   │       ├── auth.go
│   │       ├── cors.go
│   │       └── logger.go
│   │
│   ├── auth/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── handler.go
│   │
│   ├── fridge/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── handler.go
│   │
│   ├── vision/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── worker.go
│   │
│   ├── recommendation/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── external.go
│   │
│   ├── cookbook/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── handler.go
│   │
│   └── notification/
│       ├── model.go
│       ├── repository.go
│       ├── service.go
│       ├── handler.go
│       └── scheduler.go
│
└── pkg/
    ├── config/
    │   └── config.go
    ├── database/
    │   └── mongo.go
    ├── jwt/
    │   └── jwt.go
    ├── gcs/
    │   └── client.go
    └── vertexai/
        └── client.go
```

---

## 도메인 간 의존 관계

```
auth ──────────────────────────────── (독립)
fridge ────→ auth (user_id 검증)
vision ────→ auth, fridge (confirm 시 재료 등록)
recommendation → auth, fridge (재료 목록 + dietary_prefs)
cookbook ──→ auth, recommendation (Recipe 참조)
notification → auth, fridge (expiry_date 참조)
```

```
[cmd/server/main.go]
     │
     ├── pkg/config      환경변수
     ├── pkg/database    MongoDB 연결
     ├── pkg/jwt         토큰 유틸
     ├── pkg/gcs         GCS 클라이언트
     ├── pkg/vertexai    Vision AI 클라이언트
     │
     └── internal/server/router.go
              │
              ├── internal/auth/handler
              ├── internal/fridge/handler
              ├── internal/vision/handler
              ├── internal/recommendation/handler
              ├── internal/cookbook/handler
              └── internal/notification/handler
```

---

## 파일 역할 요약

### 공통 파일 패턴 (각 도메인)

| 파일              | 역할                                                   |
|-------------------|--------------------------------------------------------|
| `model.go`        | MongoDB 문서 구조체, 요청/응답 DTO 정의                |
| `repository.go`   | MongoDB CRUD, 쿼리 로직 (DB 접근만 담당)               |
| `service.go`      | 비즈니스 로직, 도메인 규칙 적용                        |
| `handler.go`      | HTTP 핸들러: 요청 파싱 → service 호출 → 응답 직렬화    |

### 도메인별 추가 파일

| 파일                           | 위치               | 역할                              |
|--------------------------------|--------------------|-----------------------------------|
| `worker.go`                    | vision/            | Vertex AI 비동기 처리 고루틴      |
| `external.go`                  | recommendation/    | YouTube/Google 외부 API 연동      |
| `scheduler.go`                 | notification/      | 유통기한 cron 배치 잡             |
| `router.go`                    | server/            | 전체 라우트 등록 + 미들웨어 체인  |
| `middleware/auth.go`           | server/            | JWT 파싱 → Context에 user_id 주입 |

---

## pkg/ 상세

| 패키지          | 책임                                              |
|-----------------|---------------------------------------------------|
| `config`        | 환경변수 로드 (`.env` 또는 OS 환경변수)           |
| `database`      | MongoDB 클라이언트 싱글톤, 컬렉션 헬퍼            |
| `jwt`           | Access/Refresh 토큰 생성·검증·파싱               |
| `gcs`           | GCS 버킷 업로드, Signed URL 생성                  |
| `vertexai`      | Cloud Vision API 요청 래핑, 응답 파싱             |

---

## 환경변수 목록 (config)

```
# Server
PORT=8080

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB=fridge_master

# JWT
JWT_SECRET=<secret>
JWT_ACCESS_TTL=3600      # 1시간 (초)
JWT_REFRESH_TTL=604800   # 7일 (초)

# GCS
GCS_BUCKET=fridge-master-images
GOOGLE_APPLICATION_CREDENTIALS=/path/to/key.json

# Vertex AI / Cloud Vision
VISION_PROJECT_ID=<gcp-project>
VISION_LOCATION=asia-northeast3

# External API
YOUTUBE_API_KEY=<key>
GOOGLE_SEARCH_API_KEY=<key>
GOOGLE_SEARCH_CX=<cx>
```
