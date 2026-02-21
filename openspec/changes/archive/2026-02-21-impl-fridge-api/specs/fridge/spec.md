## ADDED Requirements

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
