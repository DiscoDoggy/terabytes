from abc import ABC, abstractmethod
from data_loaders import *
from data_transformers import *
from data_exporter import *

class Pipeline(ABC):
    @abstractmethod
    def run(self):
        pass

class RSSBlogPipeline(Pipeline):
    def __init__(self, rss_feed_link : str):
        self.rss_feed_link = rss_feed_link
    
    def run(self):
        rss_data_loader = RSSDataLoader(self.rss_feed_link)
        parsed_rss = rss_data_loader.load()

        blog_snippet_loader = BlogSnippetLoader(parsed_rss)
        rss_blog_snippets = blog_snippet_loader.load()

        deduplication_transformer = RSSBlogDeduplicationTransformer(rss_blog_snippets)
        deduped_blog_snippets = deduplication_transformer.transform()

        blog_content_loader = BlogContentLoader(deduped_blog_snippets)
        blog_content = blog_content_loader.load()

        blog_content_transformer = BlogHTMLParseTransformer(blog_content)
        cleaned_blog_content = blog_content_transformer.transform()

        blog_exporter = RSSBlogExporter(cleaned_blog_content)
        blog_exporter.export()

class RSSFeedPipeline(Pipeline):
    def __init__(self, rss_feed_link : str):
        self.rss_feed_link=rss_feed_link        

    def run(self):
        rss_data_loader = RSSDataLoader(self.rss_feed_link)
        parsed_rss = rss_data_loader.load()

        rss_info_loader = RSSInfoLoader(parsed_rss)
        rss_info = rss_info_loader.load()

        rss_info_exporter = RSSFeedExporter(rss_info)
        rss_info_exporter.export()



        
