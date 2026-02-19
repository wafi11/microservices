use std::collections::HashSet;
use std::env::args;

use rdkafka::consumer::{Consumer, StreamConsumer};
use rdkafka::producer::{FutureProducer, FutureRecord};
use rdkafka::util::Timeout;
use rdkafka::Message;

mod kafka;
mod types;
mod validate;
use crate::kafka::connection::check_kafka_connection;
use crate::kafka::consumer::create_consumer;
use crate::kafka::producer::create_producer;
use crate::validate::{filter_nulls,normalize,validate_json};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let broker = args().nth(1).unwrap_or("192.168.1.21:9092".to_string());

    if let Err(e) = check_kafka_connection(&broker).await {
        eprintln!("{}", e);
        std::process::exit(1);
    }

    let consumer: StreamConsumer = create_consumer(&broker)?;
    let producer: FutureProducer = create_producer(&broker)?;

    consumer.subscribe(&["google-trends"])?;
    println!("🚀 Rust cleaner aktif: google-trends → google-trends-clean\n");

    // Deduplikasi: simpan timestamp yang sudah diproses
    let mut seen_timestamps: HashSet<String> = HashSet::new();

    loop {
        match consumer.recv().await {
            Ok(msg) => {
                let key = msg
                    .key()
                    .and_then(|k| std::str::from_utf8(k).ok())
                    .unwrap_or("unknown")
                    .to_string();

                let payload = match msg.payload().and_then(|p| std::str::from_utf8(p).ok()) {
                    Some(p) => p.to_string(),
                    None => {
                        eprintln!("⚠️  Payload kosong, skip");
                        continue;
                    }
                };

                // Hanya proses interest_over_time
                if key != "interest_over_time" {
                    continue;
                }

                // 1. Validasi struktur JSON
                let record = match validate_json(&payload) {
                    Some(r) => r,
                    None => {
                        eprintln!("❌ Validasi gagal, skip");
                        continue;
                    }
                };

                // 2. Deduplikasi timestamp
                if seen_timestamps.contains(&record.timestamp) {
                    println!("⏭️  Duplikat timestamp {}, skip", record.timestamp);
                    continue;
                }
                seen_timestamps.insert(record.timestamp.clone());

                // 3. Filter nilai 0 / null
                let record = match filter_nulls(record) {
                    Some(r) => r,
                    None => {
                        println!("⚠️  Semua nilai 0, skip");
                        continue;
                    }
                };

                // 4. Normalisasi 0-100
                let cleaned = normalize(record);

                // Kirim ke topic clean
                let clean_payload = serde_json::to_string(&cleaned)?;
                producer
                    .send(
                        FutureRecord::to("google-trends-clean")
                            .key("interest_over_time")
                            .payload(clean_payload.as_bytes()),
                        Timeout::Never,
                    )
                    .await
                    .expect("Gagal kirim ke google-trends-clean");

                println!("✅ Cleaned & forwarded → {}", cleaned.timestamp);
            }
            Err(e) => eprintln!("❌ Error: {}", e),
        }
    }
}
