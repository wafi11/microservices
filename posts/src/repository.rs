use sqlx::{Error as SqlxError};
use crate::config::DbPool;
use crate::types::Post;

// Repository function - pure database logic
pub async fn get_all_posts(pool: &DbPool) -> Result<Vec<Post>, SqlxError> {
    let posts = sqlx::query_as::<_, Post>(
        "SELECT id, title, description FROM posts"
    )
        .fetch_all(pool)
        .await?;

    Ok(posts)
}

pub async fn get_post_by_id(pool: &DbPool, id: i32) -> Result<Post, SqlxError> {
    let post = sqlx::query_as::<_, Post>(
        "SELECT id, title, description FROM posts WHERE id = $1"
    )
        .bind(id)
        .fetch_one(pool)
        .await?;

    Ok(post)
}

pub async fn create_post(
    pool: &DbPool,
    title: String,
    description: String
) -> Result<Post, SqlxError> {
    let post = sqlx::query_as::<_, Post>(
        "INSERT INTO posts (title, description) VALUES ($1, $2) RETURNING id, title, description"
    )
        .bind(title)
        .bind(description)
        .fetch_one(pool)
        .await?;

    Ok(post)
}

pub async fn update_post(
    pool: &DbPool,
    id: i32,
    title: String,
    description: String,
) -> Result<Post, SqlxError> {
    let post = sqlx::query_as::<_, Post>(
        "UPDATE posts SET title = $1, description = $2 WHERE id = $3 RETURNING id, title, description"
    )
        .bind(title)
        .bind(description)
        .bind(id)
        .fetch_one(pool)
        .await?;

    Ok(post)
}

pub async fn delete_post(pool: &DbPool, id: i32) -> Result<u64, SqlxError> {
    let result = sqlx::query("DELETE FROM posts WHERE id = $1")
        .bind(id)
        .execute(pool)
        .await?;

    Ok(result.rows_affected())
}