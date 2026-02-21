---
name: commit
description: 변경된 파일을 분석하여 .gitmessage 컨벤션에 맞는 커밋 메시지를 작성하고 커밋합니다. Use when committing changes, staging files, or when the user says "커밋", "commit", or "변경사항 정리".
---

# Git Commit Skill

프로젝트의 `.gitmessage` 컨벤션에 따라 변경 파일을 분석하고 구조화된 커밋 메시지를 작성하여 커밋합니다.

## 사전 확인

1. 현재 디렉토리가 git 저장소인지 확인합니다.
2. `git status`로 변경사항이 있는지 확인합니다.
3. 변경사항이 없으면 "커밋할 변경사항이 없습니다"라고 안내합니다.

## 변경사항 분석

### 1단계: 변경 파일 파악

```bash
git status --short
git diff --staged --stat
git diff --stat
```

- staged 변경이 있으면 staged 변경만 대상으로 합니다.
- staged 변경이 없으면 전체 변경사항을 보여주고, 어떤 파일을 커밋할지 사용자에게 확인합니다.

### 2단계: 변경 내용 상세 분석

```bash
git diff --staged    # staged인 경우
git diff             # unstaged인 경우
```

변경 내용을 읽고 다음을 파악합니다:
- 어떤 파일이 추가/수정/삭제되었는지
- 변경의 목적과 맥락 (기능 추가? 버그 수정? 리팩토링?)
- 영향 범위 (scope)

### 3단계: 위험 요소 점검

커밋 전 다음 항목을 검사합니다:

- `TODO`, `FIXME`, `HACK`, `XXX` 등 임시 마커
- `console.log`, `print()`, `debugger` 등 디버깅 코드
- 주석 처리된 코드 블록
- 하드코딩된 테스트 데이터, API 키, 시크릿
- `.env`, 인증 정보 등 민감한 파일

발견 시 사용자에게 경고하고 계속 진행할지 확인합니다.

## 커밋 메시지 작성

### 형식

프로젝트 `.gitmessage` 컨벤션을 따릅니다:

```
type(scope): 한 줄 요약

Why:
- (왜 이 변경이 필요한가?)

What:
- (무엇이 어떻게 바뀌었나?)

Verify:
- (어떻게 검증했나? 예: pnpm lint / pnpm typecheck / pnpm test:e2e)

Refs:
- (관련 이슈/문서/링크)

AI:
- (AI 도구를 사용했다면: 무엇을 위임했고, 무엇을 직접 검증했나?)
```

### type 분류 기준

| type | 사용 시점 |
|------|----------|
| `feat` | 새로운 기능 추가 |
| `fix` | 버그 수정 |
| `docs` | 문서 변경 (README, 주석 등) |
| `style` | 코드 포맷, 세미콜론 등 (동작 변경 없음) |
| `refactor` | 리팩토링 (기능 변경 없음) |
| `test` | 테스트 추가/수정 |
| `chore` | 빌드, 설정, 패키지 등 |
| `ci` | CI/CD 설정 변경 |
| `perf` | 성능 개선 |

### scope 결정 기준

변경된 파일 경로에서 모듈/기능 단위를 추출합니다:
- `src/auth/login.ts` → `auth`
- `src/components/Header.tsx` → `components`
- `docker-compose.yml` → `infra`
- `README.md` → `docs`
- 여러 영역에 걸치면 가장 핵심적인 scope 하나를 선택하거나 생략합니다.

### 각 섹션 작성 가이드

**한 줄 요약**: 50자 이내, 명령형("추가", "수정", "제거"), 마침표 없음

**Why**: 이 변경이 필요한 이유. "왜?"에 대한 답.
- 비즈니스 요구사항, 버그 리포트, 기술적 필요성 등

**What**: 구체적으로 무엇이 바뀌었는지.
- 파일 단위가 아닌 논리적 변경 단위로 기술
- 핵심 변경만 간결하게 (파일 목록을 나열하지 않음)

**Verify**: 검증 방법.
- 실행한 테스트, 린트, 타입체크 명령어
- 수동 확인한 사항
- 검증하지 못한 항목이 있으면 솔직하게 기재

**Refs**: 관련 참조.
- GitHub 이슈: `#123`
- 외부 링크, 문서, Slack 스레드 등
- 없으면 `- N/A`

**AI**: AI 도구 사용 여부. 해커톤에서의 투명성을 위한 섹션.
- Claude Code로 생성/수정한 부분
- 사람이 직접 검토/수정한 부분
- AI를 사용하지 않았으면 `- N/A`

## 실행 절차

1. 변경사항을 분석합니다.
2. 위험 요소를 점검합니다.
3. 커밋 메시지 초안을 사용자에게 보여줍니다.
4. 사용자가 확인하면 (또는 수정 요청을 반영한 후) 커밋합니다.

```bash
# staged 변경이 없으면 사용자 확인 후 staging
git add <files>

    # 커밋 실행
    git commit -m "type(scope): 한 줄 요약

    Why:
    - 이유

    What:
    - 변경 내용

    Verify:
    - 검증 방법

    Refs:
    - 참조

    AI:
    - AI 사용 내역"
    ```

    5. 커밋 결과를 보여줍니다:

    ```bash
    git log --oneline -1
    ```

    ## 인자 처리

    - `/commit` — 변경사항 분석 후 메시지 자동 생성
    - `/commit 로그인 기능 완성` — 힌트를 바탕으로 메시지 작성
    - `/commit --all` — 모든 변경사항을 staging 후 커밋
    - `/commit --amend` — 직전 커밋 메시지를 수정

    ## 예시

    ### 입력
    ```
    /commit
    ```

    ### 출력 (사용자 확인용 초안)
    ```
    feat(auth): 카카오 소셜 로그인 구현

    Why:
    - 해커톤 MVP에서 간편 로그인이 필수 요구사항

    What:
    - KakaoStrategy OAuth2 인증 플로우 추가
    - /api/auth/kakao 콜백 엔드포인트 구현
    - 로그인 성공 시 JWT 토큰 발급 로직 적용

    Verify:
    - pnpm typecheck 통과
    - 카카오 개발자 콘솔 테스트 앱으로 로그인 플로우 수동 확인

    Refs:
    - #12

    AI:
    - Claude Code로 KakaoStrategy 보일러플레이트 생성
    - 토큰 발급 로직 및 에러 핸들링은 직접 검토 후 수정
    ```

    ## 주의사항

    - 커밋 메시지는 반드시 사용자에게 보여주고 확인을 받은 후 커밋합니다.
    - `--amend`는 이미 push된 커밋에는 사용하지 않도록 경고합니다.
    - 변경이 너무 크면 (10개 파일 이상 또는 500줄 이상) 커밋을 분리할 것을 제안합니다.
    - AI 섹션은 해커톤 공정성 가이드에 따라 솔직하게 작성합니다.