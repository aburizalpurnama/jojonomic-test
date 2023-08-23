CREATE TABLE "tbl_harga" (
  "id" varchar(32) PRIMARY KEY,
  "topup" bigint,
  "buyback" bigint,
  "date" int,
  "admin_id" varchar(32)
);

CREATE TABLE "tbl_rekening" (
  "id" varchar(32) PRIMARY KEY,
  "norek" varchar(10),
  "saldo" float,
  "updated_at" int
);

CREATE TABLE "tbl_transaksi" (
  "id" varchar(32) PRIMARY KEY,
  "date" int,
  "type" varchar(10),
  "rekening_id" varchar(32),
  "gram" float,
  "harga_topup" bigint,
  "harga_buyback" bigint,
  "saldo" float
);

CREATE INDEX ON "tbl_rekening" ("norek");

CREATE INDEX ON "tbl_transaksi" ("rekening_id");

ALTER TABLE "tbl_transaksi" ADD FOREIGN KEY ("rekening_id") REFERENCES "tbl_rekening" ("id");
