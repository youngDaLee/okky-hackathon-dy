---
name: create-issue
description: GitHub 이슈를 생성합니다. 해커톤 협업을 위해 작업을 이슈로 분리하고 담당자를 지정할 때 사용합니다. Use when creating tasks, breaking down features, or managing hackathon work items.
---

# GitHub 이슈 생성 Skill

바이브코딩 해커톤에서 팀 협업을 위한 GitHub 이슈를 생성합니다.

## 사전 확인

1. 현재 디렉토리가 git 저장소인지 확인합니다.
2. `gh` CLI가 설치되어 있고 인증되어 있는지 확인합니다.
- 없으면 사용자에게 `gh auth login` 안내를 합니다.
3. 현재 리포지토리의 remote origin URL을 확인하여 GitHub 리포지토리 정보를 파악합니다.

## 입력 처리

사용자가 `$ARGUMENTS`로 이슈 내용을 전달합니다.

- **인자가 있는 경우**: 해당 내용을 분석하여 이슈 제목과 본문을 구성합니다.
- **인자가 없는 경우**: 사용자에게 다음을 질문합니다.
1. 어떤 작업인가요? (기능 구현, 버그 수정, 문서 작성, 설정/인프라, 디자인 등)
2. 작업 내용을 간단히 설명해주세요.
3. 담당자가 있나요? (GitHub username)
4. 우선순위는? (high, medium, low)

## 이슈 생성 규칙

### 제목 형식

카테고리 prefix를 붙여 한눈에 파악할 수 있게 합니다:

```
[feat] 로그인 페이지 UI 구현
[fix] 회원가입 시 이메일 중복 체크 오류
[docs] README에 프로젝트 설명 추가
[infra] Docker Compose 설정 및 배포 환경 구성
[design] 메인 화면 와이어프레임 작성
[refactor] API 응답 구조 통일
[test] 결제 모듈 단위 테스트 추가
```

### 본문 템플릿

```markdown
## 📋 작업 설명
<!-- 이 작업이 무엇인지 명확하게 기술 -->

## ✅ 완료 조건 (Definition of Done)
- [ ] 체크리스트 항목 1
- [ ] 체크리스트 항목 2
- [ ] 체크리스트 항목 3

## 🔗 관련 사항
- 관련 이슈: (있으면 #이슈번호)
- 참고 자료: (있으면 링크)

## 📌 비고
- 우선순위: `high` / `medium` / `low`
- 예상 소요: (해커톤 내 예상 시간)
```

### 라벨 자동 지정

작업 유형에 따라 라벨을 자동으로 지정합니다. 리포지토리에 라벨이 없으면 생성합니다:

| 카테고리 | 라벨 | 색상 |
|---------|------|------|
| feat | `feature` | `#0E8A16` |
| fix | `bug` | `#D73A4A` |
| docs | `documentation` | `#0075CA` |
| infra | `infra` | `#E4E669` |
| design | `design` | `#C5DEF5` |
| refactor | `refactor` | `#FEF2C0` |
| test | `test` | `#BFD4F2` |

우선순위 라벨도 함께 지정합니다:

| 우선순위 | 라벨 | 색상 |
|---------|------|------|
| high | `priority: high` | `#B60205` |
| medium | `priority: medium` | `#FBCA04` |
| low | `priority: low` | `#0E8A16` |

## 실행 절차

1. 이슈 제목과 본문을 구성합니다.
2. 필요한 라벨이 리포지토리에 있는지 확인하고, 없으면 `gh label create`로 생성합니다.
3. `gh issue create` 명령어로 이슈를 생성합니다:

```bash
gh issue create \
--title "[카테고리] 이슈 제목" \
--body "본문 내용" \
--label "라벨1,라벨2" \
--assignee "담당자"
```

4. 생성된 이슈 URL을 사용자에게 보여줍니다.

## 일괄 생성 모드

사용자가 여러 작업을 한 번에 전달하면 (예: "로그인, 회원가입, 마이페이지 기능을 각각 이슈로 만들어줘"), 각 작업을 개별 이슈로 분리하여 순차적으로 생성합니다.

## 예시

### 단일 이슈 생성
```
/create-issue 카카오 소셜 로그인 기능 구현
```

### 상세 지정
```
/create-issue [feat] 카카오 소셜 로그인 구현 - 담당자: @dayoung - 우선순위: high
```

### 일괄 생성
```
/create-issue 다음 작업들을 이슈로 만들어줘: 1) 로그인 페이지 2) 회원가입 API 3) DB 스키마 설계
```

## 주의사항

- 해커톤 특성상 이슈는 간결하고 실행 가능하게 작성합니다.
- 하나의 이슈는 한 사람이 해커톤 시간 내에 완료할 수 있는 크기로 분리합니다.
- 이슈 간 의존 관계가 있으면 본문의 "관련 사항"에 명시합니다.