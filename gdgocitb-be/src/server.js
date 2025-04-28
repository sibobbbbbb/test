import app from './app.js';
import { config, connectDB, disconnectDB } from './config/index.js';
import logger from './utils/logger.js';

// Connect to PostgreSQL
connectDB()
  .then(() => {
    const PORT = config.port || 5000;

    const server = app.listen(PORT, () => {
      logger.info(`Server running in ${config.nodeEnv} mode on port ${PORT}`);
    });

    // Handle unhandled promise rejections
    process.on('unhandledRejection', (err) => {
      logger.error(`Unhandled Rejection: ${err.message}`);
      
      // Disconnect from DB before closing server
      disconnectDB().then(() => {
        server.close(() => process.exit(1));
      });
    });

    // Handle graceful shutdown
    process.on('SIGTERM', () => {
      logger.info('SIGTERM received, shutting down gracefully');
      
      server.close(() => {
        logger.info('Process terminated');
        disconnectDB();
      });
    });
  })
  .catch(err => {
    logger.error(`Database connection failed: ${err.message}`);
    process.exit(1);
  });

export default app;