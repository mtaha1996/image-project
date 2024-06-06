start:
	docker-compose -f ./docker-compose.yml --env-file compose-prod.env up --build -d 


stop:
	docker-compose -f ./docker-compose.yml down


rerun:
	make stop
	make start