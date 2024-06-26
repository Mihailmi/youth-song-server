package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"song.com/v1/entity"
)

type Song struct {
	pptx    string
	obj_id  string `bson:"_id"`
	version int32
}

func writeDB() {
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://uru3xmclir6cuzp73mkz:61L8uYL9CudlbOL9PYLR@n1-c2-mongodb-clevercloud-customers.services.clever-cloud.com:27017,n2-c2-mongodb-clevercloud-customers.services.clever-cloud.com:27017/bvm47vjpnnaq4zo?replicaSet=rs0"))
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	databases, err := client.ListDatabases(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(databases)

}

func sendEmail(text string, sender string) {
	from := "mysenderg@gmail.com"
	password := "cyqorpgewebmbsue"

	// Receiver email address.
	to := []string{
		// "miha65079@gmail.com",
		sender,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// message := []byte(fmt.Sprintf("%s", text))
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: Confirming your account\r\n\r\n"+
		"%s\r\n", from, sender, text))
	fmt.Println(msg)

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)

	if err != nil {
		log.Fatal(400, err)
		return
	}

	log.Println(200, "Email Sent Successfully!")
}

func usersList(ctxx *gin.Context) string {
	// mongo
	result := ""
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	users := client.Database("nto2023").Collection("users")

	colresusers, err := users.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var logins_of_users []bson.M
	if err = colresusers.All(ctx, &logins_of_users); err != nil {
		log.Fatal(err)
	}

	// for _, user := range logins_of_users {
	// 	ctxx.JSON(200, user)
	// }
	ctxx.JSON(200, logins_of_users)
	return result
}

func hello(ctx *gin.Context) {
	// writeDB()

	// from := "mysenderg@gmail.com"
	// password := "cyqorpgewebmbsue"

	// // Receiver email address.
	// to := []string{
	// 	"mihailkorcik@gmail.com",
	// }

	// // smtp server configuration.
	// smtpHost := "smtp.gmail.com"
	// smtpPort := "587"

	// // Message
	// message := []byte(usersList(ctx))
	// fmt.Println(message)
	// // Authentication
	// auth := smtp.PlainAuth("", from, password, smtpHost)

	// // Sending email
	// err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)

	// if err != nil {
	// 	ctx.JSON(400, err)
	// 	return
	// }

	ctx.JSON(200, "Hello!")
}

func addUser(c *gin.Context) {
	var newUser entity.User

	// err := c.BindJSON(&newUser)
	newUser.Email = c.PostForm("email")
	if len(strings.Split(newUser.Email, "@")) != 2 {
		c.JSON(400, "Wrong email")
	}

	newUser.Name = c.PostForm("name")
	newUser.Surname = c.PostForm("surname")
	newUser.Email = c.PostForm("email")
	newUser.Password = c.PostForm("password")
	newUser.Role = c.PostForm("role")
	newUser.JobTitle = c.PostForm("jobTilte")
	newUser.Login = c.PostForm("login")
	newUser.IsActive = false
	fmt.Println(newUser)

	// mongo
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err3 := client.Connect(ctx)
	if err3 != nil {
		log.Fatal(err3)
	}

	defer client.Disconnect(ctx)
	err3 = client.Ping(ctx, readpref.Primary())
	if err3 != nil {
		log.Fatal(err3)
	}

	collection := client.Database("nto2023").Collection("users")
	colres, err3 := collection.Find(ctx, bson.M{"login": newUser.Login})
	if err3 != nil {
		log.Fatal(err3)
	}
	var logins []bson.M
	if err3 = colres.All(ctx, &logins); err3 != nil {
		log.Fatal(err)
	}
	fmt.Println("len of logins ", len(logins))
	if len(logins) > 0 {
		c.JSON(400, "User with this login exists")
		return
	}

	result, err := collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		log.Fatal(err)
	}
	text := string("Activate your account https://nto.onrender.com/confirmation/" + strings.Split(fmt.Sprintf("%s", result.InsertedID)[10:], "\"")[0])
	sendEmail(text, newUser.Email)
	// fmt.Println(text)

	c.JSON(200, "Successfully added user")

}

