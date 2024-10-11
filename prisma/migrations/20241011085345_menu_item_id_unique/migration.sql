/*
  Warnings:

  - A unique constraint covering the columns `[menu_item_id]` on the table `outlet_menu_items` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "outlet_menu_items_menu_item_id_key" ON "outlet_menu_items"("menu_item_id");
