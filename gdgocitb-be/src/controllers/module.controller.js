import { prisma } from '../config/index.js';
import logger from '../utils/logger.js';
import { successResponse, errorResponse } from '../utils/response.js';

// Get all modules
export const getModules = async (req, res) => {
  try {
    const modules = await prisma.module.findMany();
    return successResponse(res, 200, "Success", modules);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

// Get single module
export const getModule = async (req, res) => {
  try {
    const moduleData = await prisma.module.findUnique({
      where: { id: parseInt(req.params.id) },
    });
    if (!moduleData) {
      return res.status(404).json({ message: "Module not found" });
    }
    return successResponse(res, 200, "Success", moduleData);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

// Create new module and add it to the corresponding Path's modules relation
export const createModule = async (req, res) => {
    try {
        const { pathId, moduleName, description, video, order } = req.body;
        // Validasi bahwa field pathId (yang seharusnya berisi nama Path) ada
        console.log("Path ID:", pathId);
        // if (!pathId || pathId.trim() === "") {
        //     throw new Error("Path Name is required");
        // }
        console.log("Body request:", req.body);
  
      // Cari record Path berdasarkan nama yang diberikan
      logger.info("Path record found");
      const pathRecord = await prisma.path.findUnique({
        where: { id: parseInt(pathId) }
      });
      if (!pathRecord) {
        throw new Error(`Path not found with name: ${pathId}`);
      }
  
      // Buat module baru dan hubungkan dengan Path yang ditemukan
      const newModule = await prisma.module.create({
        data: {
          moduleName,
          description,
          video,
          order,
          path: { connect: { id: pathRecord.id } }  // hubungkan ke Path menggunakan id yang valid
        },
      });
        return successResponse(res, 201, 'Module created successfully', newModule);
    } catch (error) {
        return errorResponse(res, 500, error.message);
    }
};

// Update module
export const updateModule = async (req, res) => {
  try {
    const updatedModule = await prisma.module.update({
      where: { id: parseInt(req.params.id) },
      data: req.body,
    });
    if (!updatedModule) {
      return res.status(404).json({ message: "Module not found" });
    }
    return successResponse(res, 200, 'Success', updatedModule);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

// Delete module
export const deleteModule = async (req, res) => {
  try {
    const deletedModule = await prisma.module.delete({
      where: { id: parseInt(req.params.id) },
    });
    if (!deletedModule) {
      return res.status(404).json({ message: "Module not found" });
    }
    return successResponse(res, 200, 'Success', deletedModule);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

// Get module lectures
export const getModuleLectures = async (req, res) => {
  try {
    const moduleData = await prisma.module.findUnique({
      where: { id: parseInt(req.params.id) },
      include: { lectures: true },
    });
    if (!moduleData) {
      return res.status(404).json({ message: "Module not found" });
    }
    return successResponse(res, 200, 'Success', moduleData.lectures);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

// Get module problem sets
export const getModuleProblemSets = async (req, res) => {
  try {
    const moduleData = await prisma.module.findUnique({
      where: { id: parseInt(req.params.id) },
      include: { problemSets: true },
    });
    if (!moduleData) {
      return res.status(404).json({ message: "Module not found" });
    }
    return successResponse(res, 200, 'Success', moduleData.problemSets);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};