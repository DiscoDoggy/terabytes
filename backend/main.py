from fastapi import FastAPI
from .accounts.controller import router as accounts_router
app = FastAPI()

app.include_router(accounts_router)

app.get("/")
def read_root():
    return {"hello" : "world"}