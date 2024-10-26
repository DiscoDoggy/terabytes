import xml.etree.ElementTree as ET
from readability import Document
from bs4 import BeautifulSoup
import feedparser
import requests
import logging
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

        with open("parsed_response.txt", "a", encoding="utf-8") as file:
            for child in essential_content_body.descendants:
                if child.name == "p":
                    file.writelines("\n-----THIS IS A PARAGRAPH NORMAL CONTENT-----\n")
                    file.writelines(child.text)
                elif child.name == "h2" or child.name == "h3" or child.name == "h4" or child.name == "h5" or child.name == "h6":
                    file.writelines("\n-----This is a header-----\n")
                    file.writelines(child.text)
                elif child.name == "img":
                    file.writelines("\n-----This is an image------\n")
                    file.writelines(child["src"])
                elif child.name == "li":
                    file.writelines("\n-------THIS IS A LIST ITEM------\n")
                    file.writelines(child.text)

# class BlogContent:
#     def __init__(self, content : str)

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

                blog_parser.parse_blog_post(blog_response)
        
driver = Driver()
driver.run()








    