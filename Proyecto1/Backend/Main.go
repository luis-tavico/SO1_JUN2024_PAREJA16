package main

import (
	"Backend/Database"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	if err := Database.Connect(); err != nil {
		log.Fatal("Error en", err)
	}

	getData()

	if err := app.Listen(":8000"); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 500)
}

func getData() {
	for range time.Tick(time.Second * 1) {

		// Obtener datos de la RAM
		free_ram_percentage, err := getRAMdata()
		if err != nil {
			fmt.Println("Error al obtener datos de la RAM: ", err)
			//http.Error(w, "Error al obtener datos de la RAM", http.StatusInternalServerError)
			return
		}

		// Obtener datos de la CPU
		free_cpu_percentage, err := getCPUdata()
		if err != nil {
			fmt.Println("Error al obtener datos de la CPU: ", err)
			//http.Error(w, "Error al obtener datos de la CPU", http.StatusInternalServerError)
			return
		}

		// Obtener los procesos del sistema
		processes, err := getProcesses()
		if err != nil {
			fmt.Println("Error al obtener los procesos del sistema: ", err)
			//http.Error(w, "Error al obtener los procesos del sistema", http.StatusInternalServerError)
			return
		}

		used_ram_percentage := 100 - free_ram_percentage
		used_cpu_percentage := 100 - free_cpu_percentage

		fmt.Println("RAM")
		fmt.Println("Porcentaje de RAM libre:", free_ram_percentage)
		fmt.Println("Porcentaje de RAM usada:", used_ram_percentage)
		fmt.Println("CPU")
		fmt.Println("Porcentaje de CPU libre:", free_cpu_percentage)
		fmt.Println("Porcentaje de CPU usada:", used_cpu_percentage)
		fmt.Println("PID de los procesos padres:")

		// obtener todos los pid de los procesos padres, no los hijos
		var pids []int
		for _, process := range processes.Processes {
			pids = append(pids, process.Pid)
		}

		// Convertir la estructura a formato JSON
		jsonData, err := json.Marshal(pids)
		if err != nil {
			fmt.Println("Error al convertir los datos a JSON: ", err)
			//http.Error(w, "Error al convertir los datos a JSON", http.StatusInternalServerError)
			return
		}

		print(string(jsonData))

		// Crear una instancia de Model.Data con los porcentajes
		/*
			data := Model.Data{
				Used_percentage: used_ram_per_str,
				Free_percentage: free_ram_per_str,
			}
			// Insertar los datos en la base de datos
			Controller.InsertData("ram", data)
		*/
	}
}

// Estructura para los datos del sistema
type SystemData struct {
	RAM_percentage int `json:"ram_percentage"`
	CPU_percentage int `json:"cpu_percentage"`
}

// Estructura para los procesos hijo
type ProcesoHijo struct {
	Pid      int    `json:"pid"`
	Name     string `json:"name"`
	State    int    `json:"state"`
	PidPadre int    `json:"pidPadre"`
}

// Estructura para los procesos padre
type ProcesoPadre struct {
	Pid   int           `json:"pid"`
	Name  string        `json:"name"`
	User  int           `json:"user"`
	State int           `json:"state"`
	Ram   int           `json:"ram"`
	Child []ProcesoHijo `json:"child"`
}

// Definicion de estructura para el objeto JSON
type ProcessData struct {
	Processes []ProcesoPadre `json:"processes"`
}

// Funcion para obtener datos de la RAM
func getRAMdata() (int, error) {
	cmd := exec.Command("sh", "-c", "cat /proc/ram_so1_jun2024")
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		return 0, err
	}

	// Convertir la salida a formato JSON
	var data SystemData
	err = json.Unmarshal(stdout, &data)
	if err != nil {
		return 0, err
	}

	return data.RAM_percentage, nil
}

// Funcion para obtener datos de la CPU
func getCPUdata() (int, error) {
	cmd := exec.Command("sh", "-c", "cat /proc/cpu_so1_1s2024")
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		return 0, err
	}

	// Convertir la salida a formato JSON
	var data SystemData
	err = json.Unmarshal(stdout, &data)
	if err != nil {
		return 0, err
	}

	return data.CPU_percentage, nil
}

// Funcion para obtener los procesos del sistema
func getProcesses() (ProcessData, error) {
	cmd := exec.Command("sh", "-c", "cat /proc/cpu_so1_1s2024")
	stdout, err := cmd.CombinedOutput()

	if err != nil {
		return ProcessData{}, err
	}
	// Convertir la salida a formato JSON
	var data ProcessData
	err = json.Unmarshal(stdout, &data)
	if err != nil {
		return ProcessData{}, err
	}

	return data, nil
}
