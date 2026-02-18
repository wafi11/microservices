use axum::extract::{Request, State};
use axum::http::StatusCode;
use axum::middleware::Next;
use axum::response::{IntoResponse, Response};
use std::sync::Arc;
use tokio::sync::Mutex;
use std::collections::HashMap;
use std::time::{Duration, Instant};

#[derive(Clone)]
pub struct RateLimiter {
    pub requests: Arc<Mutex<HashMap<String, (u32, Instant)>>>,
    pub max_requests: u32,
    pub window: Duration,
}

pub async fn route_limiter(
    State(limiter): State<RateLimiter>,
    req: Request,
    next: Next,
) -> Response {

    // get identifier (IP, API key, dll)
    let key = req
        .headers()
        .get("x-api-key")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("anonymous")
        .to_string();

    let mut map = limiter.requests.lock().await;
    let now = Instant::now();

    let entry = map.entry(key).or_insert((0, now));

    // Reset window if through time
    if now.duration_since(entry.1) > limiter.window {
        *entry = (0, now);
    }

    if entry.0 >= limiter.max_requests {
        return StatusCode::TOO_MANY_REQUESTS.into_response();
    }

    // if not through → increment + next
    entry.0 += 1;
    drop(map); // off lock before await

    next.run(req).await
}