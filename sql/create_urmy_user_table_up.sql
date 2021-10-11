
CREATE TABLE IF NOT EXISTS urmyusers (
			loginuuid char(40) PRIMARY KEY NOT NULL,
			notificationtoken  char(100),
			lastlogin DATE
);

CREATE TABLE IF NOT EXISTS urmyuserinfo (
		loginuuid char(40) PRIMARY KEY NOT NULL,
		password  char(20),
		email   char(30),
		phoneno  char(20),
		phonecode char(6),
		name char(100),
		hometown char(30),
		country char(5),
		birthday TIMESTAMPTZ,
		gender BOOLEAN
);

CREATE TABLE IF NOT EXISTS urmyusersagreement (
			loginuuid char(40) PRIMARY KEY NOT NULL,
			isoverage   BOOLEAN,
			urmyaccount  BOOLEAN,
			urmyoverallservice BOOLEAN,
			urmynotiad  BOOLEAN,
			urmypersonaldataacc BOOLEAN,
			urmylocation BOOLEAN,
			urmyprofileadditional BOOLEAN,
			createdAt DATE
);

SET TIMEZONE='ASIA/SEOUL';

CREATE USER koreaogh WITH SUPERUSER CREATEROLE CREATEDB REPLICATION LOGIN INHERIT BYPASSRLS PASSWORD 'ogh1898';
CREATE DATABASE urmydb with owner koreaogh encoding 'UTF8';
INSERT INTO urmyusers (loginId, password, nickname, name, phoneNo, gender, birthday, createdat) VALUES ('tgja1075', 'qwer1234', 'whitewhale', 'johnoh', '01095801075', true, '910807', current_timestamp)