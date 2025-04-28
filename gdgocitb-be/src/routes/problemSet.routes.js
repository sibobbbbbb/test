import express from 'express';
import { protect, authorize } from '../middleware/auth.middleware.js';
import { handleSingleUpload, handleMultipleUpload } from '../middleware/upload.middleware.js';

const router = express.Router();

// Import controllers once they are implemented
import { 
  getProblemSets, 
  getProblemSet, 
  createProblemSet, 
  updateProblemSet, 
  deleteProblemSet,
  submitProblemSet,
  gradeProblemSet
} from '../controllers/problemSet.controller.js';

// Define routes
// Get all problem sets
router.get('/', protect, getProblemSets);

// Get single problem set
router.get('/:id', protect, getProblemSet);

// Get problem set by module ID
router.get('/module/:id', protect, getProblemSets);

// Create new problem set - Admin only
// Memperbolehkan upload video tutorial untuk problem set
router.post('/', 
  protect, 
  authorize('Curriculum Admin'),
  handleSingleUpload(['video', 'document'], 'video', 'problem-sets/videos'),
  createProblemSet
);

// Update problem set - Admin only
router.put('/:id',
  protect, 
  authorize('Curriculum Admin'),
  handleSingleUpload(['video', 'document'], 'video', 'problem-sets/videos'),
  updateProblemSet
);

// Delete problem set - Admin only
router.delete('/:id', 
  protect, 
  authorize('Curriculum Admin'), 
  deleteProblemSet
);

// Submit problem set solution
// Memperbolehkan upload file/image/link solusi
router.post('/:id/submit', 
  protect,
  handleSingleUpload(['image', 'document', 'archive', 'code'], 'submission', 'problem-sets/submissions'),
  submitProblemSet
);

// Grade problem set - Admin only
router.post('/:id/grade/:userId',
  protect,
  authorize('Curriculum Admin'),
  gradeProblemSet
);

export default router;