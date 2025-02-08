import uvicorn
from fastapi import FastAPI
from schema import schema
from strawberry.fastapi import GraphQLRouter

app = FastAPI()
graphql_app = GraphQLRouter(schema)
app.include_router(graphql_app, prefix="/graphql")


@app.get("/")
def health_check():
    return {"status": "ok"}


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=5000)
