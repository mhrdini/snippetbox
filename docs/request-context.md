# Request Context

Every `http.Request` that our middleware and handlers process has a `context.Context` object
embedded in it.

It stores information during the lifetime during the request.

It is a common use-case to use the context to pass information between middleware and handlers.

A few things to note:

- We don't update the context for a request directly. We create a new copy of the `http.Request` object with a new context with it.
- Use context values only for request-scoped data that transits processes and APIs.

## Authentication

We can use it to check if a user is authenticated once in some middleware and share/publish the result of
this check to all our other middleware and handlers.

## "Updating" the request context

Using `context.WithValue()` and
`http.Request.WithContext()` in tandem to update the context, for example:

```golang
// create a copy of the context with a new key-value pair, where the value is of type any
ctx = context.WithValue(r.Context(), "isAuthenticated", true)
// create new copy of the request containing the new context
r = r.WithContext(ctx)
```

## Retrieving context values

Because context values are of type `any`, we need to assert them to their original type when
retrieving them:

```golang
isAuthenticated, ok := r.Context().Value("isAuthenticated").(bool)
if !ok {
  return errors.New("could not convert value to bool")
}
```

## Avoiding key collision

It is good practice to create your own custom type which is used for context keys that we create, in case of collisions with context keys set by third-party packages.

```golang
type contextKey string
const isAuthenticatedContextKey = contextKey("isAuthenticated")
```
