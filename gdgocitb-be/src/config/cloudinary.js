import { v2 as cloudinary } from 'cloudinary';
import { CloudinaryStorage } from 'multer-storage-cloudinary';
import multer from 'multer';
import { config } from './environment.js';
import logger from '../utils/logger.js';

// Konfigurasi Cloudinary
cloudinary.config({
  cloud_name: config.cloudinary.cloudName,
  api_key: config.cloudinary.apiKey,
  api_secret: config.cloudinary.apiSecret,
  secure: true
});

// Multer configuration
const storage = new CloudinaryStorage({
  cloudinary: cloudinary,
  params: {
    folder: config.cloudinary.folder,
    allowedFormats: ['jpg', 'png', 'jpeg', 'gif', 'pdf', 'doc', 'docx', 'ppt', 'pptx', 'xls', 'xlsx', 'zip', 'rar', 'mp4', 'webm'],
    resource_type: 'auto'
  }
});

// Instance multer untuk upload file ke Cloudinary
const upload = multer({ storage: storage });

// Fungsi untuk upload file buffer ke Cloudinary
const uploadBuffer = async (buffer, options = {}) => {
  return new Promise((resolve, reject) => {
    const uploadOptions = {
      folder: config.cloudinary.folder,
      ...options
    };

    const uploadStream = cloudinary.uploader.upload_stream(
      uploadOptions,
      (error, result) => {
        if (error) {
          logger.error(`Error uploading to Cloudinary: ${error.message}`);
          return reject(error);
        }
        return resolve(result);
      }
    );

    import('streamifier').then(streamifier => {
      streamifier.default.createReadStream(buffer).pipe(uploadStream);
    });
  });
};

// Fungsi untuk menghapus file dari Cloudinary
const deleteFile = async (publicId) => {
  try {
    const result = await cloudinary.uploader.destroy(publicId);
    logger.info(`File deleted from Cloudinary: ${publicId}`);
    return result;
  } catch (error) {
    logger.error(`Error deleting file from Cloudinary: ${error.message}`);
    throw error;
  }
};

export {
  cloudinary,
  upload,
  uploadBuffer,
  deleteFile
};