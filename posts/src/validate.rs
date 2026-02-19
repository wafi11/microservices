use std::f32::NEG_INFINITY;

use crate::types::{InterestOverTime,RelatedQueries,RelatedQuery,CleanedInterest};



pub fn validate_json(payload: &str) -> Option<InterestOverTime> {
    match serde_json::from_str::<InterestOverTime>(payload) {
        Ok(r) => {
            if r.timestamp.is_empty() || r.fetched_at.is_empty() || r.data.is_empty() {
                eprintln!("❌ Field kosong");
                None
            } else {
                Some(r)
            }
        }
        Err(e) => {
            eprintln!("❌ JSON invalid: {}", e);
            None
        }
    }
}

pub fn filter_nulls(mut record: InterestOverTime) -> Option<InterestOverTime> {
    record.data.retain(|_, v| *v > 0.0);
    if record.data.is_empty() { None } else { Some(record) }
}

pub fn normalize(record: InterestOverTime) -> CleanedInterest {
    let max = record.data.values().cloned().fold(f64::NEG_INFINITY, f64::max);
    let already_normalized = max <= 100.0;

    let data = if already_normalized {
        record.data.clone()
    } else {
        record.data.iter().map(|(k, v)| (k.clone(), (v / max) * 100.0)).collect()
    };

    CleanedInterest {
        msg_type: record.msg_type,
        fetched_at: record.fetched_at,
        timestamp: record.timestamp,
        data,
        is_normalized: !already_normalized,
    }
}