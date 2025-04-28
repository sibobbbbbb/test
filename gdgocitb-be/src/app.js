import express from 'express';
import cors from 'cors';
import morgan from 'morgan';
import cookieParser from 'cookie-parser';
import { config } from './config/index.js';
import routes from './routes/index.js';
import errorMiddleware from './middleware/error.middleware.js';
import logger from './utils/logger.js';
import helmet from 'helmet';

// Initialize express
const app = express();

// Middleware
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Security middleware
app.use(helmet({
    crossOriginOpenerPolicy: { policy: 'same-origin-allow-popups' }
  })
);

// Cookie parser middleware
app.use(cookieParser());


// CORS with credentials support for cookies
app.use(cors({
  origin: function(origin, callback) {
    callback(null, true);
  },
  credentials: true
}));

// Logging middleware
if (config.nodeEnv === 'development') {
  app.use(morgan('dev'));
} else {
  app.use(morgan('combined', {
    stream: {
      write: (message) => logger.info(message.trim())
    }
  }));
}

// API Routes
app.use('/api-GDGoC-ITB', routes);

// Home route
app.get('/', (req, res) => {
  res.json({ message: 'Welcome to LMS GDGoC ITB API' });
});

// 404 handler
app.use((req, res) => {
  res.status(404).json({ message: 'Route not found' });
});

// Error handling middleware
app.use(errorMiddleware);

export default app;