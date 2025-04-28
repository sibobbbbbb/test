import Module from '../models/Module.js';
import { successResponse, errorResponse } from '../utils/response.js';
import { handleSingleUpload, handleMultipleUpload } from '../middleware/upload.middleware.js';
import ProblemSet from '../models/ProblemSet.js';

export const getProblemSets = async (req, res) => {
    try {
        const problemSets = await ProblemSet.find();
        return successResponse(res, problemSets);
    } catch (error) {
        return errorResponse(res, error.message);
    }
}

export const getProblemSet = async (req, res) => {
    try {
        const problemSetId = req.params.id;
        const module = await Module.findOne({ 'problemSets._id': problemSetId }).populate('problemSets');
        if (!module) {
            return res.status(404).json({ message: "Problem Set not found" });
        }
        const problemSet = module.problemSets.id(problemSetId);
        return successResponse(res, problemSet);
    } catch (error) {
        return errorResponse(res, error.message);
    }
}

export const createProblemSet = async (req, res) => {
    try {
        const moduleId = req.params.id;
        const { title, description, video, order } = req.body;
        const newProblemSet = { title, description, video, order };
        
        const module = await Module.findById(moduleId);
        if (!module) {
            return res.status(404).json({ message: "Module not found" });
        }
        
        module.problemSets.push(newProblemSet);
        await module.save();
        
        return res.status(201).json({ success: true, data: newProblemSet });
    } catch (error) {
        return errorResponse(res, error.message);
    }
}

export const updateProblemSet = async (req, res) => {
    try {
        const moduleId = req.params.id;
        const problemSetId = req.params.problemSetId;
        const { title, description, video, order } = req.body;
        
        const module = await Module.findById(moduleId);
        if (!module) {
            return res.status(404).json({ message: "Module not found" });
        }
        
        const problemSet = module.problemSets.id(problemSetId);
        if (!problemSet) {
            return res.status(404).json({ message: "Problem Set not found" });
        }
        
        problemSet.title = title || problemSet.title;
        problemSet.description = description || problemSet.description;
        problemSet.video = video || problemSet.video;
        problemSet.order = order || problemSet.order;
        
        await module.save();
        
        return successResponse(res, problemSet);
    } catch (error) {
        return errorResponse(res, error.message);
    }
}

export const deleteProblemSet = async (req, res) => {
    try {
        const moduleId = req.params.id;
        const problemSetId = req.params.problemSetId;
        
        const module = await Module.findById(moduleId);
        if (!module) {
            return res.status(404).json({ message: "Module not found" });
        }
        
        const problemSet = module.problemSets.id(problemSetId);
        if (!problemSet) {
            return res.status(404).json({ message: "Problem Set not found" });
        }
        
        problemSet.remove();
        await module.save();
        
        return successResponse(res, problemSet);
    } catch (error) {
        return errorResponse(res, error.message);
    }
}

export const submitProblemSet = async (req, res) => {
    try {
        const problemSetId = req.params.id;
        const { solution } = req.body;
        
        // Logic to handle submission (e.g., save to database, send email, etc.)
        
        return res.status(201).json({ success: true, message: "Problem Set submitted successfully", data: solution });
    } catch (error) {
        return errorResponse(res, error.message);
    }
}

export const gradeProblemSet = async (req, res) => {
    try {
        const problemSetId = req.params.id;
        const userId = req.params.userId;
        const { grade } = req.body;
        
        // Logic to handle grading (e.g., save to database, send email, etc.)
        
        return res.status(200).json({ success: true, message: "Problem Set graded successfully", data: { userId, problemSetId, grade } });
    } catch (error) {
        return errorResponse(res, error.message);
    }
}