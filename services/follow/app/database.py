from motor.motor_asyncio import AsyncIOMotorClient
from pymongo import ASCENDING
import os

MONGO_URI = os.getenv("MONGO_URI", "mongodb://localhost:27017")
DB_NAME = os.getenv("MONGO_DB", "followdb")

client = AsyncIOMotorClient(MONGO_URI)
db = client[DB_NAME]

# Index za brže pretraživanje
async def init_db():
    await db.follows.create_index([("follower_id", ASCENDING)])
    await db.follows.create_index([("followed_id", ASCENDING)])
