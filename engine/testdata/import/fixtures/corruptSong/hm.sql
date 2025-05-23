PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE Information ( 
	id INTEGER PRIMARY KEY AUTOINCREMENT, 
	uuid TEXT, 
	schemaVersionMajor INTEGER, 
	schemaVersionMinor INTEGER, 
	schemaVersionPatch INTEGER, 
	currentPlayedIndiciator INTEGER, 
	lastRekordBoxLibraryImportReadCounter INTEGER
);
INSERT INTO Information VALUES(1,'d1650737-af46-43b8-a4ac-3642fe419e4b',3,0,1,1685741846,NULL);
CREATE TABLE AlbumArt ( 
	id INTEGER PRIMARY KEY AUTOINCREMENT, 
	hash TEXT, 
	albumArt BLOB 
);
CREATE TABLE Track ( 
	id INTEGER PRIMARY KEY AUTOINCREMENT, 
	playOrder INTEGER, 
	length INTEGER, 
	bpm INTEGER, 
	year INTEGER, 
	path TEXT, 
	filename TEXT, 
	bitrate INTEGER, 
	bpmAnalyzed REAL, 
	albumArtId INTEGER, 
	fileBytes INTEGER, 
	title TEXT, 
	artist TEXT, 
	album TEXT, 
	genre TEXT, 
	comment TEXT, 
	label TEXT, 
	composer TEXT, 
	remixer TEXT, 
	key INTEGER, 
	rating INTEGER, 
	albumArt TEXT, 
	timeLastPlayed DATETIME, 
	isPlayed BOOLEAN, 
	fileType TEXT, 
	isAnalyzed BOOLEAN, 
	dateCreated DATETIME, 
	dateAdded DATETIME, 
	isAvailable BOOLEAN, 
	isMetadataOfPackedTrackChanged BOOLEAN, 
	isPerfomanceDataOfPackedTrackChanged BOOLEAN, 
	playedIndicator INTEGER, 
	isMetadataImported BOOLEAN, 
	pdbImportKey INTEGER, 
	streamingSource TEXT, 
	uri TEXT, 
	isBeatGridLocked BOOLEAN, 
	originDatabaseUuid TEXT, 
	originTrackId INTEGER, 
	streamingFlags INTEGER, 
	explicitLyrics BOOLEAN, 
	lastEditTime DATETIME, 
	CONSTRAINT C_originDatabaseUuid_originTrackId UNIQUE (originDatabaseUuid, originTrackId), 
	CONSTRAINT C_path UNIQUE (path), 
	FOREIGN KEY (albumArtId) REFERENCES AlbumArt (id) ON DELETE RESTRICT 
);
CREATE TABLE PerformanceData ( 
	trackId INTEGER PRIMARY KEY, 
	trackData BLOB, 
	overviewWaveFormData BLOB, 
	beatData BLOB, 
	quickCues BLOB, 
	loops BLOB, 
	thirdPartySourceId INTEGER, 
	activeOnLoadLoops INTEGER, 
	FOREIGN KEY(trackId) REFERENCES Track(id) ON DELETE CASCADE ON UPDATE CASCADE 
);
CREATE TABLE Playlist ( 
	id INTEGER PRIMARY KEY AUTOINCREMENT, 
	title TEXT, 
	parentListId INTEGER, 
	isPersisted BOOLEAN, 
	nextListId INTEGER, 
	lastEditTime DATETIME, 
	isExplicitlyExported BOOLEAN, 
	CONSTRAINT C_NAME_UNIQUE_FOR_PARENT UNIQUE (title, parentListId), 
	CONSTRAINT C_NEXT_LIST_ID_UNIQUE_FOR_PARENT UNIQUE (parentListId, nextListId) 
);
CREATE TABLE Historylist ( 
	id INTEGER PRIMARY KEY AUTOINCREMENT, 
	sessionId TEXT, 
	title TEXT, 
	startTime DATETIME, 
	timezone TEXT, 
	originDriveName TEXT, 
	originDatabaseUuid TEXT, 
	originListId INTEGER, 
	isDeleted BOOLEAN, 
	editTime DATETIME, 
	CONSTRAINT C_UNIQUE_ORIGIN_UUID_AND_LIST_ID UNIQUE (originDatabaseUuid, originListId) 
);
CREATE TABLE PlaylistEntity ( 
	id INTEGER PRIMARY KEY AUTOINCREMENT, 
	listId INTEGER, 
	trackId INTEGER, 
	databaseUuid TEXT, 
	nextEntityId INTEGER, 
	membershipReference INTEGER, 
	CONSTRAINT C_NAME_UNIQUE_FOR_LIST UNIQUE (listId, databaseUuid, trackId), 
	FOREIGN KEY (listId) REFERENCES Playlist (id) ON DELETE CASCADE 
);
CREATE TABLE HistorylistEntity ( 
	id INTEGER PRIMARY KEY AUTOINCREMENT, 
	listId INTEGER, 
	trackId INTEGER, 
	startTime DATETIME, 
	FOREIGN KEY (listId) REFERENCES Historylist (id) ON DELETE CASCADE, 
	FOREIGN KEY (trackId) REFERENCES Track (id) ON DELETE CASCADE 
);
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('Information',1);
CREATE INDEX index_AlbumArt_hash ON AlbumArt (hash);
CREATE INDEX index_Track_filename ON Track (filename);
CREATE INDEX index_Track_albumArtId ON Track (albumArtId);
CREATE INDEX index_Track_uri ON Track (uri);
CREATE INDEX index_Track_title ON Track(title);
CREATE INDEX index_Track_length ON Track(length);
CREATE INDEX index_Track_rating ON Track(rating);
CREATE INDEX index_Track_year ON Track(year);
CREATE INDEX index_Track_dateAdded ON Track(dateAdded);
CREATE INDEX index_Track_genre ON Track(genre);
CREATE INDEX index_Track_artist ON Track(artist);
CREATE INDEX index_Track_album ON Track(album);
CREATE INDEX index_Track_key ON Track(key);
CREATE INDEX index_Track_bpmAnalyzed ON Track(CAST(bpmAnalyzed + 0.5 AS int));
CREATE TRIGGER trigger_after_insert_Track_check_id 
AFTER INSERT ON Track 
	WHEN NEW.id <= (SELECT seq FROM sqlite_sequence WHERE name = 'Track') 
