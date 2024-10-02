/*
  Warnings:

  - Changed the type of `outlet_id` on the `user_outlets` table. No cast exists, the column would be dropped and recreated, which cannot be done if there is data, since the column is required.

*/
-- AlterTable
ALTER TABLE "user_outlets" DROP COLUMN "outlet_id",
ADD COLUMN     "outlet_id" INTEGER NOT NULL;
