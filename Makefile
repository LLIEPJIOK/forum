up:
	make up --directory=de_airflow
	make up --directory=forum
	docker network create airflow-network
	docker network connect airflow-network de_airflow-airflow-worker-1
	docker network connect airflow-network de_airflow-airflow-webserver-1
	docker network connect airflow-network de_airflow-airflow-scheduler-1
	docker network connect airflow-network de_airflow-airflow-triggerer-1
	docker network connect airflow-network forum-forum-1

down:
	make down --directory=de_airflow
	make down --directory=forum
	docker network remove airflow-network