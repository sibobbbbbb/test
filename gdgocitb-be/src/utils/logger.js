import winston from 'winston';

const { format, createLogger, transports } = winston;
const { combine, timestamp, printf, colorize } = format;

// Define custom format
const logFormat = printf(({ level, message, timestamp }) => {
  return `${timestamp} [${level}]: ${message}`;
});

// Create Winston logger
const logger = createLogger({
  level: 'info',
  format: combine(
    timestamp({ format: 'YYYY-MM-DD HH:mm:ss' }),
    format.errors({ stack: true }),
    logFormat
  ),
  defaultMeta: { service: 'lms-gdgoc-itb' },
  transports: [
    // Write all logs to console
    new transports.Console({
      format: combine(
        colorize(),
        logFormat
      )
    }),
    // Write all logs with level 'error' and below to error.log
    new transports.File({ filename: 'logs/error.log', level: 'error' }),
    // Write all logs to combined.log
    new transports.File({ filename: 'logs/combined.log' })
  ]
});

// If not in production, log to console with colors
if (process.env.NODE_ENV !== 'production') {
  logger.add(new transports.Console({
    format: combine(
      colorize(),
      logFormat
    )
  }));
}

export default logger;