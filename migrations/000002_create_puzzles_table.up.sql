create table puzzle_collections(
  id serial primary key,
  slug varchar(20) not null,
  created_at timestamp not null default CURRENT_TIMESTAMP
);

create table puzzles(
  id serial primary key,
  fen varchar(90) not null,
  puzzle_collection_id int not null references puzzle_collections(id),
  created_at timestamp not null default CURRENT_TIMESTAMP
);

create table puzzle_solutions(
  id serial primary key,
  solution_internal varchar(4) not null,
  solution_pretty varchar(20) not null, 
  puzzle_id int not null references puzzles(id),
  created_at timestamp not null default CURRENT_TIMESTAMP,
  unique(puzzle_id, solution_internal)
);

--create table puzzle_guesses (
--  id serial primary key,
--  guess varchar(20) not null,
--  correct boolean not null,
--  puzzle_id int not null references puzzles(id),
--  player_id int not null references players(id),
--  created_at timestamp not null default CURRENT_TIMESTAMP
--  unique(puzzle_id, puzzle_session_id, guess)
--);
