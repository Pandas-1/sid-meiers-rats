CREATE TABLE "CityDetails"(
    "UserID" INTEGER NOT NULL,
    "Resource1" BIGINT NOT NULL,
    "Resource2" BIGINT NOT NULL,
    "MaxResource1" BIGINT NOT NULL,
    "MaxResource2" BIGINT NOT NULL,
    "MaxTroopArmySize" BIGINT NOT NULL,
    "LastUpdated" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "MaxDefenceBuildings" BIGINT NOT NULL,
    "MaxResourceBuildings" BIGINT NOT NULL
);
ALTER TABLE
    "CityDetails" ADD PRIMARY KEY("UserID");
CREATE TABLE "Users"(
    "UserID" INTEGER NOT NULL,
    "Username" TEXT NOT NULL,
    "Trophies" INTEGER NOT NULL,
    "BaseLevel" INTEGER NOT NULL,
    "Xp" INTEGER NOT NULL,
    "LastPlayed" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "PasswordHash" TEXT NOT NULL
);
ALTER TABLE
    "Users" ADD PRIMARY KEY("UserID");
ALTER TABLE
    "Users" ADD CONSTRAINT "users_username_unique" UNIQUE("Username");
CREATE TABLE "TroopDetails"(
    "TroopID" INTEGER NOT NULL,
    "Name" TEXT NOT NULL,
    "BaseCost" INTEGER NOT NULL,
    "TroopAttackPower" INTEGER NOT NULL,
    "BuildingAttackPower" INTEGER NOT NULL,
    "Defence" INTEGER NOT NULL,
    "Range" INTEGER NOT NULL,
    "AttributeStrength" INTEGER NOT NULL,
    "AttributeWeakness" INTEGER NOT NULL,
    "TroopArmySpace" INTEGER NOT NULL,
    "MovementSpeed" INTEGER NOT NULL,
    "MaxLevel" INTEGER[] NOT NULL
);
ALTER TABLE
    "TroopDetails" ADD PRIMARY KEY("TroopID");
ALTER TABLE
    "TroopDetails" ADD CONSTRAINT "troopdetails_name_unique" UNIQUE("Name");
CREATE TABLE "BuildingDetails"(
    "BuildingID" INTEGER NOT NULL,
    "Name" TEXT NOT NULL,
    "Typeint" CHAR(255) NOT NULL,
    "Production" INTEGER NOT NULL,
    "Scaling" INTEGER NOT NULL,
    "HealthBar" INTEGER NOT NULL,
    "Width" INTEGER NOT NULL,
    "Height" INTEGER NOT NULL,
    "DefenceAttack" INTEGER NOT NULL,
    "DefenceRange" INTEGER NOT NULL,
    "MaxLevel" INTEGER[] NOT NULL,
    "CostResource1" INTEGER NOT NULL,
    "CostResource2" INTEGER NOT NULL
);
ALTER TABLE
    "BuildingDetails" ADD PRIMARY KEY("BuildingID");
ALTER TABLE
    "BuildingDetails" ADD CONSTRAINT "buildingdetails_name_unique" UNIQUE("Name");
CREATE TABLE "UserBattleHistory"(
    "UserID" INTEGER NOT NULL,
    "NumberOfBattles" INTEGER NOT NULL,
    "BattlesWon" INTEGER NOT NULL,
    "BattlesLost" INTEGER NOT NULL,
    "Trophies" INTEGER NOT NULL
);
ALTER TABLE
    "UserBattleHistory" ADD PRIMARY KEY("UserID");
CREATE TABLE "Battles"(
    "BattleID" BIGINT NOT NULL,
    "AttackerID" INTEGER NOT NULL,
    "DefenderID" INTEGER NOT NULL,
    "Resource1Won" INTEGER NOT NULL,
    "Resource2Won" INTEGER NOT NULL,
    "VictoryPercentage" INTEGER NOT NULL
);
ALTER TABLE
    "Battles" ADD PRIMARY KEY("BattleID");
CREATE TABLE "ArmyDetails"(
    "UserID" INTEGER NOT NULL,
    "TroopUnitsUsed" INTEGER NOT NULL,
    "ArmyComposition" jsonb NOT NULL,
    "CreatedOn" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "ArmyDetails" ADD PRIMARY KEY("UserID");
CREATE TABLE "UserBuildings"(
    "InstanceID" INTEGER NOT NULL,
    "UserID" INTEGER NOT NULL,
    "BuildingID" INTEGER NOT NULL,
    "Level" BIGINT NOT NULL,
    "GridX" INTEGER NOT NULL,
    "GridY" INTEGER NOT NULL
);
ALTER TABLE
    "UserBuildings" ADD PRIMARY KEY("InstanceID");
CREATE TABLE "UserTroopDetails"(
    "UserID" INTEGER NOT NULL,
    "TroopID" INTEGER NOT NULL,
    "TroopLevel" INTEGER NOT NULL
);
ALTER TABLE
    "UserTroopDetails" ADD PRIMARY KEY("UserID", "TroopID");
ALTER TABLE
    "CityDetails" ADD CONSTRAINT "citydetails_userid_foreign" FOREIGN KEY("UserID") REFERENCES "Users"("UserID");
ALTER TABLE
    "UserTroopDetails" ADD CONSTRAINT "usertroopdetails_troopid_foreign" FOREIGN KEY("TroopID") REFERENCES "TroopDetails"("TroopID");
ALTER TABLE
    "UserBuildings" ADD CONSTRAINT "userbuildings_buildingid_foreign" FOREIGN KEY("BuildingID") REFERENCES "BuildingDetails"("BuildingID");
ALTER TABLE
    "ArmyDetails" ADD CONSTRAINT "armydetails_userid_foreign" FOREIGN KEY("UserID") REFERENCES "Users"("UserID");
ALTER TABLE
    "UserBuildings" ADD CONSTRAINT "userbuildings_userid_foreign" FOREIGN KEY("UserID") REFERENCES "Users"("UserID");
ALTER TABLE
    "UserTroopDetails" ADD CONSTRAINT "usertroopdetails_userid_foreign" FOREIGN KEY("UserID") REFERENCES "Users"("UserID");
ALTER TABLE
    "UserBattleHistory" ADD CONSTRAINT "userbattlehistory_userid_foreign" FOREIGN KEY("UserID") REFERENCES "Users"("UserID");
ALTER TABLE
    "Battles" ADD CONSTRAINT "battles_attackerid_foreign" FOREIGN KEY("AttackerID") REFERENCES "Users"("UserID");
ALTER TABLE
    "Battles" ADD CONSTRAINT "battles_defenderid_foreign" FOREIGN KEY("DefenderID") REFERENCES "Users"("UserID");