# Health Check para los servicios en JAVA

## Descripción

Este proyecto contiene un servicio de salud para los servicios en JAVA, el cual se encarga de monitorizar el funcionamiento de uno o varios .JAR que se encuentren en ejecución.

## Uso

Para utilizar este servicio, se debe ejecutar el binario con el siguiente comando:

```bash
./healthcheck 8080:ruta/a/mi-aplicacion1.jar:/endpoint1,8081:ruta/a/mi-aplicacion2.jar:/endpoint2
```

Donde:

-   **8080**: Puerto en el que se ejecuta el .JAR de la aplicación que se desaea monitorizar.

-   **ruta/a/mi-aplicacion.jar**: es la rut absoluta más el nombre del .JAR de la aplicación que se desea monitorizar.

-   **/endpoint**: path de la aplicación el .JAR de la aplicación que se desea monitorizar (la misma debe devolver una status code 200 para ser válida).

Cuando se ejecuta el binario, se empieza a monitorizar y si la aplicación no responde, se levanta el servicio.

## Compilación

Para compilar el proyecto, se debe ejecutar el siguiente comando:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o healthcheck -v main.go
```
