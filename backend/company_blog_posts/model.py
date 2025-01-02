from pydantic import BaseModel
from uuid import UUID

class BlogSnippetResponse(BaseModel):
    id : int
    blog_name : str
    title : str
    description : str | None = None
    cover_image_src : str | None = None
