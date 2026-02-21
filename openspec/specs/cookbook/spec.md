# Cookbook Spec

## 개요

"나만의 요리책" 도메인.
사용자가 마음에 든 레시피를 저장하고, 카테고리별로 분류하며, 만족도를 기록하는 개인 레시피 보관함.

---

## 핵심 엔티티

### SavedRecipe

| 필드           | 타입       | 설명                                            |
|----------------|------------|-------------------------------------------------|
| id             | ObjectID   | 저장 레코드 고유 ID                             |
| user_id        | ObjectID   | 소유자                                          |
| recipe_id      | *ObjectID  | 내부 레시피 ID (INTERNAL 레시피인 경우)         |
| recipe_snapshot| RecipeSnap | 저장 시점의 레시피 스냅샷 (외부 링크 포함)      |
| label          | string     | 사용자 정의 카테고리 라벨                       |
| note           | string     | 메모 (최대 500자)                               |
| rating         | *int       | 만족도 1~5 (nil이면 미평가)                     |
| saved_at       | time.Time  | 저장 일시                                       |
| updated_at     | time.Time  | 수정 일시                                       |

### RecipeSnap (저장 시점 스냅샷, 내장 도큐먼트)

| 필드           | 타입     | 설명                              |
|----------------|----------|-----------------------------------|
| title          | string   | 레시피명                          |
| source_url     | string   | 원본 링크                         |
| thumbnail_url  | string   | 썸네일 URL                        |
| source_type    | string   | INTERNAL / EXTERNAL               |
| main_ingredient| string   | 주재료                            |
| category       | string   | 레시피 카테고리                   |

> 레시피 원본이 삭제되더라도 저장된 정보를 유지하기 위해 스냅샷 방식 사용

### Label (계산 집계, 별도 컬렉션 없음)

- 사용자별 label 값의 distinct 목록을 SavedRecipe에서 집계
- 사전 정의 없이 사용자가 자유롭게 생성

---

## 핵심 동작 (Behaviors)

### 레시피 저장
- 추천 화면 또는 레시피 상세에서 "저장" 버튼
- label 미지정 시 기본값 "미분류"
- 동일 recipe (내부 ID 또는 source_url 기준) 중복 저장 방지

### 목록 조회
- 기본 정렬: saved_at 내림차순 (최근 저장 순)
- 라벨 필터링
- 검색: 레시피 제목 prefix 검색

### 레시피 수정
- label, note, rating 수정 가능

### 레시피 삭제
- 단건 삭제
- 라벨별 일괄 삭제 지원

### 라벨 목록 조회
- 해당 사용자의 모든 label 목록 + 각 label의 저장 수 반환
- 예: [{"label": "주말요리", "count": 5}, {"label": "다이어트", "count": 3}]

---

## API 엔드포인트

```
GET    /api/v1/cookbook              저장 레시피 목록 (필터/정렬)
POST   /api/v1/cookbook              레시피 저장
GET    /api/v1/cookbook/:id          저장 레시피 단건 조회
PATCH  /api/v1/cookbook/:id         라벨/메모/평점 수정
DELETE /api/v1/cookbook/:id          단건 삭제
GET    /api/v1/cookbook/labels       사용자 라벨 목록 + 개수
```

---

## 제약 조건

- 사용자당 저장 레시피 최대 500개 (MVP, 이후 조정)
- label 최대 길이: 20자
- note 최대 길이: 500자
- rating: 1 이상 5 이하 정수 (nil이면 미평가 상태 유지)

---

## 의존성

- **auth**: user_id
- **recommendation**: SavedRecipe 생성 시 Recipe 정보 참조
