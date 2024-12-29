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
        