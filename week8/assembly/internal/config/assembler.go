package config

type AssemblerConfig struct {
	MinBuildTimeSec int `yaml:"min_build_time_sec" env:"ASSEMBLER_MIN_BUILD_TIME_SEC" env-default:"5"`
	MaxBuildTimeSec int `yaml:"max_build_time_sec" env:"ASSEMBLER_MAX_BUILD_TIME_SEC" env-default:"15"`
}
