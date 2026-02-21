# Fridge Spec

## 개요

사용자의 냉장고 재료 목록 관리 도메인.
재료 CRUD, 유통기한 추적, 소비 우선순위 강조를 담당한다.
재료 등록 경로는 수동 입력 / 영수증 OCR / 냉장고 사진 인식 3가지.

---

## 핵심 엔티티

### Ingredient

| 필드         | 타입       | 설명                                              |
|--------------|------------|---------------------------------------------------|
| id           | ObjectID   | 재료 고유 ID                                      |
| user_id      | ObjectID   | 소유자 (User.id)                                  |
| name         | string     | 재료명 (예: 계란, 당근)                           |
| category     | string     | 카테고리 (enum, 아래 참조)                        |
| quantity     | float64    | 수량                                              |
| unit         | string     | 단위 (예: 개, g, ml, 묶음)                       |
| expiry_date  | *time.Time | 유통기한 (nil 허용 — 유통기한 없는 재료)          |
| source       | string     | 등록 경로: manual / receipt / vision              |
| added_at     | time.Time  | 등록 일시                                         |
| updated_at   | time.Time  | 수정 일시                                         |

### Category (enum)

```
VEGETABLE  채소
FRUIT      과일
MEAT       육류
SEAFOOD    해산물
DAIRY      유제품/달걀
GRAIN      곡류/면류
CONDIMENT  양념/소스
FROZEN     냉동식품
OTHER      기타
```

### ExpiryStatus (계산 필드, DB 저장 안 함)

| 상태      | 조건                        | UI 강조     |
|-----------|-----------------------------|-------------|
| URGENT    | 오늘 포함 D-0 ~ D-2         | 빨간색      |
| SOON      | D-3 ~ D-5                   | 주황색      |
| NORMAL    | D-6 이후                    | 기본색      |
| NO_EXPIRY | expiry_date == nil          | 회색        |

---

## 핵심 동작 (Behaviors)

### 재료 등록
- 수동: name, category, quantity, unit, expiry_date 입력
- Vision 경유: VisionJob 결과를 받아 자동 등록 (vision 도메인에서 호출)
- 동일 name + user_id 재료가 이미 존재하면 수량 합산 여부를 클라이언트에서 선택

### 재료 목록 조회
- 기본 정렬: expiry_date 오름차순 (임박 재료 상단)
- 필터: category, expiry_status
- 검색: name prefix 검색

### 재료 수정
- name, category, quantity, unit, expiry_date 수정 가능

### 재료 삭제
- 단건 삭제
- 요리 완료 후 "사용한 재료 선택 삭제" 일괄 처리 지원

### 유통기한 집계
- GET /fridge/summary: 전체 재료 수, URGENT/SOON 재료 수 반환
- 알림 도메인(notification)에서 이 데이터를 풀링

---

## API 엔드포인트

```
GET    /api/v1/fridge              내 재료 목록 조회 (필터/정렬)
POST   /api/v1/fridge              재료 추가
GET    /api/v1/fridge/:id          재료 단건 조회
PATCH  /api/v1/fridge/:id         재료 수정
DELETE /api/v1/fridge/:id          재료 삭제
DELETE /api/v1/fridge              다수 재료 일괄 삭제 (body: ids[])
GET    /api/v1/fridge/summary      재료 현황 요약
```

---

## 제약 조건

- 재료명: 최대 50자
- 수량: 0 초과 (0이 되면 삭제 유도 또는 자동 삭제 정책 선택 필요 — 구현 시 결정)
- 사용자당 재료 최대 200개 (MVP 기준, 이후 조정 가능)
- expiry_date는 오늘 이전 날짜 입력 가능 (이미 만료된 재료도 등록 가능 — 사용자 주의 메시지 표시)

---

## 의존성

- **auth**: user_id 기반 소유권 검증
- **vision**: VisionJob 완료 시 이 도메인의 재료 등록 API 호출
- **recommendation**: 레시피 매칭 시 이 도메인의 재료 목록 사용
- **notification**: 유통기한 데이터 참조

---

## 검증된 요구사항 (impl-fridge-api)

### Requirement: ExpiryStatus 계산
시스템은 재료의 expiry_date와 오늘 날짜를 비교하여 ExpiryStatus를 SHALL 계산해야 한다.
expiry_date가 nil이면 NO_EXPIRY, D-0~D-2(만료 포함)면 URGENT, D-3~D-5면 SOON, D-6 이후면 NORMAL.

#### Scenario: 오늘 만료 재료
- **WHEN** expiry_date가 오늘 날짜이면
- **THEN** ExpiryStatus는 URGENT를 반환해야 한다

#### Scenario: 이미 만료된 재료
- **WHEN** expiry_date가 오늘보다 이전이면
- **THEN** ExpiryStatus는 URGENT를 반환해야 한다

#### Scenario: D-3 재료
- **WHEN** expiry_date가 오늘로부터 3일 후이면
- **THEN** ExpiryStatus는 SOON을 반환해야 한다

#### Scenario: 유통기한 없는 재료
- **WHEN** expiry_date가 nil이면
- **THEN** ExpiryStatus는 NO_EXPIRY를 반환해야 한다

---

### Requirement: 재료 등록 200개 제한
시스템은 사용자당 재료 최대 200개 제한을 SHALL 적용해야 한다.

#### Scenario: 200개 초과 시도
- **WHEN** 재료가 이미 200개인 사용자가 재료 추가를 요청하면
- **THEN** HTTP 400과 FRIDGE_LIMIT_EXCEEDED 에러를 반환해야 한다

---

### Requirement: 중복 재료 감지
시스템은 동일 user_id + name 재료 중복 등록 시 SHALL 409를 반환해야 한다.

#### Scenario: 동일 이름 재료 재등록
- **WHEN** 이미 존재하는 재료와 같은 name으로 POST /fridge를 호출하면
- **THEN** HTTP 409와 DUPLICATE_INGREDIENT 에러, 기존 재료의 id를 반환해야 한다

---

### Requirement: 재료 목록 유통기한 오름차순 정렬
GET /fridge 응답은 expiry_date 오름차순으로 SHALL 정렬되어야 한다.
expiry_date가 nil인 재료는 항상 마지막에 위치한다.

#### Scenario: 유통기한 정렬
- **WHEN** GET /fridge를 호출하면
- **THEN** expiry_date가 가장 이른 재료가 가장 먼저 반환되어야 한다

#### Scenario: nil 재료 정렬
- **WHEN** expiry_date가 nil인 재료와 있는 재료가 함께 존재하면
- **THEN** nil 재료는 목록의 마지막에 위치해야 한다

---

### Requirement: 재료 소유권 검증
시스템은 다른 사용자의 재료 접근 시도를 SHALL 차단해야 한다.

#### Scenario: 타인 재료 조회 시도
- **WHEN** 다른 사용자의 재료 id로 GET /fridge/:id를 호출하면
- **THEN** HTTP 404를 반환해야 한다 (소유권 노출 방지)

---

### Requirement: 일괄 삭제 silent ignore
DELETE /fridge (bulk)는 존재하지 않는 id가 포함되어도 SHALL 에러 없이 처리해야 한다.

#### Scenario: 존재하지 않는 id 포함 삭제
- **WHEN** 일부 존재하지 않는 id를 포함하여 DELETE /fridge를 호출하면
- **THEN** 존재하는 재료만 삭제하고 HTTP 200과 실제 삭제 수를 반환해야 한다
