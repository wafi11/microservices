import json
import threading
import logging
from collections import defaultdict
from datetime import datetime

import numpy as np
import pandas as pd
from kafka import KafkaConsumer
from dash import Dash, dcc, html, Input, Output
import plotly.graph_objects as go
from plotly.subplots import make_subplots

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

KAFKA_BROKER = "192.168.1.21:9092"
TOPIC = "google-trends-clean"
KEYWORDS = ["AI", "ChatGPT"]
REFRESH_INTERVAL = 5000  # ms

# ── Shared state (thread-safe append) ─────────────────────────────────────
history: dict[str, list] = defaultdict(list)
timestamps: list = []
lock = threading.Lock()


# ── Kafka consumer thread ─────────────────────────────────────────────────

def kafka_thread():
    consumer = KafkaConsumer(
        TOPIC,
        bootstrap_servers=KAFKA_BROKER,
        value_deserializer=lambda v: json.loads(v.decode('utf-8')),
        auto_offset_reset='earliest',
        group_id='trends-dashboard'
    )
    logger.info("Kafka consumer thread started")
    for msg in consumer:
        data = msg.value
        if data.get("type") != "interest_over_time":
            continue
        with lock:
            timestamps.append(data.get("timestamp"))
            for kw in KEYWORDS:
                if kw in data.get("data", {}):
                    history[kw].append(data["data"][kw])


# ── Forecasting ───────────────────────────────────────────────────────────

def forecast(values: list, steps: int = 10) -> list:
    if len(values) < 10:
        return []
    x = np.arange(len(values))
    y = np.array(values)
    coeffs = np.polyfit(x, y, 1)
    poly = np.poly1d(coeffs)
    future_x = np.arange(len(values), len(values) + steps)
    predicted = np.clip(poly(future_x), 0, 100)
    return [round(float(v), 2) for v in predicted]


# ── Dash App ──────────────────────────────────────────────────────────────

app = Dash(__name__)

COLORS = {
    "bg":        "#0a0e1a",
    "panel":     "#0f1627",
    "border":    "#1e2d4a",
    "ai":        "#00d4ff",
    "chatgpt":   "#7c3aed",
    "forecast":  "#f59e0b",
    "text":      "#e2e8f0",
    "subtext":   "#64748b",
    "grid":      "#1a2540",
}

