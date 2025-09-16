from fastapi import FastAPI
from routes.FollowRoutes import router as follow_router

app = FastAPI(title="Follower Microservice")

app.include_router(follow_router)

@app.get("/")
def root():
    return {"message": "Follower Microservice is running"}
