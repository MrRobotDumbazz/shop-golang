build: 
	docker-compose up 
restart:
	docker-compose down 
	docker image rm shop-golang_web
	docker-compose up 