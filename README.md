# Alertmanager and karma hands-on – Companion Code

Companion code for the book *Alertmanager and karma hands-on* (German edition: *Alertmanager und
karma zum Anfassen*) by Thomas Zachmann (IT Drills series, itdrills.com / it-trainingslager.de).

Book page: https://itdrills.com/books/karma · German edition: https://it-trainingslager.de/buecher/karma

```bash
git clone https://github.com/thomaszachmann/karma-zum-anfassen.git ~/am-lab/code
```

## Contents

| Folder | Contents |
|---|---|
| `notes-api/` | notes API in Go with `/metrics` and built-in faults (`NOTES_DELAY`, `NOTES_ERROR_RATE`) |
| `alarm-log/` | small webhook receiver for Alertmanager, one line per notification |
| `loesung/endstand/` | all files of the `~/am-lab` lab after Day 15 |
| `loesung/pruefung/` | the same after the final exam (Appendix B) |

The solutions contain the address of the test VM `172.16.0.61`. Substitute your own.
All passwords and keys are lab values only (`lab-…-123`) and must never be used in
production.

## Versions

Prometheus 3.16.0, Alertmanager 0.35.0, karma v0.134, Postgres 18.7, Mailpit 1.33.1,
Caddy 2.11.9, k3s v1.37.2+k3s1, kube-prometheus-stack 93.3.0.

## License

The code is licensed under the MIT license (see `LICENSE`). The text of the book is
excluded.
