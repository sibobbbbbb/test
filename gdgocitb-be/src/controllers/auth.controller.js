import { loginWithGoogle, checkMemberWhitelist } from '../services/auth.service.js';
import { successResponse, errorResponse, unauthorizedResponse } from '../utils/response.js';
import logger from '../utils/logger.js';
import { OAuth2Client } from 'google-auth-library';
import { config } from '../config/index.js';
import User from '../models/User.js';

// Cookie settings yang sudah ditingkatkan
const COOKIE_OPTIONS = {
  httpOnly: true,
  secure: process.env.NODE_ENV === 'production',
  sameSite: 'strict',
  maxAge: 7 * 24 * 60 * 60 * 1000, // 7 days
  path: '/'
};

/**
 * @desc    Login atau register dengan Google
 * @route   POST /api/auth/google
 * @access  Public
 */
const googleAuth = async (req, res) => {
  try {
    const { idToken } = req.body;
    
    if (!idToken) {
      return errorResponse(res, 400, 'Google ID token is required');
    }
    
    // Get payload from Google token
    let payload;
    try {
      payload = await getPayloadFromGoogleToken(idToken);
    } catch (error) {
      logger.error(`Google token verification failed: ${error.message}`);
      return errorResponse(res, 401, 'Invalid Google token');
    }
    
    logger.info(`Google auth request for: ${payload.email}`);
    
    // Check if email is in member whitelist
    const isMember = await checkMemberWhitelist(payload.email);
    logger.info(`User ${payload.email} is ${isMember ? 'a member' : 'not a member'}`);
    if (!isMember) {
      // Check if email is registered as buddy
      const buddy = await User.findOne({ email: payload.email, access: 'Buddy' });
      logger.info(`Buddy ${buddy ? 'found' : 'not found'}`);
      if (!buddy) {
        // Not a member and not a buddy - need to register first
        return errorResponse(res, 403, 'Email not in whitelist. Please register as Buddy first.');
      }
      
      // Is a buddy, proceed with login
    }
    
    // Login or register with Google
    const { user, token } = await loginWithGoogle(payload);
    
    logger.info(`Login successful for ${user.email} (${user.access})`);
    
    // Set cookie dengan token JWT
    res.cookie('gdgoc_auth_token', token, COOKIE_OPTIONS);
    
    // Tambahkan cookie untuk user data (non-sensitive)
    const userData = {
      id: user._id,
      name: user.name,
      email: user.email,
      access: user.access
    };
    
    // Set user data cookie (accessible to JavaScript)
    res.cookie('gdgoc_user', JSON.stringify(userData), {
      ...COOKIE_OPTIONS,
      httpOnly: false
    });
    
    return successResponse(res, 200, 'Login successful', {
      user: userData,
      token
    });
  } catch (error) {
    logger.error(`Error in googleAuth: ${error.message}`);
    return errorResponse(res, 500, 'Login failed', error.message);
  }
};

/**
 * @desc    Register as Buddy
 * @route   POST /api/auth/register-buddy
 * @access  Public
 */
const registerBuddyUser = async (req, res) => {
  try {
    const { name, email } = req.body;
    
    logger.info(`Registering buddy: ${email}`);
    
    // Validasi input
    if (!name || !email) {
      return errorResponse(res, 400, 'Please provide name and email');
    }
    
    // Check if email is already in member whitelist
    const isMember = await checkMemberWhitelist(email);
    if (isMember) {
      return errorResponse(res, 400, 'This email is already registered as a Member. Please login directly.');
    }
    
    // Check if email is already registered as buddy
    const existingBuddy = await User.findOne({ email, access: 'Buddy' });
    if (existingBuddy) {
      return errorResponse(res, 400, 'This email is already registered as a Buddy. Please login directly.');
    }
    
    // Create new buddy user
    const buddy = new User({
      name,
      email,
      access: 'Buddy'
    });
    
    await buddy.save();
    logger.info(`Buddy registered successfully: ${email}`);
    
    return successResponse(res, 201, 'Registration successful. You can now login with Google.', {
      email: buddy.email
    });
  } catch (error) {
    logger.error(`Error in registerBuddyUser: ${error.message}`);
    return errorResponse(res, 500, 'Registration failed', error.message);
  }
};

/**
 * @desc    Get currently authenticated user
 * @route   GET /api/auth/me
 * @access  Private
 */
const getMe = async (req, res) => {
  try {
    const user = req.user;
    
    logger.info(`getMe called for user: ${user.email}`);
    
    // Refresh cookies to extend session
    const token = user.getSignedJwtToken();
    res.cookie('gdgoc_auth_token', token, COOKIE_OPTIONS);
    
    // Refresh user data cookie
    const userData = {
      id: user._id,
      name: user.name,
      email: user.email,
      access: user.access
    };
    
    res.cookie('gdgoc_user', JSON.stringify(userData), {
      ...COOKIE_OPTIONS,
      httpOnly: false
    });
    
    logger.info('User cookies refreshed');
    
    return successResponse(res, 200, 'User fetched successfully', {
      user: userData,
      token // Include token in response
    });
  } catch (error) {
    logger.error(`Error in getMe: ${error.message}`);
    return errorResponse(res, 500, 'Failed to fetch user');
  }
};

/**
 * @desc    Register new admin
 * @route   POST /api/auth/register-admin
 * @access  Private/Admin
 */
const registerAdminUser = async (req, res) => {
  try {
    // Pastikan yang request adalah admin
    if (!['Curriculum Admin', 'Professional Development Admin', 'Technical Admin'].includes(req.user.access)) {
      return unauthorizedResponse(res, 'Not authorized to register admin');
    }
    
    const { name, email, adminType } = req.body;
    
    // Validasi input
    if (!name || !email || !adminType) {
      return errorResponse(res, 400, 'Please provide name, email, and admin type');
    }
    
    // Check if email is already registered
    const existingUser = await User.findOne({ email });
    if (existingUser) {
      return errorResponse(res, 400, 'Email already registered');
    }
    
    // Create new admin user
    const admin = new User({
      name,
      email,
      access: adminType
    });
    
    await admin.save();
    
    // Generate token
    const token = admin.getSignedJwtToken();
    
    return successResponse(res, 201, 'Admin registered successfully', {
      admin: {
        id: admin._id,
        name: admin.name,
        email: admin.email,
        access: admin.access
      },
      token
    });
  } catch (error) {
    logger.error(`Error in registerAdminUser: ${error.message}`);
    return errorResponse(res, 500, 'Failed to register admin', error.message);
  }
};

/**
 * @desc    Logout user (clear cookie on client)
 * @route   POST /api/auth/logout
 * @access  Private
 */
const logout = (req, res) => {
  logger.info(`Logout request for user: ${req.user ? req.user.email : 'unknown'}`);
  
  // Clear cookies with same path setting
  res.clearCookie('gdgoc_auth_token', { path: '/' });
  res.clearCookie('gdgoc_user', { path: '/' });
  
  logger.info('Cookies cleared');
  
  return successResponse(res, 200, 'Logged out successfully');
};

// Helper untuk mendapatkan payload dari Google token
const getPayloadFromGoogleToken = async (idToken) => {
  const client = new OAuth2Client(config.google.clientId);
  const ticket = await client.verifyIdToken({
    idToken,
    audience: config.google.clientId
  });
  
  return ticket.getPayload();
};

export {
  googleAuth,
  registerBuddyUser,
  getMe,
  registerAdminUser,
  logout
};