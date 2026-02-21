## Context

냉털마스터 Go 백엔드는 아직 코드가 없는 상태이다.
`openspec/modules/be/` 문서에 6개 도메인(auth, fridge, vision, recommendation, cookbook, notification)의
패키지 구조, 인터페이스 시그니처, 구조체 필드, MongoDB 인덱스 설계가 상세히 정의되어 있다.
이 Change는 그 문서를 기반으로 코드 뼈대(skeleton)만 생성한다 — 기능 구현 없이.

## Goals / Non-Goals

**Goals:**
- `backend/` 디렉토리 아래 Go 프로젝트 레이아웃 확립
- 각 파일에 package 선언, 빈 구조체, 빈 인터페이스, 빈 함수 시그니처 작성
- `go.mod` 의존성 선언 (실제 다운로드 포함)
- 팀원이 어떤 파일에서 무엇을 구현해야 하는지 즉시 파악 가능한 상태

**Non-Goals:**
- 어떠한 비즈니스 로직도 구현하지 않음
- DB 연결, HTTP 서버 실제 기동 없음
- 테스트 파일 작성 없음

## Decisions

### 1. 패키지 구조: Domain-flat

```
internal/<domain>/{model,repository,service,handler}.go
```

- **왜**: 해커톤 속도 최적화. 도메인별 파일을 열면 전체 흐름 즉시 파악.
- **대안**: Layered-global (handler/, service/, repository/ 분리) — 레이어 간 파일 점프가 많아 개발 속도 저하.

### 2. HTTP 프레임워크: Gin

- **왜**: Go 표준 생태계에서 가장 널리 쓰이며, 미들웨어 체인이 직관적.
- **대안**: Echo, Fiber — 성능 차이 미미, Gin이 더 풍부한 예제/레퍼런스.

### 3. 인터페이스 정의 위치: 같은 패키지 내

각 도메인의 `repository.go`, `service.go`에 인터페이스를 직접 정의.

- **왜**: 별도 `port/` 패키지 없이 심플하게 유지. 테스트 시 mock 생성 용이.
- **대안**: Hexagonal 구조 — 해커톤에는 과잉.

### 4. 진입점 DI: main.go에서 수동 조립

`cmd/server/main.go`에서 pkg 초기화 → repository → service → handler 순서로 수동 주입.

- **왜**: Wire/fx 등 DI 프레임워크 러닝커브 없이 명시적 흐름 유지.

### 5. Go 모듈 경로

```
module github.com/okky-hackathon/fridge-master-backend
go 1.22
```

## Risks / Trade-offs

- **빈 함수 반환값**: 모든 함수가 `nil, nil` 또는 빈 값을 반환 — 컴파일은 되지만 런타임 동작 없음. 이후 구현 시 반드시 채워야 함.
- **순환 참조 위험**: vision → fridge, notification → fridge 의존 시 import cycle 가능. service 레이어에서 직접 import 대신 인터페이스 주입으로 해소.
- **go.sum 없음**: 초기 `go mod tidy` 필요. 인터넷 연결 필요.

## Open Questions

- `pkg/vertexai/` 구현 시 Vertex AI SDK vs Cloud Vision API 직접 호출 중 어느 쪽을 사용할지 (스켈레톤에는 무관, 구현 Change에서 결정)
- 외부 YouTube API 캐시 레이어 — 메모리 캐시(sync.Map) vs Redis (MVP는 메모리 캐시로 결정 예정)
