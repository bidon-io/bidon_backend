class CreateBidONMigrations < ActiveRecord::Migration[7.0]
  # rubocop:disable Rails/SquishedSQLHeredocs
  def up
    safety_assured do
      execute <<~SQL
        CREATE TABLE bidon_migrations (
          id serial PRIMARY KEY,
          version_id bigint NOT NULL,
          is_applied boolean NOT NULL,
          tstamp timestamp DEFAULT now()
        );
        INSERT INTO bidon_migrations (version_id, is_applied) VALUES (20231120143101, true);
      SQL
    end
  end

  def down
    safety_assured do
      execute <<~SQL
        DROP TABLE bidon_migrations;
      SQL
    end
  end
  # rubocop:enable Rails/SquishedSQLHeredocs
end
