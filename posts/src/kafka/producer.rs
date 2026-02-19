use rdkafka::config::ClientConfig;
use rdkafka::producer::{FutureProducer};

pub  fn create_producer(bootstrap_server: &str) -> Result<FutureProducer, Box<dyn std::error::Error>> {
    Ok(ClientConfig::new()
        .set("bootstrap.servers", bootstrap_server)
        .set("queue.buffering.max.ms", "0")
        .create()?)
}

