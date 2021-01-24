package images

import (
	"cloud.google.com/go/storage"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"github.com/anthonyhawkins/savorbook/database"
	"github.com/anthonyhawkins/savorbook/middleware"
	"github.com/anthonyhawkins/savorbook/responses"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
	"io"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

var (
	storageClient *storage.Client
)

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func UploadImage(c *fiber.Ctx) error {

	response := new(responses.StandardResponse)
	response.Success = false
	db := database.GetDB()
	userId := middleware.AuthedUserId(c.Locals("user"))

	/**
	Read in the uploaded file and generate a unqiue name
	*/
	fileHeader, err := c.FormFile("image")
	if err != nil {
		response.Message = "Unable to Upload Image"
		response.Errors = append(response.Errors, response.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	//generate non-sequential-image-id
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	numStamp := strconv.Itoa(rangeIn(100000000000, 999999999999))
	h := sha1.New()
	h.Write([]byte(timeStamp))
	hexHash := hex.EncodeToString(h.Sum(nil))

	imageName := hexHash + "-" + numStamp + ".jpg"

	//TODO - validate image, resize etc.
	// validate image type via mimeType
	// jpg / png / gif
	// validate image size
	//

	/**
	Upload the image to Google Cloud Storage
	*/
	ctx := context.Background()
	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile(".googlekey.json"))
	if err != nil {
		response.Message = "Unable to Upload Image"
		//TODO respond with actual error for now
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	//create new storage-writer
	sw := storageClient.Bucket("savorbook-dev").Object(imageName).NewWriter(ctx)

	//copy file into storage-write aka upload the file
	file, err := fileHeader.Open()
	defer file.Close()
	if _, err := io.Copy(sw, file); err != nil {
		response.Message = "Unable to Upload Image"
		//TODO respond with actual error for now
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	if err := sw.Close(); err != nil {
		response.Message = "Unable to Upload Image"
		//TODO respond with actual error for now
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	resource, err := url.Parse("/" + "savorbook-dev" + "/" + sw.Attrs().Name)
	if err != nil {
		response.Message = "Unable to Upload Image"
		//TODO respond with actual error for now
		response.Errors = append(response.Errors, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	/**
	Create and Save the DB entry
	*/
	var image = new(Image)
	image.UserID = userId
	image.Name = imageName
	image.Path = resource.Path
	image.Used = true

	result := db.Create(&image)

	if result.RowsAffected == 0 {
		response.Message = "Recipe Creation Failed"
		//TODO - Should error be bubbled up DB error to client?
		response.Errors = append(response.Errors, result.Error.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Success = true
	response.Message = "Image has been uploaded"
	response.Data = image
	return c.Status(fiber.StatusCreated).JSON(response)
}
