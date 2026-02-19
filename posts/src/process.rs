// use crate::types::{InterestOverTime, RelatedQueries};

// pub  fn process_message(key: &str, payload: &str) {
//     match key {
//         "interest_over_time" => {
//             match serde_json::from_str::<InterestOverTime>(payload) {
//                 Ok(data) => {
//                     println!("📊 [Interest Over Time] @ {}", data.timestamp);
//                     for (keyword, value) in &data.data {
//                         println!("   {:<12} : {}/100", keyword, value);
//                     }
//                     analyze_interest(&data);
//                     println!();
//                 }
//                 Err(e) => eprintln!("⚠️  Gagal parse interest_over_time: {}", e),
//             }
//         }

//         "related_queries_top" => {
//             match serde_json::from_str::<RelatedQueries>(payload) {
//                 Ok(data) => {
//                     println!("🔍 [Related Top] keyword: {}", data.keyword);
//                     for (i, q) in data.data.iter().take(5).enumerate() {
//                         println!("   {}. {} ({})", i + 1, q.query, q.value);
//                     }
//                     println!();
//                 }
//                 Err(e) => eprintln!("⚠️  Gagal parse related_queries_top: {}", e),
//             }
//         }

//         "related_queries_rising" => {
//             match serde_json::from_str::<RelatedQueries>(payload) {
//                 Ok(data) => {
//                     println!("📈 [Related Rising] keyword: {}", data.keyword);
//                     for (i, q) in data.data.iter().take(5).enumerate() {
//                         println!("   {}. {} ({})", i + 1, q.query, q.value);
//                     }
//                     println!();
//                 }
//                 Err(e) => eprintln!("⚠️  Gagal parse related_queries_rising: {}", e),
//             }
//         }

//         _ => {
//             println!("❓ Unknown message type: {}", key);
//         }
//     }
// }


// pub fn analyze_interest(data: &InterestOverTime) {
//     let ai = data.data.get("AI").copied().unwrap_or(0);
//     let chatgpt = data.data.get("ChatGPT").copied().unwrap_or(0);
//     let kesehatan = data.data.get("Kesehatan").copied().unwrap_or(0);

//     if ai > 80 {
//         println!("   ⚡ AI trending tinggi!");
//     }
//     if kesehatan > 80 {
//         println!("   ⚡ Kesehatan trending tinggi!");
//     }
//     if chatgpt > 50 {
//         println!("   ⚡ ChatGPT interest di atas rata-rata!");
//     }
//     if ai > chatgpt * 2 {
//         println!("   📌 AI jauh lebih populer dari ChatGPT saat ini");
//     }
// }