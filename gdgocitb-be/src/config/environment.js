import dotenv from 'dotenv';

dotenv.config();

export const config = {
  nodeEnv: process.env.NODE_ENV || 'development',
  port: process.env.PORT || 5000,
  mongodb: {
    uri: process.env.MONGODB_URI
  },
  jwt: {
    secret: process.env.JWT_SECRET || 'secret_dev_key_should_be_changed',
    expiresIn: process.env.JWT_EXPIRES_IN || '1d',
  },
  google: {
    clientId: process.env.GOOGLE_CLIENT_ID,
    clientSecret: process.env.GOOGLE_CLIENT_SECRET,
    callbackUrl: process.env.GOOGLE_CALLBACK_URL,
  },
  cloudinary: {
    cloudName: process.env.CLOUDINARY_CLOUD_NAME,
    apiKey: process.env.CLOUDINARY_API_KEY,
    apiSecret: process.env.CLOUDINARY_API_SECRET,
    folder: process.env.CLOUDINARY_FOLDER || 'lms-gdgoc-itb',
  },
  clientURL: process.env.CLIENT_URL,
};