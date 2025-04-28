import express from 'express';
import { protect, authorize } from '../middleware/auth.middleware.js';
import { handleSingleUpload, handleMultipleUpload } from '../middleware/upload.middleware.js';

const router = express.Router();

// Import controllers once they are implemented
import { 
  getLectures, 
  getLecture, 
  createLecture, 
  updateLecture, 
  deleteLecture 
} from '../controllers/lecture.controller.js';

// Define routes
// Get all lectures
router.get('/', protect, getLectures);

// Get single lecture
router.get('/:id', protect, getLecture);

// Create new lecture - Admin only
router.post('/',
  protect,
  authorize('Curriculum Admin'),
  handleMultipleUpload(['document', 'archive', 'code'], 'materials', 3, 'lectures/materials'),
  createLecture
);

// Update lecture - Admin only
router.put('/:id',
  protect,
  authorize('Curriculum Admin'),
  handleMultipleUpload(['document', 'archive', 'code'], 'materials', 3, 'lectures/materials'),
  updateLecture
);

// Delete lecture - Admin only
router.delete('/:id',
  protect,
  deleteLecture
);

// Upload slides for lecture
router.post('/:id/slides',
  protect,
  authorize('Curriculum Admin'),
  handleSingleUpload(['document'], 'slides', 'lectures/slides'),
  (req, res) => {
    // req.fileData akan berisi informasi tentang file yang diupload
    res.status(501).json({ message: 'Not implemented yet', fileData: req.fileData });
  }
);

// Upload source code for lecture
router.post('/:id/source-code',
  protect,
  authorize('Curriculum Admin'),
  handleSingleUpload(['archive', 'code'], 'sourceCode', 'lectures/source-code'),
  (req, res) => {
    // req.fileData akan berisi informasi tentang file yang diupload
    res.status(501).json({ message: 'Not implemented yet', fileData: req.fileData });
  }
);

export default router;