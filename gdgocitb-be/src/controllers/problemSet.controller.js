import { prisma } from '../config/index.js';
import { successResponse, errorResponse } from '../utils/response.js';
import logger from '../utils/logger.js';

export const getProblemSets = async (req, res) => {
    try {
        const moduleId = req.params.id; 
        
        const filters = {};
        if (moduleId) {
            filters.where = { moduleId };
        }
        
        const problemSets = await prisma.problemSet.findMany(filters);
        return successResponse(res, 200, 'Success retrieving problem sets', problemSets);
    } catch (error) {
        logger.error(`Error getting problem sets: ${error.message}`);
        return errorResponse(res, 500, error.message);
    }
};

export const getProblemSet = async (req, res) => {
    try {
        const problemSetId = req.params.id;
        const problemSet = await prisma.problemSet.findUnique({
            where: { id: problemSetId },
            include: {
                module: true,
                submissions: {
                    where: { userId: req.user.id },
                    orderBy: { submittedAt: 'desc' },
                    take: 1
                }
            }
        });
        
        if (!problemSet) {
            return errorResponse(res, 404, "Problem Set not found");
        }
        return successResponse(res, 200, 'Success retrieving problem set', problemSet);
    } catch (error) {
        logger.error(`Error getting problem set: ${error.message}`);
        return errorResponse(res, 500, error.message);
    }
};

export const createProblemSet = async (req, res) => {
    try {
        const { 
            moduleId, 
            problemSetTitle, 
            description, 
            submissionType, 
            accessLevel, 
            deadline, 
            maxGrade, 
            passingGrade, 
            isManualGrading, 
            order 
        } = req.body;
        
        // Periksa apakah module ada
        const moduleExists = await prisma.module.findUnique({
            where: { id: moduleId }
        });
        
        if (!moduleExists) {
            return errorResponse(res, 404, "Module not found");
        }
        
        // Tambahkan video file jika ada di request
        let video = null;
        if (req.fileData) {
            video = req.fileData.url;
        }
        
        const newProblemSet = await prisma.problemSet.create({
            data: {
                moduleId,
                problemSetTitle,
                description,
                video,
                submissionType,
                accessLevel: accessLevel || 'Member',
                deadline: deadline ? new Date(deadline) : null,
                maxGrade: maxGrade || 100,
                passingGrade: passingGrade || 60,
                isManualGrading: isManualGrading || false,
                order: order || 0
            }
        });
        
        return successResponse(res, 201, 'Problem set created successfully', newProblemSet);
    } catch (error) {
        logger.error(`Error creating problem set: ${error.message}`);
        return errorResponse(res, 500, error.message);
    }
};

export const updateProblemSet = async (req, res) => {
    try {
        const problemSetId = req.params.id;
        const {
            moduleId,
            problemSetTitle,
            description,
            submissionType,
            accessLevel,
            deadline,
            maxGrade,
            passingGrade,
            isManualGrading,
            order
        } = req.body;
        
        // Periksa apakah problem set ada
        const problemSetExists = await prisma.problemSet.findUnique({
            where: { id: problemSetId }
        });
        
        if (!problemSetExists) {
            return errorResponse(res, 404, "Problem Set not found");
        }
        
        // Update data
        const updateData = {
            moduleId,
            problemSetTitle,
            description,
            submissionType,
            accessLevel,
            deadline: deadline ? new Date(deadline) : problemSetExists.deadline,
            maxGrade,
            passingGrade,
            isManualGrading,
            order
        };
        
        // Filter undefined values
        Object.keys(updateData).forEach(key => {
            if (updateData[key] === undefined) {
                delete updateData[key];
            }
        });
        
        // Update video jika ada file baru
        if (req.fileData) {
            updateData.video = req.fileData.url;
        }
        
        const updatedProblemSet = await prisma.problemSet.update({
            where: { id: problemSetId },
            data: updateData
        });
        
        return successResponse(res, 200, 'Problem set updated successfully', updatedProblemSet);
    } catch (error) {
        logger.error(`Error updating problem set: ${error.message}`);
        return errorResponse(res, 500, error.message);
    }
};

export const deleteProblemSet = async (req, res) => {
    try {
        const problemSetId = req.params.id;
        
        // Periksa apakah problem set ada
        const problemSetExists = await prisma.problemSet.findUnique({
            where: { id: problemSetId }
        });
        
        if (!problemSetExists) {
            return errorResponse(res, 404, "Problem Set not found");
        }
        
        // Hapus semua submissions terkait dulu
        await prisma.problemSetSubmission.deleteMany({
            where: { problemSetId }
        });
        
        // Hapus problem set
        const deletedProblemSet = await prisma.problemSet.delete({
            where: { id: problemSetId }
        });
        
        return successResponse(res, 200, 'Problem set deleted successfully', deletedProblemSet);
    } catch (error) {
        logger.error(`Error deleting problem set: ${error.message}`);
        return errorResponse(res, 500, error.message);
    }
};

export const submitProblemSet = async (req, res) => {
    try {
        const problemSetId = req.params.id;
        const userId = req.user.id;
        const { submissionLink } = req.body;
        
        // Periksa apakah problem set ada
        const problemSet = await prisma.problemSet.findUnique({
            where: { id: problemSetId }
        });
        
        if (!problemSet) {
            return errorResponse(res, 404, "Problem Set not found");
        }
        
        // Periksa apakah sudah melewati deadline
        if (problemSet.deadline && new Date(problemSet.deadline) < new Date()) {
            return errorResponse(res, 400, "Submission deadline has passed");
        }
        
        // Periksa jenis submission
        let submissionUrl = null;
        
        if (problemSet.submissionType === 'Link') {
            // Validasi link submission
            if (!submissionLink) {
                return errorResponse(res, 400, "Submission link is required");
            }
            submissionUrl = submissionLink;
        } else if (['File', 'Image', 'GOCI'].includes(problemSet.submissionType)) {
            // Jika menggunakan file upload, periksa apakah file ada
            if (!req.fileData) {
                return errorResponse(res, 400, "File submission is required");
            }
            submissionUrl = req.fileData.url;
        }
        
        // Cek apakah sudah ada submission sebelumnya
        const existingSubmission = await prisma.problemSetSubmission.findFirst({
            where: {
                userId,
                problemSetId
            }
        });
        
        let submission;
        
        if (existingSubmission) {
            // Update submission yang sudah ada
            submission = await prisma.problemSetSubmission.update({
                where: { id: existingSubmission.id },
                data: {
                    submissionUrl,
                    submittedAt: new Date(),
                    // Jika grading otomatis, menetapkan grade default
                    grade: problemSet.isManualGrading ? existingSubmission.grade : 0,
                    gradedAt: problemSet.isManualGrading ? existingSubmission.gradedAt : null
                }
            });
            
            logger.info(`Updated submission for problem set ${problemSetId} by user ${userId}`);
        } else {
            // Buat submission baru
            submission = await prisma.problemSetSubmission.create({
                data: {
                    userId,
                    problemSetId,
                    submissionUrl,
                    grade: 0, // Default grade
                    submittedAt: new Date()
                }
            });
            
            logger.info(`Created new submission for problem set ${problemSetId} by user ${userId}`);
            
            // Update path progress jika belum ada
            await updateUserProgress(userId, problemSetId);
        }
        
        return successResponse(res, 200, 'Submission successful', submission);
        
    } catch (error) {
        logger.error(`Error submitting problem set: ${error.message}`);
        return errorResponse(res, 500, error.message);
    }
};

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