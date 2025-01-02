class MetricsCollector:
    def __init__(self):
        self.metrics = {}

    def create_metric(self, metric_name : str, initial_value):
        self.metrics[metric_name] = initial_value

    def set_metric(self, metric_name : str, value):
        self.metrics[metric_name] = value
    
    def increment_metric(self, metric_name : str, increment_value : int):
        self.metrics[metric_name] += increment_value

    def get_metric(self, metric_name : str):
        return self.metrics[metric_name]
        