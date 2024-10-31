from ..dependencies.auth import authorize_user
from ..database_connection import DatabaseConnection
from ..database_model import accounts, company_blog_site, followings
from .model import AccountFollowingsResponse 

from fastapi import APIRouter, Depends
from sqlalchemy import select, insert
import uuid
from uuid import UUID 

router = APIRouter(
    prefix="/followings"
)

db_engine = DatabaseConnection()

@router.get("/")
def test():
    return {"message" : "this is the followings endpoint"}

@router.post("/follow_company_blog/{company_blog_id}")
def follow_company_blog(company_blog_id : UUID, account_info = Depends(authorize_user)):
    account_id = account_info["account_id"]
    new_follow_id = uuid.uuid4()

    query = (
        insert(followings)
        .values(id=new_follow_id, account_id=account_id, following_type="scraped", following_id=company_blog_id)
    )

    with db_engine.connect() as conn:
        conn.execute(query)
        conn.commit()
    
    return {"message" : "follow successful"}

@router.get("/following/{account_id}")
def get_who_account_follows(account_id : UUID) -> list[AccountFollowingsResponse]:
    query = (
        select(followings.c.following_id, company_blog_site.c.blog_name, accounts.c.username, followings.c.following_type)
        .select_from(followings)
        .join(company_blog_site, followings.c.following_id == company_blog_site.c.id)
        .join(accounts, followings.c.account_id == account_id)
        .where(accounts.c.id == account_id)
    )

    with db_engine.connect() as conn:
        results = conn.execute(query)

    if results is None:
        return []
    
    following = []
    for result in results:
        reformatted_following = {}

        reformatted_following["account_id"] = result.account_id
        reformatted_following["following_type"] = result.following_type

        if result["followings_type"] == "scraped":
            reformatted_following["username"] = result.blog_name
        else:
            reformatted_following["username"] = result.username
        
        following.append(reformatted_following)

    return following

@router.get("/followers/{account_id}")
def get_account_followers(account_id : UUID) -> list[AccountFollowingsResponse]:
    query = (
        select(followings.account_id, accounts.c.username, followings.c.following_type)
        .select_from(followings)
        .join(accounts, followings.c.following_id == account_id)
    )

    with db_engine.connect() as conn:
        results = conn.execute(query)
    if results is None:
        return []
    
    followers = []
    for result in results:
        followers.append(result)
    return followers

    