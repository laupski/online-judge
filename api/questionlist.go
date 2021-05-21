package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getQuestionListResponse struct {
	Data []getQuestionListRow `json:"data"`
}

type getQuestionListRow struct {
	Key    string `json:"key"`
	Number int    `json:"number"`
	Title  string `json:"title"`
}

func getQuestionList(c *gin.Context) {
	fmt.Println("Getting question list...")
	rows, err := PostgresConnection.Query(context.Background(), "SELECT key, number, title FROM public.questions ORDER BY number")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error fetching list",
		})
		return
	}
	defer rows.Close()

	var response getQuestionListResponse

	for rows.Next() {
		var key string
		var number int
		var title string
		err = rows.Scan(&key, &number, &title)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error parsing rows",
			})
			return
		}
		fmt.Printf("key: %v, number: %v, title: %v\n", key, number, title)
		datarow := &getQuestionListRow{
			Key:    key,
			Number: number,
			Title:  title,
		}
		response.Data = append(response.Data, *datarow)
	}

	c.JSON(http.StatusOK, response)
}
