import feedparser
import requests
from all_rss_links_complete import all_rss_links_complete as RSS_FEED_LINKS
from readability import Document

# link = "https://github.blog/engineering.atom"
# link = "http://tech.finn.no/atom.xml"
# link = "https://www.hostinger.com/blog/feed"
# response = requests.get(link)
# parsed_response = feedparser.parse(link)

# print(f"{parsed_response["channel"]}")
# print(f"Number of entries: {len(parsed_response["items"])}")
# print(f"Number of entries: {len(parsed_response.entries)}")
# for i in range(len(parsed_response["items"])):
#     print(parsed_response["items"][i]["title"])
#     print(parsed_response["items"][i]["link"])

#     # blog_tag_terms = []
#     all_tags = parsed_response["items"][i]["tags"]
#     for tag in all_tags:
#         # blog_tag_terms.append(tag["term"])
#         print(tag["term"])
#     print("\n")

# def test_for_broken_resources():
#     with open("non_working_feeds.txt", "w") as file:
#         for rss_feed_name in RSS_FEED_LINKS:
#             link = RSS_FEED_LINKS[rss_feed_name]
            
#             try:
#                 parsed_response = feedparser.parse(link)

#                 if len(parsed_response["channel"]) == 0:
#                     raise ValueError("Link did not work")
#             except:
#                 print(f"{link} | Reason: Link did not work/site was down")
#                 file.write(f"{link} | Reason: Link did not work/site was down\n")
#                 continue

#             if "channel" not in parsed_response:
#                 print(f"{link} | Reason: Channel tag or atom equivalent was not present")
#                 file.write(f"{link} | Reason: Channel tag or atom equivalent was not present\n")
#             elif "items" not in parsed_response or len(parsed_response.entries) == 0:
#                 print(f"{link} | Reason: Entries or items not present")
#                 file.write(f"{link} | Reason: entries or items not present\n")
            
#             channel = parsed_response["channel"]
#             if "link" not in channel:
#                 print(f"{link} | no link")
# def test_working_links():
#     #things i need to  check existence of 
#     # Description/subtitle
#     # categories
#     pass

# test_for_broken_resources()


