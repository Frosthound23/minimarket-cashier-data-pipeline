from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.routers.analytics_router import router as analytics_router

app = FastAPI(
    title="Minimarket Analytics API",
    description="Analytics API for multi-tenant minimarket sales data",
    version="1.0.0",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(analytics_router)


@app.get("/")
def root():
    return {
        "message": "Minimarket Analytics API",
        "docs": "/docs",
    }