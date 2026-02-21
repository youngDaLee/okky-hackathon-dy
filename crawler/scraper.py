import json
import logging
import re
import time

import requests
from bs4 import BeautifulSoup

from config import BASE_URL, LIST_URL, PAGES_TO_CRAWL, REQUEST_DELAY, MAX_RETRIES

logger = logging.getLogger(__name__)

# ── Cooking-time text → minutes mapping ─────────────────────────────
COOKING_TIME_MAP = {
    "5분 이내": 5,
    "10분 이내": 10,
    "15분 이내": 15,
    "20분 이내": 20,
    "30분 이내": 30,
    "60분 이내": 60,
    "90분 이내": 90,
    "2시간 이상": 120,
}

# ── Difficulty text → normalised value ──────────────────────────────
DIFFICULTY_MAP = {
    "아무나": "EASY",
    "초급": "EASY",
    "중급": "MEDIUM",
    "고급": "HARD",
    "신의경지": "HARD",
}


def _get(url: str) -> requests.Response | None:
    """GET with retries and delay."""
    for attempt in range(1, MAX_RETRIES + 1):
        try:
            resp = requests.get(url, timeout=10, headers={
                "User-Agent": "FridgeMaster-Crawler/1.0",
            })
            resp.raise_for_status()
            time.sleep(REQUEST_DELAY)
            return resp
        except requests.RequestException as exc:
            logger.warning("Attempt %d failed for %s: %s", attempt, url, exc)
            if attempt < MAX_RETRIES:
                time.sleep(REQUEST_DELAY * attempt)
    logger.error("All %d attempts failed for %s", MAX_RETRIES, url)
    return None


# ── Step 1: collect recipe URLs from list pages ────────────────────
def collect_recipe_urls() -> list[str]:
    """Crawl list pages (order=reco) and return up to ~100 recipe detail URLs."""
    urls: list[str] = []
    for page in range(1, PAGES_TO_CRAWL + 1):
        logger.info("Fetching list page %d/%d", page, PAGES_TO_CRAWL)
        resp = _get(f"{LIST_URL}?order=reco&page={page}")
        if resp is None:
            continue
        soup = BeautifulSoup(resp.text, "html.parser")
        for a_tag in soup.select("a.common_sp_link"):
            href = a_tag.get("href", "")
            if re.match(r"^/recipe/\d+$", href):
                urls.append(f"{BASE_URL}{href}")
    unique = list(dict.fromkeys(urls))  # deduplicate, keep order
    logger.info("Collected %d unique recipe URLs", len(unique))
    return unique[:100]


# ── Step 2: extract recipe data from detail page ───────────────────
def _extract_jsonld(soup: BeautifulSoup) -> dict | None:
    """Find Schema.org Recipe JSON-LD in the page."""
    for script in soup.find_all("script", type="application/ld+json"):
        try:
            data = json.loads(script.string or "")
            if isinstance(data, dict) and data.get("@type") == "Recipe":
                return data
            if isinstance(data, list):
                for item in data:
                    if isinstance(item, dict) and item.get("@type") == "Recipe":
                        return item
        except json.JSONDecodeError:
            continue
    return None


def _parse_cooking_time(text: str) -> int:
    """Convert cooking-time text to minutes."""
    text = text.strip()
    if text in COOKING_TIME_MAP:
        return COOKING_TIME_MAP[text]
    m = re.search(r"(\d+)\s*분", text)
    if m:
        return int(m.group(1))
    m = re.search(r"(\d+)\s*시간", text)
    if m:
        return int(m.group(1)) * 60
    return 0


def _parse_difficulty(text: str) -> str:
    return DIFFICULTY_MAP.get(text.strip(), "EASY")


def _extract_meta(soup: BeautifulSoup) -> tuple[int, str, str]:
    """Extract cooking_time, difficulty, and category from HTML meta area."""
    cooking_time = 0
    difficulty = "EASY"
    category = "기타"

    info_spans = soup.select(".view2_summary_info span")
    for span in info_spans:
        text = span.get_text(strip=True)
        if "분" in text or "시간" in text:
            cooking_time = _parse_cooking_time(text)
        elif text in DIFFICULTY_MAP:
            difficulty = _parse_difficulty(text)

    cat_tag = soup.select_one(".view2_summary_info2")
    if cat_tag:
        cat_text = cat_tag.get_text(strip=True)
        if cat_text:
            category = cat_text

    return cooking_time, difficulty, category


def _extract_tags(soup: BeautifulSoup) -> list[str]:
    """Extract hashtags from the page."""
    tags: list[str] = []
    for a_tag in soup.select("a.view_tag"):
        tag = a_tag.get_text(strip=True).lstrip("#")
        if tag:
            tags.append(tag)
    return tags


def scrape_recipe(url: str) -> dict | None:
    """Scrape a single recipe detail page and return structured data."""
    resp = _get(url)
    if resp is None:
        return None

    soup = BeautifulSoup(resp.text, "html.parser")
    jsonld = _extract_jsonld(soup)

    if jsonld is None:
        logger.warning("No JSON-LD found for %s — skipping", url)
        return None

    title = jsonld.get("name", "").strip()
    description = jsonld.get("description", "").strip()
    raw_ingredients = [ing.strip() for ing in jsonld.get("recipeIngredient", []) if ing.strip()]
    thumbnail = jsonld.get("image", "")
    if isinstance(thumbnail, list):
        thumbnail = thumbnail[0] if thumbnail else ""

    cooking_time, difficulty, category = _extract_meta(soup)
    tags = _extract_tags(soup)

    if not title:
        logger.warning("Empty title for %s — skipping", url)
        return None

    return {
        "title": title,
        "description": description,
        "raw_ingredients": raw_ingredients,
        "source_url": url,
        "source_type": "10000recipe",
        "thumbnail_url": thumbnail,
        "cooking_time_min": cooking_time,
        "difficulty": difficulty,
        "category": category,
        "tags": tags,
    }
