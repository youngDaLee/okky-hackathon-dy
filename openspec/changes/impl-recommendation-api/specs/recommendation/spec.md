## ADDED Requirements

### Requirement: 3-Tier 매칭 알고리즘
시스템은 사용자 냉장고 재료와 레시피 재료를 비교하여 Tier를 SHALL 분류해야 한다.
재료 이름은 대소문자 무시 및 공백 제거 후 비교한다.

#### Scenario: Tier 1 분류
- **WHEN** 레시피의 모든 required_ingredients가 사용자 재료 목록에 포함되면
- **THEN** 해당 레시피는 Tier 1로 분류되어야 한다

#### Scenario: Tier 2 분류
- **WHEN** 레시피의 main_ingredient가 사용자 재료에 있고 match_score가 0.6 이상이면
- **THEN** 해당 레시피는 Tier 2로 분류되어야 한다

#### Scenario: Tier 3 분류
- **WHEN** match_score가 0.3 이상 0.6 미만이면
- **THEN** 해당 레시피는 Tier 3으로 분류되어야 한다

#### Scenario: 매칭 제외
- **WHEN** match_score가 0.3 미만이면
- **THEN** 해당 레시피는 추천 결과에서 제외되어야 한다

---

### Requirement: urgency_bonus 적용
시스템은 유통기한 URGENT 재료가 사용된 레시피를 SHALL 상단 노출해야 한다.

#### Scenario: urgency_bonus 활성화
- **WHEN** 레시피의 required_ingredients 중 하나라도 사용자의 URGENT(D-0~D-2) 재료와 일치하면
- **THEN** urgency_bonus는 true이고 해당 Tier 내 최상단에 정렬되어야 한다

#### Scenario: urgency_bonus 비활성화
- **WHEN** 레시피의 재료가 URGENT 재료와 겹치지 않으면
- **THEN** urgency_bonus는 false이고 match_score 기준으로 정렬되어야 한다

---

### Requirement: 내부 DB 부족 시 외부 보완
시스템은 내부 레시피가 충분하지 않을 때 외부 검색으로 SHALL 보완해야 한다.

#### Scenario: Tier 1 외부 보완
- **WHEN** Tier 1 내부 결과가 3개 미만이면
- **THEN** 주재료 기반 YouTube 검색으로 최대 5개 외부 레시피를 보완해야 한다

#### Scenario: Tier 2 외부 보완
- **WHEN** Tier 2 내부 결과가 5개 미만이면
- **THEN** 주재료 기반 YouTube 검색으로 최대 3개 외부 레시피를 보완해야 한다

#### Scenario: 외부 API 오류 시 graceful fallback
- **WHEN** 외부 검색 API 호출이 실패하거나 API 키가 없으면
- **THEN** 오류 없이 내부 결과만 반환해야 한다

---

### Requirement: 외부 검색 결과 캐싱
시스템은 동일한 외부 검색 쿼리 결과를 SHALL 1시간 동안 캐시해야 한다.

#### Scenario: 캐시 히트
- **WHEN** 동일 쿼리로 1시간 이내 재요청이 오면
- **THEN** 외부 API를 호출하지 않고 캐시된 결과를 반환해야 한다

#### Scenario: 캐시 만료
- **WHEN** 캐시 만료 후 동일 쿼리가 오면
- **THEN** 외부 API를 새로 호출하여 결과를 갱신해야 한다

---

### Requirement: 오늘의 추천
시스템은 URGENT 재료 기반 Tier 1 추천을 SHALL 오늘의 추천으로 제공해야 한다.

#### Scenario: URGENT 재료 있을 때
- **WHEN** GET /recommendations/today를 호출하고 URGENT 재료가 존재하면
- **THEN** URGENT 재료로 만들 수 있는 Tier 1 레시피를 우선 반환해야 한다

#### Scenario: URGENT 재료 없을 때
- **WHEN** GET /recommendations/today를 호출하고 URGENT 재료가 없으면
- **THEN** 전체 재료 기반 Tier 1/2 추천 결과를 반환해야 한다

---

### Requirement: 빈 냉장고 처리
시스템은 재료가 0개인 사용자의 추천 요청에 SHALL 빈 결과를 반환해야 한다.

#### Scenario: 재료 없음
- **WHEN** 냉장고 재료가 0개인 사용자가 GET /recommendations를 호출하면
- **THEN** HTTP 200과 빈 items 배열을 반환해야 한다
