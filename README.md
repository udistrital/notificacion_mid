# notificacion_mid
Servicio para envío de notificaciónes por difusión en AWS SNS.

Para instrucciones de implementación del sistema de notificaciones consulta la [documentacion](https://drive.google.com/file/d/1ZlXdKAyooyftamyCZX8Btmd5nhDR5YEX/view?usp=drivesdk)

## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones
* [Golang](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/golang.md)
* [BeeGo](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/beego.md)
* [Docker](https://docs.docker.com/engine/install/ubuntu/)
* [Docker Compose](https://docs.docker.com/compose/)
* [AWS SNS](https://aws.amazon.com/es/sns/)

### Variables de Entorno
```shell
NOTIFICACION_MID_HTTP_PORT = [Puerto de ejecucion]
```
**NOTA:** Las variables se pueden ver en el fichero conf/app.conf y están identificadas.


### Ejecución del Proyecto
```shell
#1. Obtener el repositorio con Go
go get github.com/udistrital/notificacion_mid

#2. Moverse a la carpeta del repositorio
cd $GOPATH/src/github.com/udistrital/notificacion_mid

# 3. Moverse a la rama **develop**
git pull origin develop && git checkout develop

# 4. alimentar todas las variables de entorno que utiliza el proyecto.
notificacion_mid_HTTP_PORT=8080 CONFIGURACION_SERVICE=127.0.0.1:27017 notificacion_mid_SOME_VARIABLE=some_value bee run
```

### Ejecución Dockerfile
```shell
# docker build --tag=notificacion_mid . --no-cache
# docker run -p 80:80 notificacion_mid
```

### Ejecución docker-compose
```shell
#1. Clonar el repositorio
git clone -b develop https://github.com/udistrital/notificacion_mid

#2. Moverse a la carpeta del repositorio
cd notificacion_mid

#3. Crear un fichero con el nombre **custom.env**
touch custom.env

#4. Crear la network **back_end** para los contenedores
docker network create back_end

#5. Ejecutar el compose del contenedor
docker-compose up --build

#6. Comprobar que los contenedores estén en ejecución
docker ps
```

Pruebas unitarias
```shell
# Not Data
```
## Estado CI

| Develop | Relese 0.0.1 | Master |
| -- | -- | -- |
| [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/notificacion_mid/status.svg?ref=refs/heads/develop)](https://hubci.portaloas.udistrital.edu.co/udistrital/notificacion_mid) |  [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/notificacion_mid/status.svg?ref=refs/heads/release/0.0.1)](https://hubci.portaloas.udistrital.edu.co/udistrital/notificacion_mid) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/notificacion_mid/status.svg)](https://hubci.portaloas.udistrital.edu.co/udistrital/notificacion_mid) |


## Licencia

This file is part of notificacion_mid.

notificacion_mid is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

notificacion_mid is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with notificacion_mid. If not, see https://www.gnu.org/licenses/.
