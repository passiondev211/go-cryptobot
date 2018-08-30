# Cryptobot

See `docs/1_poc_requirements.md` for the full description.

[API description](https://documenter.getpostman.com/view/2186314/cryptobot/7Lq7gG5)

[Postman collection](https://www.getpostman.com/collections/643163e560085d98f9a7)

## Project Structure

    app/      - Go application domain layer
    db/       - Go package for DB migrations handling
    docs/     - project documentation
    fe/       - front-end
    schemas/  - JSON schemas defining the protocol between front-end and back-end
    vendor/   - Go 3rd-party libraries

## Back-End Development

### Prerequisites   

1. ensure `make` are installed on your system
2. install Docker and Docker Compose: https://www.docker.com/products/overview
3. install `glide` for package management: `go get -u github.com/Masterminds/glide`

Now you can clone the application and prepare it for running:

    mkdir -p $GOPATH/src/fkaller
    cd $GOPATH/src/fkaller
    git clone git@bitbucket.org:fkaller/cryptobot.git
    cd $GOPATH/src/fkaller/cryptobot
    cp config.sample.json config.json

### Running

As simple as `make up`.

If you're running the project on the production environment, use `make prod` to enable https. Note that `make prod`
won't work locally due to inability to get a LetsEncrypt certificate. All debugging must be done on the server.

### Introducing Changes

#### Front-End

Follow `fe/README.md` for instructions, build the project and put static files into `./dist`.

#### Back-End

New changes to the Go source code may be introduced in the following way:

1. run the app in the terminal: `make up`
2. make changes to the source code
3. `Ctrl+C` in the terminal
4. `make up` again and check your changes

### Database migrations

In this api implementing db migrations through makefile.   
run `make migrate` in working directory with parametres `type` and `n`.   
*type* required. Must be: 
  
 - up - go to later migration    
 - down - go to earlier migration     
 - current - show current last active migration name.      

*n* optional. Determines, how many steps upping/downing. 
If *n* not selected, then api will be applied to top/bottom.  

Example input: `make migrate type=up n=1`

Migration files are located at ./db/migrations   
Migration files type is SQL.    
The migration files are created manually, according to the following template:
    
    -- +migrate Up
    -- SQL in section 'Up' is executed when this migration is applied
    CREATE TABLE people (id int);
    
    
    -- +migrate Down
    -- SQL section 'Down' is executed when this migration is rolled back
    DROP TABLE people;
Here comments play the role of a separator, for parsing up and down code of migration.   
> **WARNING!**   
> use `make migrate type=current` for to be sure that migration you want to edit is not used.   
> If migration is used - go down to a more recent version.   

Recommended to follow the following file name: 
    
    0001_create_table.sql
    0002_insert_record.sql
    0003_add_role.sql
