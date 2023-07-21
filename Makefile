# .PHONY:

stop:
	docker-compose down

run: 
	docker-compose up

# build-image:
# 	docker build -t telegram-bot-lift-kz:0.1 .

# start-container:
# 	docker run --env-file .env -p 80:80 telegram-bot-lift-kz:0.1


# docker build -t telegram-bot-lift-kz:v0.1 .	