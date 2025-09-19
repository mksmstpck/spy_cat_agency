DROP TRIGGER IF EXISTS trg_auto_complete_mission ON targets;
DROP FUNCTION IF EXISTS auto_complete_mission_when_all_targets_done;

DROP TRIGGER IF EXISTS trg_ensure_all_targets_completed ON missions;
DROP FUNCTION IF EXISTS ensure_all_targets_completed_before_marking_mission;

DROP TRIGGER IF EXISTS trg_prevent_notes_update_if_completed ON targets;
DROP FUNCTION IF EXISTS prevent_notes_update_if_completed;

DROP TRIGGER IF EXISTS trg_prevent_add_target_to_completed_mission ON targets;
DROP FUNCTION IF EXISTS prevent_add_target_to_completed_mission;

DROP TRIGGER IF EXISTS trg_prevent_delete_completed_target ON targets;
DROP FUNCTION IF EXISTS prevent_delete_completed_target;

DROP TRIGGER IF EXISTS trg_targets_count_after_update ON targets;
DROP TRIGGER IF EXISTS trg_targets_count_after_delete ON targets;
DROP TRIGGER IF EXISTS trg_targets_count_after_insert ON targets;
DROP FUNCTION IF EXISTS ensure_targets_count_bounds;

DROP TRIGGER IF EXISTS trg_prevent_delete_assigned_mission ON missions;
DROP FUNCTION IF EXISTS prevent_delete_assigned_mission;

DROP TRIGGER IF EXISTS targets_touch_updated_at ON targets;
DROP TRIGGER IF EXISTS missions_touch_updated_at ON missions;
DROP TRIGGER IF EXISTS cats_touch_updated_at ON cats;
DROP FUNCTION IF EXISTS touch_updated_at;

DROP TABLE IF EXISTS targets;
DROP INDEX IF EXISTS uq_missions_assigned_cat_active;
DROP TABLE IF EXISTS missions;
DROP INDEX IF EXISTS idx_cats_name;
DROP TABLE IF EXISTS cats;
DROP TABLE IF EXISTS breeds;
