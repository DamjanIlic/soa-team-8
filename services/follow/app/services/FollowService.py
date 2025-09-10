from repo.FollowRepo import FollowRepo

class FollowService:
    def __init__(self):
        self.repo = FollowRepo()

    def follow_user(self, follower_id: str, following_id: str):
        self.repo.create_follow(follower_id, following_id)
        return {"message": f"{follower_id} now follows {following_id}"}

    def unfollow_user(self, follower_id: str, following_id: str):
        self.repo.remove_follow(follower_id, following_id)
        return {"message": f"{follower_id} unfollowed {following_id}"}

    def get_following(self, user_id: str):
        return self.repo.get_following(user_id)

    def get_followers(self, user_id: str):
        return self.repo.get_followers(user_id)

    def get_recommendations(self, user_id: str):
        return self.repo.get_recommendations(user_id)
