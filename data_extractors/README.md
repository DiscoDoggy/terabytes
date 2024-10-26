# How Blog Extraction Will Happen

In order to keep uptodate with blogs, will download and parse many RSS XML feeds form different engineering blogs. These RSS feeds are in a standard format such that they easily provide key details such as: 
1. title 
2. company
3. Link to the whole article
4. Blog name
5. Description of the article
5. A short snippet of the article

Once we find the link to the whole article, we need to scrape the entire article for its content. 

## Scraping 
Engineering blogs all have different structures, but it would be inefficient to write over 200 different implementations to scrape for each different blog. The HTML of each blog is different. From the engineering blog full article we need to store a few pieces of information:
1. Authors
2. Actual body of the article,
    a. This is a hard problem. Because a blog post is not just text, 
    it is also made up of images which must be rendered in correct position when rendered in app. We also have heading types such as h1, h2, ... p. all of which correspond to differnt font sizes. How do we encode where and what order thigns appear on page in app
3. Image links and positionality

## Processing 
* Tag articles with their topic using LLM?

I somehow need to find a way to represent different structures. For example, images and their captions should not be embedded directly into the article body text. I need to figure out how to store where an image is located on the page.