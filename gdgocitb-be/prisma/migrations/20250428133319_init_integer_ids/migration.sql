-- CreateEnum
CREATE TYPE "UserAccess" AS ENUM ('Member', 'Buddy', 'CurriculumAdmin', 'ProfessionalDevelopmentAdmin', 'TechnicalAdmin');

-- CreateEnum
CREATE TYPE "SubmissionType" AS ENUM ('File', 'Link', 'Image', 'GOCI');

-- CreateEnum
CREATE TYPE "AccessLevel" AS ENUM ('Member', 'Buddy');

-- CreateEnum
CREATE TYPE "CertificateType" AS ENUM ('Member', 'Buddy');

-- CreateEnum
CREATE TYPE "EventType" AS ENUM ('Professional', 'Community');

-- CreateEnum
CREATE TYPE "EventCategory" AS ENUM ('GoogleDeveloperGroups', 'Scholarship', 'Internship', 'Exchange');

-- CreateEnum
CREATE TYPE "EventAccessLevel" AS ENUM ('Member', 'MemberAndBuddy');

-- CreateTable
CREATE TABLE "users" (
    "id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "access" "UserAccess" NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "paths" (
    "id" SERIAL NOT NULL,
    "pathName" TEXT NOT NULL,
    "description" TEXT,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "paths_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "modules" (
    "id" SERIAL NOT NULL,
    "pathId" INTEGER NOT NULL,
    "moduleName" TEXT NOT NULL,
    "description" TEXT,
    "video" TEXT,
    "order" INTEGER NOT NULL DEFAULT 0,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "modules_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "lectures" (
    "id" SERIAL NOT NULL,
    "moduleId" INTEGER NOT NULL,
    "title" TEXT NOT NULL,
    "notes" TEXT,
    "slides" TEXT,
    "sourceCode" TEXT,
    "order" INTEGER NOT NULL DEFAULT 0,
    "accessLevel" "AccessLevel" NOT NULL DEFAULT 'Member',
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "lectures_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "problem_sets" (
    "id" SERIAL NOT NULL,
    "moduleId" INTEGER NOT NULL,
    "problemSetTitle" TEXT NOT NULL,
    "description" TEXT,
    "video" TEXT,
    "submissionType" "SubmissionType" NOT NULL,
    "accessLevel" "AccessLevel" NOT NULL DEFAULT 'Member',
    "deadline" TIMESTAMP(3),
    "maxGrade" INTEGER NOT NULL DEFAULT 100,
    "passingGrade" INTEGER NOT NULL DEFAULT 60,
    "isManualGrading" BOOLEAN NOT NULL DEFAULT false,
    "order" INTEGER NOT NULL DEFAULT 0,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "problem_sets_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "problem_set_submissions" (
    "id" SERIAL NOT NULL,
    "userId" INTEGER NOT NULL,
    "problemSetId" INTEGER NOT NULL,
    "submissionUrl" TEXT,
    "grade" INTEGER NOT NULL DEFAULT 0,
    "submittedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "gradedAt" TIMESTAMP(3),

    CONSTRAINT "problem_set_submissions_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "path_progresses" (
    "id" SERIAL NOT NULL,
    "userId" INTEGER NOT NULL,
    "pathId" INTEGER NOT NULL,

    CONSTRAINT "path_progresses_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "module_progresses" (
    "id" SERIAL NOT NULL,
    "pathProgressId" INTEGER NOT NULL,
    "moduleId" INTEGER NOT NULL,
    "completed" BOOLEAN NOT NULL DEFAULT false,
    "completedAt" TIMESTAMP(3),

    CONSTRAINT "module_progresses_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "certificates" (
    "id" SERIAL NOT NULL,
    "userId" INTEGER NOT NULL,
    "pathId" INTEGER NOT NULL,
    "certificateType" "CertificateType" NOT NULL,
    "name" TEXT NOT NULL,
    "issueDate" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "certificateUrl" TEXT,
    "certificateId" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "certificates_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "events" (
    "id" SERIAL NOT NULL,
    "title" TEXT NOT NULL,
    "subtitle" TEXT,
    "category" "EventCategory" NOT NULL,
    "description" TEXT,
    "location" TEXT,
    "date" TIMESTAMP(3) NOT NULL,
    "timeStart" TEXT NOT NULL,
    "timeEnd" TEXT NOT NULL,
    "accessLevel" "EventAccessLevel" NOT NULL DEFAULT 'Member',
    "capacity" INTEGER,
    "eventType" "EventType" NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "events_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "event_attendees" (
    "id" SERIAL NOT NULL,
    "userId" INTEGER NOT NULL,
    "eventId" INTEGER NOT NULL,
    "rsvp" BOOLEAN NOT NULL DEFAULT false,
    "attended" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "event_attendees_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "users_email_key" ON "users"("email");

-- CreateIndex
CREATE UNIQUE INDEX "problem_set_submissions_userId_problemSetId_key" ON "problem_set_submissions"("userId", "problemSetId");

-- CreateIndex
CREATE UNIQUE INDEX "path_progresses_userId_pathId_key" ON "path_progresses"("userId", "pathId");

-- CreateIndex
CREATE UNIQUE INDEX "module_progresses_pathProgressId_moduleId_key" ON "module_progresses"("pathProgressId", "moduleId");

-- CreateIndex
CREATE UNIQUE INDEX "certificates_certificateId_key" ON "certificates"("certificateId");

-- CreateIndex
CREATE UNIQUE INDEX "event_attendees_userId_eventId_key" ON "event_attendees"("userId", "eventId");

-- AddForeignKey
ALTER TABLE "modules" ADD CONSTRAINT "modules_pathId_fkey" FOREIGN KEY ("pathId") REFERENCES "paths"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "lectures" ADD CONSTRAINT "lectures_moduleId_fkey" FOREIGN KEY ("moduleId") REFERENCES "modules"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "problem_sets" ADD CONSTRAINT "problem_sets_moduleId_fkey" FOREIGN KEY ("moduleId") REFERENCES "modules"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "problem_set_submissions" ADD CONSTRAINT "problem_set_submissions_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "problem_set_submissions" ADD CONSTRAINT "problem_set_submissions_problemSetId_fkey" FOREIGN KEY ("problemSetId") REFERENCES "problem_sets"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "path_progresses" ADD CONSTRAINT "path_progresses_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "path_progresses" ADD CONSTRAINT "path_progresses_pathId_fkey" FOREIGN KEY ("pathId") REFERENCES "paths"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "module_progresses" ADD CONSTRAINT "module_progresses_pathProgressId_fkey" FOREIGN KEY ("pathProgressId") REFERENCES "path_progresses"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "module_progresses" ADD CONSTRAINT "module_progresses_moduleId_fkey" FOREIGN KEY ("moduleId") REFERENCES "modules"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "certificates" ADD CONSTRAINT "certificates_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "certificates" ADD CONSTRAINT "certificates_pathId_fkey" FOREIGN KEY ("pathId") REFERENCES "paths"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "event_attendees" ADD CONSTRAINT "event_attendees_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "event_attendees" ADD CONSTRAINT "event_attendees_eventId_fkey" FOREIGN KEY ("eventId") REFERENCES "events"("id") ON DELETE CASCADE ON UPDATE CASCADE;
