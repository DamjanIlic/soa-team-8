from config import driver

class FollowRepo:
    def __init__(self):
        self.driver = driver

    def create_follow(self, follower_id: str, following_id: str):
        with self.driver.session() as session:
            session.run(
                "MERGE (a:User {id: $follower_id}) "
                "MERGE (b:User {id: $following_id}) "
                "MERGE (a)-[:FOLLOWS]->(b)",
                follower_id=follower_id, following_id=following_id
            )

    def remove_follow(self, follower_id: str, following_id: str):
        with self.driver.session() as session:
            session.run(
                "MATCH (a:User {id: $follower_id})-[f:FOLLOWS]->(b:User {id: $following_id}) "
                "DELETE f",
                follower_id=follower_id, following_id=following_id
            )

    def get_following(self, user_id: str):
        with self.driver.session() as session:
            result = session.run(
                "MATCH (a:User {id: $user_id})-[:FOLLOWS]->(b:User) "
                "RETURN b.id AS id",
                user_id=user_id
            )
            return [record["id"] for record in result]

    def get_followers(self, user_id: str):
        with self.driver.session() as session:
            result = session.run(
                "MATCH (b:User)-[:FOLLOWS]->(a:User {id: $user_id}) "
                "RETURN b.id AS id",
                user_id=user_id
            )
            return [record["id"] for record in result]

    def get_recommendations(self, user_id: str):
        with self.driver.session() as session:
            result = session.run(
                """
                MATCH (me:User {id: $user_id})-[:FOLLOWS]->(:User)-[:FOLLOWS]->(rec:User)
                WHERE NOT (me)-[:FOLLOWS]->(rec) AND rec.id <> $user_id
                RETURN DISTINCT rec.id AS id
                """,
                user_id=user_id
            )
            return [record["id"] for record in result]
