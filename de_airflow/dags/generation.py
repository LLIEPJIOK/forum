from datetime import datetime
import json
import random
import requests
from airflow import DAG
from airflow.models import Variable
from airflow.operators.python import PythonOperator
import psycopg2


API_KEY = Variable.get("API_KEY")


default_args = {
    'owner': 'matthew',
    'start_date': datetime(2023, 9, 3, 10, 00)
}


def get_user_data():
    data = requests.get("https://randomuser.me/api/")
    data = data.json()['results'][0]
    data = format_user(data)
    return data


def format_user(data):
    final_user_data = {}
    final_user_data['nickname'] = data['name']['first'] + '_X_' +data['name']['last']
    final_user_data['email'] = data['email']
    final_user_data['password'] = get_password()
    return final_user_data

def get_post():
    api_url = 'https://api.api-ninjas.com/v1/facts'
    response = requests.get(api_url, headers={'X-Api-Key': API_KEY})
    if response.status_code == requests.codes.ok:
        post = response.json()[0]['fact']
        return post
    else:
        print("Error:", response.status_code, response.text)


def get_password():
    length = random.randint(8, 14)
    api_url = 'https://api.api-ninjas.com/v1/passwordgenerator?length={}'.format(length)
    response = requests.get(api_url, headers={'X-Api-Key': API_KEY})
    if response.status_code == requests.codes.ok:
        password = response.json()['random_password']
        return password
    else:
        print("Error:", response.status_code, response.text)

def batch_generate_users(batch: int, batch_iteration: int):
    batched_users = []
    for i in range(batch_iteration):
        print(f"batch number {i}")
        for x in range(batch):
            data = get_user_data()
            batched_users.append(data)
        for user in batched_users:
            requests.post(url="http://forum-forum-1:8000/user", data=json.dumps(user))            
        batched_users = []


def create_postgress_connection():
    conn = psycopg2.connect(
    host="forum-db-1",
    port="5430",
    database="forumdb",
    user="postgres",
    password="some_password"
)
    return conn

def format_post(id):
    post = {}
    text = get_post()
    author_id = id
    post['content'] = text
    post['author_id'] = author_id
    return post

def get_all_users_id(connection):
    cur = connection.cursor()
    cur.execute("SELECT array_agg(id) FROM users;")
    result = cur.fetchone()[0]
    if len(result) == 0:
        cur.close()
        connection.close()
        return
    elif len(result) > 1 and len(result) < 5:
        random_ids=random.sample(result,1)
        post = format_post(random_ids[0])
        requests.post("http://forum-forum-1:8000/post", data=json.dumps(post))
    elif len(result) > 10:
        random_ids=random.sample(result,8)
        for i in random_ids:
            post = format_post(i)
            requests.post("http://forum-forum-1:8000/post", data=json.dumps(post))
    cur.close()
    connection.close()


def generate_posts():
    get_all_users_id(create_postgress_connection())
    
        



with DAG('generate_users',
         default_args=default_args,
         schedule_interval='@once',
         catchup=False) as dag:
    generate_users=PythonOperator(
        task_id="generate_users",
        python_callable=batch_generate_users,
        op_args=[5, 2]
    )
    generate_posts=PythonOperator(
        task_id="generate_posts",
        python_callable=generate_posts
    )


generate_users >> generate_posts