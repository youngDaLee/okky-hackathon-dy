## ADDED Requirements

### Requirement: Go 모듈 초기화
백엔드 프로젝트는 `backend/` 디렉토리를 루트로 하는 Go 모듈을 SHALL 포함해야 한다.
모듈 경로는 `github.com/okky-hackathon/fridge-master-backend`, Go 버전은 1.22 이상이어야 한다.

#### Scenario: go.mod 존재
- **WHEN** `backend/` 디렉토리를 확인하면
- **THEN** `go.mod` 파일이 존재하고 module 경로와 Go 버전이 선언되어 있어야 한다

---

### Requirement: 진입점 파일
서버는 `cmd/server/main.go`를 진입점으로 SHALL 가져야 한다.
스켈레톤 단계에서는 `package main`과 빈 `main()` 함수만 존재한다.

#### Scenario: main.go 컴파일
- **WHEN** `go build ./cmd/server/` 를 실행하면
- **THEN** 컴파일 오류 없이 바이너리가 생성되어야 한다

---

### Requirement: 도메인 패키지 레이아웃
`internal/` 아래 6개 도메인 패키지(auth, fridge, vision, recommendation, cookbook, notification)가 SHALL 존재해야 한다.
각 패키지는 `model.go`, `repository.go`, `service.go`, `handler.go` 파일을 포함해야 한다.

#### Scenario: 도메인 파일 존재
- **WHEN** `internal/<domain>/` 디렉토리를 확인하면
- **THEN** model.go, repository.go, service.go, handler.go 4개 파일이 모두 존재해야 한다

#### Scenario: 패키지 선언 일치
- **WHEN** 각 파일을 열면
- **THEN** 파일 최상단에 `package <domain>` 선언이 있어야 한다

---

### Requirement: 특수 파일 포함 (vision, recommendation, notification)
vision 패키지는 `worker.go`를, recommendation 패키지는 `external.go`를, notification 패키지는 `scheduler.go`를 SHALL 추가로 포함해야 한다.

#### Scenario: vision worker 파일 존재
- **WHEN** `internal/vision/` 디렉토리를 확인하면
- **THEN** `worker.go` 파일이 존재하고 `package vision`으로 선언되어 있어야 한다

#### Scenario: recommendation external 파일 존재
- **WHEN** `internal/recommendation/` 디렉토리를 확인하면
- **THEN** `external.go` 파일이 존재하고 `package recommendation`으로 선언되어 있어야 한다

#### Scenario: notification scheduler 파일 존재
- **WHEN** `internal/notification/` 디렉토리를 확인하면
- **THEN** `scheduler.go` 파일이 존재하고 `package notification`으로 선언되어 있어야 한다

---

### Requirement: 서버 및 미들웨어 패키지
`internal/server/` 패키지는 `router.go`를, `internal/server/middleware/` 패키지는 `auth.go`, `cors.go`, `logger.go`를 SHALL 포함해야 한다.

#### Scenario: server 패키지 파일 존재
- **WHEN** `internal/server/` 디렉토리를 확인하면
- **THEN** router.go 파일이 존재해야 한다

#### Scenario: middleware 파일 존재
- **WHEN** `internal/server/middleware/` 디렉토리를 확인하면
- **THEN** auth.go, cors.go, logger.go 3개 파일이 모두 존재해야 한다

---

### Requirement: 공유 유틸리티 패키지
`pkg/` 아래 config, database, jwt, gcs, vertexai 5개 패키지가 SHALL 존재해야 하며, 각각 단일 `.go` 파일을 포함해야 한다.

#### Scenario: pkg 패키지 파일 존재
- **WHEN** `pkg/<name>/` 디렉토리를 확인하면
- **THEN** 각 패키지마다 해당 `.go` 파일(config.go, mongo.go, jwt.go, client.go)이 존재해야 한다

---

### Requirement: 전체 프로젝트 컴파일 가능
스켈레톤 상태에서 `go build ./...`를 실행하면 컴파일 오류가 SHALL 없어야 한다.
(기능 구현 없이 빈 시그니처만 있어도 컴파일 가능한 상태를 유지해야 한다.)

#### Scenario: 전체 빌드 성공
- **WHEN** `backend/` 루트에서 `go build ./...`를 실행하면
- **THEN** 컴파일 오류 없이 종료되어야 한다
