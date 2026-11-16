package main

// setupTelemetry richtet OpenTelemetry ein. Noch tut es nichts: Ohne Provider sind
// alle Spans leer und kosten fast nichts. An Tag 10 ersetzt du diese Datei.

import "context"

func setupTelemetry(ctx context.Context) (shutdown func(context.Context) error, err error) {
	return func(context.Context) error { return nil }, nil
}
