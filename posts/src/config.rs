use std::env;
use std::time::Duration;
use sqlx::{Pool, Postgres, postgres::PgPoolOptions};

pub type DbPool = Pool<Postgres>;

pub async fn create_pool() -> Result<DbPool, sqlx::Error> {
    dotenv::dotenv().ok();

    let database_url = env::var("DATABASE_URL")
        .expect("DATABASE_URL must be set");

    let pool = PgPoolOptions::new()
        .max_connections(50)
        .min_connections(5)
        .acquire_timeout(Duration::from_secs(3))
        .idle_timeout(Duration::from_secs(600))
        .max_lifetime(Duration::from_secs(1800))
        .connect(&database_url)
        .await?;

    Ok(pool)
}