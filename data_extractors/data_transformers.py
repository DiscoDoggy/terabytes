from abc import ABC, abstractmethod
import logging
from rss_info import RSSInfo
from blog_snippet import RSSBlogSnippet
from database_connection import DatabaseConnection
from database_model import company_blog_site, company_blog_posts
from sqlalchemy import select

db_engine = DatabaseConnection()

class DataTransformer(ABC):
    @abstractmethod
    def transform(self):
        pass

class RSSBlogDeduplicationTransformer(DataTransformer):
    def __init__(self, blog_snippets : list[RSSBlogSnippet]):
        self.blog_snippets=blog_snippets
    
    def transform(self):
        blogs_not_in_db = []
        with db_engine.connect() as conn:
            for blog_snippet in self.blog_snippets:
                query = (
                    select(company_blog_posts.c.id)
                    .select_from(company_blog_posts)
                    .join(company_blog_site, company_blog_site.c.id == company_blog_posts.c.company_id)
                    .where(blog_snippet.title == company_blog_posts.c.title)
                    .where(blog_snippet.rss_feed_name == company_blog_site.c.blog_name)
                )

                result = conn.execute(query).first()
                if result is None:
                    blogs_not_in_db.append(blog_snippet)
        
        return blogs_not_in_db
    
class BlogHTMLParseTransformer(DataTransformer):
    def __init__(self, blog_snippets : RSSBlogSnippet):
        self.blog_snippets=blog_snippets
    def transform(self):
        

