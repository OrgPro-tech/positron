/*
  Warnings:

  - A unique constraint covering the columns `[mobile_number]` on the table `users` will be added. If there are existing duplicate values, this will fail.

*/
-- AlterTable
ALTER TABLE "users" ALTER COLUMN "mobile_number" SET DATA TYPE TEXT;

-- CreateIndex
CREATE UNIQUE INDEX "users_mobile_number_key" ON "users"("mobile_number");