BEGIN 
	SELECT RAISE(ABORT, 'Recycling deleted track id''s are not allowed'); 
END;
CREATE TRIGGER trigger_after_update_Track_check_Id 
BEFORE UPDATE ON Track 
	WHEN NEW.id <> OLD.id 
BEGIN 
	SELECT RAISE(ABORT, 'Changing track id''s are not allowed'); 
END;
CREATE TRIGGER trigger_after_insert_Track_fix_origin 
AFTER INSERT ON Track 
	WHEN IFNULL(NEW.originTrackId, 0) = 0 
	OR IFNULL(NEW.originDatabaseUuid, '') = '' 
BEGIN 
	UPDATE Track SET 
		originTrackId = NEW.id, 
		originDatabaseUuid = (SELECT uuid FROM Information) 
	WHERE track.id = NEW.id; 
END;
CREATE TRIGGER trigger_after_update_Track_fix_origin 
AFTER UPDATE ON Track 
	WHEN IFNULL(NEW.originTrackId, 0) = 0 
	OR IFNULL(NEW.originDatabaseUuid, '') = '' 
BEGIN 
	UPDATE Track SET 
		originTrackId = NEW.id, 
		originDatabaseUuid = (SELECT uuid FROM Information) 
	WHERE track.id = NEW.id; 
END;
CREATE TRIGGER trigger_after_update_only_Track_timestamp 
	AFTER UPDATE OF	length, bpm, year, filename, bitrate, bpmAnalyzed, albumArtId, 
	title, artist, album, genre, comment, label, composer, remixer, key, rating, albumArt, 
	fileType, isAnalyzed, isBeatgridLocked, explicitLyrics 
	ON Track 
	FOR EACH ROW 
BEGIN 
	UPDATE Track SET lastEditTime = strftime('%s') WHERE ROWID=NEW.ROWID; 
END;
CREATE TRIGGER trigger_after_insert_Track_insert_performance_data 
AFTER INSERT ON Track 
BEGIN 
	INSERT INTO PerformanceData(trackId) VALUES(NEW.id); 
END;
CREATE TRIGGER trigger_PerformanceData_after_update_Track_timestamp 
	AFTER UPDATE OF trackData, isAnalyzed, overviewWaveFormData, beatData, quickCues, loops, activeOnLoadLoops 
	ON PerformanceData 
	FOR EACH ROW 
BEGIN 
	UPDATE Track 
	SET lastEditTime = strftime('%s') 
	WHERE id = NEW.trackId; 
END;
CREATE TRIGGER trigger_before_insert_List 
BEFORE INSERT ON Playlist 
FOR EACH ROW BEGIN 
	UPDATE Playlist SET 
		nextListId = -(1 + nextListId) 
	WHERE nextListId = NEW.nextListId 
	AND parentListId = NEW.parentListId; 
END;
CREATE TRIGGER trigger_after_insert_List 
AFTER INSERT ON Playlist 
FOR EACH ROW BEGIN 
	UPDATE Playlist SET 
		nextListId = NEW.id 
	WHERE nextListId = -(1 + NEW.nextListId) 
	AND parentListId = NEW.parentListId; 
