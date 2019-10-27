# Project starter

Ahoy! This is the project starter. I've done some
of the hard parts for you. Good luck.

## Building the project

Run `go build -mod vendor` to compile your project and then
run `./classproject` or `./classproject` on Windows
in order to run the app. Preview the app by visiting
[http://localhost:8080](http://localhost:8080) if you're
running it locally. Or, preview using the URL provided
by cloud9 if you're running on cloud9, natch.

You will need to provide a `DATABASE_URL` environment
variable. This contains the connection information
for Postgres. Yours might look something like this

```
DATABASE_URL=postgres://testing_user:testing_pass@localhost:7000/testing_db
```

Of course, it will not be precisely that. You'll need to
figure out where and how to run [Postgres](https://www.postgresql.org/)
yourself.

In production, on Heroku, it would look different.
Take a look at `event_models.go` to see the (incomplete)
database schema that I provided you.

I like to run Postgres locally using Docker. If you're
into that, you might find something like this useful

```
docker run -p 7000:5432 -e POSTGRES_PASSWORD=testing_pass -e POSTGRES_USER=testing_user -e POSTGRES_DB=testing_db postgres
```

Also, word to the wise: you can't run docker on cloud9
because cloud9 runs in docker methinks.

Many people like [Postgres.app](https://postgresapp.com/) for macs.
I don't know of anything similar for Windows.

## Deploying to Heroku

I included a `go.mod` file and a `Procfile` so you
should have no trouble deploying to Heroku. Of course,
your team should only have one Heroku app (usually)
and you should all be "collaborators" on it. You're
welcome to deploy to other places too---I don't care.

You can read about [adding Postgres to Heroku here](https://www.heroku.com/postgres).

## What is here

| File                      | Role                                                                                                                      |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| ./README.md               | This file!                                                                                                                |
| ./server.go               | Entrypoint for the code; contains the `main` function                                                                     |
| ./routes.go               | Maps your URLs to the controllers/handlers for each URL                                                                   |
| ./event_models.go         | Defines your data structure and logic. I put in a few default events.                                                     |
| ./index_controllers.go    | Controllers related to the index (home) page                                                                              |
| ./templates.go            | Sets up your templates from which you generate HTML                                                                       |
| ./templates               | Directory with your templates. You'll need more templates ;P                                                              |
| ./templates/layout.gohtml | The "base" layout for your HTML frontend                                                                                  |
| ./templates/index.gohtml  | The template for the index (home) page                                                                                    |
| ./static.go               | Sets up the static file server (see next entry)                                                                           |
| ./staticfiles             | Directory with our "static" assets like images, css, js                                                                   |
| ./Procfile                | A file that helps heroku run your app                                                                                     |
| ./go.mod                  | [Go modules file](https://www.kablamo.com.au/blog/2018/12/10/just-tell-me-how-to-use-go-modules). Lists our dependencies. |
| ./go.sum                  | A "checksum" file that says precisely what versions of our dependencies need to be installed.                             |
| ./vendor                  | A directory containing our dependencies                                                                                   |

## Automatic reload

If you want your app to reload every time you make a
change you can do the following. First
install reflex with `go get github.com/cespare/reflex`.

Then, run

```
~/go/bin/reflex -d fancy -r'\.go' -r'\.gohtml' -s -R vendor. -- go run *.go
```

or something like that. Look at the reflex documentation. Automatic
reload is pretty rad while developing. As a general rule, developers
want to let computers do what computers are good at (tasks that can be automated)
so that they, the developers, can focus on what they are good at: the
logic of the product.

## Other info

Information about the class final project is distributed between
a few places and I apologize for this. You can find information
about the project in the following places:

- This page (which you're likely looking at in your own repo)
- The sprint-1 assignment page:
  [656](https://www.656.mba/#assignments/project-sprint-1) &
  [660](https://www.660.mba/#assignments/project-sprint-1)
- The "about" repo for the class:
  [https://github.com/yale-mgt-656-fall-2019/about/blob/master/class-project.md](https://github.com/yale-mgt-656-fall-2019/about/blob/master/class-project.md)
- The grading code:
  [https://github.com/yale-mgt-656-fall-2019/project-grading](https://github.com/yale-mgt-656-fall-2019/project-grading)
- The reference solution
  [http://project.solutions.656.mba/](http://project.solutions.656.mba/)
- Your class' Piazza page:
  [656](https://piazza.com/yale/fall2019/mgt656) &
  [660](https://piazza.com/yale/fall2019/mgt660)
- My comments in class
