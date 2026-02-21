# Notification Spec

## 개요

유통기한 알림 도메인.
냉장고 재료의 유통기한을 모니터링하고 임박 시점에 사용자에게 알린다.

---

## 알림 트리거 시점

| 타입    | 조건                      | 메시지 예시                                  |
|---------|---------------------------|----------------------------------------------|
| D_3     | expiry_date - today == 3  | "당근이 3일 후 만료돼요. 오늘 사용해볼까요?" |
| D_1     | expiry_date - today == 1  | "계란이 내일 만료됩니다!"                    |
| D_DAY   | expiry_date - today == 0  | "두부 오늘이 마지막이에요. 지금 요리하세요!" |

---

## 핵심 엔티티

### ExpiryAlert

| 필드           | 타입       | 설명                                       |
|----------------|------------|--------------------------------------------|
| id             | ObjectID   | 알림 고유 ID                               |
| user_id        | ObjectID   | 수신자                                     |
| ingredient_id  | ObjectID   | 대상 재료 (Ingredient.id)                  |
| ingredient_name| string     | 재료명 스냅샷 (삭제돼도 유지)              |
| type           | string     | D_3 / D_1 / D_DAY                          |
| scheduled_date | time.Time  | 발송 예정일 (자정 기준)                    |
| sent_at        | *time.Time | 실제 발송 일시 (nil이면 미발송)            |
| is_read        | bool       | 사용자 읽음 여부                           |
| created_at     | time.Time  | 알림 생성 일시                             |

---

## 핵심 동작 (Behaviors)

### 알림 스케줄링
- 재료 등록/수정 시 expiry_date가 있으면 D_3, D_1, D_DAY 알림 자동 생성
- 이미 과거인 시점은 스킵
- 재료 expiry_date 수정 시 기존 미발송 알림 삭제 후 재생성
- 재료 삭제 시 연관 알림 삭제

### 알림 발송 (배치 잡)
- 매일 오전 9시 실행 (cron)
- 오늘 scheduled_date인 미발송 알림 조회 → 발송 처리
- MVP: In-App 알림 (DB 저장 + 클라이언트 폴링)
- 향후: FCM/APNs Push 알림 확장 고려

### 알림 목록 조회
- 사용자의 최근 알림 목록 반환 (최대 50개)
- 읽지 않은 알림 수(badge count) 포함

### 읽음 처리
- 단건 또는 전체 읽음 처리

### 알림 설정 (선택적 MVP 포함 여부 결정 필요)
- D_3 알림 ON/OFF
- 발송 시간대 설정 (기본 오전 9시)

---

## API 엔드포인트

```
GET    /api/v1/notifications           알림 목록 조회
PATCH  /api/v1/notifications/:id/read  단건 읽음 처리
POST   /api/v1/notifications/read-all  전체 읽음 처리
GET    /api/v1/notifications/count     읽지 않은 알림 수
```

---

## 제약 조건

- 동일 재료·동일 타입 알림은 중복 생성 방지 (unique index: ingredient_id + type)
- 발송 실패 시 재시도 없음 (MVP), 이후 dead-letter 큐 도입 고려
- 알림 보존 기간: 30일 (이후 자동 삭제)
- Push 알림은 MVP 범위 외 (In-App만)

---

## 의존성

- **auth**: user_id
- **fridge**: Ingredient 등록/수정/삭제 이벤트, expiry_date
