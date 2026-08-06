-- +goose Up
-- Заполняем properties для seed-деталей из миграции 00002 в соответствии с их типом.
-- Двигатель класса B (required_strength=70) совместим с титановым корпусом (strength=120),
-- но не с алюминиевым (strength=50) — для тестирования правила совместимости.
UPDATE parts SET properties = '{"hull": {"strength": 50}}'::jsonb
    WHERE uuid = '550e8400-e29b-41d4-a716-446655440001'; -- Алюминиевый корпус

UPDATE parts SET properties = '{"hull": {"strength": 120}}'::jsonb
    WHERE uuid = '550e8400-e29b-41d4-a716-446655440002'; -- Титановый корпус

UPDATE parts SET properties = '{"engine": {"class": "C", "required_strength": 30}}'::jsonb
    WHERE uuid = '550e8400-e29b-41d4-a716-446655440003'; -- Ионный двигатель C

UPDATE parts SET properties = '{"engine": {"class": "B", "required_strength": 70}}'::jsonb
    WHERE uuid = '550e8400-e29b-41d4-a716-446655440004'; -- Ионный двигатель B

UPDATE parts SET properties = '{"shield": {"shield_type": "energy"}}'::jsonb
    WHERE uuid = '550e8400-e29b-41d4-a716-446655440005'; -- Энергетический щит

UPDATE parts SET properties = '{"weapon": {"weapon_type": "laser"}}'::jsonb
    WHERE uuid = '550e8400-e29b-41d4-a716-446655440006'; -- Лазерная пушка

UPDATE parts SET properties = '{"hull": {"strength": 90}}'::jsonb
    WHERE uuid = '550e8400-e29b-41d4-a716-446655440007'; -- Плазменный корпус

-- +goose Down
UPDATE parts SET properties = '{}'::jsonb WHERE uuid IN (
    '550e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440003',
    '550e8400-e29b-41d4-a716-446655440004',
    '550e8400-e29b-41d4-a716-446655440005',
    '550e8400-e29b-41d4-a716-446655440006',
    '550e8400-e29b-41d4-a716-446655440007'
);
