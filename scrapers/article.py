from datetime import datetime
import pickle


class Article:
    def __init__(self, link: str, title: str, description: str, content: str, time: datetime):
        self.link = link
        self.title = title
        self.description = description
        self.content = content
        self.time = time

    def __repr__(self) -> str:
        return f"title='{self.title}'\ndescription='{self.description}'\ncontent='{self.content}'\ntime={self.time}"

    def dump(self, file_name: str):
        with open(f"{file_name}.pkl", 'wb') as file:
            pickle.dump(self, file)
