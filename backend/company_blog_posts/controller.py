from fastapi import APIRouter

router = APIRouter(
    prefix="/posts"
)

@router.get("/random_posts")
def get_random_posts():
    pass

@router.get("/random_posts_by_tag")
def get_random_posts_by_tag():
    pass

#should fetching posts by follower feed go in here
#or in accounts