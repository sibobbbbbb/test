import { prisma } from '../config/index.js';
import { successResponse, errorResponse } from '../utils/response.js';

export const getPaths = async (req, res) => {
  try {
    const paths = await prisma.path.findMany();
    return successResponse(res, 200, 'Success', paths);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

export const getPath = async (req, res) => {
  try {
    const path = await prisma.path.findUnique({
      where: {
        id: parseInt(req.params.id),
      },
      include: {
        modules: {
          select: {
            moduleName: true,
            description: true,
          },
        },
      },
    });
    if (!path) {
      return errorResponse(res, 404, 'Path not found');
    }
    return successResponse(res, 200, 'Success', path);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

export const createPath = async (req, res) => {
  try {
    const { pathName, description, modules } = req.body;
    const newPath = await prisma.path.create({
      data: {
        pathName,
        description,
        modules: modules
          ? {
              connect: modules.map((moduleId) => ({ id: parseInt(moduleId) })),
            }
          : undefined,
      },
    });
    return successResponse(res, 201, 'Path created', newPath);
  } catch (error) {
    return errorResponse(res, 500, error.message);
  }
};

export const updatePath = async (req, res) => {
    try {
        const { pathName, description, modules } = req.body;
        const updatedPath = await prisma.path.update({
        where: {
            id: parseInt(req.params.id),
        },
        data: {
            pathName,
            description,
            modules: modules
            ? {
                set: [],
                connect: modules.map((moduleId) => ({ id: parseInt(moduleId) })),
                }
            : undefined,
        },
        });
        return successResponse(res, 200, 'Path updated', updatedPath);
    } catch (error) {
        return errorResponse(res, 500, error.message);
    }
}