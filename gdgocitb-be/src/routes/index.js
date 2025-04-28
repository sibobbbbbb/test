import express from 'express';
import authRoutes from './auth.routes.js';
import userRoutes from './user.routes.js';
import pathRoutes from './path.routes.js';
import moduleRoutes from './module.routes.js';
import lectureRoutes from './lecture.routes.js';
import problemSetRoutes from './problemSet.routes.js';
import professionalEventRoutes from './professionalEvent.routes.js';
import communityEventRoutes from './communityEvent.routes.js';
import certificateRoutes from './certificate.routes.js';

const router = express.Router();

// Mount routes
router.use('/auth', authRoutes);
router.use('/users', userRoutes);
router.use('/paths', pathRoutes);
router.use('/modules', moduleRoutes);
router.use('/lectures', lectureRoutes);
router.use('/problem-sets', problemSetRoutes);
router.use('/professional-events', professionalEventRoutes);
router.use('/community-events', communityEventRoutes);
router.use('/certificates', certificateRoutes);

// API health check route
router.get('/health', (req, res) => {
  res.status(200).json({
    status: 'success',
    message: 'API is running'
  });
});

export default router;