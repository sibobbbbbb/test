import User from '../models/User.js';
import logger from '../utils/logger.js';
import { config } from '../config/index.js';

/**
 * Check if email is in Member whitelist
 * @param {string} email - Email to check
 * @returns {Promise<boolean>} - True if email is in whitelist
 */
export const checkMemberWhitelist = async (email) => {
  try {
    // In a real implementation, this could check against a database of approved emails
    // For this example, we'll use a simple array of whitelisted domains and emails
    
    // Define whitelist
    const whitelistedDomains = ['itb.ac.id']; // ITB domain
    const whitelistedEmails = [
      'admin@gdgoc-itb.com',
      'test@example.com',
      '13522142@std.stei.itb.ac.id',
      '13522155@std.stei.itb.ac.id'
      // Add more whitelisted emails as needed
    ];
    
    // Check if email is directly whitelisted
    if (whitelistedEmails.includes(email.toLowerCase())) {
      return true;
    }
    
    // Check if email domain is whitelisted
    const domain = email.split('@')[1];
    if (whitelistedDomains.includes(domain.toLowerCase())) {
      return true;
    }
    
    // If we have an admin database, check there too
    const adminUser = await User.findOne({
      email,
      access: { $in: ['Curriculum Admin', 'Professional Development Admin', 'Technical Admin'] }
    });
    
    return !!adminUser;
  } catch (error) {
    logger.error(`Error checking member whitelist: ${error.message}`);
    return false;
  }
};

/**
 * Login or register user with Google
 * @param {Object} payload - Google OAuth payload
 * @returns {Promise<Object>} - User and token
 */
export const loginWithGoogle = async (payload) => {
  try {
    // Extract info from payload
    const { email, email_verified, name, picture } = payload;
    
    // Ensure email is verified
    if (!email_verified) {
      throw new Error('Email not verified by Google');
    }
    
    // Check if user exists
    let user = await User.findOne({ email });
    
    // If user doesn't exist, create a new one
    if (!user) {
      // Check if email is in member whitelist
      const isMember = await checkMemberWhitelist(email);
      
      if (!isMember) {
        throw new Error('Email not in whitelist. Please register as Buddy first.');
      }
      
      // Create new user as Member
      user = new User({
        name,
        email,
        access: 'Member'
      });
      
      await user.save();
    }
    
    // Generate JWT token
    const token = user.getSignedJwtToken();
    
    return { user, token };
  } catch (error) {
    logger.error(`Error in loginWithGoogle: ${error.message}`);
    throw error;
  }
};

/**
 * Register a new buddy
 * @param {Object} buddyData - Buddy data
 * @returns {Promise<Object>} - Registered buddy
 */
export const registerBuddy = async (buddyData) => {
  try {
    const { name, email } = buddyData;
    
    // Check if email is already in member whitelist
    const isMember = await checkMemberWhitelist(email);
    if (isMember) {
      throw new Error('This email is already registered as a Member. Please login directly.');
    }
    
    // Check if email is already registered as buddy
    const existingBuddy = await User.findOne({ email, access: 'Buddy' });
    if (existingBuddy) {
      throw new Error('This email is already registered as a Buddy. Please login directly.');
    }
    
    // Create new buddy user
    const buddy = new User({
      name,
      email,
      access: 'Buddy'
    });
    
    await buddy.save();
    
    return buddy;
  } catch (error) {
    logger.error(`Error in registerBuddy: ${error.message}`);
    throw error;
  }
};

export default {
  checkMemberWhitelist,
  loginWithGoogle,
  registerBuddy
};