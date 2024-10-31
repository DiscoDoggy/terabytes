from pydantic import BaseModel
import uuid

class BlogSnippetResponse(BaseModel):
    id : uuid
    title : str
    description : str | None
