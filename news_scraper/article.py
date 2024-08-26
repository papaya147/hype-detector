import os
import json
from datetime import datetime
from typing import List, Tuple, Dict
import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from collections import defaultdict


class Article:
    def __init__(self, url: str = "", title: str = "", description: str = "", content: str = "", cleaned_content: str = "", timestamp: datetime = None, market_timestamp: datetime = None, off_market_hours: bool = False) -> None:
        self.url = url
        self.title = title
        self.description = description
        self.content = content
        self.cleaned_content = cleaned_content
        self.timestamp = timestamp
        self.market_timestamp = market_timestamp
        self.off_market_hours = off_market_hours

    def __repr__(self) -> str:
        return (f"Article(url='{self.url}', title='{self.title}', "
                f"description='{self.description}', timestamp='{self.timestamp}', "
                f"market timestamp='{self.market_timestamp}', off market hours='{self.off_market_hours}')")

    @classmethod
    def from_dict(self, data: dict) -> "Article":
        timestamp = datetime.fromisoformat(
            data.get("timestamp")) if data.get("timestamp") else None
        market_timestamp = datetime.fromisoformat(
            data.get("market_timestamp")) if data.get("market_timestamp") else None
        return self(
            url=data.get("url", ""),
            title=data.get("title", ""),
            description=data.get("description", ""),
            content=data.get("content", ""),
            cleaned_content=data.get("cleaned_content", ""),
            timestamp=timestamp,
            market_timestamp=market_timestamp,
            off_market_hours=data.get("off_market_hours", False)
        )

    def add_content_word_dict(self) -> Dict[str, int]:
        content_words = [word for word in self.cleaned_content.split()]
        word_dict = defaultdict(int)
        for word in content_words:
            word_dict[word] += 1
        self.word_dict = dict(word_dict)


def load_articles(folder: str) -> List[Article]:
    articles = []
    for filename in os.listdir(folder):
        if filename.endswith(".json"):
            file_path = os.path.join(folder, filename)
            with open(file_path, 'r') as file:
                data = json.load(file)
                article = Article.from_dict(data)
                articles.append(article)
    return articles


def get_content(articles: List[Article]) -> List[str]:
    return [article.cleaned_content for article in articles]


def tfidf_vectorise(articles: List[Article]) -> pd.DataFrame:
    content = get_content(articles)

    vectorizer = TfidfVectorizer()
    tfidf_matrix = vectorizer.fit_transform(content)
    feature_names = vectorizer.get_feature_names_out()

    return pd.DataFrame(tfidf_matrix.toarray(), columns=feature_names)


def extract_important_words(articles: List[Article], top_n: int = 10) -> List[Tuple[str, float]]:
    tfidf_df = tfidf_vectorise(articles)

    important_words = []
    for index, row in tfidf_df.iterrows():
        sorted_terms = row.sort_values(ascending=False)
        top_terms = sorted_terms.head(top_n)
        important_words.append(list(top_terms.items()))
    return important_words
