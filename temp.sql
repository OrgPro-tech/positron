-- Create enum type for UserType
CREATE TYPE "UserType" AS ENUM ('ADMIN', 'USER');

-- Create users table
CREATE TABLE "users" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "username" VARCHAR(255) UNIQUE NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "mobile_number" INTEGER NOT NULL,
    "user_type" "UserType" NOT NULL DEFAULT 'ADMIN',
    "businessId" UUID,
    "outletId" UUID
);

-- Create user_outlets table
CREATE TABLE "user_outlets" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "business_id" UUID NOT NULL,
    "outlet_id" UUID NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

-- Create businesses table
CREATE TABLE "businesses" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "contact_person_name" VARCHAR(255) NOT NULL,
    "contact_person_email" VARCHAR(255) NOT NULL,
    "contact_person_mobile_number" INTEGER NOT NULL,
    "company_name" VARCHAR(255) NOT NULL,
    "address" TEXT NOT NULL,
    "pin" INTEGER NOT NULL,
    "city" VARCHAR(255) NOT NULL,
    "state" VARCHAR(255) NOT NULL,
    "country" VARCHAR(255) NOT NULL,
    "business_type" VARCHAR(255) NOT NULL,
    "gst" VARCHAR(255),
    "pan" VARCHAR(255),
    "bank_account_number" VARCHAR(255),
    "bank_name" VARCHAR(255),
    "ifsc_code" VARCHAR(255),
    "account_type" VARCHAR(255),
    "account_holder_name" VARCHAR(255)
);

-- Create outlets table
CREATE TABLE "outlets" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "outlet_name" VARCHAR(255) NOT NULL,
    "outlet_address" TEXT NOT NULL,
    "outlet_pin" INTEGER NOT NULL,
    "outlet_city" VARCHAR(255) NOT NULL,
    "outlet_state" VARCHAR(255) NOT NULL,
    "outlet_country" VARCHAR(255) NOT NULL,
    "business_id" UUID,
    "businessId" UUID,
    "userId" UUID,
    FOREIGN KEY ("businessId") REFERENCES "businesses" ("id")
);

-- Create user_sessions table
CREATE TABLE "user_sessions" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id" UUID UNIQUE NOT NULL,
    "access_token" TEXT NOT NULL,
    "refresh_token" TEXT NOT NULL,
    "expire_at" TIMESTAMP NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

-- Add foreign key constraints
ALTER TABLE "users" ADD CONSTRAINT "users_businessId_fkey" FOREIGN KEY ("businessId") REFERENCES "businesses" ("id");




----------------------------------------
-- Create enum type for UserType
CREATE TYPE "UserType" AS ENUM ('ADMIN', 'USER');

-- Create users table
CREATE TABLE "users" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "username" VARCHAR(255) UNIQUE NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "mobile_number" TEXT NOT NULL,
    "user_type" "UserType" NOT NULL DEFAULT 'ADMIN',
    "business_id" UUID,
    "outlet_id" UUID
);

-- Create user_outlets table
CREATE TABLE "user_outlets" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL,
    "business_id" UUID NOT NULL,
    "outlet_id" UUID NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

-- Create businesses table
CREATE TABLE "businesses" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "contact_person_name" VARCHAR(255) NOT NULL,
    "contact_person_email" VARCHAR(255) NOT NULL,
    "contact_person_mobile_number" INTEGER NOT NULL,
    "company_name" VARCHAR(255) NOT NULL,
    "address" TEXT NOT NULL,
    "pin" INTEGER NOT NULL,
    "city" VARCHAR(255) NOT NULL,
    "state" VARCHAR(255) NOT NULL,
    "country" VARCHAR(255) NOT NULL,
    "business_type" VARCHAR(255) NOT NULL,
    "gst" VARCHAR(255),
    "pan" VARCHAR(255),
    "bank_account_number" VARCHAR(255),
    "bank_name" VARCHAR(255),
    "ifsc_code" VARCHAR(255),
    "account_type" VARCHAR(255),
    "account_holder_name" VARCHAR(255)
);

-- Create outlets table
CREATE TABLE "outlets" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "outlet_name" VARCHAR(255) NOT NULL,
    "outlet_address" TEXT NOT NULL,
    "outlet_pin" INTEGER NOT NULL,
    "outlet_city" VARCHAR(255) NOT NULL,
    "outlet_state" VARCHAR(255) NOT NULL,
    "outlet_country" VARCHAR(255) NOT NULL,
    "business_id" UUID,
    FOREIGN KEY ("business_id") REFERENCES "businesses" ("id")
);

-- Create user_sessions table
CREATE TABLE "user_sessions" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id" UUID UNIQUE NOT NULL,
    "access_token" TEXT NOT NULL,
    "refresh_token" TEXT NOT NULL,
    "expire_at" TIMESTAMP NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

-- Add foreign key constraints
ALTER TABLE "users" 
    ADD CONSTRAINT "users_business_id_fkey" 
    FOREIGN KEY ("business_id") REFERENCES "businesses" ("id")
    ON DELETE SET NULL;

ALTER TABLE "users" 
    ADD CONSTRAINT "users_outlet_id_fkey" 
    FOREIGN KEY ("outlet_id") REFERENCES "outlets" ("id")
    ON DELETE SET NULL;

ALTER TABLE public.businesses ALTER COLUMN contact_person_mobile_number TYPE text USING contact_person_mobile_number::text;
