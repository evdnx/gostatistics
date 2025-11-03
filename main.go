package gostatistics

import (
	"errors"
	"math"
	"sort"
)

// ErrEmptySlice indicates that a statistical operation was attempted on an empty slice.
var ErrEmptySlice = errors.New("gostatistics: cannot operate on an empty slice")

// ErrInsufficientValues indicates that a statistical calculation requires at least two data points.
var ErrInsufficientValues = errors.New("gostatistics: need at least two values to perform calculation")

// StatisticsOptions contains options for statistical calculations.
type StatisticsOptions struct {
	// UseSampleCorrection determines whether to use Bessel's correction (n-1) for sample statistics.
	// Default is true (use n-1 in the denominator).
	UseSampleCorrection bool
}

// DefaultStatisticsOptions returns the default statistics options.
func DefaultStatisticsOptions() StatisticsOptions {
	return StatisticsOptions{
		UseSampleCorrection: true,
	}
}

// Sum calculates the sum of a slice of values.
func Sum(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, ErrEmptySlice
	}
	return sum(values), nil
}

// Mean calculates the arithmetic mean of a slice of values.
func Mean(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, ErrEmptySlice
	}

	return sum(values) / float64(len(values)), nil
}

// Min returns the smallest value in the slice.
func Min(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, ErrEmptySlice
	}

	minVal := values[0]
	for _, v := range values[1:] {
		if v < minVal {
			minVal = v
		}
	}

	return minVal, nil
}

// Max returns the largest value in the slice.
func Max(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, ErrEmptySlice
	}

	maxVal := values[0]
	for _, v := range values[1:] {
		if v > maxVal {
			maxVal = v
		}
	}

	return maxVal, nil
}

// Range returns the difference between the largest and smallest values in the slice.
func Range(values []float64) (float64, error) {
	minVal, err := Min(values)
	if err != nil {
		return 0, err
	}

	maxVal, err := Max(values)
	if err != nil {
		return 0, err
	}

	return maxVal - minVal, nil
}

// Median calculates the median value of a slice of values.
func Median(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, ErrEmptySlice
	}

	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)
	mid := len(sorted) / 2

	if len(sorted)%2 == 1 {
		return sorted[mid], nil
	}

	return (sorted[mid-1] + sorted[mid]) / 2, nil
}

// StandardDeviation calculates the standard deviation of a slice of values.
// By default, it uses Bessel's correction (n-1 in the denominator) for sample standard deviation.
func StandardDeviation(values []float64, options ...StatisticsOptions) (float64, error) {
	if len(values) < 2 {
		return 0, ErrInsufficientValues
	}

	mean, err := Mean(values)
	if err != nil {
		return 0, err
	}

	return StandardDeviationWithMean(values, mean, options...), nil
}

// StandardDeviationWithMean calculates the standard deviation using a pre-calculated mean.
// This is useful when you already have the mean and want to avoid recalculating it.
func StandardDeviationWithMean(values []float64, mean float64, options ...StatisticsOptions) float64 {
	opts := resolveOptions(options)

	if len(values) < 2 {
		return 0
	}

	var sumSquaredDiff float64
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}

	var variance float64
	if opts.UseSampleCorrection {
		variance = sumSquaredDiff / float64(len(values)-1)
	} else {
		variance = sumSquaredDiff / float64(len(values))
	}

	return math.Sqrt(variance)
}

// SampleStandardDeviation calculates the sample standard deviation (uses Bessel's correction).
func SampleStandardDeviation(values []float64) (float64, error) {
	return StandardDeviation(values, StatisticsOptions{UseSampleCorrection: true})
}

// PopulationStandardDeviation calculates the population standard deviation.
func PopulationStandardDeviation(values []float64) (float64, error) {
	return StandardDeviation(values, StatisticsOptions{UseSampleCorrection: false})
}

// Variance calculates the variance of a slice of values.
func Variance(values []float64, options ...StatisticsOptions) (float64, error) {
	if len(values) < 2 {
		return 0, ErrInsufficientValues
	}

	mean, err := Mean(values)
	if err != nil {
		return 0, err
	}

	return VarianceWithMean(values, mean, options...), nil
}

// VarianceWithMean calculates the variance using a pre-calculated mean.
func VarianceWithMean(values []float64, mean float64, options ...StatisticsOptions) float64 {
	opts := resolveOptions(options)

	if len(values) < 2 {
		return 0
	}

	var sumSquaredDiff float64
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}

	if opts.UseSampleCorrection {
		return sumSquaredDiff / float64(len(values)-1)
	}

	return sumSquaredDiff / float64(len(values))
}

// SampleVariance calculates the sample variance (uses Bessel's correction).
func SampleVariance(values []float64) (float64, error) {
	return Variance(values, StatisticsOptions{UseSampleCorrection: true})
}

// PopulationVariance calculates the population variance.
func PopulationVariance(values []float64) (float64, error) {
	return Variance(values, StatisticsOptions{UseSampleCorrection: false})
}

// DescriptiveStats contains a collection of descriptive statistics derived from a data set.
type DescriptiveStats struct {
	Count                       int
	Sum                         float64
	Mean                        float64
	Median                      float64
	Min                         float64
	Max                         float64
	Range                       float64
	SampleVariance              float64
	PopulationVariance          float64
	SampleStandardDeviation     float64
	PopulationStandardDeviation float64
}

// Describe returns a DescriptiveStats summary for the provided values.
func Describe(values []float64) (DescriptiveStats, error) {
	if len(values) == 0 {
		return DescriptiveStats{}, ErrEmptySlice
	}

	minVal, err := Min(values)
	if err != nil {
		return DescriptiveStats{}, err
	}

	maxVal, err := Max(values)
	if err != nil {
		return DescriptiveStats{}, err
	}

	median, err := Median(values)
	if err != nil {
		return DescriptiveStats{}, err
	}

	total := sum(values)
	count := len(values)
	mean := total / float64(count)

	stats := DescriptiveStats{
		Count:                       count,
		Sum:                         total,
		Mean:                        mean,
		Median:                      median,
		Min:                         minVal,
		Max:                         maxVal,
		Range:                       maxVal - minVal,
		SampleVariance:              math.NaN(),
		PopulationVariance:          math.NaN(),
		SampleStandardDeviation:     math.NaN(),
		PopulationStandardDeviation: math.NaN(),
	}

	if count == 1 {
		stats.SampleVariance = math.NaN()
		stats.SampleStandardDeviation = math.NaN()
		stats.PopulationVariance = 0
		stats.PopulationStandardDeviation = 0
		return stats, nil
	}

	if variance, err := SampleVariance(values); err == nil {
		stats.SampleVariance = variance
	}
	if variance, err := PopulationVariance(values); err == nil {
		stats.PopulationVariance = variance
	}
	if std, err := SampleStandardDeviation(values); err == nil {
		stats.SampleStandardDeviation = std
	}
	if std, err := PopulationStandardDeviation(values); err == nil {
		stats.PopulationStandardDeviation = std
	}

	return stats, nil
}

func resolveOptions(options []StatisticsOptions) StatisticsOptions {
	if len(options) == 0 {
		return DefaultStatisticsOptions()
	}
	return options[0]
}

func sum(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total
}
