package fridge

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Sentinel errors for HTTP error mapping
var (
	ErrNotFound         = errors.New("INGREDIENT_NOT_FOUND")
	ErrLimitExceeded    = errors.New("FRIDGE_LIMIT_EXCEEDED")
	ErrDuplicate        = errors.New("DUPLICATE_INGREDIENT")
)

// DuplicateError wraps ErrDuplicate with the existing ingredient's ID.
type DuplicateError struct {
	ExistingID string
}

func (e *DuplicateError) Error() string { return ErrDuplicate.Error() }

const maxIngredients = 200

// Task 2.1 — CalcExpiryStatus
func CalcExpiryStatus(expiryDate *time.Time) string {
	if expiryDate == nil {
		return ExpiryStatusNoExpiry
	}
	today := time.Now().UTC().Truncate(24 * time.Hour)
	exp := expiryDate.UTC().Truncate(24 * time.Hour)
	diff := int(exp.Sub(today).Hours() / 24)

	switch {
	case diff <= 2: // includes negative (already expired)
		return ExpiryStatusUrgent
	case diff <= 5:
		return ExpiryStatusSoon
	default:
		return ExpiryStatusNormal
	}
}

func toResponse(i *Ingredient) *IngredientResponse {
	return &IngredientResponse{
		Ingredient:   *i,
		ExpiryStatus: CalcExpiryStatus(i.ExpiryDate),
	}
}

type FridgeService interface {
	ListIngredients(ctx context.Context, userID string, filter ListFilter) ([]*IngredientResponse, int, error)
	GetIngredient(ctx context.Context, id, userID string) (*IngredientResponse, error)
	AddIngredient(ctx context.Context, userID string, req CreateIngredientReq) (*IngredientResponse, error)
	UpdateIngredient(ctx context.Context, id, userID string, req UpdateIngredientReq) (*IngredientResponse, error)
	RemoveIngredient(ctx context.Context, id, userID string) error
	BulkRemove(ctx context.Context, ids []string, userID string) (int64, error)
	GetSummary(ctx context.Context, userID string) (*FridgeSummary, error)
}

type fridgeService struct {
	repo IngredientRepository
}

func NewFridgeService(repo IngredientRepository) FridgeService {
	return &fridgeService{repo: repo}
}

// Task 3.1 — ListIngredients with sort (expiry_date ASC, nil last)
func (s *fridgeService) ListIngredients(ctx context.Context, userID string, filter ListFilter) ([]*IngredientResponse, int, error) {
	uid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, ErrNotFound
	}

	items, err := s.repo.FindAllByUserID(ctx, uid, filter)
	if err != nil {
		return nil, 0, err
	}

	// Application-level sort: expiry_date ASC, nil last
	withExpiry := make([]*Ingredient, 0, len(items))
	noExpiry := make([]*Ingredient, 0)
	for _, it := range items {
		if it.ExpiryDate == nil {
			noExpiry = append(noExpiry, it)
		} else {
			withExpiry = append(withExpiry, it)
		}
	}
	// sort withExpiry by ExpiryDate ASC
	for i := 0; i < len(withExpiry); i++ {
		for j := i + 1; j < len(withExpiry); j++ {
			if withExpiry[j].ExpiryDate.Before(*withExpiry[i].ExpiryDate) {
				withExpiry[i], withExpiry[j] = withExpiry[j], withExpiry[i]
			}
		}
	}
	sorted := append(withExpiry, noExpiry...)

	// Apply expiry_status filter after CalcExpiryStatus (can't filter in DB without aggregation)
	var resp []*IngredientResponse
	for _, it := range sorted {
		r := toResponse(it)
		if filter.ExpiryStatus != "" && r.ExpiryStatus != filter.ExpiryStatus {
			continue
		}
		resp = append(resp, r)
	}

	return resp, len(resp), nil
}

