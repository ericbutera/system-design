import uvicorn
from fastapi import Depends, FastAPI, Request
from models import SessionLocal
from schema import schema
from strawberry.fastapi import GraphQLRouter


async def get_db():
    """Dependency that provides a database session per request and ensures cleanup."""
    db = SessionLocal()
    try:
        yield db  # Session is available for the request
    finally:
        db.close()  # Ensure session is closed after request


async def get_context(request: Request, db=Depends(get_db)):
    """GraphQL context getter that includes DB session and user info."""
    user_id = request.headers.get("user-id")
    user = {"user_id": user_id} if user_id else None
    return {"user": user, "db": db}


def user_middleware(next, root, info, **kwargs):
    request = info.context.get("request")
    if request:
        user_id = request.headers.get("user-id")
        if user_id:
            # Simulate a user lookup based on user_id
            info.context["user"] = {"user_id": user_id}
        else:
            info.context["user"] = None
    return next(root, info, **kwargs)


app = FastAPI()
graphql_app = GraphQLRouter(schema, context_getter=get_context)
app.include_router(graphql_app, prefix="/graphql")


@app.get("/")
def health_check():
    return {"status": "ok"}


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=5000)
