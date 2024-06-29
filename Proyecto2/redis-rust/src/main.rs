#[macro_use] extern crate rocket;

use rocket::serde::json::Json;
use rocket::serde::{Deserialize, Serialize};
use redis::Commands;

#[derive(Deserialize, Serialize)]
struct Data {
    Pais: String,
    Texto: String
}

#[post("/set", format = "json", data = "<data>")]
async fn set_data(data: Json<Data>) -> Result<&'static str, &'static str> {
    // Crear cliente de redis
    let client = redis::Client::open("redis://redis:6379/")
        .map_err(|_| "Failed to create Redis client")?;

    // Conexion a redis
    let mut con = client.get_connection()
        .map_err(|_| "Failed to connect to Redis")?;
    
    // Insertar hash en redis
    let _: () = con.hincr(&data.Pais, &data.Texto, 1)
    .map_err(|_| "Failed to set data in Redis")?;

    // Insertar hash en redis
    let _: () = con.hincr("countries", &data.Pais, 1)
        .map_err(|_| "Failed to set data in Redis")?;

    // Insertar hash en redis
    let _: () = con.hincr("messages", &data.Texto, 1)
        .map_err(|_| "Failed to set data in Redis")?;

    // Incrementar contador total de mensajes
    let _: () = con.incr("total-messages", 1)
        .map_err(|_| "Failed to increment total messages in Redis")?;

    Ok("Data set")
}

#[launch]
fn rocket() -> _ {
    rocket::build().mount("/", routes![set_data])
}