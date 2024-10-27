from sqlalchemy import create_engine
from dotenv import load_dotenv, dotenv_values
import os

class DatabaseConnection:
    _engine = None

    def __new__(cls, *args, **kwargs):
        if not cls._engine:
            cls._engine = super().__new__(cls, *args, **kwargs)
            cls._engine = cls._engine.create_connection()      
        return cls._engine      

    def create_connection(self):
        load_dotenv()
        secrets = dotenv_values("secrets.env")

        db_username = secrets["DB_USERNAME"]
        db_url = secrets["DB_URL"]
        db_port = secrets["DB_PORT"]
        db_name = secrets["DB_NAME"]
        db_password = secrets["DB_PASSWORD"]

        engine = create_engine(f"postgresql+psycopg2://{db_username}:{db_password}@{db_url}:{db_port}/{db_name}", echo=True)

        return engine