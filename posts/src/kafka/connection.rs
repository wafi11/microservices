use rdkafka::{ClientConfig, admin::AdminClient, client::DefaultClientContext};


pub async fn check_kafka_connection(broker: &str) -> Result<(), String> {
    let admin: AdminClient<DefaultClientContext> = ClientConfig::new()
        .set("bootstrap.servers", broker)
        .set("socket.timeout.ms", "3000")
        .create()
        .map_err(|e| format!("Gagal buat client: {}", e))?;

    match admin
        .inner()
        .fetch_metadata(None, std::time::Duration::from_secs(3))
    {
        Ok(metadata) => {
            println!("✅ Kafka terhubung!");
            println!("   Broker count : {}", metadata.brokers().len());
            println!("   Topic count  : {}", metadata.topics().len());
            Ok(())
        }
        Err(e) => Err(format!("❌ Gagal connect: {}", e)),
    }
}