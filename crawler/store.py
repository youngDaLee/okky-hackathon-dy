import logging
from datetime import datetime, timezone

from pymongo import MongoClient
from pymongo.errors import PyMongoError

from config import MONGO_URI, DB_NAME, COLLECTION_NAME

logger = logging.getLogger(__name__)


def get_collection():
    """Connect to MongoDB and return the recipes collection."""
    client = MongoClient(MONGO_URI)
    db = client[DB_NAME]
    col = db[COLLECTION_NAME]
    return client, col


def ensure_indexes(col) -> None:
    """Create unique partial index on source_url for upsert deduplication.

    Partial filter ensures existing docs with empty source_url don't conflict.
    """
    col.create_index(
        "source_url",
        unique=True,
        partialFilterExpression={"source_url": {"$type": "string", "$gt": ""}},
        name="source_url_unique_nonempty",
    )
    logger.info("Ensured unique partial index on source_url")


def upsert_recipe(col, recipe: dict) -> str:
    """Upsert a recipe document by source_url. Returns 'inserted' or 'updated'."""
    now = datetime.now(timezone.utc)
    filter_doc = {"source_url": recipe["source_url"]}
    update_doc = {"$set": {**recipe, "updated_at": now}, "$setOnInsert": {"created_at": now}}
    try:
        result = col.update_one(filter_doc, update_doc, upsert=True)
        if result.upserted_id:
            return "inserted"
        return "updated"
    except PyMongoError as exc:
        logger.error("Failed to upsert %s: %s", recipe.get("source_url"), exc)
        return "error"
