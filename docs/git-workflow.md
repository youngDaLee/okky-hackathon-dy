# Git 워크플로우 규칙

Claude Code 스킬 기반 커밋 / 이슈 / PR 생성 규칙을 정리합니다.

---

## 커밋 (`/commit`)

### 커밋 메시지 형식

```
type(scope): 한 줄 요약

Why:
- 왜 이 변경이 필요한가?

What:
- 무엇이 어떻게 바뀌었나?

Verify:
- 어떻게 검증했나? (예: pnpm lint / pnpm typecheck)

Refs:
- 관련 이슈/링크 (없으면 N/A)

AI:
- AI 도구를 사용했다면 무엇을 위임했고 무엇을 직접 검증했나
```

### type 기준

| type | 사용 시점 |
|------|----------|
| `feat` | 새로운 기능 추가 |
| `fix` | 버그 수정 |
| `docs` | 문서 변경 |
| `style` | 코드 포맷 (동작 변경 없음) |
| `refactor` | 리팩토링 (기능 변경 없음) |
| `test` | 테스트 추가/수정 |
| `chore` | 빌드, 설정, 패키지 등 |
| `ci` | CI/CD 설정 변경 |
| `perf` | 성능 개선 |

### scope 기준

변경된 파일 경로에서 모듈 단위를 추출합니다.
- `src/auth/` → `auth`
- `src/components/` → `components`
- `docker-compose.yml` → `infra`
- 여러 영역에 걸치면 가장 핵심적인 scope 하나 선택 또는 생략

### 커밋 절차

1. staged 변경이 있으면 staged만 대상으로 합니다
2. staged 변경이 없으면 어떤 파일을 커밋할지 확인 후 진행합니다
3. 커밋 전 위험 요소를 점검합니다:
   - `TODO`, `FIXME`, `HACK` 등 임시 마커
   - `console.log`, `debugger` 등 디버깅 코드
   - 하드코딩된 API 키, 시크릿
   - `.env`, 인증 정보 파일
4. 메시지 초안을 보여주고 사용자 확인 후 커밋합니다

### 명령어 옵션

| 명령어 | 동작 |
|--------|------|
| `/commit` | 변경사항 분석 후 메시지 자동 생성 |
| `/commit <힌트>` | 힌트를 바탕으로 메시지 작성 |
| `/commit --all` | 모든 변경사항 staging 후 커밋 |
| `/commit --amend` | 직전 커밋 메시지 수정 (push된 커밋엔 사용 금지) |

---

## 이슈 생성 (`/create-issue`)

### 이슈 제목 형식

```
[feat] 로그인 페이지 UI 구현
[fix] 회원가입 이메일 중복 체크 오류
[docs] README에 프로젝트 설명 추가
[infra] Docker Compose 배포 환경 구성
[design] 메인 화면 와이어프레임 작성
[refactor] API 응답 구조 통일
[test] 결제 모듈 단위 테스트 추가
```

### 이슈 본문 템플릿

```markdown
## 📋 작업 설명

## ✅ 완료 조건 (Definition of Done)
- [ ] 항목 1
- [ ] 항목 2

## 🔗 관련 사항
- 관련 이슈: #이슈번호
- 참고 자료:

## 📌 비고
- 우선순위: high / medium / low
- 예상 소요:
```

### 라벨 규칙

| 카테고리 | 라벨 | 우선순위 | 라벨 |
|---------|------|---------|------|
| feat | `feature` | high | `priority: high` |
| fix | `bug` | medium | `priority: medium` |
| docs | `documentation` | low | `priority: low` |
| infra | `infra` | | |
| design | `design` | | |

리포지토리에 라벨이 없으면 자동 생성합니다.

### 주의사항

- 이슈 하나 = 한 사람이 해커톤 시간 내에 완료 가능한 단위
- 의존 관계는 본문 "관련 사항"에 명시

---

## 브랜치 규칙

### 브랜치 명명 규칙

이슈 번호를 기준으로 브랜치를 생성합니다.

```
<type>/<issue-number>-<short-description>
```

**예시:**

```
feat/12-login-page
fix/34-email-validation
docs/7-readme-update
refactor/21-api-response
```

### 브랜치 생성 절차

1. GitHub에서 이슈를 먼저 생성합니다 (`/create-issue`)
2. 이슈 번호를 확인합니다 (예: `#12`)
3. 브랜치를 생성합니다:

```bash
git checkout -b feat/12-login-page
```

4. 작업 완료 후 `/create-pr` 스킬로 PR을 생성합니다

### type 기준

커밋 type과 동일하게 사용합니다.

| type | 사용 시점 |
|------|----------|
| `feat` | 새로운 기능 |
| `fix` | 버그 수정 |
| `docs` | 문서 변경 |
| `refactor` | 리팩토링 |
| `chore` | 설정, 빌드 등 |

### 주의사항

- `main` 브랜치에서 직접 작업하지 않습니다
- 브랜치명에 이슈 번호가 포함되면 `/create-pr` 실행 시 이슈가 자동 연결됩니다
- 이슈 하나당 브랜치 하나를 원칙으로 합니다

---

## PR 생성 (`/create-pr`)

### PR 제목 형식

커밋과 동일한 `type(scope): 요약` 컨벤션을 따릅니다.

### PR 본문 템플릿

```markdown
## 📋 개요

## 🎯 관련 이슈
- closes #이슈번호

## 🔄 변경 사항

### Why
-

### What
-

### How
-

## ✅ 체크리스트
- [ ] 코드 정상 동작 확인
- [ ] 린트/타입체크 통과
- [ ] 관련 이슈 연결
- [ ] 팀원 리뷰 요청

## 🧪 테스트 방법
1.
2.

## 🤖 AI 사용 내역
-
```

### 라벨 규칙

| 조건 | 라벨 |
|------|------|
| type: `feat` | `feature` |
| type: `fix` | `bug` |
| type: `docs` | `documentation` |
| type: `chore`/`ci` | `infra` |
| 변경 파일 < 5개 | `size: S` |
| 변경 파일 5~15개 | `size: M` |
| 변경 파일 > 15개 | `size: L` |

### 이슈 자동 연결 순서

1. 브랜치명에서 추출 (`feat/12-login` → `#12`)
2. 커밋 메시지 `Refs:` 섹션에서 추출
3. `/create-pr #12` 인자에서 추출
4. 탐지 실패 시 사용자에게 질문

### 주의사항

- PR 본문은 사용자 확인 후 생성합니다
- `main` 브랜치에서 직접 실행 시 피처 브랜치로 안내합니다
- 변경 파일 20개 이상이면 PR 분리를 제안합니다