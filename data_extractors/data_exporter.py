from abc import ABC, abstractmethod
from sqlalchemy import select
from database_model import company_blog_site, company_blog_posts, all_db_tags, blog_tags
from database_connection import DatabaseConnection
from blog_snippet import RSSBlogSnippet
from rss_info import RSSInfo
from sqlalchemy.dialects.postgresql import insert
import uuid

db_engine = DatabaseConnection()

class DataExporter(ABC):
    @abstractmethod
    def export(self):
        pass

class RSSFeedExporter(DataExporter):
    def __init__(self, rss_feed_info : RSSInfo):
        self.rss_feed_info = rss_feed_info
    
    def export(self):
        new_id = uuid.uuid4()
        insert_query = (
            insert(company_blog_site)
            .values(
                id=new_id, 
                blog_name=self.rss_feed_info.rss_company_blog_name,
                feed_link=self.rss_feed_info.rss_company_blog_link,
                blog_description=self.rss_feed_info.rss_company_blog_description
            )
            .on_conflict_do_nothing(index_elements=["blog_name"])
        )

        with db_engine.connect() as conn:
            conn.execute(insert_query)
            conn.commit()

class RSSBlogExporter(DataExporter):
    def __init__(self, rss_blog_snippets : list[RSSBlogSnippet]):
        self.rss_blog_snippets = rss_blog_snippets
    
    def export(self):
        with db_engine.connect() as conn:

            def get_company_id(feed_name : str):
                result = conn.execute(
                    select(company_blog_site.c.id)
                    .where(company_blog_site.c.blog_name == feed_name)
                ).first()
                
                if result is None:
                    raise ValueError("Get company ID should not yield a \'None\' value")
                
                company_id = result.id
                return company_id
            
            for blog in self.rss_blog_snippets:
                insert_query = (
                    insert(company_blog_posts)
                    .values(
                        company_id=get_company_id(blog.rss_feed_name),
                        title=blog.title,
                        description=blog.description,
                        publication_date=blog.publication_date,
                        link=blog.link,
                        content=blog.content
                    )
                    .returning(index_elements=["id"])
                )

                blog_id = conn.execute(insert_query).first().id
                conn.commit()

                for tag in blog.tags:
                    self.insert_tag_db(blog_id, tag)


    def insert_tag_db(self, post_id, tag : str, db_conn):
        cleaned_tag = tag.strip().capitalize()

        query = (
            select(all_db_tags.c.id)
            .where(cleaned_tag == all_db_tags.c.name)
        )

        result = db_conn.execute(query).first()

        tag_id = -1
        if result is None:
            query = (
                insert(all_db_tags)
                .values(name=tag)
                .returning(index_elements=["id"])
            )

            tag_result = db_conn.execute(query).first()
            tag_id = tag_result.id
            db_conn.commit()
        else:
            tag_id = result.id
        
        insert_query = (
            insert(blog_tags)
            .values(blog_post_id=post_id, tag_id=tag_id)
        )

        db_conn.execute(insert_query)
        db_conn.commit()