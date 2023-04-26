# Project Structure

```sh
snippetbox
|
+-- cmd               # application-specific code for the executable application(s) in the project
|    |
|    +-- web          # the executable web application
|
+-- pkg               # non-application-specific, potentially reusable code like validation helpers and SQL database models for the project
|
+-- ui                # user-interface assets used by the web application\
     |
     +-- html         # html templates
     |
     +-- static       # static files like CSS and images
```

## Why are we using this structure?

1. It gives a clean separation between Go and non-Go assets.
2. It scales really nicely should we want to add another executable application to the project.
