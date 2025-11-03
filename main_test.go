package gostatistics

import (
	"errors"
	"math"
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name    string
		values  []float64
		want    float64
		wantErr error
	}{
		{
			name:    "empty slice",
			values:  []float64{},
			want:    0,
			wantErr: ErrEmptySlice,
		},
		{
			name:    "single value",
			values:  []float64{5},
			want:    5,
			wantErr: nil,
		},
		{
			name:    "multiple values",
			values:  []float64{1, 2, 3, 4, 5},
			want:    15,
			wantErr: nil,
		},
		{
			name:    "mixed values",
			values:  []float64{-10, 0, 10},
			want:    0,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sum(tt.values)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Sum() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got != tt.want {
				t.Fatalf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinMaxAndRange(t *testing.T) {
	values := []float64{5, 3, 9, -2, 7}

	t.Run("min", func(t *testing.T) {
		got, err := Min(values)
		if err != nil {
			t.Fatalf("Min() error = %v", err)
		}
		if got != -2 {
			t.Fatalf("Min() = %v, want -2", got)
		}
	})

	t.Run("max", func(t *testing.T) {
		got, err := Max(values)
		if err != nil {
			t.Fatalf("Max() error = %v", err)
		}
		if got != 9 {
			t.Fatalf("Max() = %v, want 9", got)
		}
	})

	t.Run("range", func(t *testing.T) {
		got, err := Range(values)
		if err != nil {
			t.Fatalf("Range() error = %v", err)
		}
		if got != 11 {
			t.Fatalf("Range() = %v, want 11", got)
		}
	})
}

func TestMedian(t *testing.T) {
	t.Run("odd number of values", func(t *testing.T) {
		got, err := Median([]float64{5, 1, 9})
		if err != nil {
			t.Fatalf("Median() error = %v", err)
		}
		if got != 5 {
			t.Fatalf("Median() = %v, want 5", got)
		}
	})

	t.Run("even number of values", func(t *testing.T) {
		got, err := Median([]float64{5, 1, 9, 11})
		if err != nil {
			t.Fatalf("Median() error = %v", err)
		}
		if got != 7 {
			t.Fatalf("Median() = %v, want 7", got)
		}
	})

	t.Run("error on empty slice", func(t *testing.T) {
		_, err := Median([]float64{})
		if !errors.Is(err, ErrEmptySlice) {
			t.Fatalf("Median() error = %v, want %v", err, ErrEmptySlice)
		}
	})
}

func TestMean(t *testing.T) {
	tests := []struct {
		name    string
		values  []float64
		want    float64
		wantErr bool
	}{
		{
			name:    "empty slice",
			values:  []float64{},
			want:    0,
			wantErr: true,
		},
		{
			name:    "single value",
			values:  []float64{5},
			want:    5,
			wantErr: false,
		},
		{
			name:    "multiple values",
			values:  []float64{1, 2, 3, 4, 5},
			want:    3,
			wantErr: false,
		},
		{
			name:    "negative values",
			values:  []float64{-1, -2, -3, -4, -5},
			want:    -3,
			wantErr: false,
		},
		{
			name:    "mixed values",
			values:  []float64{-10, 0, 10},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mean(tt.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mean() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && math.Abs(got-tt.want) > 1e-10 {
				t.Errorf("Mean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStandardDeviation(t *testing.T) {
	tests := []struct {
		name                string
		values              []float64
		useSampleCorrection bool
		want                float64
		wantErr             bool
	}{
		{
			name:                "empty slice",
			values:              []float64{},
			useSampleCorrection: true,
			want:                0,
			wantErr:             true,
		},
		{
			name:                "single value",
			values:              []float64{5},
			useSampleCorrection: true,
			want:                0,
			wantErr:             true,
		},
		{
			name:                "sample standard deviation",
			values:              []float64{2, 4, 4, 4, 5, 5, 7, 9},
			useSampleCorrection: true,
			want:                2.14,
			wantErr:             false,
		},
		{
			name:                "population standard deviation",
			values:              []float64{2, 4, 4, 4, 5, 5, 7, 9},
			useSampleCorrection: false,
			want:                2.0,
			wantErr:             false,
		},
		{
			name:                "identical values",
			values:              []float64{3, 3, 3, 3},
			useSampleCorrection: true,
			want:                0,
			wantErr:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := StatisticsOptions{
				UseSampleCorrection: tt.useSampleCorrection,
			}

			got, err := StandardDeviation(tt.values, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("StandardDeviation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && math.Abs(got-tt.want) > 0.01 {
				t.Errorf("StandardDeviation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariance(t *testing.T) {
	tests := []struct {
		name                string
		values              []float64
		useSampleCorrection bool
		want                float64
		wantErr             bool
	}{
		{
			name:                "empty slice",
			values:              []float64{},
			useSampleCorrection: true,
			want:                0,
			wantErr:             true,
		},
		{
			name:                "single value",
			values:              []float64{5},
			useSampleCorrection: true,
			want:                0,
			wantErr:             true,
		},
		{
			name:                "sample variance",
			values:              []float64{2, 4, 4, 4, 5, 5, 7, 9},
			useSampleCorrection: true,
			want:                4.57,
			wantErr:             false,
		},
		{
			name:                "population variance",
			values:              []float64{2, 4, 4, 4, 5, 5, 7, 9},
			useSampleCorrection: false,
			want:                4.0,
			wantErr:             false,
		},
		{
			name:                "identical values",
			values:              []float64{3, 3, 3, 3},
			useSampleCorrection: true,
			want:                0,
			wantErr:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := StatisticsOptions{
				UseSampleCorrection: tt.useSampleCorrection,
			}

			got, err := Variance(tt.values, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Variance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && math.Abs(got-tt.want) > 0.01 {
				t.Errorf("Variance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStandardDeviationWithMean(t *testing.T) {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	mean, _ := Mean(values) // Calculate the actual mean

	tests := []struct {
		name                string
		useSampleCorrection bool
		want                float64
	}{
		{
			name:                "sample standard deviation",
			useSampleCorrection: true,
			want:                2.14,
		},
		{
			name:                "population standard deviation",
			useSampleCorrection: false,
			want:                2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := StatisticsOptions{
				UseSampleCorrection: tt.useSampleCorrection,
			}

			got := StandardDeviationWithMean(values, mean, opts)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("StandardDeviationWithMean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVarianceWithMean(t *testing.T) {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	mean, _ := Mean(values) // Calculate the actual mean

	tests := []struct {
		name                string
		useSampleCorrection bool
		want                float64
	}{
		{
			name:                "sample variance",
			useSampleCorrection: true,
			want:                4.57,
		},
		{
			name:                "population variance",
			useSampleCorrection: false,
			want:                4.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := StatisticsOptions{
				UseSampleCorrection: tt.useSampleCorrection,
			}

			got := VarianceWithMean(values, mean, opts)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("VarianceWithMean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvenienceVarianceAndStd(t *testing.T) {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}

	t.Run("sample variance", func(t *testing.T) {
		got, err := SampleVariance(values)
		if err != nil {
			t.Fatalf("SampleVariance() error = %v", err)
		}
		if math.Abs(got-4.57) > 0.01 {
			t.Fatalf("SampleVariance() = %v, want 4.57", got)
		}
	})

	t.Run("population variance", func(t *testing.T) {
		got, err := PopulationVariance(values)
		if err != nil {
			t.Fatalf("PopulationVariance() error = %v", err)
		}
		if math.Abs(got-4.0) > 0.01 {
			t.Fatalf("PopulationVariance() = %v, want 4.0", got)
		}
	})

	t.Run("sample standard deviation", func(t *testing.T) {
		got, err := SampleStandardDeviation(values)
		if err != nil {
			t.Fatalf("SampleStandardDeviation() error = %v", err)
		}
		if math.Abs(got-2.14) > 0.01 {
			t.Fatalf("SampleStandardDeviation() = %v, want 2.14", got)
		}
	})

	t.Run("population standard deviation", func(t *testing.T) {
		got, err := PopulationStandardDeviation(values)
		if err != nil {
			t.Fatalf("PopulationStandardDeviation() error = %v", err)
		}
		if math.Abs(got-2.0) > 0.01 {
			t.Fatalf("PopulationStandardDeviation() = %v, want 2.0", got)
		}
	})
}

func TestDescribe(t *testing.T) {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}

	stats, err := Describe(values)
	if err != nil {
		t.Fatalf("Describe() error = %v", err)
	}

	if stats.Count != 8 {
		t.Fatalf("Describe() Count = %d, want 8", stats.Count)
	}
	if math.Abs(stats.Sum-40) > 1e-10 {
		t.Fatalf("Describe() Sum = %v, want 40", stats.Sum)
	}
	if math.Abs(stats.Mean-5) > 1e-10 {
		t.Fatalf("Describe() Mean = %v, want 5", stats.Mean)
	}
	if math.Abs(stats.Median-4.5) > 1e-10 {
		t.Fatalf("Describe() Median = %v, want 4.5", stats.Median)
	}
	if stats.Min != 2 {
		t.Fatalf("Describe() Min = %v, want 2", stats.Min)
	}
	if stats.Max != 9 {
		t.Fatalf("Describe() Max = %v, want 9", stats.Max)
	}
	if math.Abs(stats.Range-7) > 1e-10 {
		t.Fatalf("Describe() Range = %v, want 7", stats.Range)
	}
	if math.Abs(stats.SampleVariance-4.57) > 0.01 {
		t.Fatalf("Describe() SampleVariance = %v, want 4.57", stats.SampleVariance)
	}
	if math.Abs(stats.PopulationVariance-4.0) > 0.01 {
		t.Fatalf("Describe() PopulationVariance = %v, want 4.0", stats.PopulationVariance)
	}
	if math.Abs(stats.SampleStandardDeviation-2.14) > 0.01 {
		t.Fatalf("Describe() SampleStandardDeviation = %v, want 2.14", stats.SampleStandardDeviation)
	}
	if math.Abs(stats.PopulationStandardDeviation-2.0) > 0.01 {
		t.Fatalf("Describe() PopulationStandardDeviation = %v, want 2.0", stats.PopulationStandardDeviation)
	}
}

func TestDescribeSingleValue(t *testing.T) {
	stats, err := Describe([]float64{42})
	if err != nil {
		t.Fatalf("Describe() error = %v", err)
	}

	if stats.Count != 1 {
		t.Fatalf("Describe() Count = %d, want 1", stats.Count)
	}
	if stats.Sum != 42 {
		t.Fatalf("Describe() Sum = %v, want 42", stats.Sum)
	}
	if stats.Min != 42 || stats.Max != 42 {
		t.Fatalf("Describe() Min/Max = %v/%v, want 42/42", stats.Min, stats.Max)
	}
	if !math.IsNaN(stats.SampleVariance) {
		t.Fatalf("Describe() SampleVariance = %v, want NaN", stats.SampleVariance)
	}
	if !math.IsNaN(stats.SampleStandardDeviation) {
		t.Fatalf("Describe() SampleStandardDeviation = %v, want NaN", stats.SampleStandardDeviation)
	}
	if stats.PopulationVariance != 0 {
		t.Fatalf("Describe() PopulationVariance = %v, want 0", stats.PopulationVariance)
	}
	if stats.PopulationStandardDeviation != 0 {
		t.Fatalf("Describe() PopulationStandardDeviation = %v, want 0", stats.PopulationStandardDeviation)
	}
}
