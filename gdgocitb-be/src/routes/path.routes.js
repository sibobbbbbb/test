import express from 'express';
import { protect, authorize } from '../middleware/auth.middleware.js';

const router = express.Router();

// Import controllers once they are implemented
import { 
  getPaths, 
  getPath, 
  createPath, 
  // updatePath, 
  // deletePath, 
  // getPathModules
} from '../controllers/path.controller.js';

// Define routes
// Get all paths
router.get('/', protect, getPaths);

// Get single path
router.get('/:id', protect, getPath);

// Create new path - Admin only
router.post('/', protect, authorize('Curriculum Admin'), createPath);

// Update path - Admin only
router.put('/:id', protect, authorize('Curriculum Admin'), (req, res) => {
  res.status(501).json({ message: 'Not implemented yet' });
});

// Delete path - Admin only
router.delete('/:id', protect, authorize('Curriculum Admin'), (req, res) => {
  res.status(501).json({ message: 'Not implemented yet' });
});

// Get path modules
router.get('/:id/modules', protect, (req, res) => {
  res.status(501).json({ message: 'Not implemented yet' });
});

export default router;