# itemSave
_ItemSave_ Es una aplicación que consume un streaming en Redis con un mensaje del tipo {"id": "123", "site": "KKK"} para consultar datos de prueba de algunas apis de mercadolibre y complementa el objeto para almacenarlo en una base de datos MongoDB.


## Pre-requisitos 📋


```
1. Go instalado
2. IDE de tu preferencia para el lenguaje Go
3. Docker, Podman, u otro gestor de containers instalado, o un servidor de Redis y Base de datos MongoDB.
4. Si tomó la opción de Docker o Podman, puede ejecutar el siguiente comando para levantar un contenedor de Redis y MongoDB en Docker:

docker pull redis
docker run -p 6379:6379 redis --protected-mode no

docker pull mongodb/mongodb-community-server
docker run --name mongodb -d -p 27017:27017 mongodb/mongodb-community-server
```

## Variables de Entorno

```
MONGO_CONNECTION_STRING=mongodb://localhost:27017;MONGO_DATABASE_NAME=test;REDIS_HOST=localhost;REDIS_PORT=6378

  1. MONGO_CONNECTION_STRING=mongodb://localhost:27017  >> String de conexión a MongoDB, para efectos de la prueba no se tiene soporte de autenticación para la conexión.
  2. MONGO_DATABASE_NAME=test >> Nombre de la base de datos para almacenar los resultados
  3. PORT=8083 >> Puerto en el que corre la aplicación
  4. REDIS_HOST=localhost >> Server host de Redis
  5. REDIS_PORT=6379 >> Server port de Redis

NOTA: Esta versión solo se conecta a Redis sin autenticación.


  
```

## Arquitectura

![alt text](https://github.com/enavarrom/itemSave/blob/main/ItemSave_Loader.drawio.png?raw=tr)

Como se mencionó previamente esta aplicación consume un bus de eventos en Redis Stream, y luego procesa los mensajes para completar la información y guardar los datos en MongoDB.



## Ejecución de la aplicación

Se puede ejecutar la aplicación haciendo el build y luego run del archivo generado. O solo descargando el proyecto y correr el comando:

```
MONGO_CONNECTION_STRING=mongodb://localhost:27017 MONGO_DATABASE_NAME=test PORT=8083 REDIS_HOST=localhost REDIS_PORT=6379 go run main.go
```
Una vez arriba, la aplicación estará escuchando para consumir los eventos y procesar la información. Si se desea se podría levantar mas de una instancia para que los mensajes sean procesados mas rapida.
