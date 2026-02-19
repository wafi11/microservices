import json
import logging
import numpy as np
import pandas as pd
from collections import defaultdict
from kafka import KafkaConsumer
from datetime import datetime

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

KAFKA_BROKER = "192.168.1.21:9092"
TOPIC = "google-trends-clean"
KEYWORDS = ["AI", "ChatGPT"]

# Simpan history data untuk analisis
history: dict[str, list] = defaultdict(list)  # {"AI": [val1, val2, ...], ...}
timestamps: list = []


def create_consumer():
    return KafkaConsumer(
        TOPIC,
        bootstrap_servers=KAFKA_BROKER,
        value_deserializer=lambda v: json.loads(v.decode('utf-8')),
        auto_offset_reset='earliest',
        group_id='trends-analyzer'
    )


# ── Analisis ──────────────────────────────────────────────────────────────

def hitung_rata_rata(keyword: str) -> dict:
    values = history[keyword]
    if not values:
        return {}
    arr = np.array(values)
    return {
        "mean": round(float(arr.mean()), 2),
        "median": round(float(np.median(arr)), 2),
        "std": round(float(arr.std()), 2),
        "min": round(float(arr.min()), 2),
        "max": round(float(arr.max()), 2),
    }


def hitung_tren(keyword: str) -> str:
    values = history[keyword]
    if len(values) < 5:
        return "data belum cukup"

    # Gunakan linear regression sederhana
    x = np.arange(len(values))
    y = np.array(values)
    slope = np.polyfit(x, y, 1)[0]

    if slope > 1.0:
        return f"📈 naik (slope: +{slope:.2f})"
    elif slope < -1.0:
        return f"📉 turun (slope: {slope:.2f})"
    else:
        return f"➡️  stabil (slope: {slope:.2f})"


def forecast(keyword: str, steps: int = 5) -> list:
    """Simple linear forecasting"""
    values = history[keyword]
    if len(values) < 10:
        return []

    x = np.arange(len(values))
    y = np.array(values)
    coeffs = np.polyfit(x, y, 1)  # linear fit
    poly = np.poly1d(coeffs)

    future_x = np.arange(len(values), len(values) + steps)
    predicted = poly(future_x)

    # Clamp ke 0-100
    predicted = np.clip(predicted, 0, 100)
    return [round(float(v), 2) for v in predicted]


def print_analisis():
    print("\n" + "="*50)
    print(f"📊 ANALISIS  |  {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print("="*50)

    for kw in KEYWORDS:
        if not history[kw]:
            continue

        rata = hitung_rata_rata(kw)
        tren = hitung_tren(kw)
        pred = forecast(kw, steps=3)

        print(f"\n🔑 Keyword: {kw}")
        print(f"   Data points : {len(history[kw])}")
        print(f"   Rata-rata   : {rata.get('mean')} | Median: {rata.get('median')}")
        print(f"   Min/Max     : {rata.get('min')} / {rata.get('max')} | Std: {rata.get('std')}")
        print(f"   Tren        : {tren}")
        if pred:
            print(f"   Forecast    : {pred} (3 titik ke depan)")

    print("="*50 + "\n")


# ── Main ──────────────────────────────────────────────────────────────────

def main():
    consumer = create_consumer()
    logger.info(f"✅ Consumer analisis aktif, listening: {TOPIC}")

    msg_count = 0

    for msg in consumer:
        data = msg.value

        if data.get("type") != "interest_over_time":
            continue

        timestamp = data.get("timestamp")
        values = data.get("data", {})

        timestamps.append(timestamp)
        for kw in KEYWORDS:
            if kw in values:
                history[kw].append(values[kw])

        msg_count += 1
        logger.info(f"Received #{msg_count} | {timestamp} | {values}")

        if msg_count % 10 == 0:
            print_analisis()


if __name__ == "__main__":
    main()