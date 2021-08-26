CREATE TABLE pets (
		id bigserial PRIMARY KEY,
		name VARCHAR ( 50 ) UNIQUE NOT NULL,
		tag VARCHAR ( 50 )
);
