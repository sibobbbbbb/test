import { createUploadMiddleware, createMultipleUploadMiddleware, uploadToCloudinary } from '../utils/upload.js';
import { errorResponse } from '../utils/response.js';
import logger from '../utils/logger.js';

/**
 * Middleware untuk menangani upload file tunggal
 * @param {Array} allowedTypes - Array dari tipe file yang diperbolehkan ('image', 'document', 'video', 'archive', 'code')
 * @param {String} fieldName - Nama field untuk file
 * @param {String} folder - Subfolder di Cloudinary (opsional)
 * @returns {Function} - Express middleware
 */
const handleSingleUpload = (allowedTypes = ['image', 'document'], fieldName = 'file', folder = '') => {
  return async (req, res, next) => {
    try {
      // Middleware untuk memproses file
      const upload = createUploadMiddleware(allowedTypes, fieldName);

      // Jalankan middleware upload
      upload(req, res, async (err) => {
        if (err) {
          logger.error(`Upload error: ${err.message}`);
          return errorResponse(res, 400, err.message);
        }

        // Jika tidak ada file yang diupload
        if (!req.file) {
          return next();
        }

        try {
          // Upload file ke Cloudinary
          const fileData = await uploadToCloudinary(req.file, folder);
          
          // Tambahkan file info ke request
          req.fileData = fileData;
          next();
        } catch (uploadError) {
          logger.error(`Cloudinary upload error: ${uploadError.message}`);
          return errorResponse(res, 500, 'Error uploading file to cloud storage');
        }
      });
    } catch (error) {
      logger.error(`Upload middleware error: ${error.message}`);
      return errorResponse(res, 500, 'Failed to process file upload');
    }
  };
};

/**
 * Middleware untuk menangani upload multiple files
 * @param {Array} allowedTypes - Array dari tipe file yang diperbolehkan ('image', 'document', 'video', 'archive', 'code')
 * @param {String} fieldName - Nama field untuk files
 * @param {Number} maxCount - Jumlah maksimum file yang bisa diupload
 * @param {String} folder - Subfolder di Cloudinary (opsional)
 * @returns {Function} - Express middleware
 */
const handleMultipleUpload = (allowedTypes = ['image', 'document'], fieldName = 'files', maxCount = 5, folder = '') => {
  return async (req, res, next) => {
    try {
      // Middleware untuk memproses files
      const upload = createMultipleUploadMiddleware(allowedTypes, fieldName, maxCount);

      // Jalankan middleware upload
      upload(req, res, async (err) => {
        if (err) {
          logger.error(`Upload error: ${err.message}`);
          return errorResponse(res, 400, err.message);
        }

        // Jika tidak ada file yang diupload
        if (!req.files || req.files.length === 0) {
          return next();
        }

        try {
          // Upload semua file ke Cloudinary
          const uploadPromises = req.files.map(file => uploadToCloudinary(file, folder));
          const filesData = await Promise.all(uploadPromises);
          
          // Tambahkan files info ke request
          req.filesData = filesData;
          next();
        } catch (uploadError) {
          logger.error(`Cloudinary upload error: ${uploadError.message}`);
          return errorResponse(res, 500, 'Error uploading files to cloud storage');
        }
      });
    } catch (error) {
      logger.error(`Upload middleware error: ${error.message}`);
      return errorResponse(res, 500, 'Failed to process files upload');
    }
  };
};

export {
  handleSingleUpload,
  handleMultipleUpload
};