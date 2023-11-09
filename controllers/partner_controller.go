package controllers

import (
	"context"

	"github.com/syobonaction/fur_lange/configs"
	"github.com/syobonaction/fur_lange/models"
	"github.com/syobonaction/fur_lange/responses"

	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var partnerCollection *mongo.Collection = configs.GetCollection(configs.DB, "partners")

func GetPartner(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	partnerId := c.Params("partnerId")
	var partner models.Partner
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(partnerId)

	err := partnerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&partner)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PartnerResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.PartnerResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": partner}})
}

func GetAllPartners(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var partners []models.Partner
	defer cancel()

	results, err := partnerCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PartnerResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singlePartner models.Partner
		if err = results.Decode(&singlePartner); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.PartnerResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		partners = append(partners, singlePartner)
	}

	return c.Status(http.StatusOK).JSON(
		responses.PartnerResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": partners}},
	)
}
