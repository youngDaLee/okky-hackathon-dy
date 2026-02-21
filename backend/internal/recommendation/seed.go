package recommendation

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SeedRecipesIfEmpty inserts mock recipes when the collection is empty.
func SeedRecipesIfEmpty(ctx context.Context, col *mongo.Collection) {
	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("warn: recipe seed count check: %v", err)
		return
	}
	if count > 0 {
		return
	}

	now := time.Now()
	docs := []interface{}{
		Recipe{
			Title:               "김치찌개",
			Description:         "돼지고기와 잘 익은 김치로 끓이는 한국 대표 찌개",
			RequiredIngredients: []string{"김치", "돼지고기", "두부", "대파"},
			OptionalIngredients: []string{"양파", "고춧가루", "참기름"},
			MainIngredient:      "김치",
			Category:            "찌개",
			Tags:                []string{"한식", "찌개", "매운맛"},
			CookingTimeMin:      30,
			Difficulty:          "easy",
			CreatedAt:           now,
		},
		Recipe{
			Title:               "된장찌개",
			Description:         "구수한 된장과 각종 채소로 끓이는 찌개",
			RequiredIngredients: []string{"된장", "두부", "감자", "애호박", "대파"},
			OptionalIngredients: []string{"양파", "고추", "버섯"},
			MainIngredient:      "된장",
			Category:            "찌개",
			Tags:                []string{"한식", "찌개", "구수한맛"},
			CookingTimeMin:      25,
			Difficulty:          "easy",
			CreatedAt:           now,
		},
		Recipe{
			Title:               "계란볶음밥",
			Description:         "간단하게 만드는 계란볶음밥",
			RequiredIngredients: []string{"밥", "계란", "대파"},
			OptionalIngredients: []string{"당근", "햄", "참기름", "간장"},
			MainIngredient:      "계란",
			Category:            "볶음",
			Tags:                []string{"한식", "볶음밥", "간단요리"},
			CookingTimeMin:      10,
			Difficulty:          "easy",
			CreatedAt:           now,
		},
		Recipe{
			Title:               "제육볶음",
			Description:         "매콤달콤한 고추장 양념의 돼지고기 볶음",
			RequiredIngredients: []string{"돼지고기", "고추장", "양파", "대파"},
			OptionalIngredients: []string{"당근", "깻잎", "버섯"},
			MainIngredient:      "돼지고기",
			Category:            "볶음",
			Tags:                []string{"한식", "볶음", "매운맛", "밑반찬"},
			CookingTimeMin:      20,
			Difficulty:          "easy",
			CreatedAt:           now,
		},
		Recipe{
			Title:               "소고기미역국",
			Description:         "소고기와 미역으로 끓이는 영양 가득 미역국",
			RequiredIngredients: []string{"소고기", "미역", "간장", "참기름"},
			OptionalIngredients: []string{"대파", "다진마늘"},
			MainIngredient:      "미역",
			Category:            "국",
			Tags:                []string{"한식", "국", "생일", "영양식"},
			CookingTimeMin:      40,
			Difficulty:          "easy",
			CreatedAt:           now,
		},
		Recipe{
			Title:               "닭볶음탕",
			Description:         "매콤하게 조린 닭과 감자 볶음탕",
			RequiredIngredients: []string{"닭", "감자", "당근", "양파", "고추장"},
			OptionalIngredients: []string{"대파", "고춧가루", "떡"},
			MainIngredient:      "닭",
			Category:            "찜/조림",
			Tags:                []string{"한식", "매운맛", "닭요리"},
			CookingTimeMin:      45,
			Difficulty:          "medium",
			CreatedAt:           now,
		},
		Recipe{
			Title:               "참치김치볶음밥",
			Description:         "참치캔과 김치로 빠르게 만드는 볶음밥",
			RequiredIngredients: []string{"밥", "참치캔", "김치"},
			OptionalIngredients: []string{"계란", "대파", "참기름", "김"},
			MainIngredient:      "참치캔",
			Category:            "볶음",
			Tags:                []string{"한식", "볶음밥", "간단요리", "자취요리"},
			CookingTimeMin:      15,
			Difficulty:          "easy",
			CreatedAt:           now,
		},
		Recipe{
			Title:               "순두부찌개",
			Description:         "부드러운 순두부와 해물로 끓이는 얼큰한 찌개",
			RequiredIngredients: []string{"순두부", "계란", "대파", "고춧가루"},
			OptionalIngredients: []string{"바지락", "새우", "양파", "버섯"},
			MainIngredient:      "순두부",
			Category:            "찌개",
			Tags:                []string{"한식", "찌개", "매운맛", "해물"},
			CookingTimeMin:      20,
			Difficulty:          "easy",
			CreatedAt:           now,
		},
		Recipe{
			Title:               "잡채",
			Description:         "당면과 각종 채소를 볶아 만드는 잡채",
			RequiredIngredients: []string{"당면", "시금치", "당근", "양파", "간장"},
			OptionalIngredients: []string{"소고기", "버섯", "피망", "참기름"},
			MainIngredient:      "당면",
			Category:            "볶음",
			Tags:                []string{"한식", "명절", "반찬"},
			CookingTimeMin:      35,
			Difficulty:          "medium",
			CreatedAt:           now,
		},
		Recipe{
			Title:               "감자채볶음",
			Description:         "아삭하게 볶아낸 감자채 반찬",
			RequiredIngredients: []string{"감자", "소금", "식용유"},
			OptionalIngredients: []string{"대파", "당근", "참기름"},
			MainIngredient:      "감자",
			Category:            "반찬",
			Tags:                []string{"한식", "반찬", "간단요리", "밑반찬"},
			CookingTimeMin:      15,
			Difficulty:          "easy",
			CreatedAt:           now,
		},
	}

	result, err := col.InsertMany(ctx, docs)
	if err != nil {
		log.Printf("warn: recipe seed insert: %v", err)
		return
	}
	log.Printf("seeded %d mock recipes", len(result.InsertedIDs))
}