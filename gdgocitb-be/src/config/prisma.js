import { PrismaClient } from '@prisma/client';
import logger from '../utils/logger.js';

let prisma;

if (process.env.NODE_ENV === 'production') {
  prisma = new PrismaClient();
} else {
  if (!global.prisma) {
    global.prisma = new PrismaClient({
      log: [
        {
          emit: 'event',
          level: 'query',
        },
        {
          emit: 'event',
          level: 'error',
        },
        {
          emit: 'event',
          level: 'info',
        },
        {
          emit: 'event',
          level: 'warn',
        },
      ],
    });

    global.prisma.$on('query', (e) => {
      logger.debug(`Query: ${e.query}`);
      logger.debug(`Duration: ${e.duration}ms`);
    });

    global.prisma.$on('error', (e) => {
      logger.error(`Prisma Error: ${e.message}`);
    });
  }
  prisma = global.prisma;
}

// Fungsi untuk connect ke database
const connectDB = async () => {
  try {
    await prisma.$connect();
    logger.info('PostgreSQL Database Connected');
    return prisma;
  } catch (error) {
    logger.error(`Error connecting to PostgreSQL: ${error.message}`);
    process.exit(1);
  }
};

// Fungsi untuk disconnect dari database
const disconnectDB = async () => {
  try {
    await prisma.$disconnect();
    logger.info('PostgreSQL Database Disconnected');
  } catch (error) {
    logger.error(`Error disconnecting from PostgreSQL: ${error.message}`);
    process.exit(1);
  }
};

export { prisma, connectDB, disconnectDB };