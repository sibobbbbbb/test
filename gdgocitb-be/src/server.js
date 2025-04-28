import app from './app.js';
import { config } from './config/index.js';
import { connectDB } from './config/database.js';
import logger from './utils/logger.js';

// Connect to MongoDB
connectDB();

const PORT = config.port || 5000;

const server = app.listen(PORT, () => {
  logger.info(`Server running in ${config.nodeEnv} mode on port ${PORT}`);
});

// Handle unhandled promise rejections
process.on('unhandledRejection', (err) => {
  logger.error(`Unhandled Rejection: ${err.message}`);
  // Close server & exit process
  server.close(() => process.exit(1));
});

export default server;