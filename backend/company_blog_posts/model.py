from pydantic import BaseModel
from uuid import UUID

class BlogSnippetResponse(BaseModel):
    id : int
    title : str
    description : str | None
