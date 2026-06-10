import os

SECRET_KEY = os.getenv(
    "SUPERSET_SECRET_KEY",
    "superset",
)

FEATURE_FLAGS = {
    "ENABLE_TEMPLATE_PROCESSING": True,
}