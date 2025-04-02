package utils

import "time"

// CalcPause berechnet Pausen in Minuten basierend auf der Arbeitsdauer
func CalcPause(start, end time.Time) int {
    dauer := end.Sub(start).Hours()
    switch {
    case dauer >= 9:
        return 45
    case dauer >= 6:
        return 30
    default:
        return 0
    }
}