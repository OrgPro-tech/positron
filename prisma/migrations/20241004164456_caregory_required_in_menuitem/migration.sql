/*
  Warnings:

  - Made the column `category_id` on table `menu_items` required. This step will fail if there are existing NULL values in that column.

*/
-- DropForeignKey
ALTER TABLE "menu_items" DROP CONSTRAINT "menu_items_category_id_fkey";

-- AlterTable
ALTER TABLE "menu_items" ALTER COLUMN "category_id" SET NOT NULL;

-- AddForeignKey
ALTER TABLE "menu_items" ADD CONSTRAINT "menu_items_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES "categories"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
