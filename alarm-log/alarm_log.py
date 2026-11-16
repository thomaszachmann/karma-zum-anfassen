"""alarm-log: nimmt Webhooks vom Alertmanager an.

Jede Benachrichtigung wird zu einer Zeile im Log:
  Uhrzeit  Pfad  Empfänger  Status  [Gruppen-Labels]  Alarmnamen  (firing/resolved)
GET / liefert die zuletzt empfangene Benachrichtigung vollständig als JSON.
"""
import json
from datetime import datetime, timezone
from http.server import BaseHTTPRequestHandler, HTTPServer

last = b"{}"


class Handler(BaseHTTPRequestHandler):
    def do_POST(self):
        global last
        body = self.rfile.read(int(self.headers.get("Content-Length", 0)))
        last = body
        msg = json.loads(body)
        alerts = msg.get("alerts", [])
        names = ",".join(sorted({a["labels"].get("alertname", "?") for a in alerts}))
        group = " ".join(f"{k}={v}" for k, v in sorted(msg["groupLabels"].items()))
        firing = sum(1 for a in alerts if a["status"] == "firing")
        now = datetime.now(timezone.utc).strftime("%H:%M:%S")
        print(f"{now} {self.path} {msg['receiver']} {msg['status']} [{group}] {names} "
              f"({firing} firing, {len(alerts) - firing} resolved)", flush=True)
        self.send_response(200)
        self.end_headers()

    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(last)

    def log_message(self, *args):
        pass


HTTPServer(("", 5001), Handler).serve_forever()
