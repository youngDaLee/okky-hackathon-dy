"""만개의레시피 크롤러 — 인기순 상위 100개를 MongoDB에 저장."""

import logging
import sys

from scraper import collect_recipe_urls, scrape_recipe
from normalizer import normalize_ingredients
from store import get_collection, ensure_indexes, upsert_recipe

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
)
logger = logging.getLogger("crawler")


def run() -> None:
    # 1) MongoDB 연결
    client, col = get_collection()
    ensure_indexes(col)

    # 2) 목록 페이지에서 레시피 URL 수집
    urls = collect_recipe_urls()
    if not urls:
        logger.error("No recipe URLs collected — aborting")
        sys.exit(1)

    logger.info("Starting to scrape %d recipes", len(urls))

    stats = {"inserted": 0, "updated": 0, "skipped": 0, "error": 0}

    # 3) 각 레시피 상세 페이지 크롤링 → 정규화 → 저장
    for i, url in enumerate(urls, 1):
        logger.info("[%d/%d] Scraping %s", i, len(urls), url)

        data = scrape_recipe(url)
        if data is None:
            stats["skipped"] += 1
            continue

        # 재료 정규화
        raw = data.pop("raw_ingredients", [])
        normalized = normalize_ingredients(raw)

        recipe_doc = {
            **data,
            "raw_ingredients": raw,
            "required_ingredients": normalized,
            "optional_ingredients": [],
            "main_ingredient": normalized[0] if normalized else "",
        }

        if not normalized:
            logger.warning("No ingredients after normalisation for %s — storing anyway", url)

        result = upsert_recipe(col, recipe_doc)
        stats[result] = stats.get(result, 0) + 1

    client.close()

    logger.info(
        "Done! inserted=%d, updated=%d, skipped=%d, errors=%d",
        stats["inserted"], stats["updated"], stats["skipped"], stats["error"],
    )


if __name__ == "__main__":
    run()
