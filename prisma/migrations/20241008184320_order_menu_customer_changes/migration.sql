/*
  Warnings:

  - The values [Extra Hot] on the enum `SpiceLevel` will be removed. If these variants are still used in the database, this will fail.
  - You are about to drop the `OutletMenuItem` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `customer` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `ingredients` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `menu_item_ingredients` table. If the table is not empty, all the data it contains will be lost.
  - Added the required column `code` to the `menu_items` table without a default value. This is not possible if the table is not empty.
  - Added the required column `customizable` to the `menu_items` table without a default value. This is not possible if the table is not empty.
  - Added the required column `size_type` to the `menu_items` table without a default value. This is not possible if the table is not empty.
  - Added the required column `tax_percentage` to the `menu_items` table without a default value. This is not possible if the table is not empty.

*/
-- CreateEnum
CREATE TYPE "SizeType" AS ENUM ('GRAM', 'PIECE');

-- CreateEnum
CREATE TYPE "OrderStatus" AS ENUM ('NEW', 'PREPARING', 'READY');

-- AlterEnum
BEGIN;
CREATE TYPE "SpiceLevel_new" AS ENUM ('Mild', 'Medium', 'Hot', 'ExtraHot');
ALTER TABLE "menu_items" ALTER COLUMN "spice_level" TYPE "SpiceLevel_new" USING ("spice_level"::text::"SpiceLevel_new");
ALTER TYPE "SpiceLevel" RENAME TO "SpiceLevel_old";
ALTER TYPE "SpiceLevel_new" RENAME TO "SpiceLevel";
DROP TYPE "SpiceLevel_old";
COMMIT;

-- DropForeignKey
ALTER TABLE "OutletMenuItem" DROP CONSTRAINT "OutletMenuItem_menu_item_id_fkey";

-- DropForeignKey
ALTER TABLE "OutletMenuItem" DROP CONSTRAINT "OutletMenuItem_outlet_id_fkey";

-- DropForeignKey
ALTER TABLE "customer" DROP CONSTRAINT "customer_business_id_fkey";

-- DropForeignKey
ALTER TABLE "customer" DROP CONSTRAINT "customer_outlet_id_fkey";

-- DropForeignKey
ALTER TABLE "menu_item_ingredients" DROP CONSTRAINT "menu_item_ingredients_ingredient_id_fkey";

-- AlterTable
ALTER TABLE "menu_items" ADD COLUMN     "code" TEXT NOT NULL,
ADD COLUMN     "customizable" BOOLEAN NOT NULL,
ADD COLUMN     "image" TEXT,
ADD COLUMN     "size_type" "SizeType" NOT NULL,
ADD COLUMN     "tax_percentage" INTEGER NOT NULL,
ADD COLUMN     "variation" JSONB;

-- DropTable
DROP TABLE "OutletMenuItem";

-- DropTable
DROP TABLE "customer";

-- DropTable
DROP TABLE "ingredients";

-- DropTable
DROP TABLE "menu_item_ingredients";

-- DropEnum
DROP TYPE "Unit";

-- CreateTable
CREATE TABLE "outlet_menu_items" (
    "id" SERIAL NOT NULL,
    "menu_item_id" INTEGER NOT NULL,
    "outlet_id" INTEGER NOT NULL,
    "price" DECIMAL(10,2) NOT NULL,
    "is_available" BOOLEAN NOT NULL DEFAULT true,

    CONSTRAINT "outlet_menu_items_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "customers" (
    "id" SERIAL NOT NULL,
    "phone_number" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "whatsapp" BOOLEAN DEFAULT false,
    "email" TEXT,
    "address" TEXT,
    "outlet_id" INTEGER NOT NULL,
    "business_id" INTEGER NOT NULL,

    CONSTRAINT "customers_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "orders" (
    "id" SERIAL NOT NULL,
    "customer_id" INTEGER NOT NULL,
    "phone_number" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "email" TEXT,
    "address" TEXT,
    "order_id" VARCHAR(21) NOT NULL,
    "status" "OrderStatus" NOT NULL DEFAULT 'NEW',
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,
    "gst_amount" DECIMAL(10,2) NOT NULL,
    "total_amount" DECIMAL(10,2) NOT NULL,
    "net_amount" DECIMAL(10,2) NOT NULL,

    CONSTRAINT "orders_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "order_items" (
    "id" SERIAL NOT NULL,
    "item_code" TEXT NOT NULL,
    "item_description" TEXT NOT NULL,
    "variation" JSONB,
    "quantity" INTEGER NOT NULL,
    "unit_price" DECIMAL(10,2) NOT NULL,
    "net_price" DECIMAL(10,2) NOT NULL,
    "tax_precentage" INTEGER NOT NULL,
    "gst_amount" DECIMAL(10,2) NOT NULL,
    "total_amount" DECIMAL(10,2) NOT NULL,
    "order_id" INTEGER NOT NULL,

    CONSTRAINT "order_items_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "customers_phone_number_key" ON "customers"("phone_number");

-- CreateIndex
CREATE UNIQUE INDEX "orders_phone_number_key" ON "orders"("phone_number");

-- CreateIndex
CREATE UNIQUE INDEX "orders_order_id_key" ON "orders"("order_id");

-- AddForeignKey
ALTER TABLE "outlet_menu_items" ADD CONSTRAINT "outlet_menu_items_menu_item_id_fkey" FOREIGN KEY ("menu_item_id") REFERENCES "menu_items"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "outlet_menu_items" ADD CONSTRAINT "outlet_menu_items_outlet_id_fkey" FOREIGN KEY ("outlet_id") REFERENCES "outlets"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "customers" ADD CONSTRAINT "customers_outlet_id_fkey" FOREIGN KEY ("outlet_id") REFERENCES "outlets"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "customers" ADD CONSTRAINT "customers_business_id_fkey" FOREIGN KEY ("business_id") REFERENCES "businesses"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "orders" ADD CONSTRAINT "orders_customer_id_fkey" FOREIGN KEY ("customer_id") REFERENCES "customers"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "order_items" ADD CONSTRAINT "order_items_order_id_fkey" FOREIGN KEY ("order_id") REFERENCES "orders"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
