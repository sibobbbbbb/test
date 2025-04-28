import express from 'express';
import { protect, authorize } from '../middleware/auth.middleware.js';

const router = express.Router();

// Import controllers once they are implemented
import { 
  getCertificates, 
  getCertificate, 
  createCertificate, 
  updateCertificate, 
  deleteCertificate,
  verifyCertificate
} from '../controllers/certificate.controller.js';

// Define routes
// Get all certificates - Admin only
router.get('/',
  protect,
  authorize('Curriculum Admin'),
  getCertificates
);

// Get single certificate
router.get('/:id',
  protect,
  getCertificate
);

// Create new certificate - Admin only
router.post('/',
  protect,
  authorize('Curriculum Admin'),
  createCertificate
);

// Update certificate - Admin only
router.put('/:id',
  protect,
  authorize('Curriculum Admin'),
  updateCertificate
);

// Delete certificate - Admin only
router.delete('/:id',
  protect,
  authorize('Curriculum Admin'),
  deleteCertificate
);

// Verify certificate - Public route
router.get('/verify/:certificateId',
  verifyCertificate
);

export default router;