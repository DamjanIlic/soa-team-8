from repo.FollowRepo import FollowRepo
from typing import List, Dict, Optional
import requests
class FollowService:
    def __init__(self):
        self.repo = FollowRepo()
        self.stakeholder_service_url = "http://stakeholders-service:8080/api/stakeholders"
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


    def get_user_info(self, user_id: str) -> Dict:
        """
        Poziva stakeholder servis da dobije informacije o korisniku.
        Ako korisnik nije pronađen, vraća default placeholder.
        """
        try:
            resp = requests.get(f"{self.stakeholder_service_url}/user/{user_id}")
            if resp.status_code == 200:
                data = resp.json()
                return {
                    "user_id": user_id,
                    "username": data.get("username", "Nepoznat"),
                    "image": data.get("profile_image", "https://upload.wikimedia.org/wikipedia/commons/8/89/Portrait_Placeholder.png"),
                    "motto": data.get("motto", "")
                }
        except Exception as e:
            print(f"Greška prilikom poziva stakeholder servisa: {e}")
        
        # fallback placeholder
        return {
            "user_id": user_id,
            "username": "Nepoznat",
            "image": "https://upload.wikimedia.org/wikipedia/commons/8/89/Portrait_Placeholder.png",
            "motto": ""
        }

    def get_recommendations(self, user_id: str) -> List[Dict]:
        """
        Vrati listu preporuka za korisnika, sa informacijama o korisnicima.
        """
        recommended_ids = self.repo.get_recommendations(user_id)  # npr lista user_id-eva
        response = []
        user_cache = {}  # keš da ne zovemo više puta isti user

        for uid in recommended_ids:
            if uid in user_cache:
                user_info = user_cache[uid]
            else:
                user_info = self.get_user_info(uid)
                user_cache[uid] = user_info
            
            response.append(user_info)
        
        return response