CREATE TABLE IF NOT EXISTS metrics(
	id varchar(255) NOT NULL,
	mtype varchar(15) NOT NULL,
	delta bigint,
	value double precision,
	hash varchar(250),
	PRIMARY KEY (mtype, id)
)
