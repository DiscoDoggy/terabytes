from dash import Dash, dcc, html, Input, Output, callback, dash_table
from dash.exceptions import PreventUpdate
import plotly.express as px
from datetime import date
import pandas as pd

start_date = None
end_date = None
feed_name = None 

df = pd.read_csv('https://raw.githubusercontent.com/plotly/datasets/master/gapminder2007.csv')

app = Dash(__name__)

df = px.data.gapminder()
fig = px.line(df.continent)

app.layout = html.Div([
    html.H1("Data pipeline Metrics"),

    html.Div([
        dcc.Dropdown(["8th Light Insights", "Meta Engineering", "Nytimes Engineering"], "adsklf", id="feed-selection-dropdown", style={"flex" : "1", "min-width" : "200px"}),

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
        dash_table.DataTable(data=df.to_dict("records"), page_size=10, style_table={"width" : "100%"}, id="extraction-job-table"),
        dcc.Graph(id="extraction-graph", figure=fig)
    ], style={"display" : "flex","justify-content" : "center", "gap" : "20px"}),
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
    Input("extraction-date-range-picker", "start_date"),
    Input("extraction-date-range-picker", "end_date")
)
def update_extraction_import():
    #update graph, update table
    pass

if __name__ == "__main__":
    app.run(debug=True)