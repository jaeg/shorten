from locust import HttpUser, TaskSet, task, between
import json

class UserTasks(TaskSet):
  @task
  def addData(self):
    response = self.client.post("/shorten",data=json.dumps({
        "link":"test"
        }))


class WebsiteUser(HttpUser):
    host = "http://127.0.0.1:8090"
    wait_time = between(2, 5)
    tasks = [UserTasks]