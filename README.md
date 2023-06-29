
# go-game

A go implementation of a game. For now no real definition of what it is about.

# Installation

## Pre-requisites

This repository requires (and make use) of the following tools:
- [Go](https://go.dev/doc/install), version 1.20 has been used for development.
- [Postgresql](https://www.postgresql.org/download/linux/ubuntu/)
- [Docker](https://docs.docker.com/engine/install/ubuntu/)
- [Migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md)

## Code

To clone the repo and build the project from source run:
```bash
git clone git@github.com:Knoblauchpilze/go-game.git
```

Then go to the db creation [section](#create-the-db-container) and follow the instructions.

Finally go to the root of the repository and run:
```bash
make
```

You can execute the services with:
```bash
make run app_name
```

## Create the DB

### Attempt with postgres server

We first tried to do it another way than having the database in a docker container. Here are a few links that we gathered:
* This [link](https://www.cherryservers.com/blog/how-to-install-and-setup-postgresql-server-on-ubuntu-20-04) defines how to install postgres and what's installed by default.
* This [link](https://chartio.com/resources/tutorials/how-to-set-the-default-user-password-in-postgresql/) defines how to change the default password for the `postgres` user.
* This [link](https://medium.com/coding-blocks/creating-user-database-and-adding-access-on-postgresql-8bfcd2f4a91e) explains how to create a user and a database.

The issue is that it doesn't seem very clear how to do this in a programmatic way without using the psql shell. Googling yields surprisingly few results, even if it could probably be because the terms were not right.

So in the end we fell back to the database-in-a-container idiom.

### Create the DB container

- Go to [database](database) folder.
- Create the docker image for the database with:
```bash
make docker_db
```
- Run the database container from the image built in the previous step with:
```bash
make create_db
```
- Note that in case a previous operation already succeeded, a container with such a name already exists. To remove it (and then restart to the previous step), one should call:
```bash
make remove_db
```
- Initialize the database by calling the migrate target: this will create the schema associated to the data model of the application and populate the needed fields:
```bash
make migrate
```
- In case a reboot of the system happens, you can restart the db without having to recreate it with the following command:
```bash
make start_db
```

### Iterate on the DB schema

In case some new information needs to be added to the database one can use the migrations mechanism. By creating a new migration file in the relevant [directory](database/migrations) and naming accordingly (increment the number so that the `migrate` tool knows in which order migrations should be ran) it is possible to perform some modifications of the db by altering some properties. The migration should respect the existing constraints on the tables.

Once this is done one can rebuild the db by using the specific [Makefile](database/Makefile) target which will only apply the migrations not yet persisted in the db schema with:
```bash
make migrate
```

The migrations are designed in a way that each one can be applied sequentially and can also be rolled back: this is accomplished by having a `XYZ.up.sql` file and a `XYZ.down.sql` file. Any operation performed in the `up` part should have a counterpart in the `down` part to allow a roll back. Typically if a `CREATE TABLE` statement is issued in the `up`, a `DROP TABLE` should be in the `down` file.

### Managing the DB

If the db container has been stopped for some reasons one can relaunch it through the `make start_db` command.

One can also directly connect to the db using the `make connect` command. The password to do so can be found in the configuration files. SQL queries can then be performed in the `psql` editor.

In case a full rebuild of the database container is needed, go to the [database](database) folder and proceed to launch the following commands:
```bash
make remove_db
make docker_db
make create_db
make migrate
```

Most of the time though it is sufficient to migrate the database all the way down and rebuild it again with:
```bash
make demigrate
make migrate
```

If all the migrations are correctly set up this will allow to effectively remove all data from the database and start fresh.
