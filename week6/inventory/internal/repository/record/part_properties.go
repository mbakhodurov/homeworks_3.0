package record

// PartPropertiesRecord — структура для десериализации JSONB из колонки properties.
// Живёт в repository-слое: json-теги не принадлежат доменной модели.
// Ровно одно поле non-nil — Pointer Union, аналогично доменной модели.
type PartPropertiesRecord struct {
	Hull   *HullPropertiesRecord   `json:"hull,omitempty"`
	Engine *EnginePropertiesRecord `json:"engine,omitempty"`
	Shield *ShieldPropertiesRecord `json:"shield,omitempty"`
	Weapon *WeaponPropertiesRecord `json:"weapon,omitempty"`
}

// HullPropertiesRecord — JSONB-структура свойств корпуса.
type HullPropertiesRecord struct {
	Strength int `json:"strength"`
}

// EnginePropertiesRecord — JSONB-структура свойств двигателя.
type EnginePropertiesRecord struct {
	Class            string `json:"class"`
	RequiredStrength int    `json:"required_strength"`
}

// ShieldPropertiesRecord — JSONB-структура свойств щита.
type ShieldPropertiesRecord struct {
	ShieldType string `json:"shield_type"`
}

// WeaponPropertiesRecord — JSONB-структура свойств оружия.
type WeaponPropertiesRecord struct {
	WeaponType string `json:"weapon_type"`
}
