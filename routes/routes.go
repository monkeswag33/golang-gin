package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gin-gonic/gin"
)

func genUpdateSQL(jsonMap map[string]interface{}, id string) (string, []interface{}) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}
	var sql string = "UPDATE mytable SET "
	var i int = 0
	var values []interface{} = make([]interface{}, 0, len(jsonMap))
	for key, value := range jsonMap {
		i++
		sql += key + "=$" + strconv.Itoa(i)
		if i != len(jsonMap) {
			sql += ", "
		}
		values = append(values, value)
	}
	sql += " WHERE id=$" + strconv.Itoa(i+1)
	values = append(values, idInt)
	sql += " RETURNING *"
	return sql, values
}

func GetHandler(ctx *gin.Context) {
	var users []*User = make([]*User, 0)
	pgxscan.Select(Context, DbPool, &users, "SELECT * FROM mytable")
	ctx.JSON(http.StatusOK, users)
}

func PostHandler(ctx *gin.Context) {
	var user User
	if err := ctx.BindJSON(&user); err != nil {
		fmt.Println(err)
		log.Fatal("Error converting JSON body to struct")
	}

	createdRow, _ := DbPool.Query(Context, "INSERT INTO mytable (firstname, lastname) VALUES ($1, $2) RETURNING *", user.Firstname, user.Lastname)
	var createdUser User
	if err := pgxscan.ScanOne(&createdUser, createdRow); err != nil {
		log.Fatal("Error processing rows")
	}
	ctx.JSON(http.StatusCreated, createdUser)
}

func UpdateHandler(ctx *gin.Context) {
	var mapbody map[string]interface{} = make(map[string]interface{})
	if err := ctx.BindJSON(&mapbody); err != nil {
		fmt.Println(err)
		log.Fatal("Error while converting JSON body to struct")
	}
	sql, values := genUpdateSQL(mapbody, ctx.Param("id"))
	updatedRow, err := DbPool.Query(Context, sql, values...)
	if err != nil {
		log.Fatal(err)
	}
	var updatedUser User
	if err := pgxscan.ScanOne(&updatedUser, updatedRow); err != nil {
		log.Fatal(err)
	}
	ctx.JSON(http.StatusOK, updatedUser)
}

func DeleteHandler(ctx *gin.Context) {
	deletedRow, _ := DbPool.Query(Context, "DELETE FROM mytable WHERE id=$1 RETURNING *", ctx.Param("id"))
	var deletedUser User
	pgxscan.ScanOne(&deletedUser, deletedRow)
	ctx.JSON(http.StatusOK, deletedUser)
}
