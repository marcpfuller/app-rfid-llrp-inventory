//
// Copyright (C) 2020 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package inventory

var (
	assetTracking = MobilityProfile{
		Slope:         -0.008,
		Threshold:     6.0,
		HoldoffMillis: 500.0,
	}

	retailGarment = MobilityProfile{
		Slope:         -0.0005,
		Threshold:     6.0,
		HoldoffMillis: 60000.0,
	}

	// this will clone it
	defaultProfile = assetTracking

	mobilityProfiles = map[string]MobilityProfile{
		"default":        defaultProfile,
		"asset_tracking": assetTracking,
		"retail_garment": retailGarment,
	}
)

// MobilityProfile defines the parameters of the weighted slope formula used in calculating a tag's location.
// Tag location is determined based on the quality of tag reads associated with a sensor/antenna averaged over time.
// For a tag to move from one location to another, the other location must be either a better signal or be more recent.
type MobilityProfile struct {
	// Slope (dBm per millisecond): Used to determine the weight applied to older RSSI values
	Slope float64
	// Threshold (dBm) RSSI threshold that must be exceeded for the tag to move from the previous sensor
	Threshold float64
	// HoldoffMillis (milliseconds) Amount of time in which the weight used is just the threshold, effectively the slope is not used
	HoldoffMillis float64
	// b = y - (m*x)
	yIntercept float64
}

// b = y - (m*x)
func (profile *MobilityProfile) calculateYIntercept() {
	profile.yIntercept = profile.Threshold - (profile.Slope * profile.HoldoffMillis)
}

// loadMobilityProfile will attempt to load a mobility profile based on defaults and user's configuration
func loadMobilityProfile(cfg ApplicationSettings) MobilityProfile {
	profile := MobilityProfile{
		Slope:         cfg.MobilityProfileSlope,
		Threshold:     cfg.MobilityProfileThreshold,
		HoldoffMillis: cfg.MobilityProfileHoldoffMillis,
	}
	profile.calculateYIntercept()
	return profile
}

// ComputeWeight computes the weight to be applied to a value based on the time it was read vs the reference timestamp.
func (profile *MobilityProfile) ComputeWeight(referenceTimestamp int64, lastRead int64) float64 {
	// y = mx + b
	weight := (profile.Slope * float64(referenceTimestamp-lastRead)) + profile.yIntercept

	// check if weight needs to be capped at threshold ceiling
	if weight > profile.Threshold {
		weight = profile.Threshold
	}

	return weight
}
