# Default variables
CONFIG_FOLDER ?= users
DB_NAME ?= user_service_db
DB_HOST ?= localhost
DB_PORT ?= 5500
DB_USER ?= user_service_administrator
DB_PASSWORD ?= Ww76hQWbbt7zi2ItM6cNo4YYT
MIGRATION_STEPS ?= 1

# Generate the connection yml file for this db
setup:
	mkdir -p ../../configs/${CONFIG_FOLDER}
	cp ../../configs/db-template-dev.yml ../../configs/${CONFIG_FOLDER}/db-dev.yml
	sed -i "s/DB_NAME/${DB_NAME}/g" ../../configs/${CONFIG_FOLDER}/db-dev.yml
	sed -i "s/DB_USER/${DB_USER}/g" ../../configs/${CONFIG_FOLDER}/db-dev.yml
	sed -i "s/DB_PASSWORD/${DB_PASSWORD}/g" ../../configs/${CONFIG_FOLDER}/db-dev.yml

# Build the docker image for the db.
docker_db:
	docker build -t user_service_db_image .

# Remove existing docker image of the db.
remove_db: stop_db
	docker rm user_service_db_container
	docker image rm user_service_db_image

# Start of the (non existing) docker image for the db.
create_db:
	docker run -d --name user_service_db_container -p ${DB_PORT}:5432 user_service_db_image

# Start the (existing) docker image for the db.
start_db:
	docker start user_service_db_container

# Stop the (running) docker image fot the db.
stop_db:
	docker stop user_service_db_container

# Connect to the db in the docker image.
connect:
	psql -U ${DB_USER} -d ${DB_NAME} -h ${DB_HOST} -p ${DB_PORT}

# Apply the next migration to the db.
.PHONY: migrate
migrate:
	migrate -path ${PWD}/migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable up
	# Target allowing to migrate the db to latest defined step.

# Apply the next MIGRATION_STEPS migration(s) to the db.
migrateO:
	migrate -path ${PWD}/migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable up ${MIGRATION_STEPS}

# Revert all migrations of the db.
demigrate:
	migrate -path ${PWD}/migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable down

# Revert the previous MIGRATION_STEPS migration(s) of the db.
demigrateO:
	migrate -path ${PWD}/migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable down ${MIGRATION_STEPS}
