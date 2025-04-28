import mongoose from 'mongoose';
import { config } from './environment.js';
import logger from '../utils/logger.js';

export const connectDB = async () => {
  try {
    const conn = await mongoose.connect(config.mongodb.uri);
    logger.info(`MongoDB Connected: ${conn.connection.host}`);
    return conn;
  } catch (error) {
    logger.error(`Error connecting to MongoDB: ${error.message}`);
    process.exit(1);
  }
};

export const disconnectDB = async () => {
  try {
    await mongoose.disconnect();
    logger.info('MongoDB Disconnected');
  } catch (error) {
    logger.error(`Error disconnecting from MongoDB: ${error.message}`);
    process.exit(1);
  }
};