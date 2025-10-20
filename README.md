# friendlymongo

> **MongoDB for Golang made easy (maybe).**

[![Go Report Card](https://goreportcard.com/badge/github.com/pmatteo/friendlymongo)](https://goreportcard.com/report/github.com/pmatteo/friendlymongo)
![GitHub Release](https://img.shields.io/github/v/release/pmatteo/friendlymongo?include_prereleases\&display_name=tag\&style=flat-square)

---

## 📦 Installation

```bash
go get -u github.com/pmatteo/friendlymongo
```

> **Note:**
> `friendlymongo` is tested with Go versions `1.18` and `1.22`, using different MongoDB versions.

---

## 🚀 Usage

A simple example is available [here](https://github.com/pmatteo/friendlymongo/tree/main/_examples/simple).

---

### 🧩 Client

`friendlymongo` provides an easy way to manage a `mongo.Client` instance as a **singleton**.

* Use `SetInstance(url)` to initialize it.
* Retrieve it anywhere using `GetInstance()`.

The wrapper type `MongoClient` exposes convenient helpers:

* `Client()` → returns the `mongo.Client` instance
* `Database(dbName)` → returns a new `mongo.Database` instance

#### Example

```go
i := friendlymongo.SetInstance("mongodb://username:password@localhost:27017")
c := i.Client()
db := i.Database("user")
```

or equivalently:

```go
friendlymongo.SetInstance("mongodb://username:password@localhost:27017")
c := friendlymongo.GetInstance().Client()
db := friendlymongo.GetInstance().Database("user")
```

---

### 🧱 Model

The `Model` interface defines lifecycle hooks for any struct intended for repository usage.

```go
type Model interface {
    OnCreate()
    OnUpdate()
    OnReplace()
}
```

A convenient `BaseModel` is included, which automatically handles:

* MongoDB `ObjectID`
* `created_at` and `updated_at` timestamps

#### Example

```go
type UserProfile struct {
    *friendlymongo.BaseModel
    Name      string
    Surname   string
    Email     string
    BirthDate time.Time
}
```

---

### 🗂 Repository

The core purpose of `friendlymongo` is to simplify MongoDB access in Go, eliminating repetitive boilerplate.

The **repository pattern** is implemented using **Go generics**, providing out-of-the-box CRUD operations for your model type.

#### Example

```go
// Retrieve a BaseRepository instance for the UserProfile type
repo := friendlymongo.NewBaseRepository(db, "userProfile", &UserProfile{})

user := NewUserProfile("John", "Doe", "john.doe@test.com", birthday)

// Insert the user into the database
if err := repo.InsertOne(context.Background(), user); err != nil {
    panic(err)
}
```

---

### 🧮 Pipeline Stage Builder

`BaseRepository` provides an `Aggregate()` method that accepts a `mongo.Pipeline`.

To simplify pipeline creation, `friendlymongo` includes a **stage builder** — a fluent, structured way to define aggregation stages.

#### Example

```go
pipeline := friendlymongo.
    NewStageBuilder().
    Match("name_filter", bson.M{"name": "John"}).
    Lookup("roles_lookup", "user_role", "_id", "fk", "role").
    Match("filter_admin", bson.M{"role.name": "admin"}).
    Build()
```

The builder includes helpers for several common stages.
You can also add custom stages using the `AddStage()` method.

---

### ⚙️ Operators

`friendlymongo` includes simplified helpers for some MongoDB operators such as `$push` and `$map`.
More will be added over time.

#### Example

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

Would you like me to make it **ready for pkg.go.dev** formatting (i.e. Markdown tuned for Go documentation sites)?
That would adjust headings, example indentation, and link formatting for GoDoc-style readability.
