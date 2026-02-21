## ADDED Requirements

### Requirement: 목록 페이지에서 레시피 URL 수집
크롤러는 만개의레시피 목록 페이지(`/recipe/list.html?order=reco`)에서
인기순 상위 100개 레시피의 상세 페이지 URL을 수집하여야(SHALL) 한다.
페이지당 약 25개이므로 4페이지를 순회한다.

#### Scenario: 정상적으로 100개 URL 수집
- **WHEN** 크롤러가 `order=reco`로 page 1~4를 요청하면
- **THEN** 최대 100개의 `/recipe/{id}` 형태 URL 목록을 반환한다

#### Scenario: 네트워크 오류 시 재시도
- **WHEN** 목록 페이지 요청이 실패하면
- **THEN** 최대 3회 재시도하고, 여전히 실패하면 에러를 로깅하고 수집된 URL만으로 진행한다

### Requirement: 상세 페이지에서 레시피 데이터 추출
크롤러는 각 레시피 상세 페이지에서 Schema.org JSON-LD를 파싱하여
레시피 정보를 추출하여야(SHALL) 한다.

#### Scenario: JSON-LD가 존재하는 레시피
- **WHEN** 상세 페이지에 `@type: Recipe` JSON-LD가 있으면
- **THEN** title, description, recipeIngredient, image, datePublished를 추출한다

#### Scenario: JSON-LD가 없는 레시피
- **WHEN** 상세 페이지에 JSON-LD가 없으면
- **THEN** 해당 레시피를 skip하고 warning을 로깅한다

### Requirement: 재료명 정규화
크롤러는 원문 재료에서 수량과 단위를 제거하여 정규화된 재료명을 생성하여야(SHALL) 한다.
원문 재료는 `raw_ingredients`에, 정규화된 이름은 `required_ingredients`에 저장한다.

#### Scenario: 수량+단위 패턴 제거
- **WHEN** 원문 재료가 "달걀 2개"이면
- **THEN** 정규화 결과는 "달걀"이다

#### Scenario: 수량 표현 단어 제거
- **WHEN** 원문 재료가 "소금 약간"이면
- **THEN** 정규화 결과는 "소금"이다

#### Scenario: 분수 표현 처리
- **WHEN** 원문 재료가 "양파 1/2개"이면
- **THEN** 정규화 결과는 "양파"이다

#### Scenario: 괄호 부가설명 제거
- **WHEN** 원문 재료가 "간장 (진간장)"이면
- **THEN** 정규화 결과는 "간장"이다

### Requirement: MainIngredient 결정
크롤러는 정규화된 재료 목록의 첫 번째 항목을 `main_ingredient`로 설정하여야(SHALL) 한다.

#### Scenario: 재료가 1개 이상인 경우
- **WHEN** 정규화된 재료 목록이 ["달걀", "대파", "간장"]이면
- **THEN** `main_ingredient`는 "달걀"이다

#### Scenario: 재료가 없는 경우
- **WHEN** 재료 목록이 비어있으면
- **THEN** `main_ingredient`는 빈 문자열이고, 해당 레시피를 warning 로깅한다

### Requirement: 조리시간 파싱
크롤러는 만개의레시피의 조리시간 텍스트를 분 단위 정수로 변환하여야(SHALL) 한다.

#### Scenario: 분 단위 텍스트
- **WHEN** 조리시간이 "5분 이내"이면
- **THEN** `cooking_time_min`은 5이다

#### Scenario: 시간 단위 텍스트
- **WHEN** 조리시간이 "2시간 이상"이면
- **THEN** `cooking_time_min`은 120이다

#### Scenario: 조리시간 정보 없음
- **WHEN** 조리시간 정보가 없으면
- **THEN** `cooking_time_min`은 0이다

### Requirement: MongoDB upsert
크롤러는 `source_url` 기준으로 MongoDB에 upsert하여야(SHALL) 한다.
신규 레시피는 insert, 기존 레시피는 update로 처리한다.

#### Scenario: 신규 레시피 삽입
- **WHEN** `source_url`이 DB에 존재하지 않는 레시피를 저장하면
- **THEN** 새 document가 insert된다

#### Scenario: 기존 레시피 갱신
- **WHEN** `source_url`이 이미 DB에 존재하는 레시피를 저장하면
- **THEN** 기존 document가 update된다

### Requirement: Rate Limiting
크롤러는 요청 간 1~2초 딜레이를 적용하여야(SHALL) 한다.

#### Scenario: 연속 요청 간격
- **WHEN** 크롤러가 상세 페이지를 연속으로 요청하면
- **THEN** 각 요청 사이에 최소 1초의 간격이 있다
