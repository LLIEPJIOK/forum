up:
	make up --directory=de_airflow
	make up --directory=forum
	make up --directory=da_clickhouse
	docker network create airflow-network
	docker network connect airflow-network de_airflow-airflow-worker-1
	docker network connect airflow-network de_airflow-airflow-webserver-1
	docker network connect airflow-network de_airflow-airflow-scheduler-1
	docker network connect airflow-network de_airflow-airflow-triggerer-1
	docker network connect airflow-network forum-forum-1
	docker network connect airflow-network forum-db-1
	docker network connect airflow-network da_clickhouse-clickhouse-1

down:
	make down --directory=de_airflow
	make down --directory=forum
	make down --directory=da_clickhouse
	docker network remove airflow-network