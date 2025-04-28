import logger from '../utils/logger.js';
import { errorResponse } from '../utils/response.js';

/**
 * Error handling middleware
 * @param {Object} err - Error object
 * @param {Object} req - Express request object
 * @param {Object} res - Express response object
 * @param {Function} next - Express next function
 * @returns {Object} - Error response
 */
const errorMiddleware = (err, req, res, next) => {
  // Log error
  logger.error(`${err.name}: ${err.message}\nStack: ${err.stack}`);

  // Default error status and message
  let statusCode = 500;
  let message = 'Server Error';
  let errors = null;

  // Handle specific error types
  if (err.name === 'ValidationError') {
    // Mongoose validation error
    statusCode = 400;
    message = 'Validation Error';
    errors = Object.values(err.errors).map((error) => error.message);
  } else if (err.name === 'CastError' && err.kind === 'ObjectId') {
    // Mongoose objectId cast error
    statusCode = 404;
    message = 'Resource not found';
  } else if (err.code === 11000) {
    // Mongoose duplicate key error
    statusCode = 400;
    message = 'Duplicate field value entered';
    errors = err.keyValue;
  } else if (err.name === 'JsonWebTokenError') {
    // JWT error
    statusCode = 401;
    message = 'Invalid token';
  } else if (err.name === 'TokenExpiredError') {
    // JWT expired error
    statusCode = 401;
    message = 'Token expired';
  } else if (err.statusCode) {
    // Custom error with statusCode
    statusCode = err.statusCode;
    message = err.message;
    errors = err.errors;
  }

  // Send error response
  return errorResponse(res, statusCode, message, errors);
};

export default errorMiddleware;