import express from 'express';
import { protect, authorize } from '../middleware/auth.middleware.js';

const router = express.Router();

import { 
  getModules, 
  getModule, 
  createModule, 
  updateModule, 
  deleteModule, 
  getModuleLectures,
  getModuleProblemSets
} from '../controllers/module.controller.js';

// Define routes
// Get all modules
router.get('/', getModules);

// Get single module
router.get('/:id', protect, getModule);

// Create new module - Admin only
router.post('/', protect, authorize('Curriculum Admin'), createModule);

// Update module - Admin only
router.put('/:id', protect, authorize('Curriculum Admin'), updateModule);

// Delete module - Admin only
router.delete('delModule/:id', protect, authorize('Curriculum Admin'), deleteModule);

// Get module lectures
router.get('module/:id/lectures', protect, getModuleLectures);

// Get module problem sets
router.get('module/:id/problem-sets', protect, getModuleProblemSets);

export default router; 