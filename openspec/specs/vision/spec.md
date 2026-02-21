# Vision Spec

## 개요

AI 이미지 인식 도메인.
Vertex AI를 통해 사용자가 촬영한 영수증/냉장고 사진에서 식재료를 자동 추출한다.
추출 결과는 사용자 확인 후 fridge 도메인으로 전달된다.

---

## 처리 유형

| 타입     | 입력           | Vertex AI 모델          | 출력                     |
|----------|----------------|-------------------------|--------------------------|
| RECEIPT  | 영수증 사진    | Document AI / OCR       | 품목명, 수량 추출        |
| FRIDGE   | 냉장고 사진    | Vision Object Detection | 식재료 객체 탐지         |

---

## 핵심 엔티티

### VisionJob

| 필드              | 타입           | 설명                                          |
|-------------------|----------------|-----------------------------------------------|
| id                | ObjectID       | 작업 고유 ID                                  |
| user_id           | ObjectID       | 요청자                                        |
| type              | string         | RECEIPT / FRIDGE                              |
| image_url         | string         | 업로드된 이미지 GCS URL                       |
| status            | string         | PENDING / PROCESSING / DONE / FAILED          |
| raw_result        | object         | Vertex AI 원본 응답 (디버깅용)                |
| extracted         | []ExtractedItem| 파싱된 결과                                   |
| error_message     | *string        | 실패 사유 (nil이면 성공)                      |
| created_at        | time.Time      | 요청 일시                                     |
| completed_at      | *time.Time     | 완료 일시                                     |

### ExtractedItem

| 필드             | 타입      | 설명                                           |
|------------------|-----------|------------------------------------------------|
| name             | string    | 추출된 재료명 (정규화 전)                      |
| normalized_name  | string    | 정규화 이름 (예: "당근 3개" → "당근")         |
| quantity         | *float64  | 추출 수량 (nil이면 미확인)                     |
| unit             | *string   | 추출 단위 (nil이면 미확인)                     |
| confidence       | float64   | 신뢰도 0.0 ~ 1.0                              |
| selected         | bool      | 사용자가 최종 선택 여부 (확인 단계에서 설정)   |

---

## 핵심 동작 (Behaviors)

### 이미지 업로드 및 작업 생성
1. 클라이언트에서 multipart/form-data로 이미지 전송
2. 서버: GCS에 이미지 업로드 → VisionJob 생성 (status: PENDING)
3. 백그라운드 워커에서 Vertex AI 호출 (비동기)
4. 완료 시 status → DONE, extracted 필드 채움

### 결과 폴링 / 웹소켓
- 클라이언트는 `GET /vision/jobs/:id` 로 상태 폴링
- (MVP) 폴링 방식 사용, 이후 웹소켓으로 전환 고려

### 사용자 확인 단계
- DONE 상태의 결과를 사용자에게 표시
- 사용자가 각 ExtractedItem의 selected 여부 및 수량/이름 수정
- 확인 완료 → POST /fridge (selected items) 호출로 냉장고에 등록

### 신뢰도 필터링
- confidence < 0.5 인 항목은 기본 selected: false (사용자가 수동으로 체크 필요)
- confidence >= 0.5 인 항목은 기본 selected: true

---

## API 엔드포인트

```
POST   /api/v1/vision/jobs              이미지 업로드 및 작업 생성
GET    /api/v1/vision/jobs/:id          작업 상태 및 결과 조회
POST   /api/v1/vision/jobs/:id/confirm  결과 확인 및 냉장고 등록 요청
```

---

## Vertex AI 통합

- **RECEIPT 타입**: Vertex AI Document AI 또는 Cloud Vision API의 TEXT_DETECTION
  - 영수증 품목 라인 파싱 → 정규화
- **FRIDGE 타입**: Cloud Vision API의 OBJECT_LOCALIZATION
  - 탐지된 객체 중 식품 관련 label만 필터링

### 재료명 정규화 규칙
- 영수증 품목명에서 가격·수량 제거 (예: "계란10구 5,400" → "계란")
- 브랜드명 제거 시도 (best-effort)
- 한글 표준 식품명으로 매핑 (내부 사전 테이블 사용)

---

## 제약 조건

- 이미지 파일 크기: 최대 10MB
- 지원 포맷: JPEG, PNG, WEBP
- 작업 보존 기간: 24시간 (이후 GCS 이미지 및 VisionJob 삭제)
- 처리 타임아웃: 30초 (초과 시 status → FAILED)

---

## 의존성

- **auth**: user_id 인증
- **fridge**: confirm 단계에서 재료 등록 API 호출
- GCS (Google Cloud Storage): 이미지 저장
- Vertex AI / Cloud Vision API: AI 분석
