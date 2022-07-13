package server

import (
	"app/app/model"
	"app/utils"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func redirect(ctx *fiber.Ctx) error {
	golyUrl := ctx.Params("id")

	goly, err := model.FindByGolyUrl((golyUrl))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not find goly in DB " + err.Error(),
		})
	}

	goly.Clicked += 1

	err = model.UpdateGol(goly)
	if err != nil {
		fmt.Printf("Error updating")
	}

	return ctx.Redirect(goly.Redirect, fiber.StatusTemporaryRedirect)
}

func getAllGolies(ctx *fiber.Ctx) error {

	golies, err := model.GetAllGolies()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting golly links " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(golies)

}

func getGoly(ctx *fiber.Ctx) error {

	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not parse id " + err.Error(),
		})
	}

	goly, err := model.GetGoly(id)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting golly " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(goly)
}

func createGoly(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")

	var goly model.Goly

	err := ctx.BodyParser(&goly)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error parsing JSON " + err.Error(),
		})
	}

	if goly.Random {
		goly.Goly = utils.RandomURL(8)
	}

	err = model.CreateGol(goly)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create Goly" + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(goly)

}

func UpdateGoly(ctx *fiber.Ctx) error {

	ctx.Accepts("application/json")

	var goly model.Goly

	err := ctx.BodyParser(&goly)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error parsing JSON " + err.Error(),
		})
	}

	err = model.UpdateGol(goly)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create Goly" + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(goly)

}

func deleteGoly(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not parse id " + err.Error(),
		})
	}

	err = model.DeleteGoly(id)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting golly " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Goly delete",
	})
}

func SetupAndListen() {
	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	router.Get("/r/:redirect", redirect)

	router.Get("/goly", getAllGolies)
	router.Get("/goly/:id", getGoly)
	router.Post("/goly", createGoly)
	router.Patch("/goly", UpdateGoly)

	err := router.Listen(":3000")

	if err != nil {
		println(err.Error())
	}

}
