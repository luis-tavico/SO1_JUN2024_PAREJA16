package main

import (
	"Backend/Controller"
	"Backend/Database"
	"Backend/Model"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	if err := Database.Connect(); err != nil {
		log.Fatal("Error en", err)
	}

	getMem()

	if err := app.Listen(":8000"); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 500)
}

func getMem() {
	// Ruta del archivo que contiene la información de la RAM
	module_ram_path := "/proc/ram_so1_jun2024"
	// Ruta del archivo que contiene la información del CPU
	//module_cpu_path := "/proc/cpu_so1_1s2024"

	// RAM
	for range time.Tick(time.Second * 1) {
		// Leer todo el contenido del archivo
		content, err := ioutil.ReadFile(module_ram_path)
		if err != nil {
			log.Fatal(err)
		}
		// Convierte los valores a string y luego separa por comas
		values := strings.Split(string(content), ",")
		// 0: free_ram, 1: used_ram, 2: total_ram
		// Convertir los valores a float64
		free_ram, err := strconv.ParseFloat(strings.TrimSpace(values[0]), 64)
		if err != nil {
			log.Fatal("Error en", err)
		}
		used_ram, err := strconv.ParseFloat(strings.TrimSpace(values[1]), 64)
		if err != nil {
			log.Fatal("Error en", err)
		}
		total_ram, err := strconv.ParseFloat(strings.TrimSpace(values[2]), 64)
		if err != nil {
			log.Fatal("Error en", err)
		}
		// Calcular porcentajes de RAM libre y ocupada
		free_ram_per := (free_ram / total_ram) * 100
		used_ram_per := (used_ram / total_ram) * 100
		// Redondear a dos decimales
		free_ram_per = math.Round(free_ram_per*100) / 100
		used_ram_per = math.Round(used_ram_per*100) / 100
		// Imprimir los porcentajes
		fmt.Printf("Porcentaje de RAM libre: %.2f\n", free_ram_per)
		fmt.Printf("Porcentaje de RAM ocupada: %.2f\n", used_ram_per)
		// Convertir el porcentaje de RAM libre a string
		free_ram_per_str := strconv.FormatFloat(free_ram_per, 'f', 2, 64)
		// Convertir el porcentaje de RAM libre a string
		used_ram_per_str := strconv.FormatFloat(used_ram_per, 'f', 2, 64)
		// Crear una instancia de Model.Data con los porcentajes
		data := Model.Data{
			Used_percentage: used_ram_per_str,
			Free_percentage: free_ram_per_str,
		}
		Controller.InsertData("ram", data)
	}

	// CPU
	// Leer todo el contenido del archivo
	/*
		content, err := ioutil.ReadFile(module_cpu_path)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(string(content))
	*/

}
