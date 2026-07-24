function "set_updated_at" {
  schema = schema.public
  return = trigger
  lang   = "PLPGSQL"
  as     = <<-SQL
    BEGIN
        NEW.updated_at = CURRENT_TIMESTAMP;
        RETURN NEW;
    END;
  SQL
}