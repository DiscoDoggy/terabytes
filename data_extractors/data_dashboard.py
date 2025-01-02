from dash import Dash, dcc, html, Input, Output, callback, dash_table
from dash.exceptions import PreventUpdate
import plotly.express as px
from datetime import date
import pandas as pd
from sqlalchemy import select
from database_model import extraction_job, company_blog_site
from database_connection import DatabaseConnection

db_engine = DatabaseConnection()

feed_names_query = (
    select(company_blog_site.c.blog_name)
)

with db_engine.connect() as conn:
    feed_name_results = conn.execute(feed_names_query)

feed_names = []
for feed in feed_name_results:
    feed_names.append(feed.blog_name)

app = Dash(__name__)

app.layout = html.Div([
    html.H1("Data pipeline Metrics"),

    html.Div([
        dcc.Dropdown(feed_names, "adsklf", id="feed-selection-dropdown", style={"flex" : "1", "min-width" : "200px"}),

        dcc.DatePickerRange(
            id="extraction-date-range-picker",
            min_date_allowed=date(2024,12,1),
            initial_visible_month=date(2024,12,1)
        ),
    ], 
    style={
        "display" : "flex",
        "justify-content" : "center"
    }
    ),

    html.Div(id="dd-output-container"),

    html.Div([
        dash_table.DataTable(page_size=10, style_table={"width" : "100%"}, id="extraction-job-table"),
        dcc.Graph(id="extraction-graph")
    ], style={"display" : "flex","justify-content" : "center", "gap" : "20px"}),
    dcc.Graph(id="import-graph")
])


@callback(
    Output("dd-output-container", "children"),
    Input("feed-selection-dropdown", "value")
)
def update_options(search_value):
    return f"you have selected {search_value}"

@callback(
    Output("extraction-job-table", "data"),
    Output("extraction-graph", "figure"),
    Output("import-graph", "figure"),
    Input("feed-selection-dropdown", "value"),
    Input("extraction-date-range-picker", "start_date"),
    Input("extraction-date-range-picker", "end_date")
)
def update_extraction_import(selected_feed, extraction_date_start, extraction_date_end):
    #update graph, update table
    jobs_table_rows_query = (
        select(
            extraction_job.c.id,
            company_blog_site.c.blog_name.label("blog_name"),
            extraction_job.c.status,
            extraction_job.c.start_time,
            extraction_job.c.end_time,
            extraction_job.c.num_extracted,
            extraction_job.c.num_imported
        )
        .select_from(extraction_job)
        .join(company_blog_site, company_blog_site.c.id == extraction_job.c.company_id)
    )

    if selected_feed:
        jobs_table_rows_query = jobs_table_rows_query.where(company_blog_site.c.blog_name == selected_feed)
    if extraction_date_start and extraction_date_end:
        jobs_table_rows_query = jobs_table_rows_query.where(extraction_job.c.start_time >= extraction_date_start)
        jobs_table_rows_query = jobs_table_rows_query.where(extraction_job.c.end_time <= extraction_date_end)
    
    compiled_query = jobs_table_rows_query.compile(db_engine,compile_kwargs={"literal_binds" : True})

    jobs_df = pd.read_sql(
        str(compiled_query), 
        db_engine,
        parse_dates={
            "start_time" : {"format" : "%m/%d/%y"},
            "end_time" : {"format" : "%m/%d/%y"}
        }
    )

    jobs_df["id"] = jobs_df["id"].astype(str)
    jobs_df["start_time"] = pd.to_datetime(jobs_df["start_time"]).dt.date.astype(str)
    jobs_df["end_time"] = pd.to_datetime(jobs_df["end_time"]).dt.date.astype(str)

    num_extractions_by_date = jobs_df.groupby("start_time")["num_extracted"].sum().reset_index()
    num_imported_by_date = jobs_df.groupby("start_time")["num_imported"].sum().reset_index()

    extraction_job_figure = px.line(num_extractions_by_date, x="start_time", y="num_extracted", title="Extractions over time")
    imported_job_figure = px.line(num_imported_by_date, x="start_time", y="num_imported", title="Imports over time")
    extraction_jobs_data = jobs_df.to_dict("records")

    return extraction_jobs_data, extraction_job_figure, imported_job_figure

if __name__ == "__main__":
    app.run(debug=True)