from fastapi import FastAPI
from .accounts.controller import router as accounts_router
from .company_blog_posts.controller import router as company_blog_posts_router
from .followings.controller import router as followings_router

app = FastAPI()

app.include_router(company_blog_posts_router)
app.include_router(accounts_router)
app.include_router(followings_router)

app.get("/")
def read_root():
    return {"hello" : "world"}