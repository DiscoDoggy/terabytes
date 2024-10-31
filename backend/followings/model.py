from pydantic import BaseModel
from uuid import UUID

class AccountFollowingsResponse(BaseModel):
    account_id : UUID
    username : str
    following_type : str
