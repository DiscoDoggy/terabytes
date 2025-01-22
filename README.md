# Terabytes
![Python](https://img.shields.io/badge/python-3670A0?style=for-the-badge&logo=python&logoColor=ffdd54) ![React Native](https://img.shields.io/badge/react_native-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB) ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![AWS](https://img.shields.io/badge/AWS-%23FF9900.svg?style=for-the-badge&logo=amazon-aws&logoColor=white) ![Plotly](https://img.shields.io/badge/Plotly-%233F4F75.svg?style=for-the-badge&logo=plotly&logoColor=white) ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white) ![Expo](https://img.shields.io/badge/expo-1C1E24?style=for-the-badge&logo=expo&logoColor=#D04A37)
## Introduction
Terabytes is tech blog social media application aimed at centralizing company tech blogs known and unknown while also supporting the sharing of user created blogs. 

Right now, Terabytes supports around 300 company blogs.

At its core, the goal of Terabytes is so the user can learn more about anything in the tech industry! Most people who read tech company blogs may only read the big company blogs such as Meta Engineering, Uber Engineering, AWS, or Door Dash.
But, there are so many great other engineering blogs out there that confront different issues than these large company blogs talk about. For example, the New York Times has their own engineering blog where one of the topics is 
how they have managed to create really great visualizations that function on multiple platforms and can be created really quickly which is really interesting to read about.

## Broad Architecture 
The appeal of Terabytes is two main features: accesing up to date tech blogs from different companies and being able to write and post your own blogs in app. 

### Company Tech Blogs - Data Pipelines
In order to support the viewing of company tech blogs, these tech blogs first need to be extracted. It would be very inefficient to write scrapers for multiple different company blogs let alone more than 300. The solution to this is to take advantage of a component most blogs have built in: RSS feeds. RSS feeds are more or less standardized XML documents that update when new blogs are uploaded to a source. They contain information such as the title, description, a few lines of the blog's content, and most importantly a direct link to the blog article. We still need to extract content from these blog articles though and different blog articles can be incredibly unstructured. The solution is a Ruby library that had been ported to python called Readablity. Most browsers and mobile browsers have a functionality called a "reader view" when viewing a article on the internet. This functionality somehow filters out all unnecessary content such as other recommended posts that may normally be at the bottom of a page, blogs, nav bars, etc. Readability can help narrow down only the relevant content on the page such as headers, images, videos, and text.

For each company blog RSS Feed what we have to do in broad strokes is:
1. Fetch the RSS feed
2. Fetch each blog article in the RSS feed
3. Parse the blog article for its content
4. Process, clean, format data
5. Ingest into Postgres

The extraction script deployed is mean't to run on a daily schedule.

Using Plotly and Dashly, I have built out a very rough dashboard to monitor extractions and database imports. 

### Backend
The backend was originally written in Python's FastAPI but I am now rewriting it in Go. Although I believe in FastAPI's ability to scale well especailly for this application that will likely never reach more than 10 people, I am doing this project to learn and one of things I have been wanting to learn is Go and how to write backends without the crux of a backend framework that constantly holds your hand. 

The backend powers the retrieval of both user created blogs and also blogs that have been extracted. 

### Frontend
The frontend is written in React Native. This is so I can  hopefully have an easier time deploying to both IOS and Android without learning both Swift and delving deeper into Kotlin. My computer is also too weak to run IDEs like Android studio 😢
