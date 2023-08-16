package config

// LoadOption option function for loadding configuration
type LoadOption func(*MixgoConfig)

// WithCodec set the codec
func WithCodec(name string) LoadOption {
	return func(mc *MixgoConfig) {
		mc.decoder = GetCodec(name)
	}
}

// WithProvider set the provider
func WithProvider(name string) LoadOption {
	return func(mc *MixgoConfig) {
		mc.p = GetProvider(name)
	}
}

// WithReciver set the reciver, which must be a pointer
func WithReciver(reciver interface{}) LoadOption {
	return func(mc *MixgoConfig) {
		mc.r = reciver
	}
}

// WithDisableWatch whether disable watcher
func WithDisableWatch(disable bool) LoadOption {
	return func(mc *MixgoConfig) {
		mc.disableWatch = disable
	}
}
