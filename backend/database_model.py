from sqlalchemy import Table, MetaData
from .database_connection import DatabaseConnection

db_engine = DatabaseConnection()
metadata_obj = MetaData()

accounts = Table("accounts", metadata_obj, autoload_with=db_engine)
followings = Table("followings", metadata_obj, autoload_with=db_engine)
sessions = Table("sessions", metadata_obj, autoload_with=db_engine)
company_blog_site = Table("company_blog_site", metadata_obj, autoload_with=db_engine)
company_blog_posts = Table("company_blog_posts", metadata_obj, autoload_with=db_engine)
company_blog_post_content = Table("company_blog_post_content", metadata_obj, autoload_with=db_engine)
