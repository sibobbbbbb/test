/**
 * Utility for standardized API responses
 */

/**
 * Success response formatter
 * @param {Object} res - Express response object
 * @param {number} statusCode - HTTP status code
 * @param {string} message - Success message
 * @param {Object|Array} data - Response data
 * @returns {Object} - Formatted response
 */
export const successResponse = (res, statusCode = 200, message = 'Success', data = {}) => {
    return res.status(statusCode).json({
      success: true,
      message,
      data
    });
  };
  
  /**
   * Error response formatter
   * @param {Object} res - Express response object
   * @param {number} statusCode - HTTP status code
   * @param {string} message - Error message
   * @param {Object|Array} errors - Error details
   * @returns {Object} - Formatted response
   */
  export const errorResponse = (res, statusCode = 500, message = 'Internal Server Error', errors = null) => {
    return res.status(statusCode).json({
      success: false,
      message,
      errors: errors || null
    });
  };
  
  /**
   * Not found response formatter
   * @param {Object} res - Express response object
   * @param {string} message - Not found message
   * @returns {Object} - Formatted response
   */
  export const notFoundResponse = (res, message = 'Resource not found') => {
    return errorResponse(res, 404, message);
  };
  
  /**
   * Validation error response formatter
   * @param {Object} res - Express response object
   * @param {Object|Array} errors - Validation errors
   * @returns {Object} - Formatted response
   */
  export const validationErrorResponse = (res, errors) => {
    return errorResponse(res, 422, 'Validation Error', errors);
  };
  
  /**
   * Unauthorized response formatter
   * @param {Object} res - Express response object
   * @param {string} message - Unauthorized message
   * @returns {Object} - Formatted response
   */
  export const unauthorizedResponse = (res, message = 'Unauthorized') => {
    return errorResponse(res, 401, message);
  };
  
  /**
   * Forbidden response formatter
   * @param {Object} res - Express response object
   * @param {string} message - Forbidden message
   * @returns {Object} - Formatted response
   */
  export const forbiddenResponse = (res, message = 'Forbidden') => {
    return errorResponse(res, 403, message);
  };
  
  export default {
    successResponse,
    errorResponse,
    notFoundResponse,
    validationErrorResponse,
    unauthorizedResponse,
    forbiddenResponse
  };