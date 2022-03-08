package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	modelMaterial "github.com/AutoResto/module/material/entity"
	modelMenu "github.com/AutoResto/module/menu/entity"
	modelRecipe "github.com/AutoResto/module/recipe/entity"

	"github.com/gin-gonic/gin"
)

// Search Menu
func SearchMenu(c *gin.Context) {
	db := Connect()
	defer db.Close()

	menuName := c.PostForm("menu_name")

	query := "SELECT m.id, m.name, m.price, m.idRecipeFK FROM `menu` m WHERE m.name LIKE '%" + menuName + "%'"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	var menu modelMenu.Menu
	var menus []modelMenu.Menu

	for rows.Next() {
		if err := rows.Scan(&menu.Id, &menu.Name, &menu.Price, &menu.Recipe.Id); err != nil {
			log.Fatal(err.Error())
		} else {
			menus = append(menus, menu)
		}
	}

	var response modelMenu.MenuResponse
	if err == nil {
		response.Message = "Get Menu Success"
		response.Data = menus
		sendMenuSuccessResponse(c, response)
	} else {
		response.Message = "Get Menu Query Error"
		sendMenuErrorResponse(c, response)
	}
}

// Search Material
func SearchMaterial(c *gin.Context) {
	db := Connect()
	defer db.Close()

	materialName := c.PostForm("material_name")

	query := "SELECT m.id, m.name, m.quantity, m.unit FROM `material` m WHERE m.name LIKE '" + materialName + "%'"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	var material modelMaterial.Material
	var materials []modelMaterial.Material

	for rows.Next() {
		if err := rows.Scan(&material.Id, &material.Name, &material.Quantity, &material.Unit); err != nil {
			log.Fatal(err.Error())
		} else {
			materials = append(materials, material)
		}
	}

	var response modelMaterial.MaterialResponse
	if err == nil {
		response.Message = "Get Material Success"
		response.Data = materials
		sendMaterialSuccessResponse(c, response)
	} else {
		response.Message = "Get Material Query Error"
		sendMaterialErrorResponse(c, response)
	}
}

// Insert Menu
func InsertMenu(c *gin.Context) {
	db := Connect()
	defer db.Close()

	name := c.PostForm("name")
	price, _ := strconv.Atoi(c.PostForm("price"))
	description := c.PostForm("description")
	materialName := c.PostForm("material")
	quantity := c.PostForm("quantity")
	unit := c.PostForm("unit")
	materialArr := strings.Split(materialName, ",")
	quantityArr := strings.Split(quantity, ",")
	unitArr := strings.Split(unit, ",")

	//insert recipe
	_, errQuery := db.Exec("INSERT INTO recipe(description) VALUES (?)",
		description,
	)

	//insert menu dengan select id dari description
	_, errQuery = db.Exec("INSERT INTO menu(name,price,idRecipeFK) VALUES (?,?,(SELECT id from recipe where recipe.description='"+description+"'));",
		name,
		price,
	)

	//mencari material dengan nama yang sama
	query := "SELECT m.id, m.name FROM `material` m WHERE m.name LIKE '" + materialName + "%'"
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	var material modelMaterial.Material
	var materials []modelMaterial.Material

	for rows.Next() {
		if err := rows.Scan(&material.Id, &material.Name); err != nil {
			log.Fatal(err.Error())
		} else {
			materials = append(materials, material)
		}
	}

	for i := 0; i < len(materialArr); i++ {

		for j := 0; j < len(materials); j++ {
			//mengecek apakah material sudah ada
			if materialArr[i] != materials[j].Name {
				// menambah material jika tidak ada
				_, errQuery = db.Exec("INSERT INTO material(name) VALUES(?)",
					materialArr[i],
				)
			}
		}

		_, errQuery = db.Exec("INSERT INTO recipedetail(quantity,idMaterialFK,idRecipeFK,unit) VALUES(?,(SELECT id FROM material WHERE material.name=?),(SELECT id FROM recipe WHERE recipe.description=?),?)",
			quantityArr[i],
			materialArr[i],
			description,
			unitArr[i],
		)
	}

	var response modelMenu.MenuResponse
	if errQuery == nil {
		response.Message = "Insert Menu Success"
		sendMenuSuccessResponse(c, response)
	} else {
		response.Message = "Insert Menu Failed Error"
		sendMenuErrorResponse(c, response)
	}
}

//Update Menu
func UpdateMenu(c *gin.Context) {
	db := Connect()
	defer db.Close()

	menuId := c.Param("menu_id")
	name := c.PostForm("name")
	price, _ := strconv.ParseFloat(c.PostForm("price"), 32)

	rows, _ := db.Query("SELECT menu.id, menu.name, menu.price FROM menu WHERE id='" + menuId + "'")
	var menu modelMenu.Menu
	for rows.Next() {
		if err := rows.Scan(&menu.Id, &menu.Name, &menu.Price); err != nil {
			log.Fatal(err.Error())
		}
	}

	// Jika kosong dimasukkan nilai lama
	if name == "" {
		name = menu.Name
	}

	if price == 0 {
		price = menu.Price
	}

	_, errQuery := db.Exec("UPDATE menu SET name = ?, price = ? WHERE id=?",
		name,
		price,
		menuId,
	)

	var response modelMenu.MenuResponse
	if errQuery == nil {
		response.Message = "Update Menu Success"
		sendMenuSuccessResponse(c, response)
	} else {
		response.Message = "Update Menu Failed Error"
		sendMenuErrorResponse(c, response)
	}
}

// Delete Menu
func DeleteMenu(c *gin.Context) {
	db := Connect()
	defer db.Close()

	menuId := c.Param("menu_id")

	_, errQuery := db.Exec("DELETE FROM menu WHERE id=?",
		menuId,
	)

	var response modelMenu.MenuResponse
	if errQuery == nil {
		response.Message = "Delete Menu Success"
		sendMenuSuccessResponse(c, response)
	} else {
		response.Message = "Delete Menu Failed Error"
		sendMenuErrorResponse(c, response)
	}
}

func sendMenuSuccessResponse(c *gin.Context, response modelMenu.MenuResponse) {
	c.JSON(http.StatusOK, response)
}

func sendMenuErrorResponse(c *gin.Context, response modelMenu.MenuResponse) {
	c.JSON(http.StatusBadRequest, response)
}

// Get Recipe
func GetAllRecipe(c *gin.Context) {
	db := Connect()
	defer db.Close()

	query := "SELECT id,description FROM recipe"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}

	var recipe modelRecipe.RecipeDetail
	var recipes []modelRecipe.RecipeDetail

	for rows.Next() {
		if err := rows.Scan(&recipe.Id, &recipe.Quantity, &recipe.Unit); err != nil {
			log.Fatal(err.Error())
		} else {
			recipes = append(recipes, recipe)
		}
	}

	var response modelRecipe.RecipeDetailResponse
	if err == nil {
		response.Message = "Delete Menu Success"
		sendRecipeDetailSuccessResponse(c, response)
	} else {
		response.Message = "Delete Menu Failed Error"
		sendRecipeDetailErrorResponse(c, response)
	}
}
