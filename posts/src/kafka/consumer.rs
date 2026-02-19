use rdkafka::{ClientConfig, consumer::StreamConsumer};
use uuid::Uuid;

pub fn create_consumer(bootstrap_server: &str) -> Result<StreamConsumer, Box<dyn std::error::Error>> {
    Ok(ClientConfig::new()
        .set("bootstrap.servers", bootstrap_server)
        .set("enable.partition.eof", "false")
        .set("group.id", format!("chat-{}", Uuid::new_v4()))
        .create()
        .expect("Failed to create client"))
}