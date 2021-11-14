EntGo example
-------------

We will design ER models and generate Ent entities by example. Let's start.
Given we have a Player model.

```go
type Player struct {
}
```

to turn `Player` mode into an `Ent`'s entity, it is just enought embed `ent.Schema`, so it becomes:

```go
type Player struct {
    ent.Schema
}
```

After we have scheme embed, we will need to define our schema parts such as `Fields`, `Relationships` and, if neccessary, some database specifics like `indexes`.
Let's add some fields to our `Player` entity, e.g. `nickname`, `email` and `scores`:

```go
func (Player) Fields() []ent.Field {
	return []ent.Field{
		field.String("nickname"),
		field.String("email"),
		field.Int("scores"),
	}
}
```

In official introduction of Entgo, it says to init model by running codegeltool. Go ahead and install it first:

```bash
go get entgo.io/ent/cmd/ent
```

After you got a tool, time to generate you first template:
```bash
go run entgo.io/ent/cmd/ent init Player
```

The `init` command will create a folder `ent` which will contain:
```
./ent
    schema/
        player.go <---- your template for Player struct
    generate.go
```

You will find it has the same structure as we describe in the beginning, plus an aadditional method:
```
func (Player) Edges() []ent.Edge {
	return nil
}
```

Now, you can copy your `Field()` method into `./ent/schema/player.go`. To verify that you did all right, run `describe` command, you should get something like bellow:
```
➜  entgo-example git:(main) go run entgo.io/ent/cmd/ent describe ./ent/schema
Player:
        +----------+--------+--------+----------+----------+---------+---------------+-----------+---------------------------+------------+
        |  Field   |  Type  | Unique | Optional | Nillable | Default | UpdateDefault | Immutable |         StructTag         | Validators |
        +----------+--------+--------+----------+----------+---------+---------------+-----------+---------------------------+------------+
        | id       | int    | false  | false    | false    | false   | false         | false     | json:"id,omitempty"       |          0 |
        | nickname | string | false  | false    | false    | false   | false         | false     | json:"nickname,omitempty" |          0 |
        | email    | string | false  | false    | false    | false   | false         | false     | json:"email,omitempty"    |          0 |
        | scores   | int    | false  | false    | false    | false   | false         | false     | json:"scores,omitempty"   |          0 |
```

We can remove `./player.go` file, as it is not gonna be needed anymore. From now on, you will edit your models in `./ent/schema` folder.

Now, you migth wonder how do we call SQL queries. To do that, you would need to generate your SQL interfaces for your schema. Invoke the following:

Full version:
```bash
go run entgo.io/ent/cmd/ent generate ./ent/schema
```

Short version:
```
go generate ./ent
```

Inspect results in `./ent`. Previously, we had there our schame file. Now, there are bunch of files, do not edit them yet. Let's understand their purpose by using it.
let's pick our `Postgres` as our RDMS engine. Steps are:

Spin up `Postgres` server in docker:
```bash
docker run -it --rm --name dev-postgres -e POSTGRES_DB=entgo_example -e POSTGRES_PASSWORD=topsecret -p 5432:5432 postgres
```

our connection setttings:
```
Host: localhost
Port: 5432
User: postgres
Password: topsecret
Database: entgo_example
```

Instead of using `database/sql` connect to our db:
```go
sql.Open("postgres", "user=postgres password=topsecret dbname=entgo_example sslmode=disable")
```

`Ent` provides out of box solution with the same interface:
```go
ent.Open("postgres", "user=postgres password=topsecret dbname=entgo_example sslmode=disable")
```

At this moment, we will be ready to use our database, but before we create our first `Player` record, let's have a look into our SQL schema with again, built-in `Ent` tool fot it:

```go
// Dump migration changes to stdout.
if err := client.Schema.WriteTo(ctx, os.Stdout); err != nil {
	log.Fatalf("failed printing schema changes: %v", err)
}
// results in:
/* 
BEGIN;
CREATE TABLE IF NOT EXISTS "players"("id" bigint GENERATED BY DEFAULT AS IDENTITY NOT NULL, "nickname" varchar NOT NULL, "email" varchar NOT NULL, "scores" bigint NOT NULL, PRIMARY KEY("id"));
COMMIT;
*/
```

if that looks good, let's run migrate again our database:
```go
if err := client.Schema.Create(context.Background()); err != nil {
	log.Fatalf("failed creating schema resources: %v", err)
}
```

let's check it in DB:
```
➜  entgo-example git:(main) ✗ psql -U postgres -d entgo_example -h localhost -p 5432
Password for user postgres: 
psql (11.5, server 14.0 (Debian 14.0-1.pgdg110+1))
WARNING: psql major version 11, server major version 14.
         Some psql features might not work.
Type "help" for help.

entgo_example=# \dt
          List of relations
 Schema |  Name   | Type  |  Owner   
--------+---------+-------+----------
 public | players | table | postgres
(1 row)
```

We moving to next thing. Creating our first entry in our `players` table.

```go
player, err := client.Player.Create().
	SetNickname("John").
	SetEmail("info@tsyren.org").
	SetScores(1).
	Save(ctx)
if err != nil {
	log.Fatalln(err)
}
// output:
// 2021/11/14 13:57:57 New player was created Player(id=1, nickname=John, email=info@tsyren.org, scores=1)
```

verifing in db:
```
entgo_example=# SELECT * FROM players;
 id | nickname |      email      | scores 
----+----------+-----------------+--------
  1 | John     | info@tsyren.org |      1
(1 row)
```

let's verify it using `Ent` query interface:

```go
player, err = client.Player.Query().Where(playerEnt.Nickname("John")).Only(ctx)
if err != nil {
	log.Fatalln(err)
}
// output:
// 2021/11/14 14:06:13 Found player was created Player(id=3, nickname=John, email=info@tsyren.org, scores=1)

```

Tsyren Ochirov (c) 2021
