/*
  Warnings:

  - You are about to drop the column `businessId` on the `users` table. All the data in the column will be lost.
  - You are about to drop the column `outletId` on the `users` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "users" DROP COLUMN "businessId",
DROP COLUMN "outletId";
