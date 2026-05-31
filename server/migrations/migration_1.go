package migrations

import (
	"database/sql"
	"fmt"

	"zyrouge.me/umi/application"
)

var MigrationV1 = Migration{
	Version: 1,
	Up: func(transaction *sql.Tx) error {
		config, err := application.GetConfig()
		if err != nil {
			return err
		}
		idType := "TEXT"
		if config.Database.Driver == application.UmiDatabaseConfigDriverPostgres {
			idType = "UUID"
		}
		statements := []string{
			// umi_user
			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS umi_user (
				id            %s NOT NULL,
				username      TEXT NOT NULL,
				email         TEXT,
				display_name  TEXT NOT NULL,
				password_hash TEXT NOT NULL,
				created_at    BIGINT NOT NULL,
				updated_at    BIGINT NOT NULL,
				CONSTRAINT pk_umi_user PRIMARY KEY (id),
				CONSTRAINT uq_umi_user_username UNIQUE (username),
				CONSTRAINT uq_umi_user_email UNIQUE (email)
			)`, idType),
			`CREATE INDEX IF NOT EXISTS idx_umi_user_display_name ON umi_user (display_name)`,
			// umi_refresh_token
			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS umi_refresh_token (
				id         %s NOT NULL,
				user_id    %s NOT NULL,
				token_hash TEXT NOT NULL,
				expires_at BIGINT NOT NULL,
				created_at BIGINT NOT NULL,
				CONSTRAINT pk_umi_refresh_token PRIMARY KEY (id),
				CONSTRAINT fk_umi_refresh_token_user_id FOREIGN KEY (user_id) REFERENCES umi_user(id) ON DELETE CASCADE,
				CONSTRAINT uq_umi_refresh_token_token_hash UNIQUE (token_hash)
			)`, idType, idType),
			`CREATE INDEX IF NOT EXISTS idx_umi_refresh_token_user_id ON umi_refresh_token (user_id)`,
			// umi_team
			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS umi_team (
				id             %s NOT NULL,
				name           TEXT NOT NULL,
				encryption_key TEXT NOT NULL,
				created_at     BIGINT NOT NULL,
				updated_at     BIGINT NOT NULL,
				CONSTRAINT pk_umi_team PRIMARY KEY (id)
			)`, idType),
			`CREATE INDEX IF NOT EXISTS idx_umi_team_name ON umi_team (name)`,
			// umi_member
			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS umi_member (
				user_id    %s NOT NULL,
				team_id    %s NOT NULL,
				role       TEXT NOT NULL,
				created_at BIGINT NOT NULL,
				updated_at BIGINT NOT NULL,
				CONSTRAINT pk_umi_member PRIMARY KEY (user_id, team_id),
				CONSTRAINT fk_umi_member_user_id FOREIGN KEY (user_id) REFERENCES umi_user(id) ON DELETE CASCADE,
				CONSTRAINT fk_umi_member_team_id FOREIGN KEY (team_id) REFERENCES umi_team(id) ON DELETE CASCADE
			)`, idType, idType),
			`CREATE INDEX IF NOT EXISTS idx_umi_member_team_id ON umi_member (team_id)`,
			// umi_service
			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS umi_service (
				id         %s NOT NULL,
				team_id    %s NOT NULL,
				name       TEXT NOT NULL,
				token_hash TEXT NOT NULL,
				created_at BIGINT NOT NULL,
				updated_at BIGINT NOT NULL,
				CONSTRAINT pk_umi_service PRIMARY KEY (id),
				CONSTRAINT fk_umi_service_team_id FOREIGN KEY (team_id) REFERENCES umi_team(id) ON DELETE CASCADE
			)`, idType, idType),
			`CREATE INDEX IF NOT EXISTS idx_umi_service_team_id ON umi_service (team_id)`,
			// umi_channel
			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS umi_channel (
				id         			%s NOT NULL,
				team_id    			%s NOT NULL,
				name       			TEXT NOT NULL,
				max_retention_days 	INTEGER NOT NULL,
				created_at BIGINT 	NOT NULL,
				updated_at BIGINT 	NOT NULL,
				CONSTRAINT pk_umi_channel PRIMARY KEY (id),
				CONSTRAINT fk_umi_channel_team_id FOREIGN KEY (team_id) REFERENCES umi_team(id) ON DELETE CASCADE
			)`, idType, idType),
			`CREATE INDEX IF NOT EXISTS idx_umi_channel_team_id ON umi_channel (team_id)`,
			// umi_event
			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS umi_event (
				id         %s NOT NULL,
				service_id %s NOT NULL,
				channel_id %s NOT NULL,
				title      TEXT NOT NULL,
				body       TEXT,
				level      TEXT,
				action_url TEXT,
				icon_url   TEXT,
				metadata   TEXT,
				created_at BIGINT NOT NULL,
				CONSTRAINT pk_umi_event PRIMARY KEY (id),
				CONSTRAINT fk_umi_event_service_id FOREIGN KEY (service_id) REFERENCES umi_service(id) ON DELETE CASCADE,
				CONSTRAINT fk_umi_event_channel_id FOREIGN KEY (channel_id) REFERENCES umi_channel(id) ON DELETE CASCADE
			)`, idType, idType, idType),
			`CREATE INDEX IF NOT EXISTS idx_umi_event_created_at ON umi_event (created_at)`,
			`CREATE INDEX IF NOT EXISTS idx_umi_event_service_id_created_at ON umi_event (service_id, created_at)`,
			`CREATE INDEX IF NOT EXISTS idx_umi_event_channel_id_created_at ON umi_event (channel_id, created_at)`,
			// umi_tag
			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS umi_tag (
				id         %s NOT NULL,
				team_id    %s NOT NULL,
				name       TEXT NOT NULL,
				created_at BIGINT NOT NULL,
				updated_at BIGINT NOT NULL,
				CONSTRAINT pk_umi_tag PRIMARY KEY (id),
				CONSTRAINT fk_umi_tag_team_id FOREIGN KEY (team_id) REFERENCES umi_team(id) ON DELETE CASCADE,
				CONSTRAINT uq_umi_tag_team_name UNIQUE (team_id, name)
			)`, idType, idType),
			`CREATE INDEX IF NOT EXISTS idx_umi_tag_team_id ON umi_tag (team_id)`,
			// umi_event_tag_map
			fmt.Sprintf(`CREATE TABLE IF NOT EXISTS umi_event_tag_map (
				event_id    %s NOT NULL,
				tag_id      %s NOT NULL,
				created_at  BIGINT NOT NULL,
				CONSTRAINT pk_umi_event_tag_map PRIMARY KEY (event_id, tag_id),
				CONSTRAINT fk_umi_event_tag_map_event_id FOREIGN KEY (event_id) REFERENCES umi_event(id) ON DELETE CASCADE,
				CONSTRAINT fk_umi_event_tag_map_tag_id FOREIGN KEY (tag_id) REFERENCES umi_tag(id) ON DELETE CASCADE
			)`, idType, idType),
			`CREATE INDEX IF NOT EXISTS idx_umi_event_tag_map_tag_id ON umi_event_tag_map (tag_id)`,
			`CREATE INDEX IF NOT EXISTS idx_umi_event_tag_map_event_id ON umi_event_tag_map (event_id)`,
		}
		for _, s := range statements {
			if _, err := transaction.Exec(s); err != nil {
				return err
			}
		}
		return nil
	},
	Down: func(transaction *sql.Tx) error {
		statements := []string{
			`DROP TABLE IF EXISTS umi_event_tag_map`,
			`DROP TABLE IF EXISTS umi_tag`,
			`DROP TABLE IF EXISTS umi_event`,
			`DROP TABLE IF EXISTS umi_channel`,
			`DROP TABLE IF EXISTS umi_service`,
			`DROP TABLE IF EXISTS umi_member`,
			`DROP TABLE IF EXISTS umi_team`,
			`DROP TABLE IF EXISTS umi_user`,
		}
		for _, s := range statements {
			if _, err := transaction.Exec(s); err != nil {
				return err
			}
		}
		return nil
	},
}
