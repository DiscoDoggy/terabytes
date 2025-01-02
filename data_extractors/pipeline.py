from abc import ABC, abstractmethod
from data_loaders import *
from data_transformers import *
from data_exporter import *
from metrics_collector import MetricsCollector

from all_rss_links_complete import all_rss_links_complete as RSS_LINKS

class Pipeline(ABC):
    @abstractmethod
    def run(self):
        pass

class RSSBlogPipeline(Pipeline):
    def __init__(self, rss_feed_link : str):
        self.rss_feed_link = rss_feed_link
        self.metrics_interface = MetricsCollector()

        self.metrics_interface.create_metric("num_extracted", None)
        self.metrics_interface.create_metric("num_deduped", None)
        self.metrics_interface.create_metric("num_errors", None)
        self.metrics_interface.create_metric("num_imported", None)

    def run(self):
        rss_data_loader = RSSDataLoader(self.rss_feed_link)
        parsed_rss = rss_data_loader.load()

        blog_snippet_loader = BlogSnippetLoader(parsed_rss)
        rss_blog_snippets = blog_snippet_loader.load()
        self.metrics_interface.set_metric("num_extracted", len(rss_blog_snippets))

        deduplication_transformer = RSSBlogDeduplicationTransformer(rss_blog_snippets)
        deduped_blog_snippets = deduplication_transformer.transform()
        self.metrics_interface.set_metric("num_deduped", len(rss_blog_snippets) - len(deduped_blog_snippets))

        blog_content_loader = BlogContentLoader(deduped_blog_snippets)
        blog_snippets_with_content, errors = blog_content_loader.load()

        blog_content_transformer = BlogHTMLParseTransformer(blog_snippets_with_content)
        cleaned_blog_content = blog_content_transformer.transform()

        blog_exporter = RSSBlogExporter(cleaned_blog_content)
        num_imported, errors = blog_exporter.export()
        self.metrics_interface.set_metric("num_imported", num_imported)

    def get_metrics(self):
        return self.metrics_interface

class RSSFeedPipeline(Pipeline):
    def __init__(self, rss_feed_link : str):
        self.rss_feed_link=rss_feed_link        

    def run(self):
        rss_data_loader = RSSDataLoader(self.rss_feed_link)
        parsed_rss = rss_data_loader.load()

        rss_info_loader = RSSInfoLoader(parsed_rss)
        rss_info = rss_info_loader.load()
        rss_info.rss_feed_link = self.rss_feed_link

        rss_info_exporter = RSSFeedExporter(rss_info)
        rss_info_exporter.export()

# for link in RSS_LINKS:
#     print(link)
#     feed_pipeline = RSSFeedPipeline(RSS_LINKS[link])
#     try:
#         feed_pipeline.run()
#     except:
#         print(f"FAILURE FOR: {link}")
#         with open("broken_links.txt", "w") as file:
#             file.write(f"FAILURE: {link} -- {RSS_LINKS[link]}")
#         continue



        
