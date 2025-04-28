import Path from "../models/Path.js";
import { successResponse, errorResponse } from '../utils/response.js';

export const getPaths = async (req, res) => {
    try {
        const paths = await Path.find();
        return successResponse(res, 200, 'Success', paths);
    } catch (error) {
        return errorResponse(res, 500, error.message);
    }
}

export const getPath = async (req, res) => {
    try {
        const path = await Path.findById(req.params.id).populate('modules', 'moduleName description');
        if (!path) {
            return res.status(404).json({ message: "Path not found" });
        }
        return successResponse(res, path);
    } catch (error) {
        return errorResponse(res, error.message);
    }
}

export const createPath = async (req, res) => {
    try {
        const { pathName, description, modules } = req.body;
        const newPath = new Path({
            pathName,
            description,
            modules
        });
        const savedPath = await newPath.save();
        return res.status(201).json({ success: true, data: savedPath });
    } catch (error) {
        return errorResponse(res, error.message);
    }
}

