package controller

import (
	"fmt"
	"log"
	"strconv"

	model "github.com/AutoResto/module/material/entity"
	"github.com/gin-gonic/gin"
)

func GetAllMaterial(c *gin.Context) {
	db := Connect()
	defer db.Close()

	query := "SELECT * FROM material"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	var material model.Material
	var materials []model.Material

	for rows.Next() {
		if err := rows.Scan(&material.Id, &material.Name, &material.Quantity, &material.Unit); err != nil {
			log.Fatal(err.Error())
		} else {
			materials = append(materials, material)
		}
	}

	var response model.MaterialResponse
	if err == nil {
		response.Message = "Get Material Success"
		response.Data = materials
		sendMaterialSuccessResponse(c, response)
	} else {
		response.Message = "Get Material Failed"
		fmt.Println(err)
		sendMaterialErrorResponse(c, response)
	}
}

func InsertMaterial(c *gin.Context) {
	db := Connect()
	defer db.Close()

	id, _ := strconv.Atoi(c.PostForm("id"))
	name := c.PostForm("name")
	quantity, _ := strconv.Atoi(c.PostForm("quantity"))
	unit := c.PostForm("unit")

	_, errQuery := db.Exec("INSERT INTO material(id,name,quantity,unit) VALUES(?,?,?,?)",
		id,
		name,
		quantity,
		unit)

	var response model.MaterialResponse

	if errQuery == nil {
		response.Message = "Insert Material Success"
		sendMaterialSuccessResponse(c, response)
	} else {
		response.Message = "Insert Material Failed"
		sendMaterialErrorResponse(c, response)
	}
}

func UpdateMaterial(c *gin.Context) {
	db := Connect()
	defer db.Close()

	idMaterial := c.Param("material_id")
	materialName := c.PostForm("name")
	materialQuantity, _ := strconv.Atoi(c.PostForm("quantity"))
	materialUnit := c.PostForm("unit")

	rows, _ := db.Query("SELECT * FROM material WHERE id = '" + idMaterial + "'")

	var material model.Material

	for rows.Next() {
		if err := rows.Scan(&material.Id, &material.Name, &material.Quantity, &material.Unit); err != nil {
			log.Fatal(err.Error())
		}
	}

	if materialName == "" {
		materialName = material.Name
	}

	if materialQuantity == materialQuantity {
		materialQuantity = material.Quantity
	}

	if materialUnit == "" {
		materialUnit = material.Unit
	}

	_, errQuery := db.Exec("UPDATE material SET name = ?,quantity = ?,unit = ? WHERE id = ?",
		materialName,
		materialQuantity,
		materialUnit,
		idMaterial,
	)

	var response model.MaterialResponse
	if errQuery == nil {
		response.Message = "Update Material Success"
		sendMaterialSuccessResponse(c, response)
	} else {
		response.Message = "Update Material Failed"
		sendMaterialErrorResponse(c, response)
	}
}

func DeleteMaterial(c *gin.Context) {
	db := Connect()
	defer db.Close()

	idmaterial := c.Param("material_id")

	_, query := db.Exec("DELETE FROM material WHERE id = ?", idmaterial)

	var response model.MaterialResponse
	if query == nil {
		response.Message = "Delete Success"
		sendMaterialSuccessResponse(c, response)

	} else {
		response.Message = "Delete Failed"
		sendMaterialErrorResponse(c, response)
	}
}
