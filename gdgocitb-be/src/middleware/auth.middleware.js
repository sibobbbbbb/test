import jwt from 'jsonwebtoken';
import { config } from '../config/index.js';
import { unauthorizedResponse, forbiddenResponse } from '../utils/response.js';
import User from '../models/User.js';

/**
 * Protect routes - verify token
 * @param {Object} req - Express request object
 * @param {Object} res - Express response object
 * @param {Function} next - Express next function
 */
export const protect = async (req, res, next) => {
  try {
    let token;

    // Cek token di Cookie terlebih dahulu
    if (req.cookies && req.cookies.gdgoc_auth_token) {
      token = req.cookies.gdgoc_auth_token;
    }
    // Jika tidak ada di cookie, cek di header Authorization
    else if (req.headers.authorization && req.headers.authorization.startsWith('Bearer')) {
      token = req.headers.authorization.split(' ')[1];
    }

    // Cek jika token ada
    if (!token) {
      return unauthorizedResponse(res, 'Not authorized, no token provided');
    }

    // Verifikasi token
    const decoded = jwt.verify(token, config.jwt.secret);

    // Cari user berdasarkan id
    const user = await User.findById(decoded.id).select('-password');

    if (!user) {
      return unauthorizedResponse(res, 'Not authorized, user not found');
    }

    // Set req.user ke user
    req.user = user;
    next();
  } catch (error) {
    return unauthorizedResponse(res, 'Not authorized, token failed');
  }
};

/**
 * Authorize by access level
 * @param  {...String} accessLevels - Required access levels (Member, Buddy, Professional Development Admin, etc)
 * @returns {Function} - Express middleware function
 */
export const authorize = (...accessLevels) => {
  return (req, res, next) => {
    if (!req.user) {
      return unauthorizedResponse(res, 'Not authorized, no user found');
    }

    if (!accessLevels.includes(req.user.access)) {
      return forbiddenResponse(res, `Access level ${req.user.access} is not authorized to access this route`);
    }

    next();
  };
};

export default {
  protect,
  authorize
};