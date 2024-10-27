import xml.etree.ElementTree as ET
from readability import Document
from bs4 import BeautifulSoup
import feedparser
import requests
import logging
from database_connection import DatabaseConnection
from sqlalchemy import text, select, insert
from database_model import company_feed_info, company_blog_posts, company_blog_post_content
import uuid

# from all_rss_feed_links import RSS_FEED_LINKS

class RSSParser:
    def __init__(self, rss_feed_link):
        self.rss_feed_link = rss_feed_link
    
    def fetch_feed(self):
        try:
            parsed_response = feedparser.parse(self.rss_feed_link)
            if "channel" not in parsed_response or len(parsed_response["channel"]) == 0:
                raise RuntimeError(f"{self.rss_feed_link} did not respond or channel tag was not present")
        except:
            logging.error(f"{self.rss_feed_link} did not respond or channel tag was not present")
            return None

        return parsed_response
    
    def parse_rss_info(self, parsed_rss):
        if parsed_rss == None:
            return
        
        rss_name = parsed_rss["channel"]["title"]
        rss_link = parsed_rss["channel"]["link"]
        
        if "description" not in parsed_rss["channel"]:
            logging.debug(f"{self.rss_feed_link} does not have feed description")
            rss_description = rss_name
        else:
            rss_description = parsed_rss["channel"]["description"]

        rss_info = RSSInfo(rss_name, rss_link, rss_description)
        
        return rss_info

    def parse_rss_blog_info(self, parsed_rss):
        
        all_blog_items = parsed_rss["items"]

        all_blog_info = []
        for blog_item in all_blog_items:
            blog_title = blog_item["title"]
            blog_link = blog_item["link"]

            if "description" not in blog_item:
                logging.debug(f"{self.rss_feed_link} : {blog_title} has no accessible description")
                blog_description = blog_title
            else:
                blog_description = blog_item["description"]
                print(blog_description)

            if "published" in blog_item:
                blog_publication_date = blog_item["published"]
            else:
                logging.debug(f"{self.rss_feed_link} : {blog_title} has no accessible publication date")
                blog_publication_date = None
            
            blog_tags = []

            if "tags" not in blog_item:
                logging.debug(f"{self.rss_feed_link} : {blog_title} has no accessible tags")
                
            else:
                for tag in blog_item["tags"]:
                    blog_tags.append(tag["term"])

            rss_blog_info = RSSBlogSnippet(
                title=blog_title,
                link=blog_link,
                publication_date=blog_publication_date,
                description=blog_description,
                tags=blog_tags
            )

            all_blog_info.append(rss_blog_info)
        
        return all_blog_info      

class RSSInfo:
    def __init__(self, rss_company_blog_name, rss_company_blog_link, rss_company_blog_description):
        self.rss_company_blog_name = rss_company_blog_name
        self.rss_company_blog_link = rss_company_blog_link
        self.rss_company_blog_description = rss_company_blog_description
    
class RSSBlogSnippet:
    def __init__(self, title : str, link : str, publication_date :str | None, description : str | None, tags : list[str]):
        self.title = title
        self.link = link
        self.publication_date = publication_date
        self.description = description
        self.tags = tags

    def to_dict(self):
        return {
            "title" : self.title,
            "link" : self.link,
            "publication_date" : self.publication_date,
            "description" : self.description,
            "tags" : self.tags
        }

class BlogParser:
    def __init__(self, blog_link):
        self.blog_link = blog_link

    def fetch_blog_html(self):
        try:
            response = requests.get(self.blog_link)
            response.raise_for_status()
        except requests.exceptions.RequestException as e:
            logging.error(f"Issue in fetching from Blog HTML: {e} blog link - {self.blog_link}")
            return None

        response.encoding = "utf-8"
        
        return response.text

    def parse_blog_post(self, html_blog_response):
        cleaned_blog_html = Document(html_blog_response)
        cleaned_blog_html = BeautifulSoup(cleaned_blog_html.summary(), "html.parser")

        essential_content_body = cleaned_blog_html.body
        all_blog_content = EssentialBlogHTMLContent()
        with open("parsed_response.txt", "a", encoding="utf-8") as file:
            for child in essential_content_body.descendants:
                tag_to_content = {}
                if child.name == "p":
                    file.writelines("\n-----THIS IS A PARAGRAPH NORMAL CONTENT-----\n")
                    file.writelines(child.text)
                    tag_to_content[child.name] = child.text
                elif child.name == "h2" or child.name == "h3" or child.name == "h4" or child.name == "h5" or child.name == "h6":
                    file.writelines("\n-----This is a header-----\n")
                    file.writelines(child.text)
                    tag_to_content[child.name] = child.text

                elif child.name == "img":
                    file.writelines("\n-----This is an image------\n")
                    file.writelines(child["src"])
                    tag_to_content[child.name] = child["src"]
                elif child.name == "li":
                    file.writelines("\n-------THIS IS A LIST ITEM------\n")
                    file.writelines(child.text)
                    tag_to_content[child.name] = child.text

                all_blog_content.add_content(tag_to_content)

        return all_blog_content            


class EssentialBlogHTMLContent:
    def __init__(self):
        """ (naturally sequential)
        [
            {tag : content},
            {tag : content},
            {tag : content},
        ]
        """
        self.content = []

    def add_content(self, tag_to_content : dict[str, str]):
        self.content.append(tag_to_content)
        
class DatabaseWriter:
    def __init__(self):
        self.engine = DatabaseConnection()
    
    def write_blog_to_db(self, blog_info : RSSBlogSnippet, blog_content : EssentialBlogHTMLContent):
        pass

    def write_company_to_db(self, feed_info: RSSInfo):
        with self.engine.connect() as conn:
            query = (
                select(company_feed_info.c.blog_name).
                where(company_feed_info.c == feed_info.rss_company_blog_name)
            )

            results = conn.execute(query)
            count = 0
            for _ in results:
                count += 1
            if not count:
                unique_id = uuid.uuid4()
                query = (
                    insert(company_feed_info).
                    values(
                        id=unique_id, 
                        blog_name=feed_info.rss_company_blog_name, 
                        feed_link=feed_info.rss_company_blog_link,
                        blog_description=feed_info.rss_company_blog_description
                    )
                )

class Driver:
    def __init__(self):
        pass
    def run(self):
        RSS_FEED_LINKS = {"Game Changer" : "https://tech.gc.com/atom.xml","Meta" : "https://engineering.fb.com/feed"}

        for rss_name in RSS_FEED_LINKS:
            rss_link = RSS_FEED_LINKS[rss_name]

            rss_parser = RSSParser(rss_link)
            rss_parser_object = rss_parser.fetch_feed()

            if rss_parser_object == None:
                continue

            rss_info = rss_parser.parse_rss_info(rss_parser_object) #returns RssInfo object
            all_rss_blog_info = rss_parser.parse_rss_blog_info(rss_parser_object) #returns a list of blogsnippets

            for blog_info in all_rss_blog_info:
                blog_link = blog_info.link
                blog_title = blog_info.title
                
                blog_parser = BlogParser(blog_link)
                blog_response = blog_parser.fetch_blog_html()
                
                if blog_response == None:
                    continue

                essential_blog_html_content = blog_parser.parse_blog_post(blog_response) #write this to the databse
        
driver = Driver()
driver.run()

# db_engine = DatabaseConnection()
# with db_engine.connect() as conn:
#     result = conn.execute(text("SELECT *"))







    