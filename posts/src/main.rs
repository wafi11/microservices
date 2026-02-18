use std::collections::HashMap;
use std::sync::Arc;
use std::time::Duration;
use tokio::sync::Mutex;
use axum::{Extension, Router};
use axum::routing::{get, patch};
use crate::middlewares::{route_limiter, RateLimiter};
mod types;
mod config;
mod repository;
mod handler;
mod middlewares;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let pool = config::create_pool().await?;
    println!("DB connection established");

    let limiter = RateLimiter {
        requests: Arc::new(Mutex::new(HashMap::new())),
        max_requests: 1,
        window: Duration::from_secs(60),
    };

    let app = Router::new()
        .route("/api", get(handler::get_root))
        .route("/api/posts", get(handler::get_posts).post(handler::create_post))
        .route("/api/posts/{id}", patch(handler::update_post).delete(handler::delete_post).get(handler::get_post))
        .layer(Extension(pool))
        .layer(axum::middleware::from_fn_with_state(limiter, route_limiter));

    let listener = tokio::net::TcpListener::bind("0.0.0.0:5000").await?;
    println!("🚀 Server running on http://0.0.0.0:5000");

    axum::serve(listener, app).await?;

    Ok(())
}