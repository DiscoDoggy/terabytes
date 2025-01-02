class RSSBlogSnippet:
    def __init__(
            self,
            title : str,
            rss_feed_name : str,
            link : str,
            publication_date : str | None,
            description : str,
            tags : list[str]
    ):
        self.title=title
        self.rss_feed_name=rss_feed_name
        self.link=link
        self.publication_date=publication_date
        self.description=description
        self.content = None
        self.tags = tags

    def print_blog_info(self):
        print(f"Blog Title: {self.title}")
        print(f"feed name: {self.rss_feed_name}")
        print(f"Link: {self.link}")
        print(f"content status: {self.content}")