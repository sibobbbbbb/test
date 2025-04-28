import multer from 'multer';
import path from 'path';
import { cloudinary, uploadBuffer } from '../config/index.js';
import logger from './logger.js';

// Definisi jenis file yang diperbolehkan
const fileTypes = {
  image: ['image/jpeg', 'image/png', 'image/gif', 'image/webp'],
  document: ['application/pdf', 'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
             'application/vnd.ms-excel', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
             'application/vnd.ms-powerpoint', 'application/vnd.openxmlformats-officedocument.presentationml.presentation'],
  video: ['video/mp4', 'video/webm', 'video/quicktime'],
  archive: ['application/zip', 'application/x-rar-compressed'],
  code: ['text/plain', 'application/json', 'text/html', 'text/css', 'text/javascript', 
         'application/javascript', 'application/x-python-code']
};

// Filter file berdasarkan type
const fileFilter = (allowedTypes) => (req, file, cb) => {
  // Gabungkan semua tipe file yang diperbolehkan
  let allowed = [];
  allowedTypes.forEach(type => {
    if (fileTypes[type]) {
      allowed = [...allowed, ...fileTypes[type]];
    }
  });

  if (allowed.includes(file.mimetype)) {
    cb(null, true);
  } else {
    cb(new Error(`Tipe file tidak diperbolehkan. Hanya menerima ${allowedTypes.join(', ')}`), false);
  }
};

// Storage sementara untuk multer
const storage = multer.memoryStorage();

// Generator untuk middleware upload
const createUploadMiddleware = (allowedTypes = ['image', 'document'], fieldName = 'file') => {
  return multer({
    storage: storage,
    fileFilter: fileFilter(allowedTypes),
    limits: {
      fileSize: 10 * 1024 * 1024 // 10MB
    }
  }).single(fieldName);
};

// Middleware untuk multiple file upload
const createMultipleUploadMiddleware = (allowedTypes = ['image', 'document'], fieldName = 'files', maxCount = 5) => {
  return multer({
    storage: storage,
    fileFilter: fileFilter(allowedTypes),
    limits: {
      fileSize: 10 * 1024 * 1024 // 10MB
    }
  }).array(fieldName, maxCount);
};

// Upload file ke Cloudinary
const uploadToCloudinary = async (file, folder = '') => {
  try {
    const folderPath = folder ? `${cloudinary.folder}/${folder}` : cloudinary.folder;
    
    // Tentukan resource_type berdasarkan mimetype
    let resourceType = 'auto';
    if (fileTypes.image.includes(file.mimetype)) {
      resourceType = 'image';
    } else if (fileTypes.video.includes(file.mimetype)) {
      resourceType = 'video';
    } else {
      resourceType = 'raw';
    }

    const result = await uploadBuffer(file.buffer, {
      folder: folderPath,
      resource_type: resourceType,
      public_id: `${Date.now()}-${path.parse(file.originalname).name}`.replace(/\s+/g, '-')
    });

    return {
      publicId: result.public_id,
      url: result.secure_url,
      format: result.format,
      resourceType: result.resource_type
    };
  } catch (error) {
    logger.error(`Error uploading file to Cloudinary: ${error.message}`);
    throw error;
  }
};

export {
  createUploadMiddleware,
  createMultipleUploadMiddleware,
  uploadToCloudinary,
  fileTypes
};