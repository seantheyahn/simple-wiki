create or replace function set_update_ts () returns trigger AS' BEGIN NEW.updated_at = NOW();
                            RETURN NEW;
                                    END;
                                    'LANGUAGE 'plpgsql' IMMUTABLE CALLED ON NULL INPUT SECURITY INVOKER;