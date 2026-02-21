## Context

fridge 도메인 스켈레톤이 이미 존재한다 (`internal/fridge/{model,repository,service,handler}.go`).
MongoDB driver v2(`go.mongodb.org/mongo-driver/v2`), Gin 프레임워크 사용.
Auth 미들웨어가 `gin.Context`에 `"userID"` 키로 문자열 user_id를 주입한다고 가정.

## Goals / Non-Goals

**Goals:**
- 7개 API 엔드포인트 완전 동작
- ExpiryStatus 계산 로직 (URGENT/SOON/NORMAL/NO_EXPIRY)
- MongoDB 인덱스 생성 (user_id + expiry_date, text index)
- 사용자당 200개 제한, 중복 재료 409 처리

**Non-Goals:**
- JWT 검증 로직 (미들웨어 미구현, Context에 userID 있다고 가정)
- 단위 테스트

## Decisions

### 1. ExpiryStatus 계산 기준

```
today := time.Now().Truncate(24 * time.Hour)
diff := expiryDate.Truncate(24h).Sub(today).Hours() / 24

diff == 0, 1, 2  → URGENT  (D-0 ~ D-2, 오늘 포함)
diff == 3, 4, 5  → SOON    (D-3 ~ D-5)
diff >= 6        → NORMAL
nil              → NO_EXPIRY
```

만료된 재료(diff < 0)는 URGENT로 처리 — 이미 지난 것도 최우선 소비 유도.

### 2. 목록 정렬 전략

MongoDB sort 우선순위:
1. `expiry_date: 1` (nil last — aggregation으로 처리)
2. `added_at: -1`

nil expiry_date를 마지막에 보내려면 `$addFields`로 정렬용 필드를 추가하거나
application 레벨에서 정렬. **MVP: 애플리케이션 레벨 정렬** 선택 (단순성 우선).

### 3. 중복 재료 처리

`name + user_id` unique index 대신 **서비스 레이어에서 검사** 후 409 반환.
이유: 같은 이름 재료라도 expiry_date가 다른 경우 별도 관리하고 싶을 수 있음 (MVP는 단순 중복 차단).

### 4. 200개 제한

재료 등록 전 `count` 조회 후 제한 초과 시 400 반환.

### 5. GET /fridge/summary normal 카운트

spec에는 `total`, `urgent`, `soon`만 있으나 api.md에는 `normal`, `no_expiry`도 있음.
→ api.md 기준으로 5개 필드 모두 반환.

## Risks / Trade-offs

- **정렬 성능**: 애플리케이션 정렬은 데이터 많을 때 비효율. MVP 200개 제한으로 수용.
- **중복 체크 Race Condition**: 동시 등록 시 중복 통과 가능. MVP 범위에서 허용.
