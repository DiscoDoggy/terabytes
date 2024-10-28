from sqlalchemy import Table, MetaData
from database_connection import DatabaseConnection

db_engine = DatabaseConnection()
metadata_obj = MetaData()

company_blog_site = Table("company_blog_site", metadata_obj, autoload_with=db_engine)
company_blog_posts = Table("company_blog_posts", metadata_obj, autoload_with=db_engine)
company_blog_post_content = Table("company_blog_post_content", metadata_obj, autoload_with=db_engine)
