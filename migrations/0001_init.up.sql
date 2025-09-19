CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE breeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    api_id TEXT UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE cats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    years_experience INT NOT NULL CHECK (years_experience >= 0),
    breed_id UUID REFERENCES breeds(id) ON DELETE SET NULL,
    salary NUMERIC(12,2) NOT NULL CHECK (salary >= 0),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
CREATE INDEX idx_cats_name ON cats (name);

CREATE TABLE missions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    assigned_cat_id UUID REFERENCES cats(id) ON DELETE SET NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
CREATE UNIQUE INDEX uq_missions_assigned_cat_active
    ON missions(assigned_cat_id)
    WHERE assigned_cat_id IS NOT NULL AND completed = FALSE;

CREATE TABLE targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mission_id UUID NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    country TEXT NOT NULL,
    notes TEXT DEFAULT '',
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT chk_name_not_empty CHECK (char_length(name) > 0)
);
CREATE INDEX idx_targets_mission_id ON targets (mission_id);

CREATE OR REPLACE FUNCTION touch_updated_at() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END;
$$;

CREATE TRIGGER cats_touch_updated_at BEFORE UPDATE ON cats
    FOR EACH ROW EXECUTE FUNCTION touch_updated_at();

CREATE TRIGGER missions_touch_updated_at BEFORE UPDATE ON missions
    FOR EACH ROW EXECUTE FUNCTION touch_updated_at();

CREATE TRIGGER targets_touch_updated_at BEFORE UPDATE ON targets
    FOR EACH ROW EXECUTE FUNCTION touch_updated_at();

CREATE OR REPLACE FUNCTION prevent_delete_assigned_mission() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF OLD.assigned_cat_id IS NOT NULL THEN
        RAISE EXCEPTION 'Cannot delete mission %: assigned to cat %', OLD.id, OLD.assigned_cat_id;
    END IF;
    RETURN OLD;
END;
$$;
CREATE TRIGGER trg_prevent_delete_assigned_mission
    BEFORE DELETE ON missions
    FOR EACH ROW EXECUTE FUNCTION prevent_delete_assigned_mission();

CREATE OR REPLACE FUNCTION ensure_targets_count_bounds() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    t_count INT;
    mission UUID;
BEGIN
    IF TG_OP = 'DELETE' THEN
        mission := OLD.mission_id;
    ELSE
        mission := COALESCE(NEW.mission_id, OLD.mission_id);
    END IF;

    SELECT COUNT(*) INTO t_count FROM targets WHERE mission_id = mission;

    IF t_count < 1 THEN
        RAISE EXCEPTION 'Mission % must have at least 1 target (current: %)', mission, t_count;
    ELSIF t_count > 3 THEN
        RAISE EXCEPTION 'Mission % cannot have more than 3 targets (current: %)', mission, t_count;
    END IF;

    RETURN NULL;
END;
$$;
CREATE TRIGGER trg_targets_count_after_insert
    AFTER INSERT ON targets
    FOR EACH ROW EXECUTE FUNCTION ensure_targets_count_bounds();
CREATE TRIGGER trg_targets_count_after_delete
    AFTER DELETE ON targets
    FOR EACH ROW EXECUTE FUNCTION ensure_targets_count_bounds();
CREATE TRIGGER trg_targets_count_after_update
    AFTER UPDATE ON targets
    FOR EACH ROW
    WHEN (OLD.mission_id IS DISTINCT FROM NEW.mission_id)
    EXECUTE FUNCTION ensure_targets_count_bounds();

CREATE OR REPLACE FUNCTION prevent_delete_completed_target() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF OLD.completed THEN
        RAISE EXCEPTION 'Cannot delete target %: it is completed', OLD.id;
    END IF;
    RETURN OLD;
END;
$$;
CREATE TRIGGER trg_prevent_delete_completed_target
    BEFORE DELETE ON targets
    FOR EACH ROW EXECUTE FUNCTION prevent_delete_completed_target();

CREATE OR REPLACE FUNCTION prevent_add_target_to_completed_mission() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    m_completed BOOLEAN;
BEGIN
    SELECT completed INTO m_completed FROM missions WHERE id = NEW.mission_id;
    IF m_completed THEN
        RAISE EXCEPTION 'Cannot add target to mission %: mission completed', NEW.mission_id;
    END IF;
    RETURN NEW;
END;
$$;
CREATE TRIGGER trg_prevent_add_target_to_completed_mission
    BEFORE INSERT ON targets
    FOR EACH ROW EXECUTE FUNCTION prevent_add_target_to_completed_mission();

CREATE OR REPLACE FUNCTION prevent_notes_update_if_completed() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    m_completed BOOLEAN;
BEGIN
    IF TG_OP = 'UPDATE' AND NEW.notes IS DISTINCT FROM OLD.notes THEN
        IF OLD.completed THEN
            RAISE EXCEPTION 'Cannot update notes for completed target %', OLD.id;
        END IF;
        SELECT completed INTO m_completed FROM missions WHERE id = OLD.mission_id;
        IF m_completed THEN
            RAISE EXCEPTION 'Cannot update notes for target %: mission % is completed', OLD.id, OLD.mission_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$;
CREATE TRIGGER trg_prevent_notes_update_if_completed
    BEFORE UPDATE ON targets
    FOR EACH ROW EXECUTE FUNCTION prevent_notes_update_if_completed();

CREATE OR REPLACE FUNCTION ensure_all_targets_completed_before_marking_mission() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    incomplete_count INT;
BEGIN
    IF TG_OP = 'UPDATE' AND NEW.completed = TRUE AND OLD.completed = FALSE THEN
        SELECT COUNT(*) INTO incomplete_count FROM targets WHERE mission_id = NEW.id AND completed = FALSE;
        IF incomplete_count > 0 THEN
            RAISE EXCEPTION 'Cannot complete mission %: % targets still incomplete', NEW.id, incomplete_count;
        END IF;
    END IF;
    RETURN NEW;
END;
$$;
CREATE TRIGGER trg_ensure_all_targets_completed
    BEFORE UPDATE ON missions
    FOR EACH ROW EXECUTE FUNCTION ensure_all_targets_completed_before_marking_mission();

CREATE OR REPLACE FUNCTION auto_complete_mission_when_all_targets_done() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    incomplete_count INT;
BEGIN
    IF OLD.completed = FALSE AND NEW.completed = TRUE THEN
        SELECT COUNT(*) INTO incomplete_count FROM targets WHERE mission_id = NEW.mission_id AND completed = FALSE;
        IF incomplete_count = 0 THEN
            UPDATE missions SET completed = TRUE, updated_at = now() WHERE id = NEW.mission_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$;
CREATE TRIGGER trg_auto_complete_mission
    AFTER UPDATE ON targets
    FOR EACH ROW
    WHEN (OLD.completed = FALSE AND NEW.completed = TRUE)
    EXECUTE FUNCTION auto_complete_mission_when_all_targets_done();

INSERT INTO breeds (api_id, name) VALUES
    ('abyssinian', 'Abyssinian'),
    ('beng', 'Bengal'),
    ('siamese', 'Siamese')
ON CONFLICT DO NOTHING;
