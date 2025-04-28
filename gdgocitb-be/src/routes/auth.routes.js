import express from 'express';
import { protect } from '../middleware/auth.middleware.js';
import { 
  googleAuth, 
  getMe, 
  registerAdminUser, 
  registerBuddyUser,
  logout 
} from '../controllers/auth.controller.js';

const router = express.Router();

// Google OAuth route
router.post('/google', googleAuth);

// Register as buddy route
router.post('/register-buddy', registerBuddyUser);

// Get current user
router.get('/me', protect, getMe);

// Register admin
router.post('/register-admin', protect, registerAdminUser);

// Logout
router.post('/logout', protect, logout);

export default router;