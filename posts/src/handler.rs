    use axum::{
        Extension,
        Json,
        http::StatusCode,
        extract::Path,
    };
    use chrono::Utc;
    use chrono_tz::Asia::Jakarta;
    use crate::config::DbPool;
    use crate::types::{ApiResponse, CreatePost, Post,RootResponse};
    use crate::repository;

    pub async fn get_root() -> Result<Json<RootResponse>, StatusCode> {
        let jakarta_time = Utc::now().with_timezone(&Jakarta);

        let response = RootResponse {
            timestamp: jakarta_time.format("%Y-%m-%d %H:%M:%S").to_string(),
            message: "Web Api Version 0.1".to_string(),
            status: 200,
        };

        Ok(Json(response))
    }

    // GET /posts
    pub async fn get_posts(
        Extension(pool): Extension<DbPool>,
    ) -> Result<Json<ApiResponse<Vec<Post>>>, StatusCode> {

        let posts = repository::get_all_posts(&pool)
            .await
            .map_err(|e| {
                eprintln!("❌ Database error: {:?}", e);
                StatusCode::INTERNAL_SERVER_ERROR
            })?;

        Ok(Json(ApiResponse::success(200, Some(posts))))
    }

    // GET /posts/:id
    pub async fn get_post(
        Extension(pool): Extension<DbPool>,
        Path(id): Path<i32>,
    ) -> Result<Json<ApiResponse<Post>>, StatusCode> {

        let post = repository::get_post_by_id(&pool, id)
            .await
            .map_err(|e| {
                eprintln!("❌ Database error: {:?}", e);
                StatusCode::NOT_FOUND
            })?;

        Ok(Json(ApiResponse::success(200, Some(post))))
    }

    // POST /posts
    pub async fn create_post(
        Extension(pool): Extension<DbPool>,
        Json(payload): Json<CreatePost>,
    ) -> Result<Json<ApiResponse<Post>>, StatusCode> {

        let post = repository::create_post(&pool, payload.title, payload.description)
            .await
            .map_err(|e| {
                eprintln!("❌ Database error: {:?}", e);
                StatusCode::INTERNAL_SERVER_ERROR
            })?;

        Ok(Json(ApiResponse::success(201, Some(post))))
    }

    // PUT /posts/:id
    pub async fn update_post(
        Extension(pool): Extension<DbPool>,
        Path(id): Path<i32>,
        Json(payload): Json<CreatePost>,
    ) -> Result<Json<ApiResponse<Post>>, StatusCode> {

        let post = repository::update_post(&pool, id, payload.title, payload.description)
            .await
            .map_err(|e| {
                eprintln!("❌ Database error: {:?}", e);
                StatusCode::NOT_FOUND
            })?;

        Ok(Json(ApiResponse::success(200, Some(post))))
    }

    // DELETE /posts/:id
    pub async fn delete_post(
        Extension(pool): Extension<DbPool>,
        Path(id): Path<i32>,
    ) -> Result<Json<ApiResponse<()>>, StatusCode> {

        let rows_affected = repository::delete_post(&pool, id)
            .await
            .map_err(|e| {
                eprintln!("❌ Database error: {:?}", e);
                StatusCode::INTERNAL_SERVER_ERROR
            })?;

        if rows_affected == 0 {
            return Err(StatusCode::NOT_FOUND);
        }

        Ok(Json(ApiResponse::success(200, None)))
    }