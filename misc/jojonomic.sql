CREATE TABLE "tbl_harga" (
  "id" varchar(32) PRIMARY KEY,
  "topup" bigint NOT NULL,
  "buyback" bigint NOT NULL,
  "date" int NOT NULL,
  "admin_id" varchar(32) NOT NULL
);

CREATE TABLE "tbl_rekening" (
  "id" varchar(32) PRIMARY KEY,
  "norek" varchar(10) NOT NULL,
  "saldo" float NOT NULL,
  "updated_at" int
);

CREATE TABLE "tbl_transaksi" (
  "id" varchar(32) PRIMARY KEY,
  "date" int NOT NULL,
  "type" varchar(10) NOT NULL,
  "rekening_id" varchar(32) NOT NULL,
  "gram" float NOT NULL,
  "harga_topup" bigint NOT NULL,
  "harga_buyback" bigint NOT NULL,
  "saldo" float NOT NULL
);

CREATE INDEX ON "tbl_rekening" ("norek");

CREATE INDEX ON "tbl_transaksi" ("rekening_id");

ALTER TABLE "tbl_transaksi" ADD FOREIGN KEY ("rekening_id") REFERENCES "tbl_rekening" ("id");
