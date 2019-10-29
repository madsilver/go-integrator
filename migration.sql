-- Function integration_trigger
CREATE OR REPLACE FUNCTION integration_trigger() RETURNS trigger AS $$  
BEGIN  
    PERFORM pg_notify('realtime_location_record', row_to_json(NEW)::text);
    RETURN NEW;
END;  
$$ LANGUAGE plpgsql;

-- Trigger tg_realtime_locations_integration
CREATE TRIGGER tg_realtime_locations_integration 
AFTER INSERT OR UPDATE ON public.realtime_locations  
FOR EACH ROW 
EXECUTE PROCEDURE integration_trigger();
