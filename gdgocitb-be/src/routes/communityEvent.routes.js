import express from 'express';
import { protect, authorize } from '../middleware/auth.middleware.js';

const router = express.Router();

// Import controllers
import { 
  getCommunityEvents, 
  getCommunityEvent, 
  createCommunityEvent, 
  updateCommunityEvent, 
  deleteCommunityEvent,
  rsvpCommunityEvent,
  cancelRsvpCommunityEvent,
  markAttendance
} from '../controllers/communityEvent.controller.js';

// Define routes
// Get all community events
router.get('/', protect, getCommunityEvents);

// Get single community event
router.get('/:id', protect, getCommunityEvent);

// Create new community event - Admin only
router.post('/', protect, authorize('Technical Admin'), createCommunityEvent);

// Update community event - Admin only
router.put('/:id', protect, authorize('Technical Admin'), updateCommunityEvent);

// Delete community event - Admin only
router.delete('/:id', protect, authorize('Technical Admin'), deleteCommunityEvent);

// RSVP to community event
router.post('/:id/rsvp', protect, rsvpCommunityEvent);

// Cancel RSVP to community event
router.delete('/:id/rsvp', protect, cancelRsvpCommunityEvent);

// Mark attendance - Admin only
router.post('/:id/attendance/:userId', protect, authorize('Technical Admin'), markAttendance);

export default router;