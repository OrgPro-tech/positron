/*
  Warnings:

  - You are about to drop the column `WhatsApp` on the `customer` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "customer" DROP COLUMN "WhatsApp",
ADD COLUMN     "whatsapp" BOOLEAN DEFAULT false;
