## MODIFIED Requirements

### Requirement: Recipe (내부 DB)

Recipe 엔티티에 `raw_ingredients` 필드를 추가한다.
크롤링된 원문 재료(수량 포함)를 저장하여 프론트엔드에서 레시피 상세 표시 시 활용한다.

| 필드                  | 타입       | 설명                                          |
|-----------------------|------------|-----------------------------------------------|
| id                    | ObjectID   | 레시피 고유 ID                                |
| title                 | string     | 레시피명                                      |
| description           | string     | 간단 설명                                     |
| required_ingredients  | []string   | 필수 재료명 목록 (정규화된 이름)              |
| optional_ingredients  | []string   | 선택 재료명 목록                              |
| **raw_ingredients**   | **[]string** | **원문 재료 목록 (수량/단위 포함, 표시용)**  |
| main_ingredient       | string     | 대표 주재료 1개                               |
| source_type           | string     | INTERNAL / EXTERNAL                           |
| source_url            | string     | 원본 링크 (YouTube, 블로그 등)                |
| thumbnail_url         | string     | 썸네일 이미지 URL                             |
| category              | string     | 한식 / 중식 / 일식 / 양식 / 간식 / 기타      |
| tags                  | []string   | 검색 태그 (예: 초간단, 10분완성, 다이어트)    |
| cooking_time_min      | int        | 조리 시간 (분)                                |
| difficulty            | string     | EASY / MEDIUM / HARD                          |
| created_at            | time.Time  |                                               |

#### Scenario: Recipe에 raw_ingredients 포함
- **WHEN** 크롤링된 레시피를 저장할 때
- **THEN** `raw_ingredients`에 원문 재료 목록이 저장되고, `required_ingredients`에는 정규화된 재료명이 저장된다

#### Scenario: 기존 레시피 하위 호환
- **WHEN** `raw_ingredients` 필드가 없는 기존 레시피를 조회할 때
- **THEN** `raw_ingredients`는 null 또는 빈 배열로 반환되며 에러가 발생하지 않는다
