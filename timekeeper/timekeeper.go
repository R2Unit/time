package timekeeper

import (
	"log"
	"sync"
	"time"
)

const ntpEpochOffset = 2208988800 // Seconds from 1900 to 1970

// DriftAwareTimeKeeper maintains high-precision time with drift correction and time zone support
type DriftAwareTimeKeeper struct {
	startTime   time.Time
	startNtp    uint64
	timeOffset  float64
	driftRate   float64
	mutex       sync.Mutex
	lastUpdated time.Time
	timeZone    *time.Location
}

// NewDriftAwareTimeKeeper initializes the TimeKeeper with time zone support
func NewDriftAwareTimeKeeper(timeZoneName string) *DriftAwareTimeKeeper {
	location, err := time.LoadLocation(timeZoneName)
	if err != nil {
		log.Fatalf("Failed to load time zone '%s': %v", timeZoneName, err)
	}

	now := time.Now()
	ntpTime := uint64(now.Unix()+ntpEpochOffset)<<32 | uint64((now.Nanosecond()*(1<<32))/1e9)

	return &DriftAwareTimeKeeper{
		startTime:   now,
		startNtp:    ntpTime,
		timeOffset:  0.0,
		driftRate:   0.0,
		lastUpdated: now,
		timeZone:    location,
	}
}

// GetCurrentNtpTime returns the current NTP time with drift correction
func (tk *DriftAwareTimeKeeper) GetCurrentNtpTime() uint64 {
	tk.mutex.Lock()
	defer tk.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tk.startTime).Seconds()
	driftAdjustment := elapsed * tk.driftRate
	return tk.startNtp + uint64((elapsed+tk.timeOffset+driftAdjustment)*(1<<32))
}

// GetLocalTime returns the current time in the configured time zone
func (tk *DriftAwareTimeKeeper) GetLocalTime() time.Time {
	tk.mutex.Lock()
	defer tk.mutex.Unlock()

	now := time.Now()
	return now.In(tk.timeZone)
}

// UpdateDriftRate adjusts the drift rate dynamically
func (tk *DriftAwareTimeKeeper) UpdateDriftRate(referenceNtpTime uint64) {
	tk.mutex.Lock()
	defer tk.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tk.lastUpdated).Seconds()
	if elapsed > 0 {
		localNtpTime := tk.GetCurrentNtpTime()
		drift := float64(referenceNtpTime-localNtpTime) / elapsed
		tk.driftRate += drift * 0.1 // Smoothing factor
		tk.lastUpdated = now
	}
}

// TimeZoneName returns the name of the configured time zone
func (tk *DriftAwareTimeKeeper) TimeZoneName() string {
	return tk.timeZone.String()
}
