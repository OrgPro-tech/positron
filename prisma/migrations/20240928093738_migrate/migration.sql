/*
  Warnings:

  - You are about to drop the column `businessId` on the `outlets` table. All the data in the column will be lost.
  - You are about to drop the column `userId` on the `outlets` table. All the data in the column will be lost.
  - You are about to drop the column `businessId` on the `users` table. All the data in the column will be lost.
  - You are about to drop the column `outletId` on the `users` table. All the data in the column will be lost.

*/
-- DropForeignKey
ALTER TABLE "outlets" DROP CONSTRAINT "outlets_businessId_fkey";

-- DropForeignKey
ALTER TABLE "users" DROP CONSTRAINT "users_businessId_fkey";

-- AlterTable
ALTER TABLE "outlets" DROP COLUMN "businessId",
DROP COLUMN "userId",
ALTER COLUMN "business_id" DROP NOT NULL;

-- AlterTable
ALTER TABLE "users" DROP COLUMN "businessId",
DROP COLUMN "outletId",
ADD COLUMN     "business_id" TEXT,
ADD COLUMN     "outlet_id" TEXT;

-- AddForeignKey
ALTER TABLE "users" ADD CONSTRAINT "users_business_id_fkey" FOREIGN KEY ("business_id") REFERENCES "businesses"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "users" ADD CONSTRAINT "users_outlet_id_fkey" FOREIGN KEY ("outlet_id") REFERENCES "outlets"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "outlets" ADD CONSTRAINT "outlets_business_id_fkey" FOREIGN KEY ("business_id") REFERENCES "businesses"("id") ON DELETE SET NULL ON UPDATE CASCADE;
