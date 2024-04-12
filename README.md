# friendlymongo

MongoDB and Golang made easy (maybe).

[![Go Report Card](https://goreportcard.com/badge/github.com/pmatteo/friendlymongo)](https://goreportcard.com/report/github.com/pmatteo/friendlymongo)

![GitHub Release](https://img.shields.io/github/v/release/pmatteo/friendlymongo?include_prereleases&display_name=tag&style=flat-square)

## Installation

`go get -u github.com/pmatteo/friendlymongo`

Note that `friendlymongo` is tested on Go `1.18` and `1.22` using different vversion of MongoDB.

## Usage

You can find a simple exmaple [here](https://github.com/pmatteo/friendlymongo/tree/main/_examples/simple).

### Client

`friendlymongo` has a simple way to handle mongo.Client instance as a singelton. You can set the instance using `setInstance(url)` method and then get it everywhere with `GetInstance()`.
It uses a wrapper class called `MongoClient`, to access the client instance use `GetInstance().Client()` or to get a new mongo.Database instance use `GetInstance().Database(dbName)`.

```go
i := friendlymongo.SetInstance("mongodb://username:password@localhost:27017")
c := i.Client()
db := i.Database("user")
```

or

```go
friendlymongo.SetInstance("mongodb://username:password@localhost:27017")
c := friendlymongo.GetInstance().Client()
db := friendlymongo.GetInstance().Database("user")
```

### Model

`Model` is an interface that defines a series of method that every struct must implemets to being able to be used by
the repository implmenetation.

```go
type Model interface {
    OnCreate()
    OnUpdate()
    OnReplace()
}
```

The library also comes with a simple BaseModel whcih already handels ObjectID, created and updated timestamp that can be used.

```go
type UserProfile struct {
    *friendlymongo.BaseModel
    // Other fields 
    Name      string
    Surname   string
    Email     string
    BirthDate time.Time
}
```

### Repository

The main goal of the package is allow basic Mongo functionalities access without tons of boilarplate code.

The implementation of the repository pattern using go generics is a way to have some operations mapped and accessible without any effort other than define the data structure you have to deal with.

```go
// retrieve a BaseRepository instance fot the userProfile type
repo := friendlymongo.NewBaseRepository(db, "userProfile", &userProfile{})

...

user := NewUserProfile("John","Doe","john.doe@test.com",birthday)
// Insert the user into the database
err := repo.InsertOne(context.Background(), user)
if err != nil {
    panic(err)
}
```

### Pipeline Stage Builder

BaseRepository offers an `Aggregate` method for Mongo's aggregation pipelines feature. It requires an instance of `mongo.Pipeline` as argument.

For some basic (or even not) pipeline, `friendlymongo` implements a simple stage builder that ìwould help developers create their stages in a more structured way and readable way.

```go
pipeline := friendlymongo.
    NewStageBuilder().
    Match("name_filter", bson.M{"name": "John"}).
    Lookup("roles_lookup", "user_role","_id","fk","role").
    Match("filter_admin", bson.M{"role.name": "admin"}).
    Build()
```

At the moment the stage builder already implements a method for a subset of possible stages, other can be added using the `AddStage` method.

### Operators

`friendlymongo` offers a semplification for some mongo operators like `$push` or `$map`. More will be added in the future.

```go
fm.NewStageBuilder().
    Lookup("product_lookup", "product", "products", "ean", "products").
    Unwind("unwind_products", "$products").
    Group("category_id_group", bson.M{
        "_id": bson.M{
            "orderId":  "$_id",
            "category": "$products.category",
        },
        "status":   fm.First("$status"),
        "products": fm.Push("$products"),
    }).
    Group("productsByCategory", bson.M{
        "_id":                "$_id.orderId",
        "status":             fm.First("$status"),
        "productsByCategory": fm.Push("category", "$_id.category", "products", "$products"),
    }).
    Project("final_project", bson.M{
        "_id":    1,
        "status": 1,
        "grouped_products": fm.ArrayToObject(
            fm.Map("$productsByCategory", "cat", "$$cat.category", "$$cat.products"),
        ),
    })
```
