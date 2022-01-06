package controllers

import (
	"fiber-mongo-api/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
GetUsers | @Desc: Get all users |
@Method: GET |
@Route: "api/v1/users" |
@Auth: Public
*/
func GetUsers(c *fiber.Ctx) error {
	filter := bson.D{{}}
	cursor, err := models.UserCollection.Find(c.Context(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	var users []models.User = make([]models.User, 0)

	if err := cursor.All(c.Context(), &users); err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(users)
}

/*
GetUser | @Desc: Get user by id |
@Method: GET |
@Route: "api/v1/users/:id" |
@Auth: Public
*/
func GetUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": idParam + " is not a valid id!"})
	}

	filer := bson.D{{Key: "_id", Value: userID}}
	userRecord := models.UserCollection.FindOne(c.Context(), filer)
	if userRecord.Err() != nil {
		return c.Status(400).JSON(fiber.Map{"message": "No user with id: " + idParam + " was found!"})
	}

	user := &models.User{}
	userRecord.Decode(user)

	return c.JSON(user)
}

/*
CreateUser | @Desc: Create new user |
@Method: POST |
@Route: "api/v1/users" |
@Auth: Public
*/
func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	insertionResult, err := models.UserCollection.InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	// get the just inserted record in order to return it as response
	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := models.UserCollection.FindOne(c.Context(), filter)

	// decode the Mongo record into Employee
	createdUser := &models.User{}
	createdRecord.Decode(createdUser)

	return c.Status(201).JSON(createdUser)
}

/*
UpdateUser | @Desc: Update user by id |
@Method: PUT |
@Route: "api/v1/users/:id" |
@Auth: Private
*/
func UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": idParam + " is not a valid id!"})
	}

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	filter := bson.D{{Key: "_id", Value: userID}}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "name", Value: user.Name},
				{Key: "email", Value: user.Email},
			},
		},
	}

	after := options.After
	userRecord := models.UserCollection.FindOneAndUpdate(c.Context(), filter, update, &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	})

	if userRecord.Err() != nil {
		return c.Status(400).JSON(fiber.Map{"message": "No user with id: " + idParam + " was found!"})
	}

	updatedUser := &models.User{}
	userRecord.Decode(updatedUser)

	return c.JSON(updatedUser)
}

/*
DeleteUser | @Desc: Delete user by id |
@Method: DELETE |
@Route: "api/v1/users/:id" |
@Auth: Private
*/
func DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": idParam + " is not a valid id!"})
	}

	filer := bson.D{{Key: "_id", Value: userID}}
	userRecord := models.UserCollection.FindOneAndDelete(c.Context(), filer)
	if userRecord.Err() != nil {
		return c.Status(400).JSON(fiber.Map{"message": "No user with id: " + idParam + " was found!"})
	}

	return c.JSON(fiber.Map{"message": "User was deleted!"})
}
