from abc import ABC, abstractmethod
from rss_info import RSSInfo
from blog_snippet import RSSBlogSnippet
import feedparser
import logging
import requests

class DataLoader(ABC):
    @abstractmethod
    def load(self):
        pass

class RSSDataLoader(DataLoader):
    def __init__(self, rss_url):
        self.rss_url = rss_url
    
    def load(self):
        #bubble request errors up to the worker/pipeline level
        parsed_response = feedparser.parse(self.rss_url)
        if "channel" not in parsed_response or len(parsed_response["channel"]) == 0:
            raise RuntimeError(f"{self.rss_url} has no channel tag")
        
        return parsed_response
    
class RSSInfoLoader(DataLoader):
    def __init__(self, feedparser_obj):
        self.feedparser_rss_feed = feedparser_obj
    
    def load(self):
        if self.feedparser_rss_feed == None:
            raise ValueError("RSS Data cannot have value None")        
        
        rss_feed_name = self.feedparser_rss_feed["channel"]["title"]
        rss_feed_link = self.feedparser_rss_feed["channel"]["link"]

        if "description" not in self.feedparser_rss_feed["channel"]:
            logging.info(f"{rss_feed_name} has no field \'description\'")
            rss_description = rss_feed_name
        else:
            rss_description = self.feedparser_rss_feed["channel"]["description"]

        rss_info = RSSInfo(
            rss_company_blog_name=rss_feed_name,
            rss_company_blog_link=rss_feed_link,
            rss_company_blog_description=rss_description
        )

        return rss_info

class BlogSnippetLoader(DataLoader):
    def __init__(self, feedparser_obj):
        self.feedparser_rss_feed = feedparser_obj

    def load(self):
        feed_blogs = self.feedparser_rss_feed["items"]
        rss_feed_title = self.feedparser_rss_feed["channel"]["title"]
        all_blogs = []
        for blog_snippet in feed_blogs:
            blog_title = blog_snippet["title"]
            blog_link = blog_snippet["link"]

            if "description" not in blog_snippet:
                logging.info(f"{self.feedparser_rss_feed["channel"]["title"]}--{blog_title} : has no accessible field \'description\'")
                blog_description = blog_title
            else:
                blog_description = blog_snippet["description"]
            
            if "published" not in blog_snippet:
                logging.info(f"{self.feedparser_rss_feed["channel"]["title"]}--{blog_title} : has no accessible field \'published\'")
                blog_publication_date = None
            else:
                blog_publication_date = blog_snippet["published"]
        
            blog_tags = []
            if "tags" not in blog_snippet:
                logging.info(f"{self.feedparser_rss_feed["channel"]["title"]}--{blog_title} : has no accessible field \'tags\'")
            else:
                for tag in blog_snippet["tags"]:
                    blog_tags.append(tag["term"])

            rss_blog_snippet = RSSBlogSnippet(
                title=blog_title,
                rss_feed_name=rss_feed_title,
                link=blog_link,
                publication_date=blog_publication_date,
                description=blog_description,
                tags=blog_tags
            )

            all_blogs.append(rss_blog_snippet)
        
        return all_blogs
    
class BlogContentLoader(DataLoader):
    def __init__(self, blog_snippets : list[RSSBlogSnippet]):
        self.blog_snippets=blog_snippets
    
    def load(self):
        errors = []
        valid_blog_snippets = [] 
        for blog_snippet in self.blog_snippets:
            try:
                response = requests.get(blog_snippet.link)
                response.raise_for_status()
            except Exception as e:
                errors.append((blog_snippet.link, str(e)))
            else:
                new_blog_snippet = blog_snippet
                new_blog_snippet.content = response
                valid_blog_snippets.append(new_blog_snippet)
            
        return valid_blog_snippets, errors
