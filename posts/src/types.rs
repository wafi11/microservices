use serde::{Deserialize, Serialize};
use sqlx::FromRow;

// Serialize:   Struct Rust  →  JSON  (send ke client)
// Deserialize: JSON         →  Struct Rust  (get dari client)
#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Post {
    pub id: i32,
    pub title: Option<String>,
    pub description: Option<String>,
}

#[derive(Debug, Deserialize)]
pub struct CreatePost {
    pub title: String,
    pub description: String,
}

#[derive(Debug, Serialize)]
pub struct ApiResponse<T: Serialize> {
    pub status_code: u16,
    pub success: bool,
    pub message: String,
    pub data: Option<T>,
}

#[derive(Debug, Serialize)]
pub struct RootResponse {
    pub timestamp: String,
    pub message: String,
    pub status : i16
}

impl<T: Serialize> ApiResponse<T> {
    pub fn success(status_code: u16, data: Option<T>) -> Self {
        ApiResponse {
            status_code,
            success: true,
            message: "Success".to_string(),
            data,
        }
    }

    // pub fn error(status_code: u16, message: &str) -> Self {
    //     ApiResponse {
    //         status_code,
    //         success: false,
    //         message: message.to_string(),
    //         data: None,
    //     }
    // }
}