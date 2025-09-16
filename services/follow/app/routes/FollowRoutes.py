from fastapi import APIRouter
from model.follow_models import FollowRequest
from services.FollowService import FollowService

router = APIRouter(prefix="/api/follow", tags=["Follow"])
service = FollowService()

@router.post("/add")
def follow_user(request: FollowRequest):
    return service.follow_user(request.follower_id, request.following_id)

@router.post("/remove")
def unfollow_user(request: FollowRequest):
    return service.unfollow_user(request.follower_id, request.following_id)

@router.get("/following/{user_id}")
def get_following(user_id: str):
    return {"following": service.get_following(user_id)}

@router.get("/followers/{user_id}")
def get_followers(user_id: str):
    return {"followers": service.get_followers(user_id)}

@router.get("/recommendations/{user_id}")
def recommendations(user_id: str):
    return {"recommendations": service.get_recommendations(user_id)}
