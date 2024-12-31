from abc import ABC, abstractmethod
import logging
from rss_info import RSSInfo
from blog_snippet import RSSBlogSnippet
from database_connection import DatabaseConnection
from database_model import company_blog_site, company_blog_posts
from sqlalchemy import select
from readability import Document

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
    def __init__(self, blog_snippets : list[RSSBlogSnippet]):
        self.blog_snippets=blog_snippets
    
    def transform(self):
        new_blog_snippets = []
        for blog_snippet in self.blog_snippets:
            new_blog_snippet = self.transform_helper(blog_snippet)
            new_blog_snippets.append(new_blog_snippet)

        return new_blog_snippets
     
    def transform_helper(self, blog_snippet : RSSBlogSnippet):
        blog_snippet.print_blog_info()
        readability_processed_html = Document(blog_snippet.content.text)
        readability_processed_html = readability_processed_html.summary()
        
        new_blog_snippet = blog_snippet
        new_blog_snippet.content = readability_processed_html

        return new_blog_snippet