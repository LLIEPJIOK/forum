from datetime import datetime
import json
import random
import requests
from airflow import DAG
from airflow.models import Variable
from airflow.operators.python import PythonOperator
from clickhouse_driver import Client




default_args = {
    'owner': 'matthew',
    'start_date': datetime(2023, 9, 3, 10, 00)
}

def create_ch_client():
    ch_client = Client(host='da_clickhouse-clickhouse-1:8123')

def load_users_data_to_clickhouse():



with DAG('generate_users',
         default_args=default_args,
         schedule_interval='@once',
         catchup=False) as dag:
    pass
