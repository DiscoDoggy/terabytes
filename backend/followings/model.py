from pydantic import BaseModel
import uuid

class AccountFollowingsResponse(BaseModel):
    account_id : uuid
    username : str
    following_type : str
