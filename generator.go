package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"syreclabs.com/go/faker"
)

type Model struct {
	Name   string   `json:"name"`
	Schema []Schema `json:"schema"`
	Seeder Seeder   `json:"seeder"`
}

type Schema struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	NotNull    bool   `json:"not_null"`
	PrimaryKey bool   `json:"primary_key"`
}

type Seeder struct {
	GeneratorSeeder GeneratorSeeder `json:"generator_seeder"`
	ManualSeeder    ManualSeeder    `json:"manual_seeder"`
}

type ManualSeeder struct {
	TargetSchema string   `json:"target_schema"`
	Values       []string `json:"values"`
}

type GeneratorSeeder struct {
	TargetSchema []TargetSchema `json:"target_schema"`
	TotalSeed    int            `json:"total_seed"`
}

type TargetSchema struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Constant string `json:"constant"`
}

func main() {
	var model Model
	files, err := ioutil.ReadDir("./json")
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		model = readFileJson(file.Name())
		dropIfExist := dropTableIfExist(model)
		createTable := createTable(model)
		seeder := generateSeeder(model)

		sqlContent := dropIfExist + createTable + seeder

		generateSQLFile(model.Name, sqlContent)
	}

}

func generateSQLFile(sqlName string, sqlContent string) {
	sqlFile := []byte(sqlContent)
	sqlPath := fmt.Sprintf("./migration/%s.sql", sqlName)
	err := ioutil.WriteFile(sqlPath, sqlFile, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func generateSeeder(model Model) string {
	var sqlSeeds string
	for i := 0; i < model.Seeder.GeneratorSeeder.TotalSeed; i++ {
		var Fakers = map[string]string{
			"name":     strings.Replace(faker.Name().Name(), "'", "", -1),
			"email":    faker.Internet().Email(),
			"constant": "bismillah",
		}

		var columns []string
		var values []string
		for _, seed := range model.Seeder.GeneratorSeeder.TargetSchema {

			value := Fakers[seed.Type]

			values = append(values, value)
			columns = append(columns, seed.Name)
		}
		sqlSeed := fmt.Sprintf("INSERT INTO %s(%s) VALUES('%s');\n", model.Name, strings.Join(columns, ","), strings.Join(values, "','"))
		sqlSeeds = sqlSeeds + sqlSeed
	}
	return sqlSeeds
}

func dropTableIfExist(model Model) string {
	dropIfExist := fmt.Sprintf("DROP TABLE IF EXISTS %s;\n", model.Name)
	return dropIfExist
}

func createTable(model Model) string {
	var columns []string
	for _, column := range model.Schema {
		columnWithOption := column.Name + " " + column.Type + " " + sqlPrimaryKey(column.PrimaryKey) + sqlNotNull(column.NotNull)
		columns = append(columns, columnWithOption)
	}
	createTable := fmt.Sprintf("CREATE TABLE %s(%s);\n", model.Name, strings.Join(columns, ",\n"))
	return createTable
}

func sqlPrimaryKey(status bool) string {
	if status {
		return " PRIMARY KEY"
	}
	return ""
}

func sqlNotNull(status bool) string {
	if status {
		return " NOT NULL"
	}
	return ""
}

func readFileJson(fileName string) Model {
	filePath := fmt.Sprintf("./json/%s", fileName)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}
	var JsonData Model
	err = json.Unmarshal(data, &JsonData)
	if err != nil {
		fmt.Println(err)
	}

	return JsonData
}
