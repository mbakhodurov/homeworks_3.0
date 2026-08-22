package model

// PartProperties — типоспецифичные свойства детали (Pointer Union Value Object).
// Ровно одно поле non-nil — определяется типом детали.
type PartProperties struct {
	hull   *HullProperties
	engine *EngineProperties
	shield *ShieldProperties
	weapon *WeaponProperties
}

// Hull возвращает свойства корпуса, если деталь — корпус.
func (p *PartProperties) Hull() *HullProperties { return p.hull }

// Engine возвращает свойства двигателя, если деталь — двигатель.
func (p *PartProperties) Engine() *EngineProperties { return p.engine }

// Shield возвращает свойства щита, если деталь — щит.
func (p *PartProperties) Shield() *ShieldProperties { return p.shield }

// Weapon возвращает свойства оружия, если деталь — оружие.
func (p *PartProperties) Weapon() *WeaponProperties { return p.weapon }