app.layout = html.Div(style={
    "backgroundColor": COLORS["bg"],
    "minHeight": "100vh",
    "fontFamily": "'JetBrains Mono', 'Fira Code', monospace",
    "color": COLORS["text"],
    "padding": "24px",
}, children=[

    # Google Fonts
    html.Link(rel="stylesheet", href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@300;400;600;700&display=swap"),

    # Header
    html.Div(style={
        "display": "flex",
        "justifyContent": "space-between",
        "alignItems": "center",
        "borderBottom": f"1px solid {COLORS['border']}",
        "paddingBottom": "16px",
        "marginBottom": "24px",
    }, children=[
        html.Div(children=[
            html.Div("TRENDS MONITOR", style={
                "fontSize": "11px",
                "letterSpacing": "4px",
                "color": COLORS["subtext"],
                "marginBottom": "4px",
            }),
            html.H1("Google Trends Dashboard", style={
                "margin": 0,
                "fontSize": "24px",
                "fontWeight": "700",
                "color": COLORS["text"],
            }),
        ]),
        html.Div(id="live-clock", style={
            "fontSize": "12px",
            "color": COLORS["ai"],
            "letterSpacing": "2px",
        }),
    ]),

    # Stats cards
    html.Div(id="stats-cards", style={
        "display": "grid",
        "gridTemplateColumns": "repeat(4, 1fr)",
        "gap": "16px",
        "marginBottom": "24px",
    }),

    # Charts row 1: Line + Bar
    html.Div(style={
        "display": "grid",
        "gridTemplateColumns": "2fr 1fr",
        "gap": "16px",
        "marginBottom": "16px",
    }, children=[
        html.Div(style={"backgroundColor": COLORS["panel"], "borderRadius": "12px", "padding": "20px", "border": f"1px solid {COLORS['border']}"}, children=[
            html.Div("INTEREST OVER TIME", style={"fontSize": "10px", "letterSpacing": "3px", "color": COLORS["subtext"], "marginBottom": "12px"}),
            dcc.Graph(id="line-chart", config={"displayModeBar": False}, style={"height": "280px"}),
        ]),
        html.Div(style={"backgroundColor": COLORS["panel"], "borderRadius": "12px", "padding": "20px", "border": f"1px solid {COLORS['border']}"}, children=[
            html.Div("KEYWORD COMPARISON", style={"fontSize": "10px", "letterSpacing": "3px", "color": COLORS["subtext"], "marginBottom": "12px"}),
            dcc.Graph(id="bar-chart", config={"displayModeBar": False}, style={"height": "280px"}),
        ]),
    ]),

    # Charts row 2: Forecast
    html.Div(style={"backgroundColor": COLORS["panel"], "borderRadius": "12px", "padding": "20px", "border": f"1px solid {COLORS['border']}", "marginBottom": "16px"}, children=[
        html.Div("FORECAST — LINEAR PROJECTION (NEXT 10 POINTS)", style={"fontSize": "10px", "letterSpacing": "3px", "color": COLORS["subtext"], "marginBottom": "12px"}),
        dcc.Graph(id="forecast-chart", config={"displayModeBar": False}, style={"height": "220px"}),
    ]),

    # Auto-refresh interval
    dcc.Interval(id="interval", interval=REFRESH_INTERVAL, n_intervals=0),
])


# ── Callbacks ─────────────────────────────────────────────────────────────

def make_stat_card(label, value, color, unit=""):
    return html.Div(style={
        "backgroundColor": COLORS["panel"],
        "border": f"1px solid {COLORS['border']}",
        "borderTop": f"2px solid {color}",
        "borderRadius": "12px",
        "padding": "16px 20px",
    }, children=[
        html.Div(label, style={"fontSize": "10px", "letterSpacing": "2px", "color": COLORS["subtext"], "marginBottom": "8px"}),
        html.Div(f"{value}{unit}", style={"fontSize": "28px", "fontWeight": "700", "color": color}),
    ])


@app.callback(
    Output("stats-cards", "children"),
    Output("line-chart", "figure"),
    Output("bar-chart", "figure"),
    Output("forecast-chart", "figure"),
    Output("live-clock", "children"),
    Input("interval", "n_intervals"),
)
def update_dashboard(_):
    with lock:
        ts = list(timestamps)
        hist = {kw: list(history[kw]) for kw in KEYWORDS}

    now = datetime.now().strftime("%Y-%m-%d  %H:%M:%S")

    # ── Stats cards ───────────────────────────────────────────────────────
    ai_vals = hist.get("AI", [])
    gpt_vals = hist.get("ChatGPT", [])

    cards = [
        make_stat_card("DATA POINTS", len(ts), COLORS["ai"]),
        make_stat_card("AI — CURRENT", ai_vals[-1] if ai_vals else 0, COLORS["ai"], "/100"),
        make_stat_card("CHATGPT — CURRENT", gpt_vals[-1] if gpt_vals else 0, COLORS["chatgpt"], "/100"),
        make_stat_card("AI AVG (7d)", round(float(np.mean(ai_vals)), 1) if ai_vals else 0, COLORS["forecast"], "/100"),
    ]

    chart_layout = dict(
        paper_bgcolor="rgba(0,0,0,0)",
        plot_bgcolor="rgba(0,0,0,0)",
        font=dict(family="JetBrains Mono", color=COLORS["text"], size=11),
        margin=dict(l=10, r=10, t=10, b=10),
        xaxis=dict(showgrid=True, gridcolor=COLORS["grid"], zeroline=False, tickfont=dict(size=10)),
        yaxis=dict(showgrid=True, gridcolor=COLORS["grid"], zeroline=False, range=[0, 105]),
        legend=dict(orientation="h", y=1.1, x=0),
        hovermode="x unified",
    )

    # ── Line chart ────────────────────────────────────────────────────────
    line_fig = go.Figure(layout=chart_layout)
    x_axis = ts if ts else list(range(max(len(ai_vals), len(gpt_vals))))

    for kw, color in [("AI", COLORS["ai"]), ("ChatGPT", COLORS["chatgpt"])]:
        vals = hist.get(kw, [])
        if vals:
            line_fig.add_trace(go.Scatter(
                x=x_axis[:len(vals)], y=vals,
                name=kw, mode="lines",
                line=dict(color=color, width=2),
                fill="tozeroy",
                fillcolor="rgba(0,212,255,0.05)" if kw == "AI" else "rgba(124,58,237,0.05)",
            ))

    # ── Bar chart ─────────────────────────────────────────────────────────
    bar_fig = go.Figure(layout={**chart_layout, "xaxis": dict(showgrid=False), "yaxis": dict(showgrid=True, gridcolor=COLORS["grid"], range=[0, 105])})
    avgs = {kw: round(float(np.mean(hist[kw])), 1) if hist[kw] else 0 for kw in KEYWORDS}
    bar_fig.add_trace(go.Bar(
        x=list(avgs.keys()),
        y=list(avgs.values()),
        marker_color=[COLORS["ai"], COLORS["chatgpt"]],
        text=[f"{v}/100" for v in avgs.values()],
        textposition="outside",
        marker_line_width=0,
    ))

    # ── Forecast chart ────────────────────────────────────────────────────
    fc_fig = go.Figure(layout=chart_layout)
    for kw, color in [("AI", COLORS["ai"]), ("ChatGPT", COLORS["chatgpt"])]:
        vals = hist.get(kw, [])
        pred = forecast(vals, steps=10)
        if not vals:
            continue

        x_hist = list(range(len(vals)))
        x_fore = list(range(len(vals), len(vals) + len(pred)))

        fc_fig.add_trace(go.Scatter(
            x=x_hist, y=vals, name=f"{kw} actual",
            mode="lines", line=dict(color=color, width=2),
        ))
        if pred:
            fc_fig.add_trace(go.Scatter(
                x=x_fore, y=pred, name=f"{kw} forecast",
                mode="lines+markers",
                line=dict(color=COLORS["forecast"], width=2, dash="dot"),
                marker=dict(size=5, color=COLORS["forecast"]),
            ))

    return cards, line_fig, bar_fig, fc_fig, now


# ── Run ───────────────────────────────────────────────────────────────────

if __name__ == "__main__":
    t = threading.Thread(target=kafka_thread, daemon=True)
    t.start()
    logger.info("🚀 Dashboard running → http://localhost:8050")
    app.run(debug=False, host="0.0.0.0", port=8050)