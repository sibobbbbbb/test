# Models Documentation

## User

**User** atribut:
- `name` (String) - Nama pengguna.
- `email` (String) - Email pengguna (harus unik).
- `googleId` (String) - ID Google pengguna (harus unik).
- `avatar` (String) - URL gambar avatar.
- `role` (String) - Salah satu dari: `admin`, `member_itb`, `buddy`. Default `buddy`.
- `eventsRegistered` (Array) - Daftar event yang didaftarkan.
- `coursesEnrolled` (Array) - Daftar kursus yang diikuti, beserta progress.
- `lastLogin` (Date) - Waktu terakhir login.
- `isActive` (Boolean) - Status aktif pengguna.
- **Timestamps** - `createdAt` dan `updatedAt`.

---

## Event

**Event** atribut:
- `name` (String) - Nama event.
- `description` (String) - Deskripsi event.
- `date` (Date) - Tanggal event berlangsung.
- `location` (String) - Lokasi event.
- `imageUrl` (String) - URL gambar event.
- `rsvpLink` (String) - Link RSVP untuk pendaftaran event.
- `eventSpeakers` (Array) - Pembicara di event.
- `attendees` (Array) - Peserta yang hadir.
- `isOpen` (Boolean) - Status pendaftaran masih dibuka/tidak.
- **Timestamps** - `createdAt` dan `updatedAt`.

---

## Speaker

**Speaker** atribut:
- `name` (String) - Nama pembicara.
- `position` (String) - Jabatan pembicara.
- `company` (String) - Perusahaan tempat bekerja.
- `imageUrl` (String) - Foto pembicara.
- **Timestamps** - `createdAt` dan `updatedAt`.

---

## Course

**Course** atribut:
- `title` (String) - Judul kursus.
- `description` (String) - Deskripsi kursus.
- `thumbnail` (String) - Thumbnail gambar kursus.
- `modules` (Array) - Modul yang terdapat dalam kursus.
- `roadmapId` (ObjectId) - Roadmap terkait.
- `enrolledUsers` (Array) - Pengguna yang terdaftar.
- `isPublished` (Boolean) - Status publikasi kursus.
- **Timestamps** - `createdAt` dan `updatedAt`.

---

## Module

**Module** atribut:
- `title` (String) - Judul modul.
- `description` (String) - Deskripsi modul.
- `courseId` (ObjectId) - Kursus terkait.
- `order` (Number) - Urutan modul dalam kursus.
- `content` (String) - Isi materi.
- `contentType` (String) - Jenis konten: `text`, `video`, `assignment`.
- `videoUrl` (String) - Jika konten berupa video.
- `assignment` (Array) - Kumpulan assignment terkait.
- `isPublished` (Boolean) - Status publikasi modul.
- `prerequisites` (Array) - Modul yang harus diselesaikan sebelum mengakses.
- **Timestamps** - `createdAt` dan `updatedAt`.

---

## Roadmap

**Roadmap** atribut:
- `title` (String) - Judul roadmap.
- `image` (String) - Thumbnail gambar roadmap.
- `steps` (Array) - Langkah-langkah dalam roadmap (setiap langkah berisi modul).
- `isPublished` (Boolean) - Status publikasi roadmap.
- `estimatedCompletionTime` (String) - Estimasi waktu penyelesaian.
- **Timestamps** - `createdAt` dan `updatedAt`.

---

## Progress

**Progress** atribut:
- `userId` (ObjectId) - Pengguna yang mengikuti kursus atau roadmap.
- `courseId` (ObjectId) - Kursus yang sedang diikuti.
- `roadmapId` (ObjectId) - Roadmap yang diikuti (opsional).
- `completedModules` (Array) - Modul yang telah diselesaikan beserta informasi tambahan (score).
- `currentModule` (ObjectId) - Modul yang sedang dikerjakan.
- `overallProgress` (Number) - Progress keseluruhan (0-100%).
- `startedAt` (Date) - Tanggal mulai.
- `lastAccessedAt` (Date) - Terakhir kali mengakses.
- `isCompleted` (Boolean) - Status penyelesaian kursus.
- `completedAt` (Date) - Tanggal selesai.
- `certificate` (Object) - Informasi sertifikat (issued atau tidak, url sertifikat).
- **Timestamps** - `createdAt` dan `updatedAt`.

---
