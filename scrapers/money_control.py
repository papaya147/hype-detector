import requests
from bs4 import BeautifulSoup
from typing import List
from datetime import datetime
from scrapers.article import Article
import os


def article_links(page_start: int, page_end: int) -> List[str]:
    article_links = []
    for i in range(page_start, page_end):
        base_url = "https://www.moneycontrol.com/news/tags/companies/news"
        url = f"{base_url}/page-{i}/"
        req = requests.request(url=url, method="GET")
        soup = BeautifulSoup(req.text, "lxml")
        articles_ul = soup.find("ul", id="cagetory")
        article_lis = articles_ul.find_all("li", class_="clearfix")
        for article_li in article_lis:
            ele = article_li.find("h2").find("a")["href"]
            article_links.append(ele)
        print(f"page {i} completed...")
    return article_links


def scrape_article(article_url: str) -> Article:
    print(f"scraping {article_url}...")
    req = requests.request(url=article_url, method="GET")
    soup = BeautifulSoup(req.text, "lxml")

    article = soup.find("div", class_="page_left_wrapper")

    title = article.find("h1").text
    desc = article.find("h2").text

    time_div = article.find("div", class_="article_schedule")
    date = time_div.find("span").text
    time = time_div.text.split("/")[-1].strip()
    datetime_obj = datetime.strptime(f"{date} {time}", "%B %d, %Y %H:%M IST")

    content_div = article.find("div", id="contentdata")
    if not content_div:
        print(f"article {article_url} contains no content, skipping...")
        return None

    ps = content_div.find_all("p", class_="")
    content = []
    for p in ps:
        content.append(p.text.strip())

    return Article(article_url, title.strip(), desc.strip(), " ".join(content), datetime_obj)


def scrape(page_start: int, page_end: int) -> List[Article]:
    articles = []
    print("scraping article list...")
    links = article_links(page_start, page_end)
    print("scraping articles...")
    for link in links:
        article = scrape_article(link)
        if article:
            articles.append(article)
    return articles


def scrape_and_save(page_start: int, page_end: int, dir_name: str):
    articles = scrape(page_start, page_end)
    os.mkdir(dir_name)
    print("dumping pkl files...")
    for article in articles:
        article.dump(f"{dir_name}/{article.time}-{article.title}")
