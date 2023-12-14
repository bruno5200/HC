package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {

	var checkInterval, requestTimeout int
	flag.IntVar(&checkInterval, "interval", 30, "Intervalo en segundos para realizar el escaneo de salud")
	flag.IntVar(&requestTimeout, "timeout", 30, "Timeout en segundos para la solicitud HTTP")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: %s [opciones] <puertos:rutas/jars:endpoints>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Opciones disponibles:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	// Verifica si se proporcionan los argumentos necesarios.
	if len(os.Args) < 1 {
		fmt.Printf("Uso: ./%s -interval 60 -timeout 10 puerto:jar:endpoint,puerto:jar:endpoint...", os.Args[0])
		return
	}

	// Parsea los argumentos en una lista de puertos, nombres de archivos JAR con rutas completas y endpoints.
	// arg := os.Args[1]
	// portJarEndpointList := strings.Split(arg, ",")

	// if len(portJarEndpointList) == 0 || !strings.Contains(arg, ":") {
	// 	fmt.Println("No se proporcionaron argumentos válidos.")
	// 	return
	// }

	portJarEndpointList := strings.Split(flag.Arg(0), ",")

	// Configura un temporizador para ejecutar el escaneo de salud de cada microservicio.
	cron := cron.New()

	for _, portJarEndpoint := range portJarEndpointList {
		puerto, rutaCompleta, endpoint, err := getPortRutaEndpoint(portJarEndpoint)

		if err != nil {
			return
		}

		_, err = cron.AddFunc(fmt.Sprintf("@every %ds", checkInterval), func() {
			checkHealth(puerto, requestTimeout, rutaCompleta, endpoint)
		})

		if err != nil {
			fmt.Printf("Error al configurar el temporizador para el puerto %d: %s\n", puerto, err)
			return
		}
	}

	cron.Start()

	// Inicializa el escaneo de salud inmediatamente.
	for _, portJarEndpoint := range portJarEndpointList {
		puerto, rutaCompleta, endpoint, err := getPortRutaEndpoint(portJarEndpoint)

		if err != nil {
			return
		}

		checkHealth(puerto, requestTimeout, rutaCompleta, endpoint)
	}

	// Mantén el programa en funcionamiento.
	select {}
}

func checkHealth(puerto, requestTimeout int, rutaCompleta, endpoint string) {
	url := fmt.Sprintf("http://localhost:%d/%s", puerto, endpoint)

	client := &http.Client{
		Timeout: time.Duration(requestTimeout) * time.Second,
	}

	// Realiza una solicitud HTTP al endpoint de verificación de salud de tu aplicación Java.
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("La aplicación Java en el puerto %d no responde: %s\n", puerto, err)

		// Reinicia la aplicación Java.
		restartJavaApplication(rutaCompleta)
	} else if resp.StatusCode != http.StatusOK {
		fmt.Printf("La aplicación Java en el puerto %d respondió con un código de estado no válido: %d\n", puerto, resp.StatusCode)

		// Reinicia la aplicación Java.
		restartJavaApplication(rutaCompleta)
	}

	// Cierra la respuesta.
	resp.Body.Close()
}

func restartJavaApplication(rutaCompleta string) {
	// Ruta completa al archivo JAR y argumentos para tu aplicación Java
	javaArgs := []string{"-Xmx512m", "-jar", rutaCompleta}

	// Ejecuta el comando Java
	cmd := exec.Command("java", javaArgs...)

	// Configura la salida estándar y la salida de errores para ver el resultado
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Ejecuta el comando
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error al reiniciar la aplicación Java: %v\n", err)
	}
}

func checkIfFileExists(rutacompleta string) bool {
	if _, err := os.Stat(rutacompleta); os.IsNotExist(err) {
		return false
	}
	return true
}

func getPortRutaEndpoint(s string) (puerto int, rutaCompleta, endpoint string, err error) {

	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		fmt.Printf("Argumento no válido: %s. Debe estar en el formato puerto:ruta-completa/nombre.jar:endpoint\n", s)
		return
	}

	if _, err = fmt.Sscan(parts[0], &puerto); err != nil {
		fmt.Printf("Error al analizar el puerto: %v\n", err)
		return
	}
	rutaCompleta = "/" + strings.Trim(parts[1], "/")

	if !checkIfFileExists(rutaCompleta) {
		fmt.Printf("El archivo %s no existe\n", rutaCompleta)
		err = fmt.Errorf("el archivo %s no existe", rutaCompleta)
		return
	}

	endpoint = parts[2]

	return
}
