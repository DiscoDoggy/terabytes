import xml.etree.ElementTree as ET
from bs4 import BeautifulSoup
from readability import Document
import requests
import html
import json
from all_rss_feed_links import RSS_FEED_LINKS

def parse_rss(rss_link : str):
    xml_response = requests.get(rss_link)
    
    xml_tree_root = ET.fromstring(xml_response.content)
    rss_channel_tag = xml_tree_root.find("channel")
    
    blog_src_title = ""
    blog_src_link = ""
    blog_src_descript = ""
    blog_posts_details = []

    for child in rss_channel_tag:
        if child.tag == "title":
            blog_src_title = child.text
        if child.tag == "link":
            blog_src_link = child.text
        if child.tag == "description":
            blog_src_descript = child.text
        if child.tag == "item":
            blog_post_info = parse_blog_post(child)
            blog_posts_details.append(blog_post_info)
        
    print(f"Blog Source: {blog_src_title}")
    print(f"blog src link: {blog_src_link}")
    print(f"blog src description: {blog_src_descript}")

    with open("rss_extraction.json", "w", encoding="utf-8") as json_file:
        json.dump(blog_posts_details, json_file, ensure_ascii=False, indent=4)

def parse_blog_post(rss_item_tag) -> dict:
    blog_post_info = {}
    categories = []

    for child in rss_item_tag:
        if child.tag == "title":
            blog_post_info["title"] = html.unescape(child.text)
        if child.tag == "link":
            blog_post_info["link"] = html.unescape(child.text)
        if child.tag == "category":
            categories.append(html.unescape(child.text))
        if child.tag == "pubDate":
            blog_post_info["publication_date"] = child.text
        if child.tag == "description":
            blog_post_info["description"] = html.unescape(child.text)
    
    blog_post_info["categories"] = categories

    return blog_post_info

def get_blog_body(blog_link : str):
    response = requests.get(blog_link)
    response.encoding = "utf-8"
    doc = Document(response.text)

    raw_html_soup = BeautifulSoup(response.text, "html.parser")
    essential_content_soup = BeautifulSoup(doc.summary(), "html.parser")
    essential_content_body = essential_content_soup.body
    
    with open("parsed_response.txt", "w", encoding="utf-8") as file:
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
        
parse_rss("https://engineering.fb.com/feed/")
get_blog_body("https://www.canva.dev/blog/engineering/endpoint-vulnerability-management-at-scale/")
# get_author("https://www.canva.dev/blog/engineering/endpoint-vulnerability-management-at-scale/")
