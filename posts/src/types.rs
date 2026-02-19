use serde::{Deserialize, Serialize};
use serde_json::Value;

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct InterestOverTime {
    #[serde(rename = "type")]
   pub msg_type: String,
   pub  fetched_at: String,
   pub  timestamp: String,
   pub  data: std::collections::HashMap<String, f64>,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct RelatedQuery {
    pub query: String,
    pub value: Value,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct RelatedQueries {
    #[serde(rename = "type")]
    pub msg_type: String,
    pub fetched_at: String,
    pub keyword: String,
    pub data: Vec<RelatedQuery>,
}

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct CleanedInterest {
    #[serde(rename = "type")]
    pub msg_type: String,
    pub fetched_at: String,
    pub timestamp: String,
    pub data: std::collections::HashMap<String, f64>,
    pub is_normalized: bool,
}