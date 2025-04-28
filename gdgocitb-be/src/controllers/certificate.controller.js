import { successResponse, errorResponse } from '../utils/response.js';
// import Certificate from '../models/Certificate.js';

export const getCertificates = async (req, res) => {
    try {
        const certificates = await Certificate.find();
        return successResponse(res, certificates);
    } catch (error) {
        return errorResponse(res, error.message);
    }
}

export const getCertificate = async (req, res) => {
    try {
        const certificate = await Certificate.findById(req.params.id);
        if (!certificate) {
            return res.status(404).json({ message: "Certificate not found" });
        }
        return successResponse(res, certificate);
    } catch (error) {
        return errorResponse(res, error.message);
    }
}

export const createCertificate = async (req, res) => {
    try {
        const { userId, moduleId, pathId } = req.body;
        const newCertificate = new Certificate({ userId, moduleId, pathId });
        const savedCertificate = await newCertificate.save();
        return res.status(201).json({ success: true, data: savedCertificate });
    } catch (error) {
        return errorResponse(res, error.message);
    }
}

export const updateCertificate = async (req, res) => {
    try {
        const updatedCertificate = await Certificate.findByIdAndUpdate(req.params.id, req.body, { new: true, runValidators: true });
        if (!updatedCertificate) {
            return res.status(404).json({ message: "Certificate not found" });
        }
        return successResponse(res, updatedCertificate);
    } catch (error) {
        return errorResponse(res, error.message);
    }
}

export const deleteCertificate = async (req, res) => {
    try {
        const deletedCertificate = await Certificate.findByIdAndDelete(req.params.id);
        if (!deletedCertificate) {
            return res.status(404).json({ message: "Certificate not found" });
        }
        return successResponse(res, deletedCertificate);
    } catch (error) {
        return errorResponse(res, error.message);
    }
}

export const verifyCertificate = async (req, res) => {
    try {
        const certificate = await Certificate.findOne({ certificateId: req.params.certificateId });
        if (!certificate) {
            return res.status(404).json({ message: "Certificate not found" });
        }
        return successResponse(res, certificate);
    } catch (error) {
        return errorResponse(res, error.message);
    }
}
