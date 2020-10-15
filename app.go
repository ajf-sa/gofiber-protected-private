package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", Protected("/areyouzombie"), func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "text/html")
		return ctx.SendString(`<h1>Welcome bro!</h1><s><a href="/logout">logout</a>`)
	})

	app.Get("/logout", func(ctx *fiber.Ctx) error {
		ctx.ClearCookie("agree")
		return ctx.Redirect("/")
	})

	app.Get("/zombienotallow", func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "text/html")
		return ctx.SendString(`<h2>Sorry mr. Zombie you'r not allow! <a href="/">I'm not ZOMBIE</a> </h2>`)
	})

	app.Post("/check", func(ctx *fiber.Ctx) error {
		type request struct {
			Agree string `json:"agree"`
		}
		var body request
		ctx.BodyParser(&body)
		if body.Agree == "Fiber is awesome" {
			ctx.Cookie(&fiber.Cookie{
				Name:     "agree",
				Value:    body.Agree,
				Expires:  time.Now().Add(24 * time.Hour),
				HTTPOnly: true,
				SameSite: "lax",
			})

			return ctx.Redirect("/")

		}
		return ctx.Redirect("/zombienotallow")

	})

	app.Get("/areyouzombie", func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "text/html")
		return ctx.SendString(`
		<h1>Are you zombie ?</h1> 
		<p style="font-size:2rem;">type '<b style="color:blue;">Fiber is awesome</b>'; <b style="color:red;"> Zombie </b> can't type correct!</p>
		<form action="/check" method="post">
		<input type="text" name="agree" />
		<input type="submit" />
		</form>
		`)
	})

	app.Listen(":3000")
}

func Protected(next ...string) fiber.Handler {
	nxt := ""
	if len(next) > 0 {
		nxt = next[0]
	}
	return func(ctx *fiber.Ctx) error {
		auth := ctx.Cookies("agree")
		if auth == "Fiber is awesome" {
			return ctx.Next()
		}
		return ctx.Redirect(nxt)
	}
}
