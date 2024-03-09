CREATE TABLE IF NOT EXISTS ReactionComment (
	Value TEXT ,
	CommentId INTEGER ,
	User_id INTEGER ,
	FOREIGN KEY (CommentId) REFERENCES Comments (Id),
FOREIGN KEY (User_id) REFERENCES User (Id)
)