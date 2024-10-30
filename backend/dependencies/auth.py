from ..database_connection import DatabaseConnection
from ..database_model import *
from fastapi import Request, HTTPException
from datetime import datetime, timezone
from database_model import sessions, accounts
from sqlalchemy import select
db_engine = DatabaseConnection()

def authorize_user(request : Request):
    session_token = request.cookies.get("session_id")

    if not session_token:
        raise HTTPException(status_code=401, detail="User is not logged in. No active session")

    curr_time = datetime.now(timezone.utc)

    query = ( #i think this is an information leak
        select(accounts.c.id, sessions.c.end_date_time)
        .select_from(sessions)
        .join(accounts)
        .where(session_token == sessions.c.id)
    )

    with db_engine.connect() as conn:
        results = conn.execute(query)

    if results is None:
        raise HTTPException(status_code=404, detail="No session with this id exists")
    for result in results:
        account_id = result.id
        session_end_time = result.end_date_time
    
    if session_end_time <= curr_time:
        raise HTTPException(status_code=401, detail="Session expired. Please login again")

    return {"account_id" : account_id}

