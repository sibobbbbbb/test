// import Lecture from '../models/Lecture.js';
import { successResponse, errorResponse } from '../utils/response.js';

export const getLectures = async (req, res) => {
    try {
        const lectures = await Lecture.find();
        return successResponse(res, lectures);
    } catch (error) {
        return errorResponse(res, error.message);
    }
};

export const getLecture = async (req, res) => {
    try {
        const lecture = await Lecture.findById(req.params.id);
        if (!lecture) {
            return res.status(404).json({ message: "Lecture not found" });
        }
        return successResponse(res, lecture);
    } catch (error) {
        return errorResponse(res, error.message);
    }
};

export const createLecture = async (req, res) => {
    try {
        const { pathId, moduleId, title, notes, order, accessLevel } = req.body;
        const newLecture = new Lecture({
            pathId,
            moduleId,
            title,
            notes,
            order,
            accessLevel
        });
        const savedLecture = await newLecture.save();
        return res.status(201).json({ success: true, data: savedLecture });
    } catch (error) {
        return errorResponse(res, error.message);
    }
};

export const updateLecture = async (req, res) => {
    try {
        const updatedLecture = await Lecture.findByIdAndUpdate(
            req.params.id,
            req.body,
            { new: true, runValidators: true }
        );
        if (!updatedLecture) {
            return res.status(404).json({ message: "Lecture not found" });
        }
        return successResponse(res, updatedLecture);
    } catch (error) {
        return errorResponse(res, error.message);
    }
};

export const deleteLecture = async (req, res) => {
    try {
        const deletedLecture = await Lecture.findByIdAndDelete(req.params.id);
        if (!deletedLecture) {
            return res.status(404).json({ message: "Lecture not found" });
        }
        return successResponse(res, deletedLecture);
    } catch (error) {
        return errorResponse(res, error.message);
    }
};