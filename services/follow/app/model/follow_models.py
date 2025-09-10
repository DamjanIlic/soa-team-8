from pydantic import BaseModel

class FollowRequest(BaseModel):
    follower_id: str
    following_id: str

class Recommendation(BaseModel):
    user_id: str
    recommended_ids: list[str]
