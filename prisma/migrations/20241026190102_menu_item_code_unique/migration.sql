/*
  Warnings:

  - A unique constraint covering the columns `[code]` on the table `menu_items` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "menu_items_code_key" ON "menu_items"("code");
