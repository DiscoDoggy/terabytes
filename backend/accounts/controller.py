from fastapi import Depends, APIRouter, HTTPException, Response, Request
from fastapi.security import HTTPBasicCredentials, HTTPBasic
from .model import SignupRequest, LoginRequest
from ..database_connection import DatabaseConnection
from sqlalchemy import insert, select, delete
from sqlalchemy.exc import IntegrityError
from ..database_model import accounts, sessions
import uuid
import bcrypt
from datetime import datetime, timedelta, timezone

router = APIRouter(
    prefix="/accounts"
)

db_engine = DatabaseConnection()
security = HTTPBasic()

@router.post("/signup")
def signup(signup_info : SignupRequest):
    username = signup_info.username
    password = signup_info.password

    password_salt = bcrypt.gensalt()
    bytes_password = password.encode("utf-8")
    hashed_password = bcrypt.hashpw(bytes_password, password_salt)
    hashed_password = hashed_password.decode()

    with db_engine.connect() as conn:
        account_id = uuid.uuid4()

        query = (
            insert(accounts)
            .values(id=account_id, username=username, password=hashed_password)
        )

        try:
            conn.execute(query)
            conn.commit()
        except IntegrityError: #i need to remember to create the unique constraint in DB
            conn.rollback()
            raise HTTPException(status_code=409, detail=f"An account with {username} username already exists.")
    
    return {"message" : "Account successfully registered"}

# def authenticate_user(credentials : HTTPBasicCredentials = Depends(security)):
#     user = users

def check_session_exists(account_id : uuid, db_engine):
    query = (
        select(sessions.c.id, sessions.c.account_id)
        .where(sessions.c.account_id == account_id)
    )

    with db_engine.connect() as conn:
        results = conn.execute(query)
    
    if results is not None:
        return True
    return False

# def authenticate_user()

@router.post("/login")
def login(response : Response, credentials:HTTPBasicCredentials=Depends(security)):
    query = (
        select(
            accounts.c.id,
            accounts.c.username,
            accounts.c.password
        )
        .where(accounts.c.username == credentials.username)
    )

    with db_engine.connect() as conn:
        results = conn.execute(query)
    
        count = 0
        for result in results:
            account_id = result.id
            password_from_db = result.password
            count += 1
        if count == 0:
            raise HTTPException(status_code=404, detail="Username and password do not match")
        
        inputted_password_encoded = credentials.password.encode("utf-8")
        password_from_db = password_from_db.encode("utf-8")

        do_passwords_match = bcrypt.checkpw(inputted_password_encoded, password_from_db)
        if not do_passwords_match:
            raise HTTPException(status_code=404, detail="Username and password do not match")

        is_logged_in = check_session_exists(account_id, db_engine)
        if is_logged_in:
            query = (
                delete(sessions)
                .where(sessions.c.account_id == account_id)
            )

        conn.execute(query)
        conn.commit()

        session_id = uuid.uuid4()
        start_date_time = datetime.now(timezone.utc)
        end_date_time = start_date_time + timedelta(days=30)

        query = (
            insert(sessions)
            .values(id=session_id, account_id=account_id, start_date_time=start_date_time, end_date_time=end_date_time)
        )

        conn.execute(query)
        conn.commit()

    response.set_cookie(key="session_id", value=session_id)
    return {"session_id" : session_id}

@router.get("/logout")
def logout(request : Request):
    session_token = request.cookies.get("session_id")
    print(f"session_token:{session_token}")
    query = (
        delete(sessions)
        .where(sessions.c.id == session_token)
    )

    with db_engine.connect() as conn:
        conn.execute(query)
        conn.commit()
        
    return {"message" : "Logout successful"}





"""
#APIS
Accounts
* POST Sign up 
    !DONE
* POST Login 
    !DONE
* GET logout
    !DONE
* UPDATE/PUT account

Blog Posts
* Get posts by accounts user follows the feed page
    ? Needs testing
* Get posts by tag
    ? Tag extraction not implemented yet
* Get posts randomly
* POST create blog post (pending feature)

Followers
* Get followers by accountId
* Get accounts user is following by accountId
* DELETE unfollow a blog or account
"""

