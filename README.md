<p align="center">
    <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white"/>
    <img src="https://img.shields.io/badge/MongoDB-4EA94B?style=for-the-badge&logo=mongodb&logoColor=white"/>
</p>

# Mango
A MongoDB schema access control framework featuring struct tags

______________________________________________________________________

## Installation
Installing Mango is simple! Simply add it to your project using go modules:
```bash
go get github.com/aphix-dev/mango
```

______________________________________________________________________

## Usage
To get started with Mango in your project, create a mango config.
```go
import "github.com/aphix-dev/mango"

var conf MangoConfig = DefaultConfig()
```
Create a struct representing a MongoDB model. For this example, we will be defining a "User" model:
```go
type User struct {
	ID              primitive.ObjectID   `json:"_id" bson:"_id"`
	Email           string               `json:"email" bson:"email"`
	Password        string               `json:"password" bson:"password"`
	FirstName       string               `json:"firstName" bson:"firstName"`
	LastName        string               `json:"lastName" bson:"lastName"`
	PurchaserSecret []primitive.ObjectID `json:"purchaserSecret" bson:"purchaseCourses,"`
}
```
Then add `access` struct tags to each field. These will be used by Mango when "trimming" data that shouldn't be there:
```go
type User struct {
	ID              primitive.ObjectID   `json:"_id" bson:"_id" access:"create"`
	Email           string               `json:"email" bson:"email" access:"create,update,priv"`
	Password        string               `json:"password" bson:"password" access:"create"`
	FirstName       string               `json:"firstName" bson:"firstName" access:"create,update,pub,priv"`
	LastName        string               `json:"lastName" bson:"lastName" access:"create,update,pub,priv"`
	PurchaserSecret []primitive.ObjectID `json:"purchaserSecret" bson:"purchaseCourses," access:"priv,purchaserOnly"`
}
```

If we want to trim the data of non-create fields, or fields that should not be set by the client upon the creation of a new User, we will use Mango's `trim` function like so:
```go
// assume userPayload (type: User) has been taken from an HTTP request body
mango.Trim(userPayload, mango.CREATE, conf)
```

Make note that Mango's default config can be extended to include custom `access` tags as well:
```go
import "github.com/aphix-dev/mango"

const PURCHASER_ONLY = iota

var conf MangoConfig = DefaultConfig().Extend(map[int]string{
	PURCHASER_ONLY: "purchaserOnly",
})
```
In this case, in the following code, only the fields tagged with `access:"purchaserOnly"`  will be returned and the other fields will be empty or nil.
```
// once again, assume userPayload (type: User) has been taken from an HTTP request body
mango.Trim(userPayload, PURCHASER_ONLY, conf)
```
Mango's default config comes with "create", "update", "pub", and "priv" tags by default. These tags are useful for general REST API use cases:
- `access:"create"` specifies that a field can be set upon creation
- `access:"update"` specifies a field that can be modified
- `access:"pub"` specifies a field that can be accessed by all users
- `access:"priv"` specifies a field that should only be accessible to the owning entity (for instance, a user's home address would likely be a field they would prefer to be kept privately rather than publicly).

-----
## Feedback
Suggestions and issues can be reported in my Discord server: https://discord.com/invite/qwz7nW43wf
