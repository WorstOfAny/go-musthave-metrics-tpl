CREATE TABLE IF NOT EXISTS metrics(
	id text PRIMARY KEY,
	mtype varchar(15) NOT NULL,
	delta int,
	value double precision,
	hash varchar(250)
)