func confirmation(c *gin.Context) {
	objectId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(objectId)
	client, err3 := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err3 != nil {
		log.Fatal(err3)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err3 = client.Connect(ctx)
	if err3 != nil {
		log.Fatal(err3)
	}

	defer client.Disconnect(ctx)
	err3 = client.Ping(ctx, readpref.Primary())
	if err3 != nil {
		log.Fatal(err3)
	}
	users := client.Database("nto2023").Collection("users")

	// checking user
	userinf, err3 := users.Find(ctx, bson.M{"_id": objectId})
	if err3 != nil {
		log.Fatal(err3)
	}
	var info_user []bson.M
	if err3 = userinf.All(ctx, &info_user); err3 != nil {
		log.Fatal(err3)
	}

	if len(info_user) != 1 {
		fmt.Fprintf(c.Writer, "No such user")
		fmt.Fprintln(c.Writer, len(info_user))
		return
	}

	result, err3 := users.UpdateOne(
		ctx,
		bson.M{"_id": objectId},
		bson.D{
			{"$set", bson.D{{"isActive", true}}},
		},
	)

	if err3 != nil {
		log.Fatal(err3)
	}

	fmt.Println(result)

	fmt.Fprintf(c.Writer, "User email confirmed!")

}

func addEvent(c *gin.Context) {
	var newEvent entity.Event

	// err := c.BindJSON(&newUser)
	newEvent.Name = c.PostForm("name")
	newEvent.Data = c.PostForm("data")
	newEvent.Description = c.PostForm("description")
	newEvent.Registration = c.PostForm("registration")
	fmt.Println(newEvent)

	// mongo
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err3 := client.Connect(ctx)
	if err3 != nil {
		log.Fatal(err3)
	}

	defer client.Disconnect(ctx)
	err3 = client.Ping(ctx, readpref.Primary())
	if err3 != nil {
		log.Fatal(err3)
	}

	collection := client.Database("nto2023").Collection("events")
	colres, err3 := collection.Find(ctx, bson.M{"name": newEvent.Name})
	if err3 != nil {
		log.Fatal(err3)
	}
	var eventsFrom []bson.M
	if err3 = colres.All(ctx, &eventsFrom); err3 != nil {
		log.Fatal(err)
	}

	fmt.Println("len of events ", len(eventsFrom))
	if len(eventsFrom) > 0 {
		c.JSON(400, "Event with this exists")
		return
	}

	result, err := collection.InsertOne(context.TODO(), newEvent)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	c.JSON(200, "Successfully added event\n")

}

func getEvents(ctxx *gin.Context) {
	// mongo

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	users := client.Database("nto2023").Collection("events")

	colresusers, err := users.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var allEvents []bson.M
	if err = colresusers.All(ctx, &allEvents); err != nil {
		log.Fatal(err)
	}

	ctxx.JSON(200, allEvents)
}

func getUsers(ctxx *gin.Context) {
	// mongo

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	users := client.Database("nto2023").Collection("users")

	colresusers, err := users.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var allUsers []bson.M
	if err = colresusers.All(ctx, &allUsers); err != nil {
		log.Fatal(err)
	}

	ctxx.JSON(200, allUsers)
}

func userUpdate(c *gin.Context) {

	var newUser entity.User

	// err := c.BindJSON(&newUser)
	newUser.Email = c.PostForm("email")
	if len(strings.Split(newUser.Email, "@")) != 2 {
		c.JSON(400, "Wrong email")
	}

	newUser.Name = c.PostForm("name")
	newUser.Surname = c.PostForm("surname")
	newUser.Email = c.PostForm("email")
	newUser.Password = c.PostForm("password")
	newUser.Role = c.PostForm("role")
	newUser.JobTitle = c.PostForm("jobTilte")
	newUser.Login = c.PostForm("login")

	objectId, err := primitive.ObjectIDFromHex(c.PostForm("id"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(objectId)
	client, err3 := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err3 != nil {
		log.Fatal(err3)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err3 = client.Connect(ctx)
	if err3 != nil {
		log.Fatal(err3)
	}

	defer client.Disconnect(ctx)
	err3 = client.Ping(ctx, readpref.Primary())
	if err3 != nil {
		log.Fatal(err3)
	}
	users := client.Database("nto2023").Collection("users")

	// checking user
	userinf, err3 := users.Find(ctx, bson.M{"_id": objectId})
	if err3 != nil {
		log.Fatal(err3)
	}
	var info_user []bson.M
	if err3 = userinf.All(ctx, &info_user); err3 != nil {
		log.Fatal(err3)
	}

	if len(info_user) != 1 {
		c.JSON(404, "No such user")
		return
	}

	result, err3 := users.UpdateOne(
		ctx,
		bson.M{"_id": objectId},
		bson.D{
			{"$set", bson.D{{"name", newUser.Name}}},
			{"$set", bson.D{{"surname", newUser.Surname}}},
			{"$set", bson.D{{"email", newUser.Email}}},
			{"$set", bson.D{{"role", newUser.Role}}},
			{"$set", bson.D{{"jobTilte", newUser.JobTitle}}},
			{"$set", bson.D{{"login", newUser.Login}}},
		},
	)

	if err3 != nil {
		log.Fatal(err3)
	}

	fmt.Println(result)

	c.JSON(200, "Successfully updated user\n")
}

func sendEmailto(c *gin.Context) {
	address := c.PostForm("address")
	text := c.PostForm("text")
	sendEmail(text, address)
	c.JSON(200, "Successfully sent email\n")
}

func findDB(ctxx *gin.Context) {
	collName := ctxx.PostForm("collection")
	fild := ctxx.PostForm("fild")
	whatToFind := ctxx.PostForm("value")
	// mongo

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	users := client.Database("nto2023").Collection(collName)

	colresusers, err := users.Find(ctx, bson.M{fild: whatToFind})
	if fild == "no" {
		colresusers, err = users.Find(ctx, bson.M{})
	}
	if err != nil {
		log.Fatal(err)
	}

	var allData []bson.M
	if err = colresusers.All(ctx, &allData); err != nil {
		log.Fatal(err)
	}

	ctxx.JSON(200, allData)
}

func updateById(c *gin.Context) {

	count := 0

	f1 := c.PostForm("f1")
	v1 := c.PostForm("v1")
	f2 := c.PostForm("f2")
	if f2 == "no" {
		count = 1
	}
	v2 := c.PostForm("v2")
	f3 := c.PostForm("f3")
	if f3 == "no" {
		count = 2
	}
	v3 := c.PostForm("v3")
	f4 := c.PostForm("f4")
	if f4 == "no" {
		count = 3
	}
	v4 := c.PostForm("v4")
	f5 := c.PostForm("f5")
	if f5 == "no" {
		count = 4
	}
	v5 := c.PostForm("v5")
	if f5 != "" && f5 != "no" {
		count = 5
	}

	collection := c.PostForm("collection")
	objectId, err := primitive.ObjectIDFromHex(c.PostForm("id"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(objectId)
	client, err3 := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err3 != nil {
		log.Fatal(err3)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err3 = client.Connect(ctx)
	if err3 != nil {
		log.Fatal(err3)
	}

	defer client.Disconnect(ctx)
	err3 = client.Ping(ctx, readpref.Primary())
	if err3 != nil {
		log.Fatal(err3)
	}
	users := client.Database("nto2023").Collection(collection)

	// checking user
	userinf, err3 := users.Find(ctx, bson.M{"_id": objectId})
	if err3 != nil {
		log.Fatal(err3)
	}
	var info_user []bson.M
	if err3 = userinf.All(ctx, &info_user); err3 != nil {
		log.Fatal(err3)
	}

	if len(info_user) != 1 {
		c.JSON(404, "No")
		return
	}

	var bsonForAdd bson.D

	bsonForAdd = bson.D{
		{Key: "$set", Value: bson.D{{Key: f1, Value: v1}}},
		{Key: "$set", Value: bson.D{{Key: f2, Value: v2}}},
		{Key: "$set", Value: bson.D{{Key: f3, Value: v3}}},
		{Key: "$set", Value: bson.D{{Key: f4, Value: v4}}},
		{Key: "$set", Value: bson.D{{Key: f5, Value: v5}}},
	}
	if count == 4 {
		bsonForAdd = bson.D{
			{Key: "$set", Value: bson.D{{Key: f1, Value: v1}}},
			{Key: "$set", Value: bson.D{{Key: f2, Value: v2}}},
			{Key: "$set", Value: bson.D{{Key: f3, Value: v3}}},
			{Key: "$set", Value: bson.D{{Key: f4, Value: v4}}},
		}
	}
	if count == 3 {
		bsonForAdd = bson.D{
			{Key: "$set", Value: bson.D{{Key: f1, Value: v1}}},
			{Key: "$set", Value: bson.D{{Key: f2, Value: v2}}},
			{Key: "$set", Value: bson.D{{Key: f3, Value: v3}}},
		}
	}
	if count == 2 {
		bsonForAdd = bson.D{
			{Key: "$set", Value: bson.D{{Key: f1, Value: v1}}},
			{Key: "$set", Value: bson.D{{Key: f2, Value: v2}}},
		}
	}
	if count == 1 {
		bsonForAdd = bson.D{
			{Key: "$set", Value: bson.D{{Key: f1, Value: v1}}},
		}
	}

	result, err3 := users.UpdateOne(
		ctx,
		bson.M{"_id": objectId},
		bsonForAdd,
	)

	if err3 != nil {
		log.Fatal(err3)
	}

	fmt.Println(result)

	c.JSON(200, "Successfully updated user\n")
}

func createEmty(c *gin.Context) {
	collName := c.PostForm("collection")

	// mongo
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err3 := client.Connect(ctx)
	if err3 != nil {
		log.Fatal(err3)
	}

	defer client.Disconnect(ctx)
	err3 = client.Ping(ctx, readpref.Primary())
	if err3 != nil {
		log.Fatal(err3)
	}

	collection := client.Database("nto2023").Collection(collName)

	result, err := collection.InsertOne(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(text)

	c.JSON(200, result.InsertedID)
}

func uploadSongToMongoDB(c *gin.Context) {

	var songToUpload Song
	// Get handler for filename, size and headers

	file, _ := c.FormFile("pptx")
	songToUpload.obj_id = c.PostForm("id")
	songToUpload.version = 1

	// songToUpload.pptx

	//mongo
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("appSongs").Collection("songs")

	// filename := filepath.Base(file.Filename)
	// c.SaveUploadedFile(file, filename)
	// file1, err := os.Open(filename)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file1.Close()

	openedFile, _ := file.Open()
	songpptx, _ := ioutil.ReadAll(openedFile)

	songToUpload.pptx = string(songpptx)

	if songToUpload.obj_id != "" {
		objectId, err := primitive.ObjectIDFromHex(songToUpload.obj_id)
		if err != nil {
			log.Fatal(err)
		}
		so, err3 := collection.Find(ctx, bson.M{"_id": objectId})
		if err3 != nil {
			log.Fatal(err)
		}
		var findsongs []bson.M
		if err3 = so.All(ctx, &findsongs); err3 != nil {
			log.Fatal(err)
		}

		if len(findsongs) != 0 {
			c.JSON(409, songToUpload.obj_id+" File already in mongodb")
			return
		}
		colres, err := collection.InsertOne(context.TODO(), bson.D{
			{Key: "_id", Value: objectId},
			{Key: "version", Value: songToUpload.version},
			{Key: "pptx", Value: songToUpload.pptx},
			{Key: "title", Value: file.Filename},
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(colres)
	} else {
		colres, err := collection.InsertOne(context.TODO(), bson.D{
			{Key: "version", Value: songToUpload.version},
			{Key: "pptx", Value: songToUpload.pptx},
			{Key: "title", Value: file.Filename},
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(colres)
	}
	// fmt.Println(songToUpload)
	c.JSON(200, "Uploaded successfully!")
}

func getSongFromMongoDB(c *gin.Context) {

	objectId, err := primitive.ObjectIDFromHex(c.PostForm("id"))
	if err != nil {
		log.Fatal(err)
	}

	//mongo
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("appSongs").Collection("songs")
	filter := bson.M{"_id": objectId}

	var result bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(204, "File isn't in mongodb")
			return
		}
		panic(err)
	}

	str := fmt.Sprintf("%v", result["pptx"])
	// fmt.Println(str)
	dst, err3 := os.Create(c.PostForm("id") + ".pptx")
	if err3 != nil {
		log.Fatal(err)
	}
	dst.Write([]byte(str))
	dst.Close()
	// result["pptx"] = ""
	simplifiedSong := map[string]interface{}{
		"id":      result["_id"],
		"version": result["version"],
		"title":   result["title"],
	}
	c.JSON(200, simplifiedSong)
}

func getListOfSongs(c *gin.Context) {

	//mongo
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("appSongs").Collection("songs")

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and build the songs array
	var songs []bson.M
	if err := cursor.All(ctx, &songs); err != nil {
		log.Fatal(err)
	}

	// Create a new array for storing the simplified song objects
	var simplifiedSongs []map[string]interface{}
	for _, song := range songs {
		simplifiedSong := map[string]interface{}{
			"id":      song["_id"].(primitive.ObjectID).Hex(),
			"version": song["version"],
			"title":   song["title"],
		}
		simplifiedSongs = append(simplifiedSongs, simplifiedSong)
	}

	c.JSON(200, simplifiedSongs)
}

func downloadSongFromMongoDB(c *gin.Context) {
	filename := c.PostForm("id") + ".pptx"

	// Check if file exists
	if _, err := os.Stat(filename); err == nil {
		c.File(filename)
		e := os.Remove(filename)
		if e != nil {
			log.Fatal(e)
		}
	} else if os.IsNotExist(err) {
		c.JSON(204, "No file")
	} else {
		c.JSON(500, "Error:"+err.Error())
	}

}

func changeSongIDInMongoDB(c *gin.Context) {

	// Parse the ID from the request
	objectID, err := primitive.ObjectIDFromHex(c.PostForm("id"))
	if err != nil {
		log.Fatal(err)
	}
	newOjectId, err := primitive.ObjectIDFromHex(c.PostForm("newId"))
	if err != nil {
		log.Fatal(err)
	}

	// Create a new MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Connect to the MongoDB server
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// Access the "appSongs" collection in the "songs" database
	collection := client.Database("appSongs").Collection("songs")
	filter := bson.M{"_id": objectID}
	oldDoc := collection.FindOne(context.Background(), filter)

	var docData bson.M
	err = oldDoc.Decode(&docData)
	if err != nil {
		c.JSON(404, "Error:"+err.Error())
		return
	}

	docData["_id"] = newOjectId

	_, err = collection.InsertOne(context.Background(), docData)
	if err != nil {
		c.JSON(409, newOjectId.String()+" File already in mongodb")
		return
	}

	// Delete the old document
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, "ID changed successfully")
}

func updatePptx(c *gin.Context) {

	// Parse the ID from the request
	objectID, err := primitive.ObjectIDFromHex(c.PostForm("id"))
	if err != nil {
		log.Fatal(err)
	}
	file, _ := c.FormFile("pptx")

	// Create a new MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Connect to the MongoDB server
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// Access the "appSongs" collection in the "songs" database
	collection := client.Database("appSongs").Collection("songs")
	filter := bson.M{"_id": objectID}

	openedFile, _ := file.Open()
	songpptx, _ := ioutil.ReadAll(openedFile)

	collection.FindOneAndUpdate(context.Background(), filter, bson.D{
		{"$set", bson.D{{"pptx", string(songpptx)}}},
		{"$inc", bson.D{{"version", 1}}},
	})

	var result bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(204, "File isn't in mongodb")
			return
		}
		c.JSON(204, err)
		return
	}
	simplifiedSong := map[string]interface{}{
		"id":      result["_id"],
		"version": result["version"],
		"title":   result["title"],
	}
	c.JSON(200, simplifiedSong)

}

func login(c *gin.Context) {

	Ulogin := c.PostForm("login")
	Upassword := c.PostForm("password")

	//mongo
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://firstuser:xwI7zM83v62q5SVj@testcluster1.1brcg.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("appSongs").Collection("users")

	filter := bson.D{
		{Key: "login", Value: Ulogin},
		{Key: "password", Value: Upassword},
	}

	fmt.Println(filter)

	so, err3 := collection.Find(ctx, filter)
	if err3 != nil {
		log.Fatal(err)
	}
	var finduser []bson.M
	if err3 = so.All(ctx, &finduser); err3 != nil {
		log.Fatal(err)
	}

	if len(finduser) == 0 {
		c.JSON(401, "Error login")
		return
	}

	c.JSON(200, "Login success!")
}

func main() {
	server := gin.Default()

	server.GET("/hello", hello)
	// server.POST("/upload", uploadFile)
	server.POST("/updateUser", userUpdate)
	server.POST("/updateById", updateById)
	server.POST("/createEmty", createEmty)
	server.POST("/uploadSong", uploadSongToMongoDB)
	server.POST("/getSong", getSongFromMongoDB)
	server.POST("/login", login)
	server.POST("/downloadSong", downloadSongFromMongoDB)
	server.POST("/changeId", changeSongIDInMongoDB)
	server.POST("/updatePptx", updatePptx)
	server.GET("/getListOfSongs", getListOfSongs)
	server.POST("/addUser", addUser)
	server.POST("/addEvent", addEvent)
	server.POST("/find", findDB)           // collection, fild, whatToFind, если fild пустое вернет всё
	server.POST("/sendEmail", sendEmailto) // address, text
	server.GET("/getEvents", getEvents)
	server.GET("/getUsers", getUsers)
	server.GET("/confirmation/:id", confirmation)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := server.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
