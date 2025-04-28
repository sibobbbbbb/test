// import Module from '../models/Module.js';
import { successResponse, errorResponse } from '../utils/response.js';


export const getModules = async (req, res) => {
    try {
        const modules = await Module.find();
        return successResponse(res, modules);
    } catch (error) {
        return errorResponse(res, error.message);
    }
};

export const getModule = async (req, res) => {
    try {
        const module = await Module.findById(req.params.id);
        if (!module) {
            return res.status(404).json({ message: "Module not found" });
        }
        return successResponse(res, module);
    } catch (error) {
        return errorResponse(res, error.message);
    }
};

export const createModule = async (req, res) => {
    try {
        const { pathId, moduleName, description, video, order } = req.body;
        const newModule = new Module({ pathId, moduleName, description, video, order });
        const savedModule = await newModule.save();
        return res.status(201).json({ success: true, data: savedModule });
    } catch (error) {
        return errorResponse(res, error.message);
    }
};

export const updateModule = async (req, res) => {
    try {
        const updatedModule = await Module.findByIdAndUpdate(req.params.id, req.body, { new: true, runValidators: true });
        if (!updatedModule) {
            return res.status(404).json({ message: "Module not found" });
        }
        return successResponse(res, updatedModule);
    } catch (error) {
        return errorResponse(res, error.message);
    }
};

export const deleteModule = async (req, res) => {
    try {
        const deletedModule = await Module.findByIdAndDelete(req.params.id);
        if (!deletedModule) {
            return res.status(404).json({ message: "Module not found" });
        }
        return successResponse(res, deletedModule);
    } catch (error) {
        return errorResponse(res, error.message);
    }
};

export const getModuleLectures = async (req, res) => {
    try {
        const module = await Module.findById(req.params.id).populate('lectures');
        if (!module) {
            return res.status(404).json({ message: "Module not found" });
        }
        return successResponse(res, module.lectures);
    } catch (error) {
        return errorResponse(res, error.message);
    }
};

export const getModuleProblemSets = async (req, res) => {
    try {
        const module = await Module.findById(req.params.id).populate('problemSets');
        if (!module) {
            return res.status(404).json({ message: "Module not found" });
        }
        return successResponse(res, module.problemSets);
    } catch (error) {
        return errorResponse(res, error.message);
    }
};