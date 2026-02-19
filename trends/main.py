import json
import time
import logging
from datetime import datetime
from pytrends.request import TrendReq
from kafka import KafkaProducer

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Config
KAFKA_BROKER = "192.168.1.21:9092"
TOPIC = "google-trends"
KEYWORDS = ["AI", "ChatGPT"]
INTERVAL_SECONDS = 60 * 60  # fetch setiap 1 jam

def create_producer():
    return KafkaProducer(
        bootstrap_servers=KAFKA_BROKER,
        value_serializer=lambda v: json.dumps(v).encode('utf-8'),
        key_serializer=lambda k: k.encode('utf-8')
    )

def fetch_trends():
    pytrends = TrendReq(hl='id-ID', tz=420, timeout=(10, 25), retries=2, backoff_factor=0.5)

    pytrends.build_payload(
        KEYWORDS,
        timeframe='now 7-d',
        geo='ID'
    )

    # Interest over time
    df = pytrends.interest_over_time()
    df = df.drop(columns=['isPartial'])

    # Related queries
    related = pytrends.related_queries()

    return df, related

def build_payload(df, related):
    records = []

    # Interest over time
    for timestamp, row in df.iterrows():
        record = {
            "type": "interest_over_time",
            "fetched_at": datetime.utcnow().isoformat(),
            "timestamp": timestamp.isoformat(),
            "data": {kw: int(row[kw]) for kw in KEYWORDS}
        }
        records.append(record)

    # Related queries
    for kw in KEYWORDS:
        top = related.get(kw, {}).get('top')
        rising = related.get(kw, {}).get('rising')

        if top is not None:
            records.append({
                "type": "related_queries_top",
                "fetched_at": datetime.utcnow().isoformat(),
                "keyword": kw,
                "data": top.to_dict(orient='records')
            })

        if rising is not None:
            records.append({
                "type": "related_queries_rising",
                "fetched_at": datetime.utcnow().isoformat(),
                "keyword": kw,
                "data": rising.to_dict(orient='records')
            })

    return records

def main():
    producer = create_producer()
    logger.info(f"Producer connected ke {KAFKA_BROKER}, topic: {TOPIC}")

    while True:
        try:
            logger.info("Fetching data dari Google Trends...")
            df, related = fetch_trends()
            records = build_payload(df, related)

            for record in records:
                key = record["type"]
                producer.send(TOPIC, key=key, value=record)
                logger.info(f"Sent → type: {record['type']}, key: {key}")

            producer.flush()
            logger.info(f"Done. Total {len(records)} records dikirim. Tunggu {INTERVAL_SECONDS}s...")
            time.sleep(INTERVAL_SECONDS)

        except Exception as e:
            logger.error(f"Error: {e}")
            time.sleep(30)  # retry setelah 30 detik kalau error

if __name__ == "__main__":
    main()