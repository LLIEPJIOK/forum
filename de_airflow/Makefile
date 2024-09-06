up:
	mkdir -p ./plugins ./logs
	echo -e "AIRFLOW_UID=$$(id -u)\nAIRFLOW_GID=0" > .env
	docker-compose up -d

down:
	docker-compose down