import re

# ── Quantity/unit patterns to strip from ingredient text ────────────

# Korean units (longest first to avoid partial match)
_UNITS = (
    "큰술", "작은술", "스푼", "숟가락", "티스푼",
    "컵", "줄기", "장", "쪽", "톨", "봉", "근", "모", "포기", "마리",
    "개", "대", "알", "캔", "팩", "병", "조각", "토막", "움큼", "줌",
    "kg", "g", "mg", "ml", "L", "cc",
    "T", "t", "ts", "Ts",  # tablespoon/teaspoon abbreviations
)

# Vague quantity words
_VAGUE = (
    "약간", "적당량", "조금", "적당히", "한줌", "반개", "반", "소량",
    "적량", "충분히", "넉넉히", "듬뿍", "톡톡", "살짝",
)

# Build regex: number-ish part (including fractions) + optional unit
_NUM = r"[\d\s/.,~\-½⅓⅔¼¾]+"
_UNIT_PAT = "|".join(re.escape(u) for u in _UNITS)
_VAGUE_PAT = "|".join(re.escape(v) for v in _VAGUE)

_STRIP_PATTERNS = [
    # parenthesised notes  e.g. (진간장)
    re.compile(r"\(.*?\)"),
    # number + unit  e.g. "2개", "200g", "1/2큰술"
    re.compile(rf"({_NUM})\s*({_UNIT_PAT})\b"),
    # trailing bare number  e.g. "달걀 2"
    re.compile(rf"\s+{_NUM}\s*$"),
    # vague words  e.g. "약간"
    re.compile(rf"\s*({_VAGUE_PAT})\s*"),
]


def normalize_ingredient(raw: str) -> str:
    """Strip quantities, units, and parenthetical notes from a raw ingredient string.

    Examples:
        "달걀 2개"      → "달걀"
        "소금 약간"      → "소금"
        "양파 1/2개"     → "양파"
        "간장 (진간장)"  → "간장"
    """
    text = raw.strip()
    for pat in _STRIP_PATTERNS:
        text = pat.sub("", text)
    return text.strip()


def normalize_ingredients(raw_list: list[str]) -> list[str]:
    """Normalize a list of raw ingredients, dropping empty results."""
    result: list[str] = []
    for raw in raw_list:
        name = normalize_ingredient(raw)
        if name:
            result.append(name)
    return result
