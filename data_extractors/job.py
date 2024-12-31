from abc import ABC, abstractmethod
from pipeline import RSSBlogPipeline
from datetime import datetime, timezone
from database_connection import DatabaseConnection
from database_model import company_blog_site, extraction_job
from sqlalchemy import select, update
from sqlalchemy.dialects.postgresql import insert
from metrics_collector import MetricsCollector

class Job(ABC):
    @abstractmethod
    def run_job(self):
        pass

class PipelineExtractionJob(ABC):
    def __init__(self, pipeline : RSSBlogPipeline): #fix with refactoring later
        self.pipeline = pipeline
        self.db_engine = DatabaseConnection()
        self.job_id = self.on_job_create_tasks()

    def run_job(self):
        with self.db_engine.connect() as conn:
            update_query = (
                update(extraction_job)
                .where(extraction_job.c.id == self.job_id)
                .values(status="RUNNING")
            )

            conn.execute(update_query)
            conn.commit()

        start_time = datetime.now(timezone.utc)
        self.pipeline.run()
        end_time = datetime.now(timezone.utc)

        metrics = self.pipeline.get_metrics()
        self.write_metrics_db(metrics=metrics, start=start_time, end=end_time)
    
    def on_job_create_tasks(self):
        company_id_query = (
            select(company_blog_site.c.id)
            .where(company_blog_site.c.feed_link == self.pipeline.rss_feed_link)
        )

        with self.db_engine.connect() as conn:
            company_id = conn.execute(company_id_query).first().id
            job_insert_query = (
                insert(extraction_job)
                .values(company_id=company_id)
                .returning(extraction_job.c.id)
            )

            job_id = conn.execute(job_insert_query).first().id
            conn.commit()
            return job_id

    def write_metrics_db(self, metrics : MetricsCollector, start, end):
        with self.db_engine.connect() as conn:
            update_job_query = (
                update(extraction_job)
                .where(extraction_job.c.id == self.job_id)
                .values(
                    status="FINISHED",
                    start_time=start,
                    end_time=end,
                    num_extracted=metrics.get_metric("num_extracted"),
                    num_imported=metrics.get_metric("num_imported")
                )
            )

            conn.execute(update_job_query)
            conn.commit()
        
link = "https://8thlight.com/insights/feed/rss.xml"
blog_pipeline = RSSBlogPipeline(link)
job = PipelineExtractionJob(blog_pipeline)
job.run_job()

    