END;
CREATE TRIGGER trigger_after_delete_List 
AFTER DELETE ON Playlist 
FOR EACH ROW BEGIN 
	UPDATE Playlist SET 
		nextListId = OLD.nextListId 
	WHERE nextListId = OLD.id; 
	DELETE FROM Playlist 
	WHERE parentListId = OLD.id; 
END;
CREATE TRIGGER trigger_after_update_isPersistParent 
AFTER UPDATE ON Playlist 
	WHEN (old.isPersisted = 0 
	AND new.isPersisted = 1) 
	OR (old.parentListId != new.parentListId 
	AND new.isPersisted = 1) 
BEGIN 
	UPDATE Playlist SET 
		isPersisted = 1 
	WHERE id IN (SELECT parentListId FROM PlaylistAllParent WHERE id=new.id); 
END;
CREATE TRIGGER trigger_after_update_isPersistChild 
AFTER UPDATE ON Playlist 
	WHEN old.isPersisted = 1 
	AND new.isPersisted = 0 
BEGIN 
	UPDATE Playlist SET 
		isPersisted = 0 
	WHERE id IN (SELECT childListId FROM PlaylistAllChildren WHERE id=new.id); 
END;
CREATE TRIGGER trigger_after_insert_isPersist 
AFTER INSERT ON Playlist 
	WHEN new.isPersisted = 1 
BEGIN 
	UPDATE Playlist SET 
		isPersisted = 1 
	WHERE id IN (SELECT parentListId FROM PlaylistAllParent WHERE id=new.id); 
END;
CREATE VIEW PlaylistAllParent AS 
WITH FindAllParent AS ( 
	SELECT id, parentListId FROM Playlist 
	UNION ALL 
	SELECT recursiveCTE.id, Plist.parentListId FROM Playlist Plist 
	INNER JOIN FindAllParent recursiveCTE 
	ON recursiveCTE.parentListId = Plist.id 
) 
SELECT * FROM FindAllParent;
CREATE VIEW PlaylistAllChildren AS 
WITH FindAllChild AS ( 
SELECT id, id as childListId FROM Playlist 
UNION ALL 
SELECT recursiveCTE.id, Plist.id FROM Playlist Plist 
INNER JOIN FindAllChild recursiveCTE 
ON recursiveCTE.childListId = Plist.parentListId 
) 
SELECT * FROM FindAllChild WHERE id <> childListId;
CREATE VIEW PlaylistPath AS 
WITH RECURSIVE Heirarchy AS 
( 
	SELECT id AS child, parentListId AS parent, title AS name, 1 AS depth FROM Playlist 
	UNION ALL 
	SELECT child, parentListId AS parent, title AS name, h.depth + 1 AS depth FROM Playlist c 
	JOIN Heirarchy h ON h.parent = c.id 
	ORDER BY depth DESC 
), 
OrderedList AS 
( 
	SELECT id , nextListId, 1 AS position 
	FROM Playlist 
	WHERE nextListId = 0 
	UNION ALL 
	SELECT c.id , c.nextListId , l.position + 1 
	FROM Playlist c 
	INNER JOIN OrderedList l 
	ON c.nextListId = l.id 
), 
NameConcat AS 
( 
	SELECT 
		child AS id, 
		GROUP_CONCAT(name ,';') || ';' AS path 
	FROM 
	( 
		SELECT child, name 
		FROM Heirarchy 
		ORDER BY depth DESC 
	) 
	GROUP BY child 
) 
SELECT 
	id, 
	path, 
	ROW_NUMBER() OVER 
	( 
		ORDER BY 
		(SELECT COUNT(*) FROM (SELECT * FROM Heirarchy WHERE child = id) ) DESC, 
		(SELECT position FROM OrderedList ol WHERE ol.id = c.id) ASC 
	) AS position 
FROM Playlist c 
LEFT JOIN NameConcat g USING (id);
CREATE TRIGGER trigger_after_update_Historylist 
AFTER UPDATE ON Historylist 
	WHEN COALESCE(NEW.title != OLD.title, OLD.title IS NULL AND NEW.title IS NOT NULL) 
BEGIN 
	UPDATE Historylist SET 
		editTime = strftime('%s','now') 
	WHERE id = NEW.id; 
END;
CREATE INDEX index_PlaylistEntity_nextEntityId_listId ON PlaylistEntity(nextEntityId, listId);
CREATE TRIGGER trigger_before_delete_PlaylistEntity 
BEFORE DELETE ON PlaylistEntity 
WHEN OLD.trackId > 0 
BEGIN 
	UPDATE PlaylistEntity SET 
		nextEntityId = OLD.nextEntityId 
	WHERE nextEntityId = OLD.id 
	AND listId = OLD.listId; 
END;
CREATE INDEX index_HistorylistEntity_listId ON HistorylistEntity (listId);
CREATE INDEX index_HistorylistEntity_trackId ON HistorylistEntity (trackId);
COMMIT;
