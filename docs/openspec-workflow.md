# OpenSpec 워크플로우 규칙

Spec-Driven 개발을 위한 OpenSpec 스킬 사용 규칙을 정리합니다.

---

## 개요

OpenSpec은 기능 변경을 **artifact** 단위로 관리합니다.
`proposal → specs → design → tasks` 순서로 문서를 작성한 뒤 구현(`apply`)하고, 검증(`verify`) 후 보관(`archive`)합니다.

```
openspec/
  changes/
    <change-name>/        # 진행 중인 변경
      proposal.md
      specs/<capability>/spec.md
      design.md
      tasks.md
    archive/
      YYYY-MM-DD-<name>/  # 완료된 변경
  specs/
    <capability>/spec.md  # 메인 스펙 (archive 시 반영)
```

---

## 워크플로우 전체 흐름

```
새 변경 시작
    │
    ├── /openspec-new-change   (단계별 진행)
    │   또는
    └── /openspec-ff-change    (한 번에 artifacts 생성)
            │
            ▼
    /openspec-apply-change     (태스크 구현)
            │
            ▼
    /openspec-verify-change    (구현 검증)
            │
            ▼
    /openspec-archive-change   (완료 보관)
```

---

## 스킬 상세 규칙

### `/openspec-new-change` — 새 변경 시작

- 입력: kebab-case 이름 또는 설명 텍스트
- `openspec new change "<name>"` 실행 후 첫 번째 artifact 템플릿만 보여줍니다
- artifact를 직접 생성하지 않으며 사용자 확인 후 대기합니다
- 이미 같은 이름의 변경이 있으면 `/openspec-continue-change` 사용을 안내합니다

### `/openspec-continue-change` — 다음 artifact 작성

- 한 번 호출 시 **artifact 1개만** 생성합니다
- 의존하는 artifact를 먼저 읽고 작성합니다
- `context`, `rules` 블록은 AI 내부 제약으로만 사용하며 출력 파일에 포함하지 않습니다
- 진행 중인 변경이 여러 개면 사용자에게 선택을 요청합니다

### `/openspec-ff-change` — 빠른 전체 artifact 생성

- `apply`에 필요한 모든 artifact를 한 번에 생성합니다
- `openspec status --json`의 `applyRequires` 배열 기준으로 완료 여부를 판단합니다
- 각 artifact 생성 후 `openspec status`로 진행 상황을 확인합니다
- 컨텍스트가 불명확하면 중간에 사용자에게 질문합니다

### `/openspec-apply-change` — 태스크 구현

- `tasks.md`의 체크박스 순서대로 구현합니다
- 태스크 완료 시 바로 `- [ ]` → `- [x]`로 표시합니다
- 다음 경우 일시 중지 후 사용자 안내:
  - 태스크 내용이 불명확한 경우
  - 구현 중 설계 이슈 발견
  - 에러 또는 블로커 발생
- 코드 변경은 태스크 범위 내 최소한으로 유지합니다

### `/openspec-verify-change` — 구현 검증

3가지 차원으로 검증하고 리포트를 생성합니다:

| 차원 | 검증 내용 | 이슈 등급 |
|------|----------|----------|
| **Completeness** | 태스크 완료 여부, 스펙 요구사항 구현 여부 | CRITICAL |
| **Correctness** | 요구사항과 구현 일치 여부, 시나리오 커버리지 | WARNING |
| **Coherence** | design.md 결정 사항 준수, 코드 패턴 일관성 | SUGGESTION |

- CRITICAL: archive 전 반드시 수정
- WARNING: 수정 권장
- SUGGESTION: 선택적 개선

### `/openspec-archive-change` — 변경 완료 보관

1. artifact 완료 여부 확인 (미완성 있으면 경고 후 확인)
2. 미완료 태스크 확인 (있으면 경고 후 확인)
3. delta specs 존재 시 main specs와 비교 후 sync 여부 결정
4. `openspec/changes/archive/YYYY-MM-DD-<name>/` 으로 이동

### `/openspec-sync-specs` — 스펙 동기화 (단독)

- `openspec/changes/<name>/specs/` → `openspec/specs/<capability>/spec.md` 반영
- archive 없이 스펙만 먼저 동기화할 때 사용합니다

---

## Spec-Driven 스키마 artifact 순서

| 순서 | artifact | 내용 |
|------|----------|------|
| 1 | `proposal.md` | Why, What Changes, Capabilities, Impact |
| 2 | `specs/<capability>/spec.md` | proposal의 Capabilities 항목당 1개 |
| 3 | `design.md` | 기술 결정, 아키텍처, 구현 접근법 |
| 4 | `tasks.md` | 구현 태스크 체크리스트 |

> `specs`는 proposal의 Capabilities 개수만큼 생성합니다.
> spec 파일명은 change 이름이 아닌 capability 이름을 사용합니다.

---

## 공통 주의사항

- `context`, `rules` 블록은 AI 제약 조건이며 artifact 파일에 포함하지 않습니다
- `openspec-continue-change`는 한 번에 artifact 1개만 생성합니다
- 변경 이름은 반드시 kebab-case여야 합니다 (예: `add-user-auth`)
- archive 전 `/openspec-verify-change`로 검증을 권장합니다
- 구현 중 설계 문제 발견 시 artifact를 수정할 수 있습니다 (단계 고정 아님)