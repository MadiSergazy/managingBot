## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
# .PHONY:
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]
## stopped docker-compose file
stop:
	docker-compose down
##run docker-compose file
run: 
	docker-compose up --build
##deleted the volume of application
volume/del: confirm
	docker-compose down -v



# build-image:
# 	docker build -t telegram-bot-lift-kz:0.1 .

# start-container:
# 	docker run --env-file .env -p 80:80 telegram-bot-lift-kz:0.1


# docker build -t telegram-bot-lift-kz:v0.1 .	