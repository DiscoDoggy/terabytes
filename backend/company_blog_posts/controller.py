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
        select(company_blog_posts.c.id, company_blog_site.c.blog_name, company_blog_posts.c.title, company_blog_posts.c.description)
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
        response_dict = {}
        response_dict["id"] = post_snippet.id
        response_dict["blog_name"] = post_snippet.blog_name
        response_dict["title"] = post_snippet.title
        response_dict["description"] = post_snippet.description

        post_snippets.append(response_dict)

    new_post_snippets = []
    with db_engine.connect() as conn:
        for post_snippet in post_snippets:
            new_post_snippet = post_snippet.copy()
            query = (
                select(company_blog_post_content.c.tag_content)
                .where(company_blog_post_content.c.company_blog_post_id == post_snippet["id"])
                .where(company_blog_post_content.c.tag_type == "img")
                .limit(1)
            )

            query_results = conn.execute(query)
            print(query_results)
            if query_results.first() is None:
                print("-----ENTE NONE------")
                new_post_snippet["cover_image_src"] = "https://t4.ftcdn.net/jpg/03/08/69/75/240_F_308697506_9dsBYHXm9FwuW0qcEqimAEXUvzTwfzwe.jpg"    
            else:
                print("------ENTER HERE-----")
                for result in query_results:
                    new_post_snippet["cover_image_src"] = result.tag_content 
            new_post_snippets.append(new_post_snippet)
        
    print(f"NEW POST \n{new_post_snippets[0]}")
    
    return new_post_snippets

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