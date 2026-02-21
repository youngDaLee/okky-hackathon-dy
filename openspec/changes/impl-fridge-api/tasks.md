## 1. Model & Repository

- [x] 1.1 `internal/fridge/repository.go` — MongoDB ingredients 컬렉션 연결 및 인덱스 생성 (user_id+expiry_date 복합, text index)
- [x] 1.2 `internal/fridge/repository.go` — `FindByUserID` 구현 (category, expiry_status, q 필터 지원)
- [x] 1.3 `internal/fridge/repository.go` — `FindByID` 구현 (ObjectID 파싱 포함)
- [x] 1.4 `internal/fridge/repository.go` — `Create` 구현 (InsertOne, 생성된 ID 반환)
- [x] 1.5 `internal/fridge/repository.go` — `Update` 구현 (UpdateOne, $set partial update)
- [x] 1.6 `internal/fridge/repository.go` — `DeleteByID` 구현 (DeleteOne)
- [x] 1.7 `internal/fridge/repository.go` — `DeleteByIDs` 구현 (DeleteMany, silent ignore)
- [x] 1.8 `internal/fridge/repository.go` — `CountByUserID` 구현 (200개 제한용)
- [x] 1.9 `internal/fridge/repository.go` — `FindByNameAndUserID` 구현 (중복 재료 감지용)

## 2. Service — ExpiryStatus 계산

- [x] 2.1 `internal/fridge/service.go` — `CalcExpiryStatus(expiryDate *time.Time) ExpiryStatus` 구현
  - nil → NO_EXPIRY
  - diff <= 2 (만료 포함) → URGENT
  - diff 3~5 → SOON
  - diff >= 6 → NORMAL

## 3. Service — 비즈니스 로직

- [x] 3.1 `internal/fridge/service.go` — `ListIngredients` 구현 (필터 적용, expiry_date 오름차순 정렬, nil last)
- [x] 3.2 `internal/fridge/service.go` — `GetIngredient` 구현 (소유권 검증 — 타인 재료 404 반환)
- [x] 3.3 `internal/fridge/service.go` — `AddIngredient` 구현 (200개 초과 → 400 FRIDGE_LIMIT_EXCEEDED, 동일 이름 → 409 DUPLICATE_INGREDIENT)
- [x] 3.4 `internal/fridge/service.go` — `UpdateIngredient` 구현 (소유권 검증 포함, ExpiryStatus 재계산)
- [x] 3.5 `internal/fridge/service.go` — `DeleteIngredient` 구현 (소유권 검증, 단건)
- [x] 3.6 `internal/fridge/service.go` — `BulkDeleteIngredients` 구현 (silent ignore, 실제 삭제 수 반환)
- [x] 3.7 `internal/fridge/service.go` — `GetSummary` 구현 (total/urgent/soon/normal/no_expiry 카운트)

## 4. Handler — HTTP 핸들러

- [x] 4.1 `internal/fridge/handler.go` — `GET /fridge` 핸들러 (query params 파싱, items+total 응답)
- [x] 4.2 `internal/fridge/handler.go` — `POST /fridge` 핸들러 (request binding, 201 응답, 에러 매핑)
- [x] 4.3 `internal/fridge/handler.go` — `GET /fridge/summary` 핸들러 (summary 5필드 응답)
- [x] 4.4 `internal/fridge/handler.go` — `GET /fridge/:id` 핸들러 (단건 조회, 404 처리)
- [x] 4.5 `internal/fridge/handler.go` — `PATCH /fridge/:id` 핸들러 (partial update binding, 200 응답)
- [x] 4.6 `internal/fridge/handler.go` — `DELETE /fridge/:id` 핸들러 (204 응답)
- [x] 4.7 `internal/fridge/handler.go` — `DELETE /fridge` 핸들러 (body binding ids, 200 deleted_count 응답)

## 5. Router & Config

- [x] 5.1 `internal/server/router.go` — fridge 라우트 7개 등록 (auth 미들웨어 적용)
- [x] 5.2 `pkg/config/config.go` — 환경변수 로드 구현 (MONGODB_URI, JWT_SECRET 등)
