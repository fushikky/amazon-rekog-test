package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			c.JSON(he.Code, he.Error())
		}
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/", func(c echo.Context) error {

		bytes, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			fmt.Printf("err: %+v\n", err)
			return err
		}
		// fmt.Printf("bytes: %+v\n", string(bytes))

		sess := session.Must(session.NewSession())
		svc := rekognition.New(sess, aws.NewConfig().WithRegion("us-east-1"))

		cparams := &rekognition.RecognizeCelebritiesInput{
			Image: &rekognition.Image{
				Bytes: bytes,
			},
		}

		cresp, err := svc.RecognizeCelebrities(cparams)
		if err != nil {
			fmt.Println(err.Error())
		}

		return c.JSON(http.StatusOK, cresp.CelebrityFaces)
	})

	e.Logger.Fatal(e.Start(":1323"))

}