// Task 3.2 — GetIngredient (ownership: 404 if not found or not owned)
func (s *fridgeService) GetIngredient(ctx context.Context, id, userID string) (*IngredientResponse, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrNotFound
	}
	uid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrNotFound
	}
	item, err := s.repo.FindByID(ctx, oid, uid)
	if err != nil || item == nil {
		return nil, ErrNotFound
	}
	return toResponse(item), nil
}

// Task 3.3 — AddIngredient (200 limit + duplicate check)
func (s *fridgeService) AddIngredient(ctx context.Context, userID string, req CreateIngredientReq) (*IngredientResponse, error) {
	uid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrNotFound
	}

	count, err := s.repo.CountByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if count >= maxIngredients {
		return nil, ErrLimitExceeded
	}

	existing, err := s.repo.FindByNameAndUserID(ctx, req.Name, uid)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, &DuplicateError{ExistingID: existing.ID.Hex()}
	}

	now := time.Now().UTC()
	ingredient := &Ingredient{
		UserID:     uid,
		Name:       req.Name,
		Category:   req.Category,
		Quantity:   req.Quantity,
		Unit:       req.Unit,
		ExpiryDate: req.ExpiryDate,
		Source:     SourceManual,
		AddedAt:    now,
		UpdatedAt:  now,
	}

	if err := s.repo.Create(ctx, ingredient); err != nil {
		return nil, err
	}
	return toResponse(ingredient), nil
}

// Task 3.4 — UpdateIngredient
func (s *fridgeService) UpdateIngredient(ctx context.Context, id, userID string, req UpdateIngredientReq) (*IngredientResponse, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrNotFound
	}
	uid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrNotFound
	}

	fields := bson.M{"updated_at": time.Now().UTC()}
	if req.Name != nil {
		fields["name"] = *req.Name
	}
	if req.Category != nil {
		fields["category"] = *req.Category
	}
	if req.Quantity != nil {
		fields["quantity"] = *req.Quantity
	}
	if req.Unit != nil {
		fields["unit"] = *req.Unit
	}
	if req.ExpiryDate != nil {
		fields["expiry_date"] = *req.ExpiryDate
	}

	updated, err := s.repo.Update(ctx, oid, uid, fields)
	if err != nil || updated == nil {
		return nil, ErrNotFound
	}
	return toResponse(updated), nil
}

// Task 3.5 — RemoveIngredient (ownership check)
func (s *fridgeService) RemoveIngredient(ctx context.Context, id, userID string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ErrNotFound
	}
	uid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return ErrNotFound
	}
	// verify ownership first
	item, err := s.repo.FindByID(ctx, oid, uid)
	if err != nil || item == nil {
		return ErrNotFound
	}
	return s.repo.Delete(ctx, oid, uid)
}

// Task 3.6 — BulkRemove (silent ignore)
func (s *fridgeService) BulkRemove(ctx context.Context, ids []string, userID string) (int64, error) {
	uid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return 0, ErrNotFound
	}
	oids := make([]bson.ObjectID, 0, len(ids))
	for _, id := range ids {
		oid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			continue // silent ignore invalid IDs
		}
		oids = append(oids, oid)
	}
	if len(oids) == 0 {
		return 0, nil
	}
	return s.repo.BulkDelete(ctx, oids, uid)
}

// Task 3.7 — GetSummary (5 fields)
func (s *fridgeService) GetSummary(ctx context.Context, userID string) (*FridgeSummary, error) {
	uid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, ErrNotFound
	}

	items, err := s.repo.FindAllByUserID(ctx, uid, ListFilter{})
	if err != nil {
		return nil, err
	}

	summary := &FridgeSummary{Total: len(items)}
	for _, it := range items {
		switch CalcExpiryStatus(it.ExpiryDate) {
		case ExpiryStatusUrgent:
			summary.Urgent++
		case ExpiryStatusSoon:
			summary.Soon++
		case ExpiryStatusNormal:
			summary.Normal++
		case ExpiryStatusNoExpiry:
			summary.NoExpiry++
		}
	}
	return summary, nil
}
