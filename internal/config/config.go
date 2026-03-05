package config

// VarDef defines a single environment variable.
type VarDef struct {
	Value     string      `yaml:"value"`
	Command   string      `yaml:"command"`
	Test      string      `yaml:"test"`
	Required  *bool       `yaml:"required"`
	Cache     bool        `yaml:"cache"`
	Secret    bool        `yaml:"secret"`
	Decrypt   *DecryptDef `yaml:"decrypt"`
	Signature *SigDef     `yaml:"signature"`
}

// IsRequired returns true if the variable is required (default: true).
func (v VarDef) IsRequired() bool {
	if v.Required == nil {
		return true
	}
	return *v.Required
}

// DecryptDef defines how to decrypt an encrypted value.
type DecryptDef struct {
	Algorithm string `yaml:"algorithm"`
	Key       string `yaml:"key"`
	Hidden    bool   `yaml:"hidden"`
}

// SigDef defines how to verify a cryptographic signature.
type SigDef struct {
	Algorithm string `yaml:"algorithm"`
	Key       string `yaml:"key"`
	Signature string `yaml:"signature"`
}

// Config is the top-level smartvar configuration.
type Config struct {
	Vars map[string]VarDef `yaml:"vars"`
}
