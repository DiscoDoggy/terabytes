from ..dependencies.auth import authorize_user
from ..database_connection import DatabaseConnection
from ..database_model import accounts, company_blog_posts, company_blog_site, company_blog_post_content, followings
from .model import BlogSnippetResponse

from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy import select, func

router = APIRouter(
    prefix="/posts"
)

db_engine = DatabaseConnection()
#we should paginate for infinite scroll but that takes time so we arent going to right now and posts that the user hasn't seen
@router.get("/following")
def get_following_posts(account_info : dict = Depends(authorize_user)) -> list[BlogSnippetResponse]:
    account_id = account_info["account_id"]

    query = (
        select(company_blog_posts.c.id, company_blog_posts.c.title, company_blog_posts.c.description)
        .select_from(company_blog_posts)
        .join(company_blog_site, company_blog_site.c.id == company_blog_posts.c.company_id)
        .join(followings, followings.c.following_id == company_blog_site.c.id)
        .where(followings.c.account_id == account_id)
    )

    with db_engine.connect() as conn:
        results = conn.execute(query)
    if results is None:
        return []
    
    post_snippets = []
    for post_snippet in results:
        post_snippets.append(post_snippet)
    
    return post_snippets

@router.get("/random_posts")
def get_random_posts() -> list[BlogSnippetResponse]:
    query = (
        select(company_blog_posts.c.id, company_blog_posts.c.title, company_blog_posts.c.description)
        .order_by(func.random())
        .limit(100)
    )

    with db_engine.connect() as conn:
        results = conn.execute(query)
    random_blog_posts = []
    for result in results:
        random_blog_posts.append(result)
    return random_blog_posts

@router.get("/random_posts_by_tag/{tag_id}")
def get_random_posts_by_tag():
    pass

#should fetching posts by follower feed go in here
#or in accounts