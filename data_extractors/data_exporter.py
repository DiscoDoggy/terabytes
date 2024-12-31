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
        # new_id = uuid.uuid4()
        insert_query = (
            insert(company_blog_site)
            .values(
                # id=new_id, 
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
        errors = []
        num_imported = 0
        with db_engine.connect() as conn:

            def get_company_id(feed_name : str):
                result = conn.execute(
                    select(company_blog_site.c.id)
                    .where(company_blog_site.c.blog_name == feed_name)
                ).first()      
                
                try:
                    company_id = result.id
                except Exception as e:
                    errors.append(f"{str(e)}")

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
                    .returning(company_blog_posts.c.id)
                )

                try:
                    blog_id = conn.execute(insert_query).first().id
                
                    for tag in blog.tags:
                        self.insert_tag_db(blog_id, tag, conn)
                except Exception as e:
                    errors.append(f"{str(e)}")
                    continue
                else:
                    num_imported += 1
            conn.commit()
            return num_imported, errors

    def insert_tag_db(self, post_id, tag : str, db_conn): 
        cleaned_tag = tag.strip().capitalize()
    
        query = (
            insert(all_db_tags)
            .values(name=cleaned_tag)
            .on_conflict_do_nothing(index_elements=["name"])
        )

        db_conn.execute(query)
        db_conn.commit()
        
        query = (
            select(all_db_tags.c.id)
            .where(all_db_tags.c.name == cleaned_tag)
        )

        result = db_conn.execute(query).first()

        query = (
            insert(blog_tags)
            .values(blog_post_id=post_id, tag_id=result.id)
        )

        db_conn.execute(query)