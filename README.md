# Alertmanager und karma zum Anfassen – Begleitcode

Begleitcode zum Buch „Alertmanager und karma zum Anfassen“ von Thomas Zachmann
(Reihe IT-Trainingslager, it-trainingslager.de).

```bash
git clone https://github.com/thomaszachmann/karma-zum-anfassen.git ~/am-lab/code
```

## Inhalt

| Ordner | Inhalt |
|---|---|
| `notes-api/` | Notiz-API in Go mit `/metrics` und eingebauten Störungen (`NOTES_DELAY`, `NOTES_ERROR_RATE`) |
| `alarm-log/` | kleiner Webhook-Empfänger für den Alertmanager, eine Zeile pro Benachrichtigung |
| `loesung/endstand/` | alle Dateien des Labs `~/am-lab` nach Tag 15 |
| `loesung/pruefung/` | dasselbe nach der Abschlussprüfung (Anhang B) |

In den Lösungen steht die Adresse der Test-VM `172.16.0.61`. Setze deine eigene ein.
Alle Passwörter und Schlüssel sind reine Lab-Werte (`lab-…-123`) und dürfen nie
produktiv verwendet werden.

## Versionen

Prometheus 3.16.0, Alertmanager 0.35.0, karma v0.134, Postgres 18.7, Mailpit 1.33.1,
Caddy 2.11.9, k3s v1.37.2+k3s1, kube-prometheus-stack 93.3.0.

## Lizenz

Der Code steht unter der MIT-Lizenz (siehe `LICENSE`). Die Texte des Buchs sind davon
ausgenommen.
