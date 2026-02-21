import os

MONGO_URI = os.getenv("MONGODB_URI", "mongodb://localhost:27017")
DB_NAME = os.getenv("MONGODB_DB", "OKKY")
COLLECTION_NAME = "recipes"

BASE_URL = "https://www.10000recipe.com"
LIST_URL = f"{BASE_URL}/recipe/list.html"

PAGES_TO_CRAWL = 4        # ~25 recipes/page → ~100 total
REQUEST_DELAY = 1.5        # seconds between requests
MAX_RETRIES = 3
